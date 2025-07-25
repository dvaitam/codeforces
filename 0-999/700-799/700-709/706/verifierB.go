package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(100) + 1
	prices := make([]int, n)
	for i := range prices {
		prices[i] = rng.Intn(1000) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, p := range prices {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", p)
	}
	sb.WriteByte('\n')
	sort.Ints(prices)
	q := rng.Intn(100) + 1
	fmt.Fprintf(&sb, "%d\n", q)
	ans := make([]int, q)
	for i := 0; i < q; i++ {
		m := rng.Intn(1000) + 1
		fmt.Fprintf(&sb, "%d\n", m)
		cnt := 0
		for _, p := range prices {
			if p <= m {
				cnt++
			} else {
				break
			}
		}
		ans[i] = cnt
	}
	return sb.String(), ans
}

func runCase(bin, input string, expected []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outScan := bufio.NewScanner(strings.NewReader(out.String()))
	outScan.Split(bufio.ScanWords)
	for i, exp := range expected {
		if !outScan.Scan() {
			return fmt.Errorf("missing output for query %d", i+1)
		}
		var got int
		fmt.Sscan(outScan.Text(), &got)
		if got != exp {
			return fmt.Errorf("query %d: expected %d got %d", i+1, exp, got)
		}
	}
	if outScan.Scan() {
		return fmt.Errorf("extra output detected")
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
