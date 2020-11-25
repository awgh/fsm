package main

import (
	"fmt"
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

	if err := fsm.Handle("I am"); err != nil {
		fmt.Println(err)
	}
	//fsm.Handle("hi")
	if err := fsm.Handle("23"); err != nil {
		fmt.Println(err)
	}
}
