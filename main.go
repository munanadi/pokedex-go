package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/munanadi/pokedex/pokecache"
	"github.com/munanadi/pokedex/pokehelp"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *pokehelp.RequestConfig, args ...[]string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    CommandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the pokedex",
			callback:    CommandExit,
		},
		"map": {
			name:        "map",
			description: "Lets you explore the map in skips of 20",
			callback:    CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "To go back 20 skips in map locations",
			callback:    CommandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Let's you explore a city area",
			callback:    CommandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Let's you catch a Pokemon",
			callback:    CommandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Let's you check on the Pokemon",
			callback:    CommandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Let's check your pokedex",
			callback:    CommandPokedex,
		},
	}
}

func main() {

	// Store values in cache for seconds specified here
	const CACHE_REFRESH_IN_SECONDS int64 = 30

	timeInterval := time.Duration(CACHE_REFRESH_IN_SECONDS) * time.Second
	cache := pokecache.NewCache(timeInterval)

	var pokedex = make(map[string]pokehelp.Pokemon)

	config := &pokehelp.RequestConfig{Next: nil, Prev: nil, Cache: cache, Pokedex: pokedex}

	for {
		fmt.Printf("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		if ok := scanner.Scan(); ok {
			for k, v := range getCommands() {
				args := strings.Split(scanner.Text(), " ")
				if k == args[0] {
					v.callback(config, args[1:])
				}
			}
		}
	}
}
