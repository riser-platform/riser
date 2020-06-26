package ui

import "io"

type FakeView struct {
	RenderHumanCalled bool
	RenderJsonCalled  bool
}

func (fake *FakeView) RenderHuman(writer io.Writer) error {
	fake.RenderHumanCalled = true
	return nil
}

func (fake *FakeView) RenderJson(writer io.Writer) error {
	fake.RenderJsonCalled = true
	return nil
}
