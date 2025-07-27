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

func run(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveA(n, m, r, c int64) int64 {
	d1 := abs(1-r) + abs(1-c)
	d2 := abs(1-r) + abs(m-c)
	d3 := abs(n-r) + abs(1-c)
	d4 := abs(n-r) + abs(m-c)
	ans := d1
	if d2 > ans {
		ans = d2
	}
	if d3 > ans {
		ans = d3
	}
	if d4 > ans {
		ans = d4
	}
	return ans
}

func runCase(bin string, n, m, r, c int64) error {
	input := fmt.Sprintf("1\n%d %d %d %d\n", n, m, r, c)
	expect := fmt.Sprintf("%d", solveA(n, m, r, c))
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	if out != expect {
		return fmt.Errorf("expected %s got %s", expect, out)
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
	total := 0

	// some edge cases
	edges := [][4]int64{
		{1, 1, 1, 1},
		{1, 5, 1, 3},
		{5, 1, 2, 1},
		{5, 5, 5, 5},
		{10, 10, 1, 10},
	}
	for _, e := range edges {
		if err := runCase(bin, e[0], e[1], e[2], e[3]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}

	for total < 100 {
		n := rng.Int63n(1000) + 1
		m := rng.Int63n(1000) + 1
		r := rng.Int63n(n) + 1
		c := rng.Int63n(m) + 1
		if err := runCase(bin, n, m, r, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d m=%d r=%d c=%d)\n", total+1, err, n, m, r, c)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}
