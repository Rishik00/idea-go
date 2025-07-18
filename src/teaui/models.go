package teaui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

type (
	errMsg error
)

type TitleModel struct {
	input 		textinput.Model
	err			error
}

type Descriptionmodel struct {
	textarea textarea.Model
	ActualString string
	err      error
}

type Choicemodel struct {
	Cursor int
	Choice string
	Choices []string
}

