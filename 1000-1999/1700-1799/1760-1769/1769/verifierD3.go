package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strings"
)

const totalCards = 36

func cardIndex(card string) (int, bool) {
	if len(card) != 2 {
		return 0, false
	}
	vals := "6789TJQKA"
	suits := "CDSH"
	v := strings.IndexByte(vals, card[0])
	s := strings.IndexByte(suits, card[1])
	if v < 0 || s < 0 {
		return 0, false
	}
	return s*9 + v, true
}

type gameSolver struct {
	aMask uint64
	bMask uint64
	memo  map[uint64]int
}

func (g *gameSolver) legalMoves(rem, played uint64) []uint64 {
	moves := make([]uint64, 0, bits.OnesCount64(rem))
	for i := 0; i < totalCards; i++ {
		bit := uint64(1) << i
		if rem&bit == 0 {
			continue
		}
		r := i % 9
		ok := false
		switch {
		case r == 3: // 9
			ok = true
		case r < 3: // 6,7,8 require next rank
			ok = (played & (uint64(1) << (i + 1))) != 0
		default: // T,J,Q,K,A require previous rank
			ok = (played & (uint64(1) << (i - 1))) != 0
		}
		if ok {
			moves = append(moves, bit)
		}
	}
	return moves
}

// solve returns game value from Alice's perspective:
// +x -> Alice wins, Bob has x cards left; -y -> Bob wins, Alice has y cards left.
func (g *gameSolver) solve(played uint64, turn int) int {
	key := (played << 1) | uint64(turn)
	if v, ok := g.memo[key]; ok {
		return v
	}

	aliceRem := g.aMask &^ played
	bobRem := g.bMask &^ played
	if aliceRem == 0 {
		v := bits.OnesCount64(bobRem)
		g.memo[key] = v
		return v
	}
	if bobRem == 0 {
		v := -bits.OnesCount64(aliceRem)
		g.memo[key] = v
		return v
	}

	rem := aliceRem
	if turn == 1 {
		rem = bobRem
	}
	moves := g.legalMoves(rem, played)
	if len(moves) == 0 {
		oppRem := bobRem
		if turn == 1 {
			oppRem = aliceRem
		}
		if len(g.legalMoves(oppRem, played)) == 0 {
			g.memo[key] = 0
			return 0
		}
		v := g.solve(played, 1-turn)
		g.memo[key] = v
		return v
	}

	if turn == 0 {
		best := -1 << 30
		for _, mv := range moves {
			val := g.solve(played|mv, 1)
			if val > best {
				best = val
			}
		}
		g.memo[key] = best
		return best
	}

	best := 1 << 30
	for _, mv := range moves {
		val := g.solve(played|mv, 0)
		if val < best {
			best = val
		}
	}
	g.memo[key] = best
	return best
}

func importance(aMask, bMask uint64) int {
	s := gameSolver{aMask: aMask, bMask: bMask, memo: make(map[uint64]int)}
	vAliceFirst := s.solve(0, 0)
	s.memo = make(map[uint64]int)
	vBobFirst := s.solve(0, 1)
	if vAliceFirst >= vBobFirst {
		return vAliceFirst - vBobFirst
	}
	return vBobFirst - vAliceFirst
}

func validateOutput(out string, k int) error {
	tokens := strings.Fields(out)
	need := 36 * k
	if len(tokens) != need {
		return fmt.Errorf("expected %d cards in output, got %d", need, len(tokens))
	}

	seenImportance := make(map[int]bool)
	for i := 0; i < k; i++ {
		start := i * 36
		var aMask, bMask uint64
		used := make(map[int]bool, 36)
		for j := 0; j < 18; j++ {
			idx, ok := cardIndex(tokens[start+j])
			if !ok {
				return fmt.Errorf("deal %d: invalid card %q", i+1, tokens[start+j])
			}
			if used[idx] {
				return fmt.Errorf("deal %d: duplicated card %q", i+1, tokens[start+j])
			}
			used[idx] = true
			aMask |= uint64(1) << idx
		}
		for j := 18; j < 36; j++ {
			idx, ok := cardIndex(tokens[start+j])
			if !ok {
				return fmt.Errorf("deal %d: invalid card %q", i+1, tokens[start+j])
			}
			if used[idx] {
				return fmt.Errorf("deal %d: duplicated card %q", i+1, tokens[start+j])
			}
			used[idx] = true
			bMask |= uint64(1) << idx
		}
		if len(used) != 36 {
			return fmt.Errorf("deal %d: must contain all 36 unique cards", i+1)
		}

		imp := importance(aMask, bMask)
		if seenImportance[imp] {
			return fmt.Errorf("deal %d: repeated importance value %d", i+1, imp)
		}
		seenImportance[imp] = true
	}
	return nil
}

func runCase(bin string, k int) error {
	in := fmt.Sprintf("%d\n", k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if err := validateOutput(out.String(), k); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD3.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for k := 1; k <= 26; k++ {
		if err := runCase(bin, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", k, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
