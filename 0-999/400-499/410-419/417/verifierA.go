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

func expected(c, d, n, m, k int) int {
	total := n*m - k
	if total < 0 {
		total = 0
	}
	maxMain := 0
	if n > 0 {
		maxMain = (total + n - 1) / n
	}
	best := int(^uint(0) >> 1)
	for x := 0; x <= maxMain; x++ {
		rem := total - x*n
		if rem < 0 {
			rem = 0
		}
		cost := x*c + rem*d
		if cost < best {
			best = cost
		}
	}
	return best
}

func runCase(bin string, c, d, n, m, k int) error {
	input := fmt.Sprintf("%d %d\n%d %d\n%d\n", c, d, n, m, k)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
	expect := expected(c, d, n, m, k)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
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
		c := rng.Intn(100) + 1
		d := rng.Intn(100) + 1
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		k := rng.Intn(n*m + 5)
		if err := runCase(bin, c, d, n, m, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: c=%d d=%d n=%d m=%d k=%d\n", i+1, err, c, d, n, m, k)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
