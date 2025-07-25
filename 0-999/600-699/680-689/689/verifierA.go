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

// solve computes the expected output for the given phone number.
func solve(n int, s string) string {
	pos := map[byte][2]int{
		'1': {0, 0}, '2': {1, 0}, '3': {2, 0},
		'4': {0, 1}, '5': {1, 1}, '6': {2, 1},
		'7': {0, 2}, '8': {1, 2}, '9': {2, 2},
		'0': {1, 3},
	}

	board := [4][3]bool{}
	for d, p := range pos {
		board[p[1]][p[0]] = true
		_ = d
	}

	if n == 1 {
		return "YES"
	}

	moves := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		a := pos[s[i-1]]
		b := pos[s[i]]
		moves[i-1] = [2]int{b[0] - a[0], b[1] - a[1]}
	}

	for d := byte('0'); d <= '9'; d++ {
		start, ok := pos[d]
		if !ok {
			continue
		}
		x, y := start[0], start[1]
		valid := true
		for _, mv := range moves {
			x += mv[0]
			y += mv[1]
			if y < 0 || y >= 4 || x < 0 || x >= 3 || !board[y][x] {
				valid = false
				break
			}
		}
		if valid && d != s[0] {
			return "NO"
		}
	}
	return "YES"
}

func runCase(bin string, n int, s string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%d\n%s\n", n, s))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solve(n, s)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) (int, string) {
	n := rng.Intn(9) + 1 // 1..9
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + rng.Intn(10))
	}
	return n, string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, s := randomCase(rng)
		if err := runCase(bin, n, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %s\n", i+1, err, n, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
