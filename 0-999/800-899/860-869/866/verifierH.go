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

const mod = 1000000007

func fact(n int) int {
	res := 1
	for i := 2; i <= n; i++ {
		res = (res * i) % mod
	}
	return res
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20)
	input := fmt.Sprintf("%d\n", n)
	expected := fmt.Sprintf("%d", fact(n))
	return input, expected
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 1; i <= cases; i++ {
		input, expect := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\n got: %s\n", i, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
