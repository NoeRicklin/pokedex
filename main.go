package main

import(
	"fmt"
	"strings"
	"bufio"
	"os"
	"time"
	"math/rand"
	"github.com/NoeRicklin/pokedex/internal/pokeapi/utils"
)

type cliCommand struct {
	name		string
	desc		string
	callback	func(...string) error
}

var pokedex map[string]Pokemon
var commands map[string]cliCommand
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
		"explore": {
			name:		"explore",
			desc:		"Shows Pokémon in area",
			callback:	cmdExplore,
		},
		"catch": {
			name:		"catch",
			desc:		"Attempt to catch a Pokémon",
			callback:	cmdCatch,
		},
		"inspect": {
			name:		"inspect",
			desc:		"Show misc. information about a caught Pokémon",
			callback:	cmdInspect,
		},
		"pokedex": {
			name:		"pokedex",
			desc:		"List caught Pokémon",
			callback:	cmdPokedex,
		},
	}

	pokedex = map[string]Pokemon{}
	s = bufio.NewScanner(os.Stdin)
	utils.SetupCache(3 * time.Second)

	var input string

	for {
		fmt.Print("Pokedex >")
		s.Scan()
		
		input = s.Text()
		cmds := cleanInput(input)

		command, ok := commands[cmds[0]]
		if !ok {
			fmt.Println("Unknown Command")
			continue
		}

		if err := command.callback(cmds[1:]...); err != nil {
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

func cmdExit(argv ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cmdHelp(argv ...string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for cmd, _ := range commands {
		fmt.Printf("%s: %s\n", commands[cmd].name, commands[cmd].desc)
	}
	return nil
}

var count int
func cmdMap(argv ...string) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?limit=20&offset=%d",
	count * 20)

	area, err := utils.GetURLBody[AreaJSON](url)
	if err != nil { return err }

	for _, location := range area.Results {
		fmt.Println(location.Name)
	}

	count++
	return nil
}

func cmdMapb(argv ...string) error {
	count--		// Cancels out prev increase of count
	count--
	if count < 0 {
		return fmt.Errorf("You're on the first page")
	}

	return cmdMap()
}

func cmdExplore(argv ...string) error {
	if len(argv) < 1 {
		return fmt.Errorf("Usage: explore <location>")
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s",
	argv[0])

	location, err := utils.GetURLBody[LocationJSON](url)
	if err != nil { return err }

	for _, encounter := range location.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}

func cmdCatch(argv ...string) error {
	if len(argv) < 1 {
		return fmt.Errorf("Usage: catch <name of Pokémon>")
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s",
	argv[0])

	pokemon, err := utils.GetURLBody[Pokemon](url)
	if err != nil { return err }

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	rng := rand.Float64()
	// Highest BE is 608 => if rng >= 608 / 800 every pokemon will be caught
	if float64(pokemon.BaseExperience) / 800.0 <= rng {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func cmdInspect(argv ...string) error {
	if len(argv) < 1 {
		return fmt.Errorf("Usage: inspect <name of Pokéemon")
	}

	pokemon, ok := pokedex[argv[0]]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}

	name		:= pokemon.Name
	height		:= pokemon.Height
	weight		:= pokemon.Weight

	stats		:= pokemon.Stats
	types		:= pokemon.Types

	fmt.Printf(`
Name: %s
Height: %d
Weight: %d
`,
	name,
	height,
	weight)

	fmt.Println("Stats:")
	for _, s := range stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}
	fmt.Println()

	return nil
}

func cmdPokedex(argv ...string) error {
	fmt.Println("Your Pokedex:")

	for name := range pokedex {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}
