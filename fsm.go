package fsm

import (
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/navossoc/bayesian"
	"gopkg.in/yaml.v2"
)

var normalizeRegexp *regexp.Regexp

func init() {
	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		log.Fatal(err)
	}
	normalizeRegexp = reg
}

// Load - makes a new finite state machine from the given config file
func Load(path string) *FSM {
	table := loadFile(path)
	return New(table)
}

func loadFile(path string) *TransitionTable {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var table TransitionTable
	if err := yaml.Unmarshal(b, &table); err != nil {
		panic(err)
	}

	for _, v := range table.Imports {
		t := loadFile(v)
		table.Transitions = append(table.Transitions, t.Transitions...)
	}
	return &table
}

// New - makes a new finite state machine from the given config
func New(table *TransitionTable) *FSM {

	states := make(map[string]*State)
	states["$start"] = &State{Name: "$start"}
	for _, transition := range table.Transitions {
		srcState, ok := states[transition.Source]
		if !ok {
			srcState = &State{Name: transition.Source}
			states[transition.Source] = srcState
		}
		dstState, ok := states[transition.Dest]
		if !ok {
			dstState = &State{Name: transition.Dest}
			states[transition.Dest] = dstState
		}

		srcState.Classes = append(srcState.Classes, Class{
			Name:   transition.Dest,
			Values: transition.On,
		})
	}

	for _, state := range states {
		// classifier is only if there is more than one option,
		// otherwise you know which state you're going to on any input
		if len(state.Classes) > 1 {
			var classes []bayesian.Class
			for _, class := range state.Classes {
				classes = append(classes, bayesian.Class(class.Name))
			}
			classifier := bayesian.NewClassifier(classes...)
			for _, class := range state.Classes {
				classifier.Learn(normalize(class.Values), bayesian.Class(class.Name))
			}
			log.Println("Loaded commands for classifier:")
			for i, c := range classes {
				log.Printf("%v\t%v\n", i, c)
			}
			//classifier.ConvertTermsFreqToTfIdf()
			state.Classifier = classifier
		} else {
			state.Classifier = nil
		}
	}
	fsm := &FSM{Transitions: table, States: states}
	fsm.Transition("$start", "")
	return fsm
}

// Handle - process input and transition
func (f *FSM) Handle(input string) (string, error) {

	if f.State.Classifier != nil {
		// todo: do I want underflow checking here or not?
		probs, likely, _ := f.State.Classifier.ProbScores(normalize([]string{input}))
		log.Printf("prob scores: %+v %+v\n", probs, likely)

		if probs[likely] <= 1.0/float64(len(probs)) {
			// no real winner
			return "", errors.New("Not sure what you're trying to say")
		}

		return f.Transition(f.State.Classes[likely].Name, input)
	}
	if len(f.State.Classes) > 0 {
		// only one place to go
		return f.Transition(f.State.Classes[0].Name, input)
	}
	// no place to go
	return f.Transition("$start", input)
}

// Transition - transition to new state
func (f *FSM) Transition(newState string, input string) (string, error) {
	prevState := ""
	if f.State != nil {
		prevState = f.State.Name
	}
	if _, ok := f.States[newState]; !ok {
		return "", errors.New("Unknown State, can't transition")
	}
	log.Println("Transitioning to ", newState)
	f.State = f.States[newState]

	// Locate and run actions (do and once)
	for _, t := range f.Transitions.Transitions {
		if (t.Source == "" || t.Source == prevState) &&
			(t.Dest == "" || t.Dest == newState) {

			// run-once actions, if state has not been entered before
			if !f.State.enteredAtLeastOnce {
				for _, s := range t.Once {
					f.Eval(s, input)
				}
			}
			for _, s := range t.Do {
				f.Eval(s, input)
			}
		}
	}
	f.State.enteredAtLeastOnce = true

	return "", nil
}

func normalize(texts []string) []string {
	//todo: handle unicode normalization
	var retval []string
	for _, v := range texts {
		retval = append(retval, strings.ToLower(normalizeRegexp.ReplaceAllString(v, "")))
	}
	return retval
}
