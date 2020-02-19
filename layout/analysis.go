package layout

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/buggaarde/layout-optimizer/text"
)

func AnalyseLayout(l Layout, f text.Frequencies) {
	fmt.Println("-- Total score --")
	fmt.Printf("   %.2f\n\n", Score(l, f.T, f.D, f.S))

	fmt.Println("-- Single press load --")

	leftHandPresses := 0.0
	rightHandPresses := 0.0

	for key, hand := range l.PositionHand {
		if hand == 0 {
			leftHandPresses += float64(f.S[key])
		}
		if hand == 1 {
			rightHandPresses += float64(f.S[key])
		}
	}

	total := leftHandPresses + rightHandPresses

	leftPercentage := 100 * leftHandPresses / total
	rightPercentage := 100 * rightHandPresses / total
	fmt.Printf("Left hand    Right hand\n  %.1f%%         %.1f%%\n\n", leftPercentage, rightPercentage)

	fmt.Println("- Single hand doubles -")
	leftHandDoubles := 0.0
	rightHandDoubles := 0.0

	for keys, freq := range f.D {
		h0, ok := l.PositionHand[keys[0]]
		if !ok {
			continue
		}

		h1, ok := l.PositionHand[keys[1]]
		if !ok {
			continue
		}
		if h0 == h1 && h0 == 0 {
			leftHandDoubles += float64(freq)
		}
		if h0 == h1 && h0 == 1 {
			rightHandDoubles += float64(freq)
		}
	}

	total = leftHandDoubles + rightHandDoubles
	leftPercentage = 100 * leftHandDoubles / total
	rightPercentage = 100 * rightHandDoubles / total
	fmt.Printf("Left hand    Right hand\n  %.1f%%         %.1f%%\n\n", leftPercentage, rightPercentage)

	fmt.Println("- Single hand doubles -")
	fmt.Println("-- with alternations --")
	leftHandAlternationDoubles := 0.0
	rightHandAlternationDoubles := 0.0

	for keys, freq := range f.T {
		h0, ok := l.PositionHand[keys[0]]
		if !ok {
			continue
		}

		h1, ok := l.PositionHand[keys[1]]
		if !ok {
			continue
		}

		h2, ok := l.PositionHand[keys[2]]
		if !ok {
			continue
		}

		firstAlternation := h0 != h1
		lastAlternation := h1 != h2
		alternationDouble := firstAlternation && lastAlternation

		if alternationDouble && h0 == 0 {
			leftHandAlternationDoubles += float64(freq)
		}

		if alternationDouble && h1 == 1 {
			rightHandAlternationDoubles += float64(freq)
		}
	}

	total = leftHandDoubles + leftHandAlternationDoubles + rightHandDoubles + rightHandAlternationDoubles
	leftPercentage = 100 * (leftHandDoubles + leftHandAlternationDoubles) / total
	rightPercentage = 100 * (rightHandDoubles + rightHandAlternationDoubles) / total
	fmt.Printf("Left hand    Right hand\n  %.1f%%         %.1f%%\n\n", leftPercentage, rightPercentage)
}

func singleHandUtilizationFile(l Layout, file string) (int, float64, int, float64) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rr := bufio.NewReader(f)

	var alternationStreaks, sameHandStreaks []float64
	var alternationStreak, sameHandStreak float64
	lastHand := 0
	for {
		r, size, err := rr.ReadRune()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		characterIsTooFancy := size > 1
		if characterIsTooFancy {
			continue
		}

		if hand, ok := l.PositionHand[strings.ToLower(string(r))]; !ok {
			continue
		} else {
			if hand != lastHand { // alternation
				alternationStreak++
				lastHand = hand
				if sameHandStreak > 0 {
					sameHandStreaks = append(sameHandStreaks, sameHandStreak)
				}
				sameHandStreak = 0
			} else { // same hand
				sameHandStreak++
				lastHand = hand
				if alternationStreak > 0 {
					alternationStreaks = append(alternationStreaks, alternationStreak)
				}
				alternationStreak = 0
			}
		}
	}

	var avgSameHandStreak, avgAlternationStreak float64
	lenAlt := len(alternationStreaks)
	lenSame := len(sameHandStreaks)

	if lenAlt != 0 {
		for _, a := range alternationStreaks {
			avgAlternationStreak += a
		}
		avgAlternationStreak /= float64(len(alternationStreaks))
	} else {
		avgAlternationStreak = 0
	}

	if lenSame != 0 {
		for _, s := range sameHandStreaks {
			avgSameHandStreak += s
		}
		avgSameHandStreak /= float64(len(sameHandStreaks))
	} else {
		avgSameHandStreak = 0
	}

	return lenAlt, avgAlternationStreak, lenSame, avgSameHandStreak
}

func singleHandUtilizationDir(l Layout, dir string) (int, float64, int, float64) {
	d, err := os.Open(dir)
	if err != nil {
		panic("couldn't open the directory")
	}

	info, err := d.Readdir(-1)
	d.Close()
	if err != nil {
		panic("error reading dir")
	}

	var alternations, sameHand int
	var avgAlternationStreak, avgSameHandStreak float64
	for _, file := range info {
		name := fmt.Sprintf("%s/%s", dir, file.Name())
		altStreaks, avgAltStreak, sameStreaks, avgSameStreak := singleHandUtilizationFile(l, name)
		// fmt.Println(altStreaks, avgAltStreak, sameStreaks, avgSameStreak)
		alternations += altStreaks
		sameHand += sameStreaks
		avgAlternationStreak += float64(altStreaks) * avgAltStreak
		avgSameHandStreak += float64(sameStreaks) * avgSameStreak
	}

	avgAlternations := avgAlternationStreak / float64(alternations)
	avgSameHand := avgSameHandStreak / float64(sameHand)

	// fmt.Println(avgAlternations, avgSameHand)

	return alternations, avgAlternations, sameHand, avgSameHand
}

func SingleHandUtilization(l Layout, dir string) {
	alts, avgAlternations, sames, avgSameHand := singleHandUtilizationDir(l, dir)
	fmt.Println(alts, avgAlternations, sames, avgSameHand)

	shu := float64(sames) * avgSameHand / (float64(sames)*avgSameHand + float64(alts)*avgAlternations)

	fmt.Printf("- Single Hand Utilization (SHU) = %.2f%%\n\n", 100*shu)
}
