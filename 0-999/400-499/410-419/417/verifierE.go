package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func isSquare(x int64) bool {
	r := int64(math.Round(math.Sqrt(float64(x))))
	return r*r == x
}

func runCase(bin string, n, m int) error {
	input := fmt.Sprintf("%d %d\n", n, m)
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	numbers := []int64{}
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		for _, f := range fields {
			var v int64
			if _, err := fmt.Sscan(f, &v); err == nil {
				numbers = append(numbers, v)
			}
		}
	}
	if len(numbers) != n*m {
		return fmt.Errorf("expected %d numbers got %d", n*m, len(numbers))
	}
	mat := make([][]int64, n)
	idx := 0
	for i := 0; i < n; i++ {
		row := make([]int64, m)
		for j := 0; j < m; j++ {
			v := numbers[idx]
			if v <= 0 || v > 1e8 {
				return fmt.Errorf("invalid value")
			}
			row[j] = v
			idx++
		}
		mat[i] = row
	}
	for i := 0; i < n; i++ {
		var sum int64
		for j := 0; j < m; j++ {
			sum += mat[i][j] * mat[i][j]
		}
		if !isSquare(sum) {
			return fmt.Errorf("row %d not square", i+1)
		}
	}
	for j := 0; j < m; j++ {
		var sum int64
		for i := 0; i < n; i++ {
			sum += mat[i][j] * mat[i][j]
		}
		if !isSquare(sum) {
			return fmt.Errorf("col %d not square", j+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		if err := runCase(bin, n, m); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d m=%d\n", i+1, err, n, m)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
