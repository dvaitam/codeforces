package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(99) + 2 // 2..100
	return fmt.Sprintf("1\n%d\n", n)
}

func isPrime(x int) bool {
	if x < 2 {
		return false
	}
	r := int(math.Sqrt(float64(x)))
	for i := 2; i <= r; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func check(out string, n int) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != n {
		return fmt.Errorf("expected %d lines, got %d", n, len(lines))
	}
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		fields := strings.Fields(lines[i])
		if len(fields) != n {
			return fmt.Errorf("line %d: expected %d numbers, got %d", i+1, n, len(fields))
		}
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[j])
			if err != nil {
				return fmt.Errorf("line %d column %d: invalid integer", i+1, j+1)
			}
			if v < 0 || v > 100000 {
				return fmt.Errorf("value %d out of range at (%d,%d)", v, i+1, j+1)
			}
			if isPrime(v) {
				return fmt.Errorf("value %d is prime at (%d,%d)", v, i+1, j+1)
			}
			mat[i][j] = v
		}
	}
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < n; j++ {
			sum += mat[i][j]
		}
		if !isPrime(sum) {
			return fmt.Errorf("row %d sum %d is not prime", i+1, sum)
		}
	}
	for j := 0; j < n; j++ {
		sum := 0
		for i := 0; i < n; i++ {
			sum += mat[i][j]
		}
		if !isPrime(sum) {
			return fmt.Errorf("column %d sum %d is not prime", j+1, sum)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		in := genCase(rng)
		got, err := run(exe, in)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		var n int
		fmt.Sscanf(in, "1\n%d\n", &n)
		if err := check(got, n); err != nil {
			fmt.Printf("wrong answer on test %d\ninput:\n%sgot:\n%s\nerror: %v", i+1, in, got, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
