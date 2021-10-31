package main

import (
	"fmt"
	"os/exec"

	"_/src/controller"
	"_/src/servicedef"

	"github.com/chebyrash/promise"
)

func awaitRestartCommand(text string, files ...controller.File) {
	promise.New(func(resolve func(v promise.Any), reject func(error)) {
		controller.Send(text, []controller.Button{
			{
				Text:         "🔄 Рестарт",
				RMKeyOnClick: true,
				OnClick: func() {
					resolve(nil)
				},
			},
		}, files...)
	}).Await()
}

func run() {
	cmd := exec.Command(servicedef.Command, servicedef.Args...)
	stdout, stdoutPipeErr := cmd.StdoutPipe()
	if stdoutPipeErr != nil {
		awaitRestartCommand(fmt.Sprintf("Не удаётся передать stdout: %v", stdoutPipeErr))
		return
	}
	stderr, stderrPipeErr := cmd.StderrPipe()
	if stderrPipeErr != nil {
		awaitRestartCommand(fmt.Sprintf("Не удаётся передать stderr: %v", stderrPipeErr))
		return
	}
	err := cmd.Start()
	if err != nil {
		awaitRestartCommand(fmt.Sprintf("Не удаётся запустить процесс сервиса: %v", err))
		return
	}
    files := readStdoutStderr(stdout, stderr)
	waitErr := cmd.Wait()
	if waitErr != nil {
		awaitRestartCommand(fmt.Sprintf("Процесс сервиса крашнулся: %v", waitErr), files...)
		return
	}
    awaitRestartCommand("Процесс сервиса успешно завершился. Этого не должно было случиться. Проверте вывод сервиса", files...)
}

func main() {
    for{
        run()
    }
}
