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

type testCaseF struct {
	n int
	k int
	a []int
}

func genTestsF() []testCaseF {
	rand.Seed(47)
	tests := make([]testCaseF, 20)
	for i := range tests {
		n := rand.Intn(3) + 2 // 2..4
		a := make([]int, n)
		for j := range a {
			a[j] = rand.Intn(10) + 1
		}
		tests[i] = testCaseF{n: n, k: rand.Intn(3) + 1, a: a}
	}
	return tests
}

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

func solveF(tc testCaseF) int {
	n := tc.n
	k := tc.k
	b := make([][]int, n)
	for i := range b {
		b[i] = make([]int, n)
		for j := 0; j < n; j++ {
			idx := (j - i + n) % n
			b[i][j] = tc.a[idx]
		}
	}
	total := n * n
	visited := make([]bool, total)
	comp := 0
	for idx := 0; idx < total; idx++ {
		if visited[idx] {
			continue
		}
		comp++
		queue := []int{idx}
		visited[idx] = true
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			x1 := v / n
			y1 := v % n
			for w := 0; w < total; w++ {
				if visited[w] {
					continue
				}
				x2 := w / n
				y2 := w % n
				if abs(x1-x2)+abs(y1-y2) <= k && gcd(b[x1][y1], b[x2][y2]) > 1 {
					visited[w] = true
					queue = append(queue, w)
				}
			}
		}
	}
	return comp
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	expected := make([]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveF(tc)
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
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if val != expected[i] {
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
