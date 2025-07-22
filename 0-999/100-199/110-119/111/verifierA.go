package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, bool) {
	n := rng.Intn(10) + 1 // 1..10
	u := rng.Intn(50) + 1
	y := u + n - 1 + rng.Intn(20)
	maxSq := int64(u*u + n - 1)
	x := rng.Int63n(maxSq) + 1
	if rng.Intn(3) == 0 { // make unsatisfiable
		if rng.Intn(2) == 0 {
			y = n - 1 // y < n -> impossible
		} else {
			x = maxSq + int64(rng.Intn(10)+1)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, y))
	feasible := y >= n && int64(y-n+1)*int64(y-n+1)+int64(n-1) >= x
	return sb.String(), feasible
}

func verify(input, output string, feasible bool) error {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	var x, y int64
	fmt.Fscan(in, &n, &x, &y)
	out := bufio.NewReader(strings.NewReader(strings.TrimSpace(output)))
	var vals []int64
	for {
		var v int64
		if _, err := fmt.Fscan(out, &v); err != nil {
			break
		}
		vals = append(vals, v)
	}
	if feasible {
		if len(vals) != n {
			return fmt.Errorf("expected %d numbers, got %d", n, len(vals))
		}
		var sum, sq int64
		for _, v := range vals {
			if v <= 0 {
				return fmt.Errorf("non-positive number %d", v)
			}
			sum += v
			sq += v * v
		}
		if sum > y {
			return fmt.Errorf("sum %d exceeds %d", sum, y)
		}
		if sq < x {
			return fmt.Errorf("sum squares %d less than %d", sq, x)
		}
	} else {
		if len(vals) != 1 || vals[0] != -1 {
			return fmt.Errorf("expected -1 for impossible case")
		}
	}
	return nil
}

func runCase(exe, input string, feasible bool) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verify(input, out.String(), feasible)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, feasible := generateCase(rng)
		if err := runCase(exe, in, feasible); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
