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

func expected(s1, s2 string) int {
	n1 := len(s1)
	n2 := len(s2)
	maxL := n1
	if n2 < maxL {
		maxL = n2
	}
	for L := 1; L <= maxL; L++ {
		cnt1 := make(map[string]int)
		for i := 0; i+L <= n1; i++ {
			cnt1[s1[i:i+L]]++
		}
		cnt2 := make(map[string]int)
		for i := 0; i+L <= n2; i++ {
			cnt2[s2[i:i+L]]++
		}
		for sub, c1 := range cnt1 {
			if c1 == 1 && cnt2[sub] == 1 {
				return L
			}
		}
	}
	return -1
}

func generateCase(rng *rand.Rand) (string, int) {
	l1 := rng.Intn(10) + 1
	l2 := rng.Intn(10) + 1
	letters := "abcdefghijklmnopqrstuvwxyz"
	b1 := make([]byte, l1)
	b2 := make([]byte, l2)
	for i := range b1 {
		b1[i] = letters[rng.Intn(len(letters))]
	}
	for i := range b2 {
		b2[i] = letters[rng.Intn(len(letters))]
	}
	s1 := string(b1)
	s2 := string(b2)
	input := fmt.Sprintf("%s\n%s\n", s1, s2)
	return input, expected(s1, s2)
}

func runCase(exe, input string, exp int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
