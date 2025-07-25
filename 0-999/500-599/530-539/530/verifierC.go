package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type pair struct{ X, Y int }

func runCandidate(bin string, input string) (string, error) {
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

func expected(a, b, c int) string {
	sols := []pair{}
	for x := 1; a*x < c; x++ {
		rem := c - a*x
		if rem%b == 0 {
			y := rem / b
			if y >= 1 {
				sols = append(sols, pair{x, y})
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d", len(sols))
	for _, p := range sols {
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d %d", p.X, p.Y)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		a := rng.Intn(20) + 1
		b := rng.Intn(20) + 1
		c := rng.Intn(20) + 1
		input := fmt.Sprintf("%d %d %d\n", a, b, c)
		want := expected(a, b, c)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\n", i+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
