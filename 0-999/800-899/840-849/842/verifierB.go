package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type piece struct {
	x, y, r int
}

type TestB struct {
	r, d   int
	pieces []piece
}

func generateTests() []TestB {
	rng := rand.New(rand.NewSource(2))
	tests := make([]TestB, 0, 120)
	for i := 0; i < 120; i++ {
		r := rng.Intn(500-1) + 1
		d := rng.Intn(r)
		n := rng.Intn(5) + 1
		pcs := make([]piece, n)
		for j := 0; j < n; j++ {
			pcs[j] = piece{
				x: rng.Intn(1001) - 500,
				y: rng.Intn(1001) - 500,
				r: rng.Intn(501),
			}
		}
		tests = append(tests, TestB{r, d, pcs})
	}
	return tests
}

func solve(t TestB) string {
	outer := float64(t.r)
	inner := float64(t.r - t.d)
	count := 0
	for _, p := range t.pieces {
		dist := math.Hypot(float64(p.x), float64(p.y))
		rr := float64(p.r)
		if dist+rr <= outer && dist-rr >= inner {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}

func (t TestB) input() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", t.r, t.d)
	fmt.Fprintf(&b, "%d\n", len(t.pieces))
	for _, p := range t.pieces {
		fmt.Fprintf(&b, "%d %d %d\n", p.x, p.y, p.r)
	}
	return b.String()
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		want := solve(t)
		input := t.input()
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != want {
			fmt.Printf("Test %d failed:\nInput:\n%sExpected: %s\nGot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}
