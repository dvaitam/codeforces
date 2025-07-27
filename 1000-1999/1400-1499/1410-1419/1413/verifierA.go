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

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveA(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	var sb strings.Builder
	nextInt := func() int64 {
		if !in.Scan() {
			return 0
		}
		v, _ := strconv.ParseInt(in.Text(), 10, 64)
		return v
	}
	t := nextInt()
	for ; t > 0; t-- {
		n := int(nextInt())
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = nextInt()
		}
		// reverse a
		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			a[i], a[j] = a[j], a[i]
		}
		res := make([]int64, n)
		for i := 0; i < n; i++ {
			if 2*i < n {
				res[i] = -a[i]
			} else {
				res[i] = a[i]
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(res[i], 10))
		}
		if t > 1 {
			sb.WriteByte('\n')
		}
	}
	return strings.TrimSpace(sb.String())
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(8/2)*2 + 2 // even between 2 and 10
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			v := int64(rng.Intn(200) - 100)
			if v == 0 {
				v = 1
			}
			a[j] = v
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(a[j], 10))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveA(input)
		tests[i] = testCase{input: input, expect: expect}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
