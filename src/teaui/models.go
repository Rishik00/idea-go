package teaui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
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

type Initmodel struct {
    title       textarea.Model
    description textarea.Model
    focusIndex  int
}

type Choicemodel struct {
	Cursor int
	Choice string
	Choices []string
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}