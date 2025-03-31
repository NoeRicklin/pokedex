package main

import(
	"encoding/json"
	"time"
	"fmt"
	"strings"
	"bufio"
	"os"
	"github.com/NoeRicklin/pokedex/internal/pokeapi/utils"
	"github.com/NoeRicklin/pokedex/internal/pokecache"
)

type cliCommand struct {
	name		string
	desc		string
	callback	func() error
}

var commands map[string]cliCommand
var c *pokecache.Cache
var s *bufio.Scanner

func main() {
	commands = map[string]cliCommand{
		"exit": {
			name:		"exit",
			desc:		"Exits the REPL",
			callback:	cmdExit,
		},
		"help": {
			name:		"help",
			desc:		"Displays the help",
			callback:	cmdHelp,
		},
		"map": {
			name:		"map",
			desc:		"Shows next 20 locations",
			callback:	cmdMap,
		},
		"mapb": {
			name:		"mapb",
			desc:		"Shows prev 20 locations",
			callback:	cmdMapb,
		},
	}

	c = pokecache.NewCache(3 * time.Second)
	s = bufio.NewScanner(os.Stdin)

	var input string

	for {
		fmt.Print("Pokedex >")
		s.Scan()
		
		input = s.Text()
		command := cleanInput(input)[0]

		if _, ok := commands[command]; !ok {
			fmt.Println("Unknown Command")
			continue
		}

		if err := commands[command].callback(); err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	trimmed := strings.Trim(text, " ")

	words := strings.Split(trimmed, " ")

	var lowerWords []string
	for _, word := range words {
		lowerWords = append(lowerWords, strings.ToLower(word))
	}

	return lowerWords
}

func cmdExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cmdHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for cmd, _ := range commands {
		fmt.Printf("%s: %s\n", commands[cmd].name, commands[cmd].desc)
	}
	return nil
}

var count int
func cmdMap() error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?limit=20&offset=%d",
	count * 20)

	var body utils.AreaJSON
	if rawBody, cacheHit := c.Get(url); cacheHit {
		if err := json.Unmarshal(rawBody, &body); err != nil {
			return err
		}
	} else {
		var err error
		body, err = utils.GetURLBody(url)
		if err != nil { return err }

		rawBody, err := json.Marshal(body)
		if err != nil { return err }
		if err := c.Add(url, rawBody); err != nil { return err }
	}

	for _, location := range body.Results {
		fmt.Println(location.Name)
	}

	count++
	return nil
}

func cmdMapb() error {
	count--		// Cancels out last increase of count
	count--
	if count < 0 {
		return fmt.Errorf("You're on the first page")
	}

	return cmdMap()
}
