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

type Bet struct {
	t   int
	w   int64
	idx []int
	fr  int
}

type InputE struct {
	n    int
	m    int
	g    int64
	s    []int
	v    []int64
	bets []Bet
}

func solveE(inp InputE) string {
	n := inp.n
	maxMask := 1 << uint(n)
	best := int64(-1 << 63)
	for mask := 0; mask < maxMask; mask++ {
		profit := int64(0)
		genders := make([]int, n)
		for i := 0; i < n; i++ {
			genders[i] = inp.s[i]
			if mask&(1<<uint(i)) != 0 {
				genders[i] ^= 1
				profit -= inp.v[i]
			}
		}
		for _, b := range inp.bets {
			ok := true
			for _, id := range b.idx {
				if genders[id] != b.t {
					ok = false
					break
				}
			}
			if ok {
				profit += b.w
			} else if b.fr == 1 {
				profit -= inp.g
			}
		}
		if profit > best {
			best = profit
		}
	}
	return fmt.Sprintf("%d", best)
}

func generateCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(3)
	g := int64(rng.Intn(5))
	s := make([]int, n)
	v := make([]int64, n)
	for i := 0; i < n; i++ {
		s[i] = rng.Intn(2)
		v[i] = int64(rng.Intn(5))
	}
	bets := make([]Bet, m)
	for i := 0; i < m; i++ {
		t := rng.Intn(2)
		w := int64(rng.Intn(6))
		k := rng.Intn(n) + 1
		perm := rng.Perm(n)[:k]
		fr := rng.Intn(2)
		idx := make([]int, k)
		for j := 0; j < k; j++ {
			idx[j] = perm[j]
		}
		bets[i] = Bet{t: t, w: w, idx: idx, fr: fr}
	}
	inp := InputE{n: n, m: m, g: g, s: s, v: v, bets: bets}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, g)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", s[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		b := bets[i]
		fmt.Fprintf(&sb, "%d %d %d", b.t, b.w, len(b.idx))
		for _, id := range b.idx {
			fmt.Fprintf(&sb, " %d", id+1)
		}
		fmt.Fprintf(&sb, " %d\n", b.fr)
	}
	expected := solveE(inp)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
