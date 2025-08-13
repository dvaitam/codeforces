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
	b := make([]int, n+1)
	used := make([]bool, n+1)
	rev := make([]int, n+1)

	for i := 1; i <= n; i++ {
		v := a[i]
		if !used[v] {
			b[i] = v
			used[v] = true
			rev[v] = i
		}
	}

	remI := make([]int, 0)
	remT := make([]int, 0)
	for i := 1; i <= n; i++ {
		if b[i] == 0 {
			remI = append(remI, i)
		}
		if !used[i] {
			remT = append(remT, i)
		}
	}

	if len(remI) == 1 {
		i0 := remI[0]
		v0 := remT[0]
		if i0 == v0 {
			j := rev[a[i0]]
			b[j] = v0
			b[i0] = a[i0]
		} else {
			b[i0] = v0
		}
	} else if len(remI) > 0 {
		arrT := append([]int(nil), remT...)
		posConf := make([]int, 0)
		for k := 0; k < len(remI); k++ {
			if remI[k] == arrT[k] {
				posConf = append(posConf, k)
			}
		}
		if len(posConf) == 1 {
			p := posConf[0]
			q := 0
			if p == q {
				q = 1
			}
			arrT[p], arrT[q] = arrT[q], arrT[p]
		} else if len(posConf) >= 2 {
			first := posConf[0]
			prev := arrT[first]
			for _, pos := range posConf[1:] {
				arrT[pos], prev = prev, arrT[pos]
			}
			arrT[first] = prev
		}

		for k := 0; k < len(remI); k++ {
			if remI[k] == arrT[k] {
				swapK := 0
				if k == 0 {
					swapK = 1
				}
				arrT[k], arrT[swapK] = arrT[swapK], arrT[k]
			}
		}

		for k := 0; k < len(remI); k++ {
			b[remI[k]] = arrT[k]
		}
	}

	matched := 0
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		if b[i] == a[i] {
			matched++
		}
		res[i-1] = b[i]
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
