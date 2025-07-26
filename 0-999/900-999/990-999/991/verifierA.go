package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expectedA(a, b, c, n int) int {
	visited := a + b - c
	if a < c || b < c || visited >= n || a > n || b > n {
		return -1
	}
	return n - visited
}

func generateCaseA(rng *rand.Rand) (string, int) {
	a := rng.Intn(101)
	b := rng.Intn(101)
	c := rng.Intn(101)
	n := rng.Intn(101) + 1
	input := fmt.Sprintf("%d %d %d %d\n", a, b, c, n)
	expected := expectedA(a, b, c, n)
	return input, expected
}

func runCaseA(bin, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(outStr)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseA(rng)
		if err := runCaseA(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
