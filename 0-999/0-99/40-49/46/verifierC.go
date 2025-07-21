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

func expected(s string) int {
	n := len(s)
	totalH := 0
	bad := make([]int, n)
	for i, ch := range s {
		if ch == 'H' {
			totalH++
		} else {
			bad[i] = 1
		}
	}
	m := totalH
	curr := 0
	for i := 0; i < m; i++ {
		curr += bad[i]
	}
	minSwaps := curr
	for i := 1; i < n; i++ {
		curr -= bad[i-1]
		add := i + m - 1
		if add >= n {
			add -= n
		}
		curr += bad[add]
		if curr < minSwaps {
			minSwaps = curr
		}
	}
	return minSwaps
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(19) + 2
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('H')
		} else {
			sb.WriteByte('T')
		}
	}
	s := sb.String()
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return input, expected(s)
}

func runCase(exe, input string, expect int) error {
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
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
