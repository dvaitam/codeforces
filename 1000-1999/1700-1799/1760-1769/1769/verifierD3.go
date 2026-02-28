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
const (
	suitStates  = 25
	totalStates = suitStates * suitStates * suitStates * suitStates
)

type stateMove struct {
	bit  uint64
	next int
}

var (
	stateMasks     [totalStates]uint64
	stateMoveCount [totalStates]uint8
	stateMoves     [totalStates][8]stateMove
)

func suitCode(l, r int) int {
	// 0 means empty; other states are intervals [l, r] containing rank 9 (index 3).
	if l == -1 {
		return 0
	}
	return 1 + l*6 + (r - 3)
}

func decodeSuit(code int) (int, int) {
	if code == 0 {
		return -1, -1
	}
	code--
	l := code / 6
	r := 3 + code%6
	return l, r
}

func initStateGraph() {
	mul := [4]int{suitStates * suitStates * suitStates, suitStates * suitStates, suitStates, 1}
	for s0 := 0; s0 < suitStates; s0++ {
		l0, r0 := decodeSuit(s0)
		for s1 := 0; s1 < suitStates; s1++ {
			l1, r1 := decodeSuit(s1)
			for s2 := 0; s2 < suitStates; s2++ {
				l2, r2 := decodeSuit(s2)
				for s3 := 0; s3 < suitStates; s3++ {
					l3, r3 := decodeSuit(s3)
					s := [4]int{s0, s1, s2, s3}
					ls := [4]int{l0, l1, l2, l3}
					rs := [4]int{r0, r1, r2, r3}
					idx := s0*mul[0] + s1*mul[1] + s2*mul[2] + s3

					var played uint64
					var cnt uint8
					for suit := 0; suit < 4; suit++ {
						l, r := ls[suit], rs[suit]
						if l == -1 {
							bit := uint64(1) << (suit*9 + 3)
							next := idx + (suitCode(3, 3)-s[suit])*mul[suit]
							stateMoves[idx][cnt] = stateMove{bit: bit, next: next}
							cnt++
							continue
						}
						for rank := l; rank <= r; rank++ {
							played |= uint64(1) << (suit*9 + rank)
						}
						if l > 0 {
							rank := l - 1
							bit := uint64(1) << (suit*9 + rank)
							next := idx + (suitCode(l-1, r)-s[suit])*mul[suit]
							stateMoves[idx][cnt] = stateMove{bit: bit, next: next}
							cnt++
						}
						if r < 8 {
							rank := r + 1
							bit := uint64(1) << (suit*9 + rank)
							next := idx + (suitCode(l, r+1)-s[suit])*mul[suit]
							stateMoves[idx][cnt] = stateMove{bit: bit, next: next}
							cnt++
						}
					}
					stateMasks[idx] = played
					stateMoveCount[idx] = cnt
				}
			}
		}
	}
}

func init() {
	initStateGraph()
}

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
	memo  [2][totalStates]int
}

func (g *gameSolver) hasMove(state int, rem uint64) bool {
	for i := 0; i < int(stateMoveCount[state]); i++ {
		if rem&stateMoves[state][i].bit != 0 {
			return true
		}
	}
	return false
}

// solve returns game value from Alice's perspective:
// +x -> Alice wins, Bob has x cards left; -y -> Bob wins, Alice has y cards left.
func (g *gameSolver) solve(state, turn int) int {
	if g.memo[turn][state] != 1<<30 {
		return g.memo[turn][state]
	}

	played := stateMasks[state]
	aliceRem := g.aMask &^ played
	bobRem := g.bMask &^ played
	if aliceRem == 0 {
		v := bits.OnesCount64(bobRem)
		g.memo[turn][state] = v
		return v
	}
	if bobRem == 0 {
		v := -bits.OnesCount64(aliceRem)
		g.memo[turn][state] = v
		return v
	}

	rem := aliceRem
	if turn == 1 {
		rem = bobRem
	}
	if !g.hasMove(state, rem) {
		oppRem := bobRem
		if turn == 1 {
			oppRem = aliceRem
		}
		if !g.hasMove(state, oppRem) {
			g.memo[turn][state] = 0
			return 0
		}
		v := g.solve(state, 1-turn)
		g.memo[turn][state] = v
		return v
	}

	if turn == 0 {
		best := -1 << 30
		for i := 0; i < int(stateMoveCount[state]); i++ {
			mv := stateMoves[state][i]
			if rem&mv.bit == 0 {
				continue
			}
			val := g.solve(mv.next, 1)
			if val > best {
				best = val
			}
		}
		g.memo[turn][state] = best
		return best
	}

	best := 1 << 30
	for i := 0; i < int(stateMoveCount[state]); i++ {
		mv := stateMoves[state][i]
		if rem&mv.bit == 0 {
			continue
		}
		val := g.solve(mv.next, 0)
		if val < best {
			best = val
		}
	}
	g.memo[turn][state] = best
	return best
}

func importance(aMask, bMask uint64) int {
	s := gameSolver{aMask: aMask, bMask: bMask}
	for turn := 0; turn < 2; turn++ {
		for i := 0; i < totalStates; i++ {
			s.memo[turn][i] = 1 << 30
		}
	}
	vAliceFirst := s.solve(0, 0)
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
