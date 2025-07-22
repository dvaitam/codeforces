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

func solveCase(n int64, t string) int64 {
	runs := map[rune]int64{'A': 0, 'B': 0, 'C': 0, 'D': 0}
	var prev rune
	var cur int64
	for i, ch := range t {
		if i == 0 || ch != prev {
			cur = 1
		} else {
			cur++
		}
		if cur > runs[ch] {
			runs[ch] = cur
		}
		prev = ch
	}
	rmin := runs['A']
	for _, c := range []rune{'B', 'C', 'D'} {
		if runs[c] < rmin {
			rmin = runs[c]
		}
	}
	if rmin <= 0 {
		rmin = 1
	}
	return (n + rmin - 1) / rmin
}

func randomString(rng *rand.Rand) string {
	letters := []byte{'A', 'B', 'C', 'D'}
	l := rng.Intn(20) + 4
	s := make([]byte, l)
	for i := 0; i < 4; i++ {
		s[i] = letters[i]
	}
	for i := 4; i < l; i++ {
		s[i] = letters[rng.Intn(4)]
	}
	rng.Shuffle(l, func(i, j int) { s[i], s[j] = s[j], s[i] })
	return string(s)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1_000_000_000) + 1
	t := randomString(rng)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n%s\n", n, t)
	ans := solveCase(n, t)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
