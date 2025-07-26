package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type testF struct {
	a []int
}

func genTestsF() []testF {
	rand.Seed(122006)
	tests := make([]testF, 100)
	for i := range tests {
		n := rand.Intn(8) + 2
		a := make([]int, n)
		pos := rand.Intn(n)
		for j := range a {
			if j == pos {
				a[j] = 1
			} else {
				a[j] = rand.Intn(9) + 2
			}
		}
		tests[i] = testF{a: a}
	}
	return tests
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveF(tc testF) (int, int) {
	n := len(tc.a)
	a := tc.a
	root := 0
	for i, v := range a {
		if v == 1 {
			root = i
			break
		}
	}
	b := make([][]int, 2)
	b[0] = make([]int, n)
	b[1] = make([]int, n)
	for i := 0; i < n; i++ {
		b[0][i] = a[(root+i)%n]
		b[1][i] = a[(root+n-i)%n]
	}
	dpl := make([][]int, 2)
	dpr := make([][]int, 2)
	ans := make([][]int, 2)
	for t := 0; t < 2; t++ {
		dpl[t] = make([]int, n)
		dpr[t] = make([]int, n)
		ans[t] = make([]int, n)
	}
	for t := 0; t < 2; t++ {
		stack := make([]int, 0, n)
		dpl[t][0] = 1
		ans[t][0] = 1
		stack = append(stack, 0)
		for i := 1; i < n; i++ {
			nex := -1
			for len(stack) > 0 && b[t][stack[len(stack)-1]] > b[t][i] {
				j := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if nex == -1 {
					dpr[t][j] = dpl[t][j]
					nex = j
				} else {
					dpr[t][j] = max(dpr[t][nex]+1, dpl[t][j])
					nex = j
				}
			}
			if nex == -1 {
				dpl[t][i] = 1
			} else {
				dpl[t][i] = dpr[t][nex] + 1
			}
			stack = append(stack, i)
			cur := len(stack) + dpl[t][i] - 1
			ans[t][i] = max(cur, ans[t][i-1])
		}
	}
	ret := n + n
	k := 0
	for i := 0; i < n; i++ {
		val := max(ans[0][i], ans[1][n-1-i])
		if val < ret {
			ret = val
			k = (root + 1 + i) % n
		}
	}
	return ret, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, len(tc.a))
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}

	expectedRet := make([]int, len(tests))
	expectedK := make([]int, len(tests))
	for i, tc := range tests {
		r, k := solveF(tc)
		expectedRet[i] = r
		expectedK[i] = k
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		r, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		k, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if r != expectedRet[i] || k != expectedK[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
