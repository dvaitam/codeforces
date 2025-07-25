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

func solve(n int) string {
	s := fmt.Sprintf("%05d", n)
	rearr := []byte{s[0], s[2], s[4], s[3], s[1]}
	var k int
	fmt.Sscanf(string(rearr), "%d", &k)
	res := 1
	for i := 0; i < 5; i++ {
		res = res * k % 100000
	}
	return fmt.Sprintf("%05d", res)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(90000) + 10000
	ans := solve(n)
	input := fmt.Sprintf("%d\n", n)
	return input, ans
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
