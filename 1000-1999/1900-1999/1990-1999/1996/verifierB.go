package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "1996B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveLocal(n, k int, grid []string) string {
	var sb strings.Builder
	for i := 0; i < n; i += k {
		for j := 0; j < n; j += k {
			sb.WriteByte(grid[i][j])
		}
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	divisors := make([]int, 0)
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			divisors = append(divisors, d)
		}
	}
	k := divisors[rng.Intn(len(divisors))]
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		grid[i] = string(row)
	}
	for i := 0; i < n; i += k {
		for j := 0; j < n; j += k {
			ch := byte('0')
			if rng.Intn(2) == 1 {
				ch = '1'
			}
			for a := i; a < i+k; a++ {
				bRow := []byte(grid[a])
				for b := j; b < j+k; b++ {
					bRow[b] = ch
				}
				grid[a] = string(bRow)
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
