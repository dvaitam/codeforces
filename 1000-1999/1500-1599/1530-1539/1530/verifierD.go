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

type testD struct {
	n int
	a []int
}

func genTestsD() []testD {
	rand.Seed(1530004)
	tests := make([]testD, 100)
	for i := range tests {
		n := rand.Intn(20) + 1
		a := make([]int, n+1)
		for j := 1; j <= n; j++ {
			a[j] = rand.Intn(n) + 1
		}
		tests[i] = testD{n: n, a: a}
	}
	return tests
}

func solveD(tc testD) (int, []int) {
	n := tc.n
	a := tc.a
	ans := make([]int, n+1)
	vis := make([]bool, n+1)
	matched := 0
	for i := 1; i <= n; i++ {
		if !vis[a[i]] {
			vis[a[i]] = true
			ans[i] = a[i]
			matched++
		}
	}
	st := make([]int, 0, n)
	for i := n; i >= 1; i-- {
		if !vis[i] {
			st = append(st, i)
		}
	}
	for i := 1; i <= n; i++ {
		if ans[i] == 0 {
			last := st[len(st)-1]
			st = st[:len(st)-1]
			ans[i] = last
		}
	}
	las := 0
	for i := 1; i <= n; i++ {
		if ans[i] == i {
			if las == 0 {
				las = i
			} else {
				ans[i], ans[las] = ans[las], ans[i]
			}
		}
	}
	if las != 0 {
		for i := 1; i <= n; i++ {
			if a[i] == a[las] {
				ans[i], ans[las] = ans[las], ans[i]
				break
			}
		}
	}
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		res[i-1] = ans[i]
	}
	return matched, res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for j := 1; j <= tc.n; j++ {
			if j > 1 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, tc.a[j])
		}
		input.WriteByte('\n')
	}

	expMatch := make([]int, len(tests))
	expAns := make([][]int, len(tests))
	for i, tc := range tests {
		m, arr := solveD(tc)
		expMatch[i] = m
		expAns[i] = arr
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s\n", err, stderr.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil || val != expMatch[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
		for j := 0; j < tests[i].n; j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			v, err := strconv.Atoi(scanner.Text())
			if err != nil || v != expAns[i][j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
