package main

import (
	"fmt"
	"os/exec"
	"time"

	"_/src/controller"
	"_/src/servicedef"

	"github.com/chebyrash/promise"
)

var state string
var startTime = ""

func awaitRestartCommand(text string, files ...controller.File) {
	promise.New(func(resolve func(v promise.Any), reject func(error)) {
		controller.Send(text, []controller.Button{
			{
				Text:         "🔄 Рестарт",
				RMMsgsOnClick: true,
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
		state = fmt.Sprintf("Не удаётся передать stdout: %v", stdoutPipeErr)
		awaitRestartCommand(state)
		return
	}
	stderr, stderrPipeErr := cmd.StderrPipe()
	if stderrPipeErr != nil {
		state = fmt.Sprintf("Не удаётся передать stderr: %v", stderrPipeErr)
		awaitRestartCommand(state)
		return
	}
	err := cmd.Start()
	if err != nil {
		state = fmt.Sprintf("Не удаётся запустить процесс сервиса: %v", err)
		awaitRestartCommand(state)
		return
	}
	state = "Запущен"
	startTime = time.Now().Format(timeFormat)
	files := readStdoutStderr(stdout, stderr)
	waitErr := cmd.Wait()
	if waitErr != nil {
		startTime = ""
		state = fmt.Sprintf("Процесс сервиса крашнулся: %v", waitErr)
		awaitRestartCommand(state, files...)
		return
	}
	startTime = ""
	state = "Процесс сервиса успешно завершился. Этого не должно было случиться. Проверте вывод сервиса"
	awaitRestartCommand(state, files...)
}

func main() {
	for{
		state = "Запуск..."
		controller.OnStatusCheck(func() (string, string, []controller.Button, []controller.File) {
			return state, startTime, []controller.Button{}, []controller.File{}
		})
		run()
	}
}
