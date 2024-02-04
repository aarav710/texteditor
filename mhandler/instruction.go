package mhandler

type Instruction struct {
	motion  *Motion
	command *Command
}

func NewInstruction(motion *Motion, command *Command) *Instruction {
	i := Instruction{motion: motion, command: command}
	return &i
}
