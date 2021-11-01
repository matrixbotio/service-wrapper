package main

import (
	"io"
	"os"
	"strings"
	"time"

	"_/src/controller"
	"_/src/servicedef"

	"github.com/chebyrash/promise"
)

var timeFormat = "2006-01-02 15:04:05"

func appendToLimitedArr(arr []string, str string, count int) []string {
	if len(arr) > count - 1 {
		arr = arr[len(arr) - count + 1:]
	}
	return append(arr, str)
}

func readData(readCloser io.ReadCloser, onNewLine func(string)) *promise.Promise {
	return promise.New(func(resolve func(v promise.Any), reject func(error)){
		lastline := ""
		buf := make([]byte, 8)
		for {
			n, err := readCloser.Read(buf)
			lastline = lastline + string(buf[:n])
			if strings.Contains(lastline, "\n") {
				splitted := strings.Split(lastline, "\n")
				linesToAdd := splitted[:len(splitted) - 1]
				for i := 0; i < len(linesToAdd); i++ {
					onNewLine(linesToAdd[i])
				}
				lastline = splitted[len(splitted) - 1]
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
			Name: servicedef.ServiceName + " — " + timestr + ".stdout.txt",
			Content: strings.Join(stdoutbuf, "\n"),
		})
	}
	if len(stderr) > 0 {
		files = append(files, controller.File{
			Name: servicedef.ServiceName + " — " + timestr + ".stderr.txt",
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
			Name: servicedef.ServiceName + " — " + time.Now().Format(timeFormat) + ".output.txt",
			Content: strings.Join(linebuf, "\n"),
		})
	}
	return files
}

func readStdoutStderr(stdout io.ReadCloser, stderr io.ReadCloser) []controller.File {
	if servicedef.SeparateStdoutStderr {
		stdoutbuf := make([]string, 0)
		stderrbuf := make([]string, 0)
		controller.OnStatusCheck(func() (string, string, []controller.File) {
			return state, startTime, filesStdoutStderr(stdoutbuf, stderrbuf)
		})
		promise.All(
			readData(stdout, func(line string) {
				os.Stdout.WriteString(line + "\n")
				stdoutbuf = appendToLimitedArr(stdoutbuf, line, servicedef.LinesToPreserve)
			}),
			readData(stderr, func(line string) {
				os.Stderr.WriteString(line + "\n")
				stderrbuf = appendToLimitedArr(stderrbuf, line, servicedef.LinesToPreserve)
			}),
		).Await()
		return filesStdoutStderr(stdoutbuf, stderrbuf)
	} else {
		linebuf := make([]string, 0)
		controller.OnStatusCheck(func() (string, string, []controller.File) {
			return state, startTime, filesLinebuf(linebuf)
		})
		promise.All(
			readData(stdout, func(line string) {
				os.Stdout.WriteString(line + "\n")
				linebuf = appendToLimitedArr(linebuf, line, servicedef.LinesToPreserve)
			}),
			readData(stderr, func(line string) {
				os.Stderr.WriteString(line + "\n")
				linebuf = appendToLimitedArr(linebuf, line, servicedef.LinesToPreserve)
			}),
		).Await()
		return filesLinebuf(linebuf)
	}
}
