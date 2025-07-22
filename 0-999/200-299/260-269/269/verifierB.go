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

type testCase struct {
	input    string
	expected string
}

func fenwickQuery(bit []int, i int) int {
	res := 0
	for i > 0 {
		if bit[i] > res {
			res = bit[i]
		}
		i &= i - 1
	}
	return res
}

func fenwickUpdate(bit []int, i, val int) {
	for i < len(bit) {
		if val > bit[i] {
			bit[i] = val
		}
		i += i & -i
	}
}

func solveCase(species []int, m int) string {
	bit := make([]int, m+2)
	best := 0
	for _, s := range species {
		cur := fenwickQuery(bit, s) + 1
		if cur > best {
			best = cur
		}
		fenwickUpdate(bit, s, cur)
	}
	return fmt.Sprintf("%d", len(species)-best)
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(15) + 1
	m := rng.Intn(n) + 1
	species := make([]int, n)
	positions := make([]float64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	pos := 0.0
	// guarantee each species appears at least once
	for i := 1; i <= m; i++ {
		species[i-1] = i
	}
	for i := m; i < n; i++ {
		species[i] = rng.Intn(m) + 1
	}
	rng.Shuffle(n, func(i, j int) { species[i], species[j] = species[j], species[i] })
	for i := 0; i < n; i++ {
		pos += rng.Float64()*10 + 0.1
		positions[i] = pos
		sb.WriteString(fmt.Sprintf("%d %.6f\n", species[i], positions[i]))
	}
	expect := solveCase(species, m)
	return testCase{input: sb.String(), expected: expect}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{}
	cases = append(cases, testCase{input: "1 1\n1 0.0\n", expected: "0"})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
