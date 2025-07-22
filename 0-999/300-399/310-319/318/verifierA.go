package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(n, k int64) int64 {
	half := (n + 1) / 2
	if k <= half {
		return 2*k - 1
	}
	return 2 * (k - half)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// deterministic edge cases
	fixed := [][2]int64{
		{1, 1},
		{2, 1},
		{2, 2},
		{5, 3},
		{10, 10},
		{1000000000000, 1},
		{1000000000000, 1000000000000},
		{999999999999, 500000000000},
	}
	idx := 0
	for ; idx < len(fixed); idx++ {
		n := fixed[idx][0]
		k := fixed[idx][1]
		input := fmt.Sprintf("%d %d\n", n, k)
		expectedOut := strconv.FormatInt(expected(n, k), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	// random cases to reach at least 100
	for ; idx < 100; idx++ {
		n := rng.Int63n(1000000000000) + 1 // up to 1e12
		k := rng.Int63n(n) + 1
		input := fmt.Sprintf("%d %d\n", n, k)
		expectedOut := strconv.FormatInt(expected(n, k), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
