package ui

import "io"

// MessageView provides a view with a human message accompanied by viewmodel for structured data output (e.g. JSON)
type MessageView struct {
	Message string
	Model   interface{}
}

func NewMessageView(message string, model interface{}) *MessageView {
	return &MessageView{message, model}
}

func (view *MessageView) RenderHuman(writer io.Writer) error {
	_, err := writer.Write([]byte(view.Message + "\n"))
	return err
}

func (view *MessageView) RenderJson(writer io.Writer) error {
	return RenderJson(view.Model, writer)
}
