package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func checkCase(m, k int, cnt []int, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("missing n")
	}
	n, err := strconv.Atoi(scanner.Text())
	if err != nil || n <= 0 {
		return fmt.Errorf("invalid n")
	}
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if !scanner.Scan() {
				return fmt.Errorf("not enough numbers")
			}
			v, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return fmt.Errorf("bad number")
			}
			if v < 0 || v > k {
				return fmt.Errorf("value out of range")
			}
			mat[i][j] = v
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extraneous output")
	}
	// count
	used := make([]int, k+1)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if mat[i][j] > 0 {
				used[mat[i][j]]++
			}
		}
	}
	for i := 1; i <= k; i++ {
		if used[i] != cnt[i-1] {
			return fmt.Errorf("count mismatch for %d", i)
		}
	}
	// check submatrices
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1; j++ {
			vals := []int{mat[i][j], mat[i][j+1], mat[i+1][j], mat[i+1][j+1]}
			nonZero := 0
			for _, v := range vals {
				if v != 0 {
					nonZero++
				}
			}
			if nonZero > 3 {
				return fmt.Errorf("more than 3 numbers in 2x2")
			}
			if mat[i][j] != 0 && mat[i+1][j+1] != 0 && mat[i][j] == mat[i+1][j+1] {
				return fmt.Errorf("diagonal repetition")
			}
			if mat[i][j+1] != 0 && mat[i+1][j] != 0 && mat[i][j+1] == mat[i+1][j] {
				return fmt.Errorf("diagonal repetition")
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 1; t <= 100; t++ {
		k := rand.Intn(4) + 1
		m := rand.Intn(15) + k
		counts := make([]int, k)
		for i := 0; i < m; i++ {
			idx := rand.Intn(k)
			counts[idx]++
		}
		input := fmt.Sprintf("1\n%d %d\n", m, k)
		for i, v := range counts {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", t, err)
			os.Exit(1)
		}
		if err := checkCase(m, k, counts, out); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:%soutput:%s\n", t, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
