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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(s string) int {
	zeros := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			zeros++
		}
	}
	ones := len(s) - zeros
	if zeros == ones {
		return zeros - 1
	}
	if zeros < ones {
		return zeros
	}
	return ones
}

func generateCase(r *rand.Rand) (string, string) {
	l := r.Intn(100) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		if r.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	s := string(b)
	expect := fmt.Sprintf("%d", solveB(s))
	input := fmt.Sprintf("1\n%s\n", s)
	return input, expect
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
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
