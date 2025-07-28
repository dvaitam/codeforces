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

func parseOutput(out string) (int, []int, error) {
	rdr := strings.NewReader(out)
	var m int
	if _, err := fmt.Fscan(rdr, &m); err != nil {
		return 0, nil, fmt.Errorf("failed to read m: %v", err)
	}
	if m < 1 || m > 25 {
		return 0, nil, fmt.Errorf("invalid m %d", m)
	}
	arr := make([]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(rdr, &arr[i]); err != nil {
			return 0, nil, fmt.Errorf("failed to read value %d: %v", i+1, err)
		}
		if arr[i] < 0 || arr[i] > 1_000_000_000 {
			return 0, nil, fmt.Errorf("value out of range")
		}
	}
	return m, arr, nil
}

func check(n, k int, arr []int) error {
	reachable := make([]bool, n+1)
	reachable[0] = true
	for _, x := range arr {
		if x > n {
			continue
		}
		for s := n; s >= 0; s-- {
			if reachable[s] && s+x <= n {
				reachable[s+x] = true
			}
		}
	}
	if k <= n && reachable[k] {
		return fmt.Errorf("found subsequence sum %d", k)
	}
	for v := 1; v <= n; v++ {
		if v == k {
			continue
		}
		if !reachable[v] {
			return fmt.Errorf("no subsequence for %d", v)
		}
	}
	return nil
}

func genCase(rng *rand.Rand) (string, int, int) {
	n := rng.Intn(20) + 2
	k := rng.Intn(n) + 1
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	return input, n, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n, k := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		_, arr, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", i+1, err, out, in)
			os.Exit(1)
		}
		if err := check(n, k, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", i+1, err, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
