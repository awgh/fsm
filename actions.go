package fsm

import (
	"errors"
	"fmt"
	"strings"
)

// ActionHandler - function pointer for this action's handler
type ActionHandler func(...string) (string, error)

// Action - named handler
type Action struct {
	Name    string
	Handler ActionHandler
}

var actionRegistry map[string]Action
var logs map[string][]string

func init() {
	actionRegistry = make(map[string]Action)
	RegisterAction(Action{Name: "print", Handler: printHandler})

	// Log Functions
	logs = make(map[string][]string)
	RegisterAction(Action{Name: "newLog", Handler: newLogHandler})
	RegisterAction(Action{Name: "printAndLog", Handler: printAndLogHandler})
	RegisterAction(Action{Name: "log", Handler: logHandler})
	RegisterAction(Action{Name: "printLog", Handler: printLogHandler})
}

// RegisterAction - register an action
func RegisterAction(action Action) {
	actionRegistry[action.Name] = action
}

func printHandler(args ...string) (string, error) {
	fmt.Println(args)
	return strings.Join(args, " "), nil
}

func newLogHandler(args ...string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Not enough arguments for newLog, needs log name")
	}
	logs[args[0]] = []string{}
	return "", nil
}

func printAndLogHandler(args ...string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Not enough arguments for printAndLog, needs log name")
	}
	logs[args[0]] = append(logs[args[0]], args[1:]...)
	fmt.Println(args[1:])
	return strings.Join(args[1:], " "), nil
}

func logHandler(args ...string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Not enough arguments for log, needs log name")
	}
	logs[args[0]] = append(logs[args[0]], args[1:]...)
	return "", nil
}

func printLogHandler(args ...string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Not enough arguments for printLog, needs log name")
	}
	s := strings.Join(logs[args[0]], "\n")
	fmt.Println(s)
	return s, nil
}

// Eval - invoke an action with arguments
func (f *FSM) Eval(call Call, input string) (string, error) {

	action, ok := actionRegistry[call.Name]
	if !ok {
		return "", errors.New("Unknown Action: " + call.Name)
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
