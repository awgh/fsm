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
	Transitions []Transition `yaml:"transitions"`
}

// Transition - single entry in the transition table
type Transition struct {
	Source string `yaml:"src"`
	Dest   string `yaml:"dst"`

	On []string `yaml:"on"`
	Do []string `yaml:"do"`
}

// State - single state
type State struct {
	Name       string
	Classes    []Class
	Classifier *bayesian.Classifier
}

// Class - named set of stings for classifier
type Class struct {
	Name   string
	Values []string
}