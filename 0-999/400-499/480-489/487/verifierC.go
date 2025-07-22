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

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

func parseOutput(out string, n int) ([]int, bool, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return nil, false, fmt.Errorf("no output")
	}
	if tokens[0] == "NO" {
		if len(tokens) != 1 {
			return nil, false, fmt.Errorf("extra tokens after NO")
		}
		return nil, false, nil
	}
	if tokens[0] != "YES" {
		return nil, false, fmt.Errorf("first token must be YES or NO")
	}
	if len(tokens) != n+1 {
		return nil, true, fmt.Errorf("expected %d numbers, got %d", n, len(tokens)-1)
	}
	seq := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return nil, true, fmt.Errorf("invalid int %q", tokens[i+1])
		}
		seq[i] = v
	}
	return seq, true, nil
}

func validSequence(n int, seq []int) error {
	if len(seq) != n {
		return fmt.Errorf("wrong length")
	}
	used := make([]bool, n+1)
	for _, v := range seq {
		if v < 1 || v > n || used[v] {
			return fmt.Errorf("not a permutation")
		}
		used[v] = true
	}
	prefUsed := make([]bool, n)
	prod := 1
	for i, v := range seq {
		if i == 0 {
			prod = v % n
		} else {
			prod = (prod * v) % n
		}
		if prefUsed[prod] {
			return fmt.Errorf("duplicate prefix value")
		}
		prefUsed[prod] = true
	}
	for i := 0; i < n; i++ {
		if !prefUsed[i] {
			return fmt.Errorf("missing prefix value %d", i)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		input := fmt.Sprintf("%d\n", n)
		expOut, err := run("487C.go", input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal reference failed on case %d: %v\n", t+1, err)
			os.Exit(1)
		}
		expYes := strings.HasPrefix(expOut, "YES")
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		seq, gotYes, err := parseOutput(got, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\noutput:\n%s\ninput:\n%s", t+1, err, got, input)
			os.Exit(1)
		}
		if gotYes != expYes {
			fmt.Fprintf(os.Stderr, "case %d: expected %v got %v\ninput:\n%s", t+1, expYes, gotYes, input)
			os.Exit(1)
		}
		if gotYes {
			if err := validSequence(n, seq); err != nil {
				fmt.Fprintf(os.Stderr, "case %d: invalid sequence: %v\noutput:\n%s\ninput:\n%s", t+1, err, got, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
