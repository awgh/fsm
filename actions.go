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
func (f *FSM) Eval(call Call, input string) error {

	action, ok := actionRegistry[call.Name]
	if !ok {
		return errors.New("Unknown Action: " + call.Name)
	}
	var args []string
	for _, v := range call.Args {
		if v == "$$" {
			args = append(args, input)
		} else {
			args = append(args, v)
		}
	}
	return action.Handler(args...)
}
