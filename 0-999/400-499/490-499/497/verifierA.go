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

type TestA struct {
	grid []string
}

func expectedA(t TestA) int {
	n := len(t.grid)
	m := len(t.grid[0])
	removed := 0
	good := make([]bool, n-1)
	for col := 0; col < m; col++ {
		remove := false
		for i := 0; i < n-1 && !remove; i++ {
			if !good[i] && t.grid[i][col] > t.grid[i+1][col] {
				remove = true
			}
		}
		if remove {
			removed++
			continue
		}
		for i := 0; i < n-1; i++ {
			if !good[i] && t.grid[i][col] < t.grid[i+1][col] {
				good[i] = true
			}
		}
	}
	return removed
}

func genCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			b[j] = byte('a' + rng.Intn(3))
		}
		grid[i] = string(b)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	exp := fmt.Sprintf("%d\n", expectedA(TestA{grid}))
	return sb.String(), exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseA(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
