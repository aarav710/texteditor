package mhandler

import "texteditor/components"

type Vselect struct {
	CommandImpl
	motion Motion
}

func NewVselect(text string, commandType commandType, editor *components.EditorModel) *Vselect {
	c := CommandImpl{commandType: commandType, text: text, editor: editor}
	vs := Vselect{CommandImpl: c}
	return &vs
}

func (vs *Vselect) getCommandType() commandType {
	return vs.commandType
}

func (vs *Vselect) execute() {
}
