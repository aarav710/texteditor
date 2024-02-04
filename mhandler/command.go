package mhandler

import "texteditor/components"

type commandType string

const (
	find    = "FIND"
	del     = "DEL"
	yank    = "YANK"
	vselect = "VSELECT"
)

type CommandImpl struct {
	commandType commandType
	text        string
	editor      *components.EditorModel
}

type Move struct {
	row, col int
}

type Command interface {
	getCommandType() commandType
	execute()
}

func CommandFactory(text string, commandType commandType, editor *components.EditorModel) Command {
	if commandType == find {
		return NewFind(text, commandType, editor)
	} else if commandType == del {
		return NewDel(text, commandType, editor)
	} else if commandType == yank {
		return NewYank(text, commandType, editor)
	} else {
		return NewVselect(text, commandType, editor)
	}
}

func ParseInput() {
}
