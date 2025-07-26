package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(n, m int64) int64 {
	if n >= 31 {
		return m
	}
	return m % (1 << uint(n))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	const cases = 100
	for i := 0; i < cases; i++ {
		n := rand.Int63n(108) + 1
		m := rand.Int63n(108) + 1
		input := fmt.Sprintf("%d\n%d\n", n, m)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			fmt.Printf("program output:\n%s\n", string(out))
			return
		}
		got := strings.TrimSpace(string(out))
		want := fmt.Sprintf("%d", solve(n, m))
		if got != want {
			fmt.Printf("case %d failed: n=%d m=%d expected %s got %s\n", i+1, n, m, want, got)
			return
		}
	}
	fmt.Printf("OK %d cases\n", cases)
}
