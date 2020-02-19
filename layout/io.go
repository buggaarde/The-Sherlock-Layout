package layout

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func Load(tomlfile string) Layout {
	var keyboard Layout
	if _, err := toml.DecodeFile(tomlfile, &keyboard); err != nil {
		log.Fatal(err)
	}

	keyboard.PositionKey = make(map[int]string)

	for key, idx := range keyboard.KeyPositions {
		keyboard.PositionKey[idx] = key
	}

	keyboard.r = rand.New(rand.NewSource(time.Now().UnixNano()))

	return keyboard
}

func Write(l Layout, file string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	f.WriteString("adjacency = [\n")
	for _, adj := range l.Adjacency {
		f.WriteString(strings.ReplaceAll(fmt.Sprintf("\t%v,\n", adj), " ", ", "))
	}
	f.WriteString("]\n\n")

	f.WriteString("[key_positions]\n")
	for i := 0; i < len(l.Adjacency); i++ {
		letter := l.PositionKey[i]
		if letter == "." || letter == "," || letter == "'" || letter == "?" {
			letter = fmt.Sprintf("\"%s\"", letter)
		}
		if letter == "æ" || letter == "ø" || letter == "å" {
			letter = fmt.Sprintf("\"%s\"", letter)
		}
		if letter == "\"" {
			letter = fmt.Sprintf("'\"'")
		}
		f.WriteString(fmt.Sprintf("%s = %d\n", letter, i))
	}
	f.WriteString("\n")

	f.WriteString("[position_hand]\n")
	for i := 0; i < len(l.Adjacency); i++ {
		letter := l.PositionKey[i]
		hand := l.PositionHand[letter]
		if letter == "." || letter == "," || letter == "'" || letter == "?" {
			letter = fmt.Sprintf("\"%s\"", letter)
		}
		if letter == "æ" || letter == "ø" || letter == "å" {
			letter = fmt.Sprintf("\"%s\"", letter)
		}
		if letter == "\"" {

			letter = fmt.Sprintf("'\"'")
		}
		f.WriteString(fmt.Sprintf("%s = %d\n", letter, hand))
	}
	f.WriteString("\n")

	f.WriteString("[convenience.singles]\n")
	for key, val := range l.Convenience.SinglePress {
		f.WriteString(fmt.Sprintf("%s = %.1f\n", key, val))
	}
	f.WriteString("\n")

	f.WriteString("[convenience.doubles]\n")
	for key, val := range l.Convenience.DoublePress {
		f.WriteString(fmt.Sprintf("\"%s\" = %.1f\n", key, val))
	}
	f.WriteString("\n")

	f.WriteString("[convenience.triples]\n")
	for key, val := range l.Convenience.TriplePress {
		f.WriteString(fmt.Sprintf("\"%s\" = %.1f\n", key, val))
	}
	f.WriteString("\n")
}

func Print(l Layout) {
	for i := 0; i < 5; i++ {
		fmt.Printf(" %s ", l.PositionKey[i])
	}
	fmt.Printf(" %s     %s ", l.PositionKey[30], l.PositionKey[32])
	for i := 5; i < 10; i++ {
		fmt.Printf(" %s ", l.PositionKey[i])
	}
	fmt.Println()

	for i := 10; i < 15; i++ {
		fmt.Printf(" %s ", l.PositionKey[i])
	}
	fmt.Printf("         ")
	for i := 15; i < 20; i++ {
		fmt.Printf(" %s ", l.PositionKey[i])
	}
	fmt.Println()

	for i := 20; i < 25; i++ {
		fmt.Printf(" %s ", l.PositionKey[i])
	}
	fmt.Printf(" %s     %s ", l.PositionKey[31], l.PositionKey[33])
	for i := 25; i < 30; i++ {
		fmt.Printf(" %s ", l.PositionKey[i])
	}
	fmt.Println()
	fmt.Println()
}
