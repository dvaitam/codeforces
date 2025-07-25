package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(a, b, c, x, y, z int) string {
	spare := 0
	need := 0
	if a >= x {
		spare += (a - x) / 2
	} else {
		need += x - a
	}
	if b >= y {
		spare += (b - y) / 2
	} else {
		need += y - b
	}
	if c >= z {
		spare += (c - z) / 2
	} else {
		need += z - c
	}
	if spare >= need {
		return "Yes"
	}
	return "No"
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateCase(rng *rand.Rand) (string, string) {
	a := rng.Intn(1_000_001)
	b := rng.Intn(1_000_001)
	c := rng.Intn(1_000_001)
	x := rng.Intn(1_000_001)
	y := rng.Intn(1_000_001)
	z := rng.Intn(1_000_001)
	input := fmt.Sprintf("%d %d %d\n%d %d %d\n", a, b, c, x, y, z)
	exp := solve(a, b, c, x, y, z)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	special := []struct {
		a, b, c, x, y, z int
	}{
		{4, 4, 0, 2, 1, 2},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{1_000_000, 1_000_000, 1_000_000, 0, 0, 0},
		{0, 1_000_000, 0, 0, 500000, 0},
		{5, 5, 5, 5, 5, 5},
		{10, 0, 0, 0, 0, 5},
		{0, 10, 0, 0, 10, 1},
		{1000000, 0, 0, 0, 0, 1000000},
		{3, 3, 3, 2, 2, 2},
	}

	for i, tc := range special {
		input := fmt.Sprintf("%d %d %d\n%d %d %d\n", tc.a, tc.b, tc.c, tc.x, tc.y, tc.z)
		expected := solve(tc.a, tc.b, tc.c, tc.x, tc.y, tc.z)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("special case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expected {
			fmt.Printf("special case %d failed: expected %s got %s\n", i+1, expected, out)
			os.Exit(1)
		}
	}

	for i := 0; i < 90; i++ {
		input, expected := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expected {
			fmt.Printf("case %d failed:\ninput:%sexpected %s got %s\n", i+1, input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
