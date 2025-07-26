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

func runCandidate(bin, input string) (string, error) {
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

func checkGrid(n int, out string) error {
	rows := strings.Split(strings.TrimSpace(out), "\n")
	if len(rows) != n {
		return fmt.Errorf("expected %d rows got %d", n, len(rows))
	}
	seen := make([]bool, n*n)
	var target int
	for i, row := range rows {
		fields := strings.Fields(row)
		if len(fields) != n {
			return fmt.Errorf("row %d: expected %d numbers got %d", i+1, n, len(fields))
		}
		xor := 0
		for _, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return fmt.Errorf("invalid integer %q", f)
			}
			if v < 0 || v >= n*n {
				return fmt.Errorf("value out of range")
			}
			if seen[v] {
				return fmt.Errorf("value %d repeated", v)
			}
			seen[v] = true
			xor ^= v
		}
		if i == 0 {
			target = xor
		} else if xor != target {
			return fmt.Errorf("row %d XOR mismatch", i+1)
		}
	}
	for j := 0; j < n; j++ {
		xor := 0
		for i := 0; i < n; i++ {
			fields := strings.Fields(rows[i])
			v, _ := strconv.Atoi(fields[j])
			xor ^= v
		}
		if xor != target {
			return fmt.Errorf("column %d XOR mismatch", j+1)
		}
	}
	for _, ok := range seen {
		if !ok {
			return fmt.Errorf("missing numbers")
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int) {
	opts := []int{4, 8, 12, 16}
	n := opts[rng.Intn(len(opts))]
	return fmt.Sprintf("%d\n", n), n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := checkGrid(n, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
