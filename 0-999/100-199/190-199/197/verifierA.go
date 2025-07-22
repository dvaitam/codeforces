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

func expected(a, b, r int) string {
	na := a / (2 * r)
	nb := b / (2 * r)
	total := na * nb
	if total%2 == 1 {
		return "First"
	}
	return "Second"
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases [][3]int
	edge := [][3]int{
		{1, 1, 1},
		{2, 2, 1},
		{100, 100, 1},
		{1, 100, 1},
		{100, 1, 1},
		{2, 1, 1},
		{1, 2, 1},
		{2, 2, 2},
		{3, 3, 2},
		{2, 3, 1},
		{3, 2, 1},
	}
	cases = append(cases, edge...)
	for i := 0; i < 100; i++ {
		a := rng.Intn(100) + 1
		b := rng.Intn(100) + 1
		r := rng.Intn(100) + 1
		cases = append(cases, [3]int{a, b, r})
	}

	for i, c := range cases {
		input := fmt.Sprintf("%d %d %d\n", c[0], c[1], c[2])
		want := expected(c[0], c[1], c[2])
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
