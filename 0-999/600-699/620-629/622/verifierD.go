package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveD(n int) []int {
	res := make([]int, 0, 2*n)
	x := (n - 1) % 2
	if x == 0 {
		x = 2
	}
	for x != n+1 {
		res = append(res, x)
		x += 2
	}
	x -= 2
	for x > 0 {
		res = append(res, x)
		x -= 2
	}
	res = append(res, n)
	x = n % 2
	if x == 0 {
		x = 2
	}
	for x != n+2 {
		res = append(res, x)
		x += 2
	}
	x -= 4
	for x > 0 {
		res = append(res, x)
		x -= 2
	}
	return res
}

func scoreArr(arr []int, n int) (int, error) {
	if len(arr) != 2*n {
		return 0, fmt.Errorf("expected %d numbers, got %d", 2*n, len(arr))
	}
	pos := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		pos[i] = make([]int, 0, 2)
	}
	for i, v := range arr {
		if v < 1 || v > n {
			return 0, fmt.Errorf("value %d out of range [1,%d]", v, n)
		}
		pos[v] = append(pos[v], i+1)
	}
	s := 0
	for i := 1; i <= n; i++ {
		if len(pos[i]) != 2 {
			return 0, fmt.Errorf("number %d appears %d times, expected 2", i, len(pos[i]))
		}
		d := pos[i][1] - pos[i][0]
		diff := d + i - n
		if diff < 0 {
			diff = -diff
		}
		s += (n - i) * diff
	}
	return s, nil
}

func parseInts(s string) ([]int, error) {
	fields := strings.Fields(s)
	res := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(4)
	var tests []int
	tests = append(tests, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	for len(tests) < 100 {
		tests = append(tests, rand.Intn(50)+1)
	}
	for i, n := range tests {
		input := fmt.Sprintf("%d\n", n)

		refArr := solveD(n)
		refScore, err := scoreArr(refArr, n)
		if err != nil {
			fmt.Printf("test %d: reference error: %v\n", i+1, err)
			os.Exit(1)
		}

		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotArr, err := parseInts(got)
		if err != nil {
			fmt.Printf("test %d failed: n=%d: cannot parse output: %v\n", i+1, n, err)
			os.Exit(1)
		}
		gotScore, err := scoreArr(gotArr, n)
		if err != nil {
			fmt.Printf("test %d failed: n=%d: invalid output: %v\n", i+1, n, err)
			os.Exit(1)
		}
		if gotScore != refScore {
			fmt.Printf("test %d failed: n=%d: score %d, want %d\ngot: %s\n", i+1, n, gotScore, refScore, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
