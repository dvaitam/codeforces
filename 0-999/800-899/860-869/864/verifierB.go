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

func solveB(s string) string {
	maxDistinct := 0
	seen := make(map[rune]bool)
	for _, ch := range s {
		if ch >= 'a' && ch <= 'z' {
			seen[ch] = true
		} else {
			if len(seen) > maxDistinct {
				maxDistinct = len(seen)
			}
			seen = make(map[rune]bool)
		}
	}
	if len(seen) > maxDistinct {
		maxDistinct = len(seen)
	}
	return fmt.Sprintf("%d", maxDistinct)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(200) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		} else {
			sb.WriteByte(byte('A' + rng.Intn(26)))
		}
	}
	s := sb.String()
	input := fmt.Sprintf("%d\n%s\n", len(s), s)
	return input, solveB(s)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
