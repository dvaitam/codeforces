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

const MOD = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func solveCase(n int, cells []cell) int64 {
	assign := make(map[int]int)
	ok := true
	for _, c := range cells {
		diff := c.a - c.b
		if diff < 0 {
			diff = -diff
		}
		val := 0
		if c.ch == 'o' {
			val = 1
		}
		if prev, found := assign[diff]; found {
			if prev != val {
				ok = false
			}
		} else {
			assign[diff] = val
		}
	}
	if !ok {
		return 0
	}
	m := len(assign)
	exp := int64(n - m)
	if exp < 0 {
		exp = 0
	}
	return modPow(2, exp)
}

type cell struct {
	a  int
	b  int
	ch byte
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n + 1)
	cells := make([]cell, 0, k)
	used := make(map[[2]int]bool)
	for len(cells) < k {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if used[[2]int{a, b}] {
			continue
		}
		used[[2]int{a, b}] = true
		ch := byte('x')
		if rng.Intn(2) == 0 {
			ch = 'o'
		}
		cells = append(cells, cell{a: a, b: b, ch: ch})
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for _, c := range cells {
		fmt.Fprintf(&sb, "%d %d %c\n", c.a, c.b, c.ch)
	}
	ans := solveCase(n, cells)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
