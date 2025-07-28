package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Test struct {
	n int
	m int
	a [][]int
}

func generateTests() []Test {
	rand.Seed(6)
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(4) + 1
		m := rand.Intn(4) + 1
		a := make([][]int, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int, m)
			for j := 0; j < m; j++ {
				a[i][j] = rand.Intn(5)
			}
		}
		tests = append(tests, Test{n, m, a})
	}
	return tests
}

func primeFactorsDistinct(x int) []int {
	res := []int{}
	for p := 2; p*p <= x; p++ {
		if x%p == 0 {
			res = append(res, p)
			for x%p == 0 {
				x /= p
			}
		}
	}
	if x > 1 {
		res = append(res, x)
	}
	return res
}

func checkRows(a [][]int, period int) bool {
	n := len(a)
	m := len(a[0])
	for i := period; i < n; i++ {
		for j := 0; j < m; j++ {
			if a[i][j] != a[i-period][j] {
				return false
			}
		}
	}
	return true
}

func checkCols(a [][]int, period int) bool {
	n := len(a)
	m := len(a[0])
	for j := period; j < m; j++ {
		for i := 0; i < n; i++ {
			if a[i][j] != a[i][j-period] {
				return false
			}
		}
	}
	return true
}

func numDivisors(x int) int {
	if x == 0 {
		return 0
	}
	res := 1
	for p := 2; p*p <= x; p++ {
		if x%p == 0 {
			e := 0
			for x%p == 0 {
				x /= p
				e++
			}
			res *= e + 1
		}
	}
	if x > 1 {
		res *= 2
	}
	return res
}

func solve(t Test) string {
	a := t.a
	n := t.n
	m := t.m
	rowPeriod := n
	for _, p := range primeFactorsDistinct(rowPeriod) {
		for rowPeriod%p == 0 && checkRows(a, rowPeriod/p) {
			rowPeriod /= p
		}
	}
	colPeriod := m
	for _, p := range primeFactorsDistinct(colPeriod) {
		for colPeriod%p == 0 && checkCols(a, colPeriod/p) {
			colPeriod /= p
		}
	}
	cntRows := numDivisors(n / rowPeriod)
	cntCols := numDivisors(m / colPeriod)
	return fmt.Sprintf("%d", cntRows*cntCols)
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var in strings.Builder
		fmt.Fprintf(&in, "%d %d\n", t.n, t.m)
		for x := 0; x < t.n; x++ {
			for y := 0; y < t.m; y++ {
				fmt.Fprintf(&in, "%d ", t.a[x][y])
			}
			in.WriteByte('\n')
		}
		expect := solve(t)
		got, err := run(binary, in.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != expect {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, in.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
	time.Sleep(0)
}
