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

func solve(x, y, z int64) (int64, int64) {
	a := x / z
	b := y / z
	total := (x + y) / z
	var transfer int64
	if total == a+b {
		transfer = 0
	} else {
		r1 := z - x%z
		r2 := z - y%z
		if r1 < r2 {
			transfer = r1
		} else {
			transfer = r2
		}
	}
	return total, transfer
}

func runCase(bin string, x, y, z int64) error {
	input := fmt.Sprintf("%d %d %d\n", x, y, z)
	exp1, exp2 := solve(x, y, z)
	expect := fmt.Sprintf("%d %d", exp1, exp2)
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewBufferString(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
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
		x := rng.Int63n(1e12)
		y := rng.Int63n(1e12)
		z := rng.Int63n(1e12) + 1
		if err := runCase(bin, x, y, z); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
