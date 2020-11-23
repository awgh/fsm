package fsm

import (
	"errors"
	"fmt"
)

// ActionHandler - function pointer for this action's handler
type ActionHandler func(...string) error

// Action - named handler
type Action struct {
	Name    string
	Handler ActionHandler
}

var actionRegistry map[string]Action

func init() {
	actionRegistry = make(map[string]Action)
	RegisterAction(Action{Name: "print", Handler: printHandler})
}

// RegisterAction - register an action
func RegisterAction(action Action) {
	actionRegistry[action.Name] = action
}

func printHandler(args ...string) error {
	fmt.Println(args)
	return nil
}

// Eval - invoke an action with arguments
func (f *FSM) Eval(call Call) error {

	action, ok := actionRegistry[call.Name]
	if !ok {
		return errors.New("Unknown Action: " + call.Name)
	}
	return action.Handler(call.Args...)
}
