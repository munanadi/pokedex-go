package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the pokedex",
			callback:    commandExit,
		},
	}
}

func main() {

	for {
		fmt.Printf("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		if ok := scanner.Scan(); ok {
			for k, v := range getCommands() {
				if k == scanner.Text() {
					v.callback()
				}
			}
		}
	}
}

func commandExit() error {
	os.Exit(0)
	return errors.New("something went wrong in exit")
}

func commandHelp() error {
	fmt.Printf(`Welcome to Pokedex
  Usage:
  
  help: Displays a help message
  exit: Exit the pokedex
`)
	return errors.New("something went wrong in help")
}
