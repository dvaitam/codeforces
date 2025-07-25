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

func solveCase(n int, s, t string) string {
	total := 0
	for i := 0; i < n; i++ {
		a := int(s[i] - '0')
		b := int(t[i] - '0')
		diff := a - b
		if diff < 0 {
			diff = -diff
		}
		if diff > 5 {
			diff = 10 - diff
		}
		total += diff
	}
	return fmt.Sprint(total)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	sb := make([]byte, n)
	tb := make([]byte, n)
	for i := 0; i < n; i++ {
		sb[i] = byte(rng.Intn(10)) + '0'
		tb[i] = byte(rng.Intn(10)) + '0'
	}
	s := string(sb)
	t := string(tb)
	input := fmt.Sprintf("%d\n%s\n%s\n", n, s, t)
	expected := solveCase(n, s, t)
	return input, expected
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, expect := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\n got: %s\ninput:\n%s", i+1, expect, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
