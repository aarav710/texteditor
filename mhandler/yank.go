package mhandler

import "texteditor/components"

type Yank struct {
	CommandImpl
	motion Motion
}

func NewYank(text string, commandType commandType, editor *components.EditorModel) *Yank {
	c := CommandImpl{commandType: commandType, text: text, editor: editor}
	y := Yank{CommandImpl: c}
	return &y
}

func (y *Yank) getCommandType() commandType {
	return y.commandType
}

func (y *Yank) execute() {
}
