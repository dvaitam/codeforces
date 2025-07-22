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

func computeCase(n, m, k int) (int, []int) {
	if n == 2 && m == 2 && k == 1 {
		return 4, []int{1, 2, 3, 4}
	}
	q := 0
	if n+n == n+m+1 && k == 1 {
		q = 1
	}
	if n+n < n+m+1 {
		q = 1
	}
	total := k*2 + q
	res := make([]int, 0, total)
	for i := 0; i < k; i++ {
		res = append(res, 1)
	}
	res = append(res, n)
	for i := 0; i < k-1; i++ {
		res = append(res, n+1)
	}
	a := 0
	if n == m && k == 1 {
		a = 1
	}
	if q == 1 {
		res = append(res, n+m-a)
	}
	return total, res
}

func buildInput(n, m, k int) string {
	return fmt.Sprintf("%d %d %d\n", n, m, k)
}

func parseInts(s string) ([]int, error) {
	scan := strings.Fields(s)
	vals := make([]int, len(scan))
	for i, f := range scan {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	return vals, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		n := rand.Intn(9) + 2   // 2..10
		m := rand.Intn(n-1) + 1 // 1..n
		k := rand.Intn(n) + 1   // 1..n
		input := buildInput(n, m, k)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
		vals, err := parseInts(out)
		if err != nil || len(vals) == 0 {
			fmt.Fprintf(os.Stderr, "case %d: bad output\n", i)
			os.Exit(1)
		}
		expZ, expArr := computeCase(n, m, k)
		if vals[0] != expZ || len(vals)-1 != expZ {
			fmt.Fprintf(os.Stderr, "case %d: expected %d numbers, got %d\n", i, expZ, len(vals)-1)
			os.Exit(1)
		}
		for j := 0; j < expZ; j++ {
			if vals[j+1] != expArr[j] {
				fmt.Fprintf(os.Stderr, "case %d: mismatch\n", i)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
