package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/quduss/pokedexcli/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
}

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var commands map[string]cliCommand

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

// commandExit prints a goodbye message and exits the program
func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// commandHelp prints all registered commands and their descriptions
func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, args []string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.Next != nil {
		url = *cfg.Next
	}

	var parsed LocationAreaResponse
	if data, ok := cfg.Cache.Get(url); ok {

		err := json.Unmarshal(data, &parsed)
		if err != nil {
			return err
		}
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cfg.Cache.Add(url, body)

		err = json.Unmarshal(body, &parsed)
		if err != nil {
			return err
		}
	}

	for _, area := range parsed.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = parsed.Next
	cfg.Previous = parsed.Previous
	return nil
}

func commandMapBack(cfg *config, args []string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	var parsed LocationAreaResponse
	if data, ok := cfg.Cache.Get(*cfg.Previous); ok {

		err := json.Unmarshal(data, &parsed)
		if err != nil {
			return err
		}
	} else {
		resp, err := http.Get(*cfg.Previous)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cfg.Cache.Add(*cfg.Previous, body)

		err = json.Unmarshal(body, &parsed)
		if err != nil {
			return err
		}
	}

	for _, area := range parsed.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = parsed.Next
	cfg.Previous = parsed.Previous
	return nil
}

// Initialize command registry after all functions are declared
func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Explore the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back to the previous 20 location areas",
			callback:    commandMapBack,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := &config{}
	cfg.Cache = pokecache.NewCache(5 * time.Second)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		command := words[0]
		args := words[1:]

		if cmd, ok := commands[command]; ok {
			err := cmd.callback(cfg, args)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
