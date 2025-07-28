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

type testcaseA struct {
	n int
	k int
}

func solveCaseA(n, k int) []string {
	if k > (n+1)/2 {
		return []string{"-1"}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		line := make([]byte, n)
		for j := 0; j < n; j++ {
			line[j] = '.'
		}
		if i%2 == 0 && i/2 < k {
			line[i] = 'R'
		}
		res[i] = string(line)
	}
	return res
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCaseA(rng *rand.Rand) (string, []string) {
	t := rng.Intn(3) + 1
	cases := make([]testcaseA, t)
	input := fmt.Sprintf("%d\n", t)
	var allOut []string
	for i := 0; i < t; i++ {
		n := rng.Intn(6) + 1
		k := rng.Intn(n + 1)
		cases[i] = testcaseA{n, k}
		input += fmt.Sprintf("%d %d\n", n, k)
		outLines := solveCaseA(n, k)
		allOut = append(allOut, outLines...)
	}
	return input, allOut
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, outLines := generateCaseA(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(gotLines) != len(outLines) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(outLines), len(gotLines), in)
			os.Exit(1)
		}
		for j := range outLines {
			if strings.TrimSpace(gotLines[j]) != outLines[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, outLines[j], gotLines[j], in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
