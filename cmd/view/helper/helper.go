package helper

import (
	"math/rand"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

type ErrMsg struct {
	err error
}

func Err(err error) ErrMsg {
	return ErrMsg{err: err}
}

func (e *ErrMsg) Error() string {
	return e.err.Error()
}

type IdModel interface {
	tea.Model
	GetId() int
}

type Id struct {
	Id int
}

func (m Id) GetId() int {
	return m.Id
}

func NewId() Id {
	return Id{Id: rand.Intn(1000000)}
}

// Internal ID management. Used during animating to ensure that frame messages
// are received only by spinner components that sent them.
var (
	lastID int
	idMtx  sync.Mutex
)

// Return the next ID we should use on the Model.
func NextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}
