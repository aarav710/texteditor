package mhandler

import "texteditor/components"

type Find struct {
	CommandImpl
}

func NewFind(text string, commandType commandType, editor *components.EditorModel) *Find {
	c := CommandImpl{commandType: commandType, text: text, editor: editor}
	find := Find{CommandImpl: c}
	return &find
}

func (f *Find) getCommandType() commandType {
	return f.commandType
}

func (f *Find) execute() {
}
