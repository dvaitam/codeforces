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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func genCases() [][2]int64 {
	rng := rand.New(rand.NewSource(1306))
	cases := make([][2]int64, 100)
	for i := range cases {
		a := rng.Int63n(1_000_000_000) + 1
		b := rng.Int63n(1_000_000_000) + 1
		cases[i] = [2]int64{a, b}
	}
	return cases
}

func runCase(bin string, ab [2]int64) error {
	input := fmt.Sprintf("%d %d\n", ab[0], ab[1])
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	want := gcd(ab[0], ab[1])
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	cases := genCases()
	for i, ab := range cases {
		if err := runCase(bin, ab); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
