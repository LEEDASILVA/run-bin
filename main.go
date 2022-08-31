package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func readOutput(output chan string, rc io.ReadCloser) {
	r := bufio.NewReader(rc)
	for {
		x, _ := r.ReadString('\n')
		output <- string(x)
	}
}

func main() {
	if len(os.Args) == 1 {
		handleError(errors.New(`to run the program you need to give the following arguments :
1. name of the command you want to execute
2. the rest of the arguments are related to the command passed in the first position`))
	}

	commandName := os.Args[1]
	fmt.Println("executing command: ", commandName)
	cmd := exec.Command(commandName, os.Args[2:]...)
	stdout, err := cmd.StdoutPipe()
	handleError(err)

	stdin, err := cmd.StdinPipe()
	handleError(err)

	err = cmd.Start()
	handleError(err)

	go func() {
		io.WriteString(stdin, "writing something in to the program\n")
	}()

	output := make(chan string)
	defer close(output)
	go readOutput(output, stdout)

	for o := range output {
		fmt.Printf("stdout from %s:%s\n", commandName, o)
	}
}
