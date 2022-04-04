package main

import (
	"errors"
	"io"
	"os"
	"strings"
	"syscall"
	"time"

	"_/src/controller"

	"github.com/chebyrash/promise"
)

var timeFormat = "2006-01-02 15:04:05"
var stoppedFile = "wrapper-process-stopped"

func appendToLimitedArr(arr []string, str string, count int) []string {
	if len(arr) > count-1 {
		arr = arr[len(arr)-count+1:]
	}
	return append(arr, str)
}

func readData(readCloser io.ReadCloser, onNewLine func(string)) *promise.Promise {
	return promise.New(func(resolve func(v promise.Any), reject func(error)) {
		lastline := ""
		buf := make([]byte, 8)
		for {
			n, err := readCloser.Read(buf)
			lastline = lastline + string(buf[:n])
			if strings.Contains(lastline, "\n") {
				splitted := strings.Split(lastline, "\n")
				linesToAdd := splitted[:len(splitted)-1]
				for i := 0; i < len(linesToAdd); i++ {
					onNewLine(linesToAdd[i])
				}
				lastline = splitted[len(splitted)-1]
			}
			if err == io.EOF {
				if lastline != "" {
					onNewLine(lastline)
				}
				break
			}
		}
		resolve(nil)
	})
}

func filesStdoutStderr(stdoutbuf []string, stderrbuf []string) []controller.File {
	timestr := time.Now().Format(timeFormat)
	files := []controller.File{}
	stdout := strings.Join(stdoutbuf, "\n")
	stderr := strings.Join(stderrbuf, "\n")
	if len(stdout) > 0 {
		files = append(files, controller.File{
			Name:    processDef.ServiceName + " — " + timestr + ".stdout.txt",
			Content: strings.Join(stdoutbuf, "\n"),
		})
	}
	if len(stderr) > 0 {
		files = append(files, controller.File{
			Name:    processDef.ServiceName + " — " + timestr + ".stderr.txt",
			Content: strings.Join(stderrbuf, "\n"),
		})
	}
	return files
}

func filesLinebuf(linebuf []string) []controller.File {
	files := []controller.File{}
	output := strings.Join(linebuf, "\n")
	if len(output) > 0 {
		files = append(files, controller.File{
			Name:    processDef.ServiceName + " — " + time.Now().Format(timeFormat) + ".output.txt",
			Content: strings.Join(linebuf, "\n"),
		})
	}
	return files
}

func sendSig(sig os.Signal) {
	if cmd == nil {
		return
	}
	err := cmd.Process.Signal(sig)
	if err != nil {
		controller.Send("Error when sending signal to process: "+err.Error(), []controller.Button{
			{
				Text:          "👌 Ok",
				RMMsgsOnClick: true,
				OnClick:       func() {},
			},
		})
	}
}

func notifyExternalProcessesStopped() {
	_, err := os.OpenFile(stoppedFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		controller.Send("Error on healthcheck notify: "+err.Error(), []controller.Button{
			{
				Text:          "👌 Ok",
				RMMsgsOnClick: true,
				OnClick:       func() {},
			},
		})
	}
}

func notifyExternalProcessesStarted() {
	err := os.Remove(stoppedFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		controller.Send("Error on healthcheck notify: "+err.Error(), []controller.Button{
			{
				Text:          "👌 Ok",
				RMMsgsOnClick: true,
				OnClick:       func() {},
			},
		})
	}
}

func isProcessStoppedByWrapper() bool {
	_, err := os.Stat(stoppedFile)
	return errors.Is(err, os.ErrNotExist)
}

func stopProcess() {
	if !running {
		return
	}
	paused = true
	running = false
	notifyExternalProcessesStopped()
	sendSig(syscall.SIGSTOP)
}

func resumeProcess() {
	if !paused {
		return
	}
	paused = false
	running = true
	sendSig(syscall.SIGCONT)
}

func termProcess() {
	resumeProcess()
	notifyExternalProcessesStopped()
	sendSig(syscall.SIGTERM)
}

func killProcess() {
	resumeProcess()
	notifyExternalProcessesStopped()
	sendSig(syscall.SIGKILL)
}

func statusButtons() []controller.Button {
	if running {
		return []controller.Button{{
			RMMsgsOnClick: true,
			Text:          "⏸ Пауза",
			OnClick:       stopProcess,
		}, {
			RMMsgsOnClick: true,
			Text:          "⏹ Остановить",
			OnClick:       termProcess,
		}, {
			RMMsgsOnClick: true,
			Text:          "⛔️ Убить",
			OnClick:       killProcess,
		}}
	}
	if paused {
		return []controller.Button{{
			RMMsgsOnClick: true,
			Text:          "▶️ Продолжить",
			OnClick:       resumeProcess,
		}, {
			RMMsgsOnClick: true,
			Text:          "⏹ Остановить",
			OnClick:       termProcess,
		}, {
			RMMsgsOnClick: true,
			Text:          "⛔️ Убить",
			OnClick:       killProcess,
		}}
	}
	return []controller.Button{}
}

func readStdoutStderr(stdout io.ReadCloser, stderr io.ReadCloser) []controller.File {
	if processDef.SeparateStdoutStderr {
		stdoutbuf := make([]string, 0)
		stderrbuf := make([]string, 0)
		controller.OnStatusCheck(func() (string, string, []controller.Button, []controller.File) {
			return state, startTime, statusButtons(), filesStdoutStderr(stdoutbuf, stderrbuf)
		})
		promise.All(
			readData(stdout, func(line string) {
				os.Stdout.WriteString(line + "\n")
				stdoutbuf = appendToLimitedArr(stdoutbuf, line, processDef.LinesToPreserve)
			}),
			readData(stderr, func(line string) {
				os.Stderr.WriteString(line + "\n")
				stderrbuf = appendToLimitedArr(stderrbuf, line, processDef.LinesToPreserve)
			}),
		).Await()
		return filesStdoutStderr(stdoutbuf, stderrbuf)
	} else {
		linebuf := make([]string, 0)
		controller.OnStatusCheck(func() (string, string, []controller.Button, []controller.File) {
			return state, startTime, statusButtons(), filesLinebuf(linebuf)
		})
		promise.All(
			readData(stdout, func(line string) {
				os.Stdout.WriteString(line + "\n")
				linebuf = appendToLimitedArr(linebuf, line, processDef.LinesToPreserve)
			}),
			readData(stderr, func(line string) {
				os.Stderr.WriteString(line + "\n")
				linebuf = appendToLimitedArr(linebuf, line, processDef.LinesToPreserve)
			}),
		).Await()
		return filesLinebuf(linebuf)
	}
}
