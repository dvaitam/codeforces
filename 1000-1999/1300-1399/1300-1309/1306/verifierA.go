package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
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

func expected(n int64) int64 {
	return n * (n + 1) / 2
}

type Case struct{ n int64 }

func genCases() []Case {
	rng := rand.New(rand.NewSource(1306))
	cases := make([]Case, 100)
	for i := range cases {
		switch i {
		case 0:
			cases[i] = Case{0}
		case 1:
			cases[i] = Case{1}
		default:
			cases[i] = Case{rng.Int63n(1_000_000)}
		}
	}
	return cases
}

func runCase(bin string, c Case) error {
	input := fmt.Sprintf("%d\n", c.n)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	want := expected(c.n)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
