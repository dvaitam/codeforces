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

func expectedSeq(n int) []int {
	if n%2 == 1 {
		return nil
	}
	a := make([]int, n+1)
	vis := make([]bool, n)
	a[n] = 0
	a[n-1] = n / 2
	vis[n/2] = true
	for i := n - 1; i > 0; i-- {
		x1 := a[i] / 2
		x2 := (a[i] + n) / 2
		if vis[x1] {
			a[i-1] = x2
			vis[x2] = true
		} else if vis[x2] {
			a[i-1] = x1
			vis[x1] = true
		} else if x1 > x2 {
			a[i-1] = x1
			vis[x1] = true
		} else {
			a[i-1] = x2
			vis[x2] = true
		}
	}
	return a
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := (rng.Intn(10) + 1) * 2
	return fmt.Sprintf("%d\n", n), expectedSeq(n)
}

func runCase(bin string, input string, exp []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(fields))
	}
	for i, f := range fields {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if v != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
