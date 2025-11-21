package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}

	freq := make([]int, 26)
	for _, ch := range s {
		freq[ch-'A']++
	}

	isVowel := make([]bool, 26)
	for _, c := range []byte{'A', 'E', 'I', 'O', 'U'} {
		isVowel[c-'A'] = true
	}

	V := 0
	for _, c := range []byte{'A', 'E', 'I', 'O', 'U'} {
		V += freq[c-'A']
	}
	Y := freq['Y'-'A']
	cntN := freq['N'-'A']
	cntG := freq['G'-'A']

	Cother := 0
	for i := 0; i < 26; i++ {
		if isVowel[i] || i == int('Y'-'A') || i == int('N'-'A') || i == int('G'-'A') {
			continue
		}
		Cother += freq[i]
	}

	Cbase := Cother + cntN + cntG
	pairMax := cntN
	if cntG < pairMax {
		pairMax = cntG
	}

	best := 0
	limitS := V + Y
	maxSByLength := len(s) / 3
	if maxSByLength < limitS {
		limitS = maxSByLength
	}

	for S := 1; S <= limitS; S++ {
		yV := 0
		if S > V {
			yV = S - V
		}
		if yV > Y {
			continue
		}
		Yrem := Y - yV
		totalSingles := Cbase + Yrem
		if totalSingles < 2*S {
			continue
		}
		extra := totalSingles - 2*S
		ngMax := pairMax
		if 2*S < ngMax {
			ngMax = 2 * S
		}
		if extra < ngMax {
			ngMax = extra
		}
		if ngMax < 0 {
			continue
		}
		length := 3*S + ngMax
		if length > best {
			best = length
		}
	}

	fmt.Println(best)
}
