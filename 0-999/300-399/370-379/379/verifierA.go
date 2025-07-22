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

func expected(a, b int) int {
	candles := a
	leftover := 0
	hours := 0
	for candles > 0 {
		hours += candles
		leftover += candles
		candles = leftover / b
		leftover %= b
	}
	return hours
}

func runCase(exe string, a, b int) error {
	input := fmt.Sprintf("%d %d\n", a, b)
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
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := expected(a, b)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	cases := [][2]int{{1, 2}, {4, 2}, {9, 3}, {10, 2}, {1000, 2}, {1000, 1000}}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a := rng.Intn(1000) + 1
		b := rng.Intn(999) + 2
		cases = append(cases, [2]int{a, b})
	}
	for idx, c := range cases {
		if err := runCase(exe, c[0], c[1]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
