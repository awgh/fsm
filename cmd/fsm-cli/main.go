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

	//fsm.Handle("hi")
	fsm.Handle("I am")

}
