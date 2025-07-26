package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod = 998244353

func solve(n int) int {
	games := n * (n - 1) / 2
	return games % mod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	const cases = 100
	for i := 0; i < cases; i++ {
		n := rand.Intn(20) + 2
		a := rand.Intn(n) + 1
		b := rand.Intn(n) + 1
		input := fmt.Sprintf("%d %d %d\n", n, a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			fmt.Printf("program output:\n%s\n", string(out))
			return
		}
		got := strings.TrimSpace(string(out))
		want := fmt.Sprintf("%d", solve(n))
		if got != want {
			fmt.Printf("case %d failed: input=%s expected %s got %s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Printf("OK %d cases\n", cases)
}
