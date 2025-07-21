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

type test struct {
	input    string
	expected string
}

func winner(names []string) string {
	counts := make(map[string]int)
	best := ""
	max := -1
	for _, n := range names {
		counts[n]++
		if counts[n] > max {
			max = counts[n]
			best = n
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, string) {
	teams := []string{"ALPHA", "BETA", "GAMMA", "DELTA", "OMEGA"}
	t1 := teams[rng.Intn(len(teams))]
	t2 := t1
	if rng.Intn(2) == 1 {
		for t2 == t1 {
			t2 = teams[rng.Intn(len(teams))]
		}
	}
	n := rng.Intn(100) + 1
	names := make([]string, n)
	var c1 int
	if t1 == t2 {
		c1 = n
	} else {
		c1 = rng.Intn(n)
		if c1 == n-c1 {
			if c1 == 0 {
				c1 = 1
			} else {
				c1++
			}
		}
	}
	for i := 0; i < c1; i++ {
		names[i] = t1
	}
	for i := c1; i < n; i++ {
		names[i] = t2
	}
	rng.Shuffle(n, func(i, j int) { names[i], names[j] = names[j], names[i] })
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, name := range names {
		fmt.Fprintln(&sb, name)
	}
	return sb.String(), winner(names)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
