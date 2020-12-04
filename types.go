package fsm

import "github.com/navossoc/bayesian"

// FSM - finite state machine struct
type FSM struct {
	Transitions *TransitionTable  // all transitions
	States      map[string]*State // all states
	State       *State            // current state
}

// TransitionTable - main structure of FSM
type TransitionTable struct {
	Imports     []string     `yaml:"imports"`
	Transitions []Transition `yaml:"transitions"`
}

// Transition - single entry in the transition table
type Transition struct {
	Source string `yaml:"src"`
	Dest   string `yaml:"dst"`

	Auto bool     `yaml:"auto"`
	On   []string `yaml:"on"`
	Do   []Call   `yaml:"do"`
	Once []Call   `yaml:"once"`
}

// State - single state
type State struct {
	Name       string
	Classes    []Class
	Classifier *bayesian.Classifier

	enteredAtLeastOnce bool
}

// Class - named set of stings for classifier
type Class struct {
	Name   string
	Values []string
}

// Call - call to a handler, name of handler and args
type Call struct {
	Name string   `yaml:"fn"`
	Args []string `yaml:"args"`
}
