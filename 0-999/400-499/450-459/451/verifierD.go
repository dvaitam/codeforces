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

func isGood(s string) bool {
	if len(s) == 0 {
		return true
	}
	comp := []byte{s[0]}
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			comp = append(comp, s[i])
		}
	}
	for i := 0; i < len(comp)/2; i++ {
		if comp[i] != comp[len(comp)-1-i] {
			return false
		}
	}
	return true
}

func solveCase(str string) (int, int) {
	even, odd := 0, 0
	n := len(str)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if isGood(str[i : j+1]) {
				if (j-i+1)%2 == 0 {
					even++
				} else {
					odd++
				}
			}
		}
	}
	return even, odd
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 'a'
		} else {
			b[i] = 'b'
		}
	}
	s := string(b)
	even, odd := solveCase(s)
	input := fmt.Sprintf("%s\n", s)
	expect := fmt.Sprintf("%d %d", even, odd)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
