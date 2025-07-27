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
)

type testCase struct {
	input  string
	expect string
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveD(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	var n int
	if in.Scan() {
		n, _ = strconv.Atoi(in.Text())
	}
	total := 2 * n
	opsType := make([]bool, total)
	opsVal := make([]int, total)
	for i := 0; i < total; i++ {
		in.Scan()
		s := in.Text()
		if s == "+" {
			opsType[i] = true
		} else {
			opsType[i] = false
			in.Scan()
			x, _ := strconv.Atoi(in.Text())
			opsVal[i] = x
		}
	}
	res := make([]int, n)
	stack := make([]int, 0, n)
	pi := n - 1
	for i := total - 1; i >= 0; i-- {
		if !opsType[i] {
			x := opsVal[i]
			if len(stack) > 0 && x > stack[len(stack)-1] {
				return "NO"
			}
			stack = append(stack, x)
		} else {
			if len(stack) == 0 {
				return "NO"
			}
			x := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res[pi] = x
			pi--
		}
	}
	if pi != -1 {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(res[i]))
	}
	return sb.String()
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(4))
	tests := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 1
		total := 2 * n
		plusCount := 0
		minusCount := 0
		ops := make([]string, total)
		for i := 0; i < total; i++ {
			if plusCount == n {
				val := rng.Intn(n) + 1
				ops[i] = fmt.Sprintf("- %d", val)
				minusCount++
			} else if minusCount == n {
				ops[i] = "+"
				plusCount++
			} else if rng.Intn(2) == 0 {
				ops[i] = "+"
				plusCount++
			} else {
				val := rng.Intn(n) + 1
				ops[i] = fmt.Sprintf("- %d", val)
				minusCount++
			}
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i := 0; i < total; i++ {
			sb.WriteString(ops[i])
			if i+1 < total {
				sb.WriteByte('\n')
			}
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveD(input)
		tests[t] = testCase{input: input, expect: expect}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			return
		}
		if out != tc.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expect, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
