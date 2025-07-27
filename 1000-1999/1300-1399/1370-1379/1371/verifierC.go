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

func expected(a, b, n, m int64) string {
	if a+b < n+m || min(a, b) < m {
		return "No"
	}
	return "Yes"
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func genCase(rng *rand.Rand) (string, []string) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	ans := make([]string, t)
	for i := 0; i < t; i++ {
		a := rng.Int63n(1_000_000_000) + 1
		b := rng.Int63n(1_000_000_000) + 1
		n := rng.Int63n(1_000_000_000) + 1
		m := rng.Int63n(1_000_000_000) + 1
		fmt.Fprintf(&sb, "%d %d %d %d\n", a, b, n, m)
		ans[i] = expected(a, b, n, m)
	}
	return sb.String(), ans
}

func runCase(bin, input string, exp []string) error {
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
		got := strings.TrimSpace(scanner.Text())
		if got != e {
			return fmt.Errorf("line %d: expected %s got %s", i+1, e, got)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
