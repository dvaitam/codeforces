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

func expected(n, r int64) int64 {
	if r < n {
		return r * (r + 1) / 2
	}
	return n*(n-1)/2 + 1
}

func genCase(rng *rand.Rand) (string, []int64) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	ans := make([]int64, t)
	for i := 0; i < t; i++ {
		n := rng.Int63n(1_000_000_000) + 1
		r := rng.Int63n(1_000_000_000) + 1
		fmt.Fprintf(&sb, "%d %d\n", n, r)
		ans[i] = expected(n, r)
	}
	return sb.String(), ans
}

func runCase(bin, input string, exp []int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	for i, e := range exp {
		if !scanner.Scan() {
			return fmt.Errorf("missing output line %d", i+1)
		}
		var got int64
		if _, err := fmt.Sscan(scanner.Text(), &got); err != nil {
			return fmt.Errorf("bad output on line %d: %v", i+1, err)
		}
		if got != e {
			return fmt.Errorf("line %d: expected %d got %d", i+1, e, got)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
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
