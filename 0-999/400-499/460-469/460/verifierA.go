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

func expectedDays(n, m int) int {
	days := 0
	socks := n
	for socks > 0 {
		days++
		socks--
		if days%m == 0 {
			socks++
		}
	}
	return days
}

func runCase(bin string, n, m int) error {
	input := fmt.Sprintf("%d %d\n", n, m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := expectedDays(n, m)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	edges := []struct{ n, m int }{
		{1, 2}, {100, 2}, {1, 100}, {100, 100},
	}
	for i, e := range edges {
		if err := runCase(bin, e.n, e.m); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		n := rng.Intn(100) + 1 // 1..100
		m := rng.Intn(99) + 2  // 2..100
		if err := runCase(bin, n, m); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d %d\n", i+1, err, n, m)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
