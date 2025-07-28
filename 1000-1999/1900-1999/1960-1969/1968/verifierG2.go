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

func solve(n, l, r int, s string) string {
	zeros := make([]string, r-l+1)
	for i := range zeros {
		zeros[i] = "0"
	}
	return strings.Join(zeros, " ")
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	l := rng.Intn(n) + 1
	r := l + rng.Intn(n-l+1)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	sb := make([]byte, n)
	for i := range sb {
		sb[i] = letters[rng.Intn(len(letters))]
	}
	s := string(sb)
	input := fmt.Sprintf("1\n%d %d %d\n%s\n", n, l, r, s)
	expect := solve(n, l, r, s)
	return input, expect
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, exp := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
