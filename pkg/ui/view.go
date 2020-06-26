package ui

import (
	"io"
)

type View interface {
	RenderHuman(io.Writer) error
	RenderJson(io.Writer) error
}
