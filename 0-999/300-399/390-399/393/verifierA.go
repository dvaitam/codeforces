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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expectedAnswer(s string) int {
	cntN, cntI, cntE, cntT := 0, 0, 0, 0
	for _, ch := range s {
		switch ch {
		case 'n':
			cntN++
		case 'i':
			cntI++
		case 'e':
			cntE++
		case 't':
			cntT++
		}
	}
	maxByN := 0
	if cntN >= 3 {
		maxByN = (cntN - 1) / 2
	}
	maxByE := cntE / 3
	res := maxByN
	if maxByE < res {
		res = maxByE
	}
	if cntI < res {
		res = cntI
	}
	if cntT < res {
		res = cntT
	}
	return res
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(100) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	s := string(b)
	return s + "\n", expectedAnswer(s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if out != fmt.Sprint(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
