package text

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type SingleFrequency = map[string]int
type DoubleFrequency = map[[2]string]int
type TripletFrequency = map[[3]string]int
type Frequencies struct {
	S SingleFrequency
	D DoubleFrequency
	T TripletFrequency
}

func AnalyseFile(file string) Frequencies {
	char0 := " "
	char1 := " "
	singleFrequency := make(SingleFrequency)
	doubleFrequency := make(DoubleFrequency)
	tripletFrequency := make(TripletFrequency)
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rr := bufio.NewReader(f)

	for {
		r, size, err := rr.ReadRune()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("couldn't read file")
		}

		characterIsTooFancy := size > 1
		if characterIsTooFancy {
			continue
		}
		char2 := strings.ToLower(string(r))

		singleFrequency[char2]++
		doubleFrequency[[2]string{char1, char2}]++
		tripletFrequency[[3]string{char0, char1, char2}]++
		char0 = char1
		char1 = char2
	}

	return Frequencies{singleFrequency, doubleFrequency, tripletFrequency}
}

func AnalyseDir(dir string) Frequencies {
	singleFrequency := make(SingleFrequency)
	doubleFrequency := make(DoubleFrequency)
	tripletFrequency := make(TripletFrequency)

	d, err := os.Open(dir)
	if err != nil {
		panic("couldn't open the directory")
	}

	info, err := d.Readdir(-1)
	d.Close()
	if err != nil {
		panic("error reading dir")
	}

	freqChan := make(chan Frequencies, len(info))

	for _, file := range info {
		go func() {
			freqChan <- AnalyseFile(fmt.Sprintf("%s/%s", dir, file.Name()))
		}()
		time.Sleep(1 * time.Millisecond) // too many open files if not
	}

	counter := 0
	for f := range freqChan {
		for s, freq := range f.S {
			singleFrequency[s] += freq
		}

		for d, freq := range f.D {
			doubleFrequency[d] += freq
		}

		for t, freq := range f.T {
			tripletFrequency[t] += freq
		}

		counter++
		if counter >= len(info) {
			close(freqChan)
		}
	}

	return Frequencies{singleFrequency, doubleFrequency, tripletFrequency}
}
