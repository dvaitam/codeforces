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

type test struct {
	n   int
	arr []int
}

func genTests() []test {
	rand.Seed(1)
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(97) + 3 // 3..99
		m := map[int]bool{}
		arr := make([]int, 0, n)
		// ensure one negative
		for {
			x := -rand.Intn(1000) - 1
			if !m[x] {
				m[x] = true
				arr = append(arr, x)
				break
			}
		}
		// ensure one positive
		for {
			x := rand.Intn(1000) + 1
			if !m[x] {
				m[x] = true
				arr = append(arr, x)
				break
			}
		}
		// ensure one zero
		if !m[0] {
			m[0] = true
			arr = append(arr, 0)
		}
		for len(arr) < n {
			x := rand.Intn(2001) - 1000
			if !m[x] {
				m[x] = true
				arr = append(arr, x)
			}
		}
		rand.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
		tests = append(tests, test{n, arr})
	}
	return tests
}

func buildInput(t test) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func signProduct(nums []int) int {
	sign := 1
	for _, v := range nums {
		if v == 0 {
			return 0
		}
		if v < 0 {
			sign = -sign
		}
	}
	return sign
}

func verifyOutput(out string, t test) bool {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	vals := []int{}
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return false
		}
		vals = append(vals, v)
	}
	if len(vals) == 0 {
		return false
	}
	idx := 0
	if idx >= len(vals) {
		return false
	}
	n1 := vals[idx]
	idx++
	if idx+n1 > len(vals) {
		return false
	}
	s1 := vals[idx : idx+n1]
	idx += n1
	if idx >= len(vals) {
		return false
	}
	n2 := vals[idx]
	idx++
	if idx+n2 > len(vals) {
		return false
	}
	s2 := vals[idx : idx+n2]
	idx += n2
	if idx >= len(vals) {
		return false
	}
	n3 := vals[idx]
	idx++
	if idx+n3 > len(vals) {
		return false
	}
	s3 := vals[idx : idx+n3]
	idx += n3
	if idx != len(vals) {
		return false
	}
	if n1+n2+n3 != t.n {
		return false
	}
	counts := make(map[int]int)
	for _, v := range s1 {
		counts[v]++
	}
	for _, v := range s2 {
		counts[v]++
	}
	for _, v := range s3 {
		counts[v]++
	}
	if len(counts) != t.n {
		return false
	}
	for _, v := range t.arr {
		if counts[v] != 1 {
			return false
		}
	}
	if signProduct(s1) >= 0 {
		return false
	}
	if signProduct(s2) <= 0 {
		return false
	}
	if signProduct(s3) != 0 {
		return false
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := buildInput(t)
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: error running candidate: %v\n", i+1, err)
			os.Exit(1)
		}
		if !verifyOutput(out, t) {
			fmt.Printf("test %d failed. input:\n%s\noutput:\n%s\n", i+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
