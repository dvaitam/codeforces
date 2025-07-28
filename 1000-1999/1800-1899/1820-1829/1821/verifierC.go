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

func solveC(s string) int {
	n := len(s)
	best := int(^uint(0) >> 1)
	for ch := byte('a'); ch <= 'z'; ch++ {
		prev := -1
		maxSeg := 0
		for i := 0; i < n; i++ {
			if s[i] == ch {
				seg := i - prev - 1
				if seg > maxSeg {
					maxSeg = seg
				}
				prev = i
			}
		}
		seg := n - prev - 1
		if seg > maxSeg {
			maxSeg = seg
		}
		ops := 0
		for v := maxSeg; v > 0; v >>= 1 {
			ops++
		}
		if ops < best {
			best = ops
		}
	}
	return best
}

func genCaseC(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	letters := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := genCaseC(rng)
		input := fmt.Sprintf("1\n%s\n", s)
		expect := fmt.Sprintf("%d", solveC(s))
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
