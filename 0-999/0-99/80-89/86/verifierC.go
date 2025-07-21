package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MOD = 1000000009

type automaton struct {
	next   [][4]int
	out    [][]int
	link   []int
	maxLen int
}

func buildAutomaton(patterns []string) *automaton {
	toIndex := func(c byte) int {
		switch c {
		case 'A':
			return 0
		case 'C':
			return 1
		case 'G':
			return 2
		case 'T':
			return 3
		}
		return 0
	}
	next := make([][4]int, 1)
	out := make([][]int, 1)
	link := []int{0}
	maxLen := 0
	for _, p := range patterns {
		if len(p) > maxLen {
			maxLen = len(p)
		}
		v := 0
		for i := 0; i < len(p); i++ {
			ci := toIndex(p[i])
			if next[v][ci] == 0 {
				next = append(next, [4]int{})
				out = append(out, nil)
				link = append(link, 0)
				next[v][ci] = len(next) - 1
			}
			v = next[v][ci]
		}
		out[v] = append(out[v], len(p))
	}
	queue := make([]int, 0, len(next))
	for c := 0; c < 4; c++ {
		if next[0][c] != 0 {
			queue = append(queue, next[0][c])
			link[next[0][c]] = 0
		}
	}
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		lv := link[v]
		if len(out[lv]) > 0 {
			out[v] = append(out[v], out[lv]...)
		}
		for c := 0; c < 4; c++ {
			u := next[v][c]
			if u != 0 {
				link[u] = next[link[v]][c]
				queue = append(queue, u)
			} else {
				next[v][c] = next[link[v]][c]
			}
		}
	}
	return &automaton{next: next, out: out, link: link, maxLen: maxLen}
}

func countStrings(n int, patterns []string) int {
	auto := buildAutomaton(patterns)
	numStates := len(auto.next)
	M := auto.maxLen
	dpCur := make([][]int, M)
	dpNext := make([][]int, M)
	for k := 0; k < M; k++ {
		dpCur[k] = make([]int, numStates)
		dpNext[k] = make([]int, numStates)
	}
	dpCur[0][0] = 1
	for pos := 0; pos < n; pos++ {
		for k := 0; k < M; k++ {
			for s := 0; s < numStates; s++ {
				dpNext[k][s] = 0
			}
		}
		for k := 0; k < M; k++ {
			for s := 0; s < numStates; s++ {
				ways := dpCur[k][s]
				if ways == 0 {
					continue
				}
				for c := 0; c < 4; c++ {
					ns := auto.next[s][c]
					maxP := 0
					for _, L := range auto.out[ns] {
						if L > maxP {
							maxP = L
						}
					}
					t := k + 1
					if maxP > t {
						maxP = t
					}
					k2 := t - maxP
					if k2 < M {
						dpNext[k2][ns] = (dpNext[k2][ns] + ways) % MOD
					}
				}
			}
		}
		dpCur, dpNext = dpNext, dpCur
	}
	ans := 0
	for s := 0; s < numStates; s++ {
		ans = (ans + dpCur[0][s]) % MOD
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(20) + 1
	m := rng.Intn(5) + 1
	patterns := make([]string, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(5) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = "ACGT"[rng.Intn(4)]
		}
		patterns[i] = string(b)
	}
	input := fmt.Sprintf("%d %d\n%s\n", n, m, strings.Join(patterns, "\n"))
	return input, countStrings(n, patterns)
}

func runCase(bin, input string, exp int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got%MOD != exp%MOD {
		return fmt.Errorf("expected %d got %d", exp%MOD, got%MOD)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
