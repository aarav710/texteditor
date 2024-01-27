package textctrl

// will need to implement decorator pattern for this

type Decorator interface {
}

type DecoratorImpl struct {
	component *Decorator
}

type Handler struct {
	currMotion string
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) AddToCurrMotion(cmd string) {
	if cmd == " " {
		return
	}
	h.currMotion += cmd
}

func (h *Handler) IsValidMotion() bool {
	return false
}

func (h *Handler) Clear() {
	h.currMotion = ""
}

func (h *Handler) ExecuteMotion() {
	if h.currMotion == "" || h.currMotion == " " {
		return
	}

}
