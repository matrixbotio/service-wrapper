package controller

import (
	"fmt"
)

var statusCheckFunc func() (string, string, []File)

type Button struct {
	Text string
	RMKeyOnClick bool
	OnClick func()
}

type File struct {
	Name string
	Content []byte
}

func Send(text string, buttons []Button, files ...File) {
	fmt.Println(text)
}

func OnStatusCheck(check func() (string, string, []File)) {
	statusCheckFunc = check
}
