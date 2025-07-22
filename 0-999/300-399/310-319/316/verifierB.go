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

type Test struct {
	n int
	x int
	a []int
}

func solve(n, x int, a []int) []int {
	next := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if a[i] != 0 {
			next[a[i]] = i
		}
	}
	var other []int
	sumOther := 0
	posInChain := 0
	for i := 1; i <= n; i++ {
		if a[i] != 0 {
			continue
		}
		cur := i
		size := 0
		containsX := false
		for cur != 0 {
			if cur == x {
				containsX = true
				posInChain = size
			}
			size++
			cur = next[cur]
		}
		if !containsX {
			other = append(other, size)
			sumOther += size
		}
	}
	dp := make([]bool, sumOther+1)
	dp[0] = true
	for _, sz := range other {
		for s := sumOther; s >= sz; s-- {
			if dp[s-sz] {
				dp[s] = true
			}
		}
	}
	var res []int
	for s := 0; s <= sumOther; s++ {
		if dp[s] {
			res = append(res, s+posInChain+1)
		}
	}
	return res
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 3 // up to 9/10
	ids := rng.Perm(n)
	a := make([]int, n+1)
	start := true
	for i, id := range ids {
		if start || rng.Float64() < 0.3 {
			a[id+1] = 0
			start = false
		} else {
			a[id+1] = ids[i-1] + 1
		}
	}
	x := ids[rng.Intn(n)] + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, x)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	ans := solve(n, x, a)
	var out strings.Builder
	for i, v := range ans {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", v)
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
