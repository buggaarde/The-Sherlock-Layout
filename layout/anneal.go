package layout

import (
	"fmt"
	"math"

	"github.com/buggaarde/The-Sherlock-Layout/text"
	"github.com/mitchellh/copystructure"
)

// Anneal the given keyboard layout
func Anneal(
	l Layout,
	initialTemperature float64,
	tf text.TripletFrequency,
	df text.DoubleFrequency,
	sf text.SingleFrequency,
) Layout {
	temp := initialTemperature

	var swapLayout Layout
	currentLayout := l
	bestLayout := l

	currentScore := Score(l, tf, df, sf)
	bestScore := currentScore

	stepsNoBest := 0

	for {
		if temp < 0 {
			break
		}

		// the possible permutations in each time step
		if l.r.Float64() > 0.2 {
			// 80% of the time, simply try to swap adjacent keys in the layout
			key1, key2 := randomNeighbours(currentLayout)
			swapLayout = swap(currentLayout, key1, key2)
		} else if l.r.Float64() > 0.33 {
			// of the remaining 20%
			// 2/3 of the time, choose three adjacent keys and cycle them
			// i.e. key1, key2, key3 -> key2, key3, key1
			key1, key2, key3 := random3Cycle(currentLayout)
			swapLayout = cycle3(currentLayout, key1, key2, key3)
		} else {
			// and the last third is used to attempt a symmetric swap
			// from left to right
			key := randomKey(currentLayout)
			swapLayout = swapLeftRight(currentLayout, key)
		}

		// compare the scores of the current and swapped layouts,
		// and accept/decline the swap based on standard swapping criteria
		swapIsAccepted, swapScore := acceptSwap(
			currentLayout, swapLayout,
			currentScore, temp,
			tf, df, sf,
		)

		if swapIsAccepted {
			currentLayout, currentScore = swapLayout, swapScore
		}

		// save the best layout if one is found
		if currentScore < bestScore-1 {
			bestLayout, bestScore = currentLayout, currentScore
			fmt.Printf("\rNew best score found at temperature %d/%d!\n", int(temp), int(initialTemperature))
			fmt.Printf("New best score is %.2f\n", bestScore)
			fmt.Println()
			Print(bestLayout)
			layoutName := fmt.Sprintf("intermediate_layouts/%d.toml", int(bestScore))
			Write(bestLayout, layoutName)
			stepsNoBest = 0
		}

		// if no new best has been found in the last 10000 swaps,
		// assume that the path through search space, given by the
		// swaps so far, have been in vain, and reset back to the
		// best known layout
		if stepsNoBest > 20000 {
			stepsNoBest = 0
			currentLayout = bestLayout

		}
		// } else if swapIsAccepted {
		// 	stepsNoBest++
		// }

		temp -= 0.25

		fmt.Printf(" -- Current temperature: %.2f        \r", temp)
	}

	return bestLayout
}

func swap(l Layout, letter1, letter2 string) Layout {
	keyPos, err := copystructure.Copy(l.KeyPositions)
	if err != nil {
		panic(err)
	}
	k := keyPos.(map[string]int)

	keyHand, err := copystructure.Copy(l.PositionHand)
	if err != nil {
		panic(err)
	}
	h := keyHand.(map[string]int)

	posKey, err := copystructure.Copy(l.PositionKey)
	if err != nil {
		panic(err)
	}
	p := posKey.(map[int]string)

	k[letter1], k[letter2] = l.KeyPositions[letter2], l.KeyPositions[letter1]
	h[letter1], h[letter2] = l.PositionHand[letter2], l.PositionHand[letter1]
	p[k[letter1]], p[k[letter2]] = letter1, letter2
	l.KeyPositions = k
	l.PositionHand = h
	l.PositionKey = p

	return l
}

func cycle3(l Layout, letter1, letter2, letter3 string) Layout {
	keyPos, err := copystructure.Copy(l.KeyPositions)
	if err != nil {
		panic(err)
	}
	k := keyPos.(map[string]int)

	keyHand, err := copystructure.Copy(l.PositionHand)
	if err != nil {
		panic(err)
	}
	h := keyHand.(map[string]int)

	posKey, err := copystructure.Copy(l.PositionKey)
	if err != nil {
		panic(err)
	}
	p := posKey.(map[int]string)

	k[letter1], k[letter2], k[letter3] = l.KeyPositions[letter2], l.KeyPositions[letter3], l.KeyPositions[letter1]
	h[letter1], h[letter2], h[letter3] = l.PositionHand[letter2], l.PositionHand[letter3], l.PositionHand[letter1]
	p[k[letter1]], p[k[letter2]], p[k[letter3]] = letter1, letter2, letter3
	l.KeyPositions = k
	l.PositionHand = h
	l.PositionKey = p

	return l
}

func swapLeftRight(l Layout, letter string) Layout {
	keyPos, err := copystructure.Copy(l.KeyPositions)
	if err != nil {
		panic(err)
	}
	k := keyPos.(map[string]int)

	keyHand, err := copystructure.Copy(l.PositionHand)
	if err != nil {
		panic(err)
	}
	h := keyHand.(map[string]int)

	posKey, err := copystructure.Copy(l.PositionKey)
	if err != nil {
		panic(err)
	}
	p := posKey.(map[int]string)

	letter1, letter2 := letter, letter
	letterIdx := k[letter]
	if letterIdx == 30 || letterIdx == 31 {
		letter2 = l.PositionKey[letterIdx+2]
	} else if letterIdx == 32 || letterIdx == 33 {
		letter2 = l.PositionKey[letterIdx-2]
	} else {
		finger := letterIdx % 10
		row := letterIdx - finger

		oppositeFinger := 10 - finger - 1
		letter2 = l.PositionKey[row+oppositeFinger]

	}

	k[letter1], k[letter2] = l.KeyPositions[letter2], l.KeyPositions[letter1]
	h[letter1], h[letter2] = l.PositionHand[letter2], l.PositionHand[letter1]
	p[k[letter1]], p[k[letter2]] = letter1, letter2
	l.KeyPositions = k
	l.PositionHand = h
	l.PositionKey = p

	return l
}

func randomKey(l Layout) string {
	return l.PositionKey[l.r.Intn(len(l.Adjacency))]
}

func randomNeighbours(l Layout) (string, string) {
	key := l.r.Intn(len(l.Adjacency))
	neigh := l.Adjacency[key][l.r.Intn(len(l.Adjacency[key]))]

	s1, s2 := l.PositionKey[key], l.PositionKey[neigh]

	return s1, s2
}

func random3Cycle(l Layout) (string, string, string) {
	key1 := l.r.Intn(len(l.Adjacency))

	key2 := l.Adjacency[key1][l.r.Intn(len(l.Adjacency[key1]))]

	key3 := key2
	for {
		if key3 != key1 && key3 != key2 {
			break
		}
		key3 = l.Adjacency[key2][l.r.Intn(len(l.Adjacency[key2]))]
	}

	s1, s2, s3 := l.PositionKey[key1], l.PositionKey[key2], l.PositionKey[key3]

	return s1, s2, s3
}

func acceptSwap(
	currentLayout, newLayout Layout,
	currentScore float64,
	temperature float64,
	tf text.TripletFrequency,
	df text.DoubleFrequency,
	sf text.SingleFrequency,
) (bool, float64) {
	sc, sn := currentScore, Score(newLayout, tf, df, sf)
	if sn < sc {
		return true, sn
	}

	prob := currentLayout.r.Float64()
	energy := math.Exp(-(sn - sc) / temperature)
	if energy > prob {
		return true, sn
	}

	return false, currentScore
}
