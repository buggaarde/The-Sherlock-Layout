package layout

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/buggaarde/layout-optimizer/text"
)

type Layout struct {
	KeyPositions map[string]int `toml:"key_positions"`
	PositionKey  map[int]string
	Adjacency    [][]int        `toml:"adjacency"`
	Convenience  Convenience    `toml:"convenience"`
	PositionHand map[string]int `toml:"position_hand"`
	r            *rand.Rand
}

type Convenience struct {
	SinglePress map[string]float64 `toml:"singles"`
	DoublePress map[string]float64 `toml:"doubles"`
	TriplePress map[string]float64 `toml:"triples"`
}

func Score(
	l Layout,
	tf text.TripletFrequency,
	df text.DoubleFrequency,
	sf text.SingleFrequency,
) float64 {

	singleScore := 0.0
	var leftHands, rightHands float64

	for key, pos := range l.KeyPositions {
		spos := fmt.Sprintf("%d", pos)
		singleScore += l.Convenience.SinglePress[spos] * float64(sf[key])

		// Estimate the distribution of the handedness of the layout
		if hand, ok := l.PositionHand[key]; !ok {
			continue
		} else if hand == 0 {
			leftHands += float64(sf[key])
		} else if hand == 1 {
			rightHands += float64(sf[key])
		}
	}

	totalHands := leftHands + rightHands
	leftP, rightP := 100*leftHands/totalHands, 100*rightHands/totalHands

	var handednessDiff float64
	if leftP > rightP {
		handednessDiff = leftP - 50
	} else {
		handednessDiff = rightP - 50
	}

	doubleScore := 0.0
	handAlternations := 0

	for pair, freq := range df {

		first, ok := l.KeyPositions[pair[0]]
		if !ok {
			continue
		}
		second, ok := l.KeyPositions[pair[1]]
		if !ok {
			continue
		}

		if l.PositionHand[pair[0]] != l.PositionHand[pair[1]] {
			handAlternations++
		}

		posPair := fmt.Sprintf("%d,%d", first, second)

		doubleScore += l.Convenience.DoublePress[posPair] * float64(freq)
	}

	tripleScore := 0.0
	for triple, freq := range tf {
		t0, ok := l.KeyPositions[triple[0]]
		if !ok {
			continue
		}
		t1, ok := l.KeyPositions[triple[1]]
		if !ok {
			continue
		}
		t2, ok := l.KeyPositions[triple[2]]
		if !ok {
			continue
		}

		trpl := fmt.Sprintf("%d,%d,%d", t0, t1, t2)
		tripleScore += float64(freq) * l.Convenience.TriplePress[trpl]

		firstHand, ok := l.PositionHand[triple[0]]
		if !ok {
			continue
		}
		secondHand, ok := l.PositionHand[triple[1]]
		if !ok {
			continue
		}

		// count doubles if there is a hand alternation in between key presses
		if firstHand == secondHand {
			continue
		}

		if firstHand != l.PositionHand[triple[2]] {
			continue
		}

		first, ok := l.KeyPositions[triple[0]]
		if !ok {
			continue
		}
		second, ok := l.KeyPositions[triple[2]]
		if !ok {
			continue
		}

		posPair := fmt.Sprintf("%d,%d", first, second)
		doubleScore += l.Convenience.DoublePress[posPair] * float64(freq)
	}

	penalty := handednessPenalty(handednessDiff)
	return penalty * (singleScore + doubleScore + tripleScore - 0.2*float64(handAlternations))
}

func handednessPenalty(diff float64) float64 {
	pen := 1.0 + 1.0/100.*penaltyFromValues(diff, [2]float64{5, 10}, [2]float64{15, 150})
	return pen
}

func penaltyFromValues(diff float64, vals0, vals1 [2]float64) float64 {
	x0, y0 := vals0[0], vals0[1]
	x1, y1 := vals1[0], vals1[1]

	c := y0 * math.Pow(y0/y1, x0/(x1-x0))
	k := math.Log(y0/y1) / (x0 - x1)

	return c * math.Exp(k*diff)
}
