package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/awgh/fsm"
)

func main() {
	var yamlFile string
	flag.StringVar(&yamlFile, "f", "fsm.yml", "YAML FSM Specification")

	reader := bufio.NewReader(os.Stdin)
	cwd, _ := os.Getwd()
	sm := fsm.Load(filepath.Join(cwd, yamlFile))

	//log.Printf("%+v\n", sm.Transitions)

	for {
		fmt.Print("\n>")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		if fsm.Normalize([]string{line})[0] == "quit" {
			break
		}
		result, err := sm.Handle(line)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	}
}
