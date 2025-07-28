package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveCase(a, b, c int) int {
	time1 := abs(a - 1)
	time2 := abs(b-c) + abs(c-1)
	if time1 < time2 {
		return 1
	} else if time2 < time1 {
		return 2
	}
	return 3
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
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

	cases := make([][3]int, 0, 120)
	// some deterministic small cases
	for a := 1; a <= 3; a++ {
		for b := 1; b <= 3; b++ {
			for c := 1; c <= 3; c++ {
				if b == c {
					continue
				}
				cases = append(cases, [3]int{a, b, c})
			}
		}
	}
	// random cases
	for len(cases) < 120 {
		a := rng.Intn(100000000) + 1
		b := rng.Intn(100000000) + 1
		c := rng.Intn(100000000) + 1
		if b == c {
			continue
		}
		cases = append(cases, [3]int{a, b, c})
	}

	for i, tc := range cases {
		input := fmt.Sprintf("1\n%d %d %d\n", tc[0], tc[1], tc[2])
		expected := strconv.Itoa(solveCase(tc[0], tc[1], tc[2]))
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
