package controller

import (
	"fmt"
)

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
