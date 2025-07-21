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

func expected(p1, p2, p3, p4, a, b int) int {
	m := p1
	if p2 < m {
		m = p2
	}
	if p3 < m {
		m = p3
	}
	if p4 < m {
		m = p4
	}
	r := b
	if r > m-1 {
		r = m - 1
	}
	if r >= a {
		return r - a + 1
	}
	return 0
}

func generateCase(rng *rand.Rand) (string, int) {
	// distinct p1..p4 in [1,1000]
	vals := make([]int, 0, 4)
	for len(vals) < 4 {
		v := rng.Intn(1000) + 1
		ok := true
		for _, x := range vals {
			if x == v {
				ok = false
				break
			}
		}
		if ok {
			vals = append(vals, v)
		}
	}
	a := rng.Intn(31416)
	b := rng.Intn(31416)
	if a > b {
		a, b = b, a
	}
	exp := expected(vals[0], vals[1], vals[2], vals[3], a, b)
	input := fmt.Sprintf("%d %d %d %d %d %d\n", vals[0], vals[1], vals[2], vals[3], a, b)
	return input, exp
}

func runCase(bin, input string, exp int) error {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
