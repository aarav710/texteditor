package mhandler

import "texteditor/components"

type Del struct {
	CommandImpl
	motion Motion
}

func NewDel(text string, commandType commandType, editor *components.EditorModel) *Del {
	c := CommandImpl{commandType: commandType, text: text, editor: editor}
	d := Del{CommandImpl: c}
	return &d
}

func (d *Del) getCommandType() commandType {
	return d.commandType
}

func (d *Del) execute() {
}
