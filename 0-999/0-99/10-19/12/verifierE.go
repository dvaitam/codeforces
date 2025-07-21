package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func checkMatrix(n int, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(output)))
	matrix := make([][]int, 0, n)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		if len(fields) != n {
			return fmt.Errorf("expected %d numbers per line", n)
		}
		row := make([]int, n)
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil || v < 0 || v >= n {
				return fmt.Errorf("invalid value %q", f)
			}
			row[i] = v
		}
		matrix = append(matrix, row)
	}
	if len(matrix) != n {
		return fmt.Errorf("expected %d lines, got %d", n, len(matrix))
	}
	for i := 0; i < n; i++ {
		if matrix[i][i] != 0 {
			return fmt.Errorf("diagonal element at %d,%d not zero", i, i)
		}
		seen := make(map[int]bool)
		for j := 0; j < n; j++ {
			if seen[matrix[i][j]] {
				return fmt.Errorf("row %d has duplicates", i)
			}
			seen[matrix[i][j]] = true
			if matrix[i][j] != matrix[j][i] {
				return fmt.Errorf("matrix not symmetric at %d,%d", i, j)
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := (rng.Intn(10) + 1) * 2 // even between 2 and 20
		input := fmt.Sprintf("%d\n", n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if err := checkMatrix(n, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", t+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
