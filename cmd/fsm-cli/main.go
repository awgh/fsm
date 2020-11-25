package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/awgh/fsm"
)

func main() {
	//reader := bufio.NewReader(os.Stdin)
	cwd, _ := os.Getwd()
	fsm := fsm.Load(filepath.Join(cwd, "fsm.yml"))

	log.Printf("%+v\n", fsm.Transitions)

	fsm.Handle("hi")
	fsm.Handle("Answer 1")
	fsm.Handle("Answer 2")
	fsm.Handle("Answer 3")
	/*
		if err := fsm.Handle("I am"); err != nil {
			fmt.Println(err)
		}

		if err := fsm.Handle("23"); err != nil {
			fmt.Println(err)
		}
	*/
}
