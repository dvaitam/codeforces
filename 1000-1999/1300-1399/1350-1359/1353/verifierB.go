package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testB struct {
	n int
	k int
	a []int
	b []int
}

func genTestsB() []testB {
	rand.Seed(2353)
	tests := make([]testB, 0, 100)
	// some edge tests
	tests = append(tests, testB{n: 1, k: 0, a: []int{1}, b: []int{2}})
	tests = append(tests, testB{n: 1, k: 1, a: []int{5}, b: []int{10}})
	tests = append(tests, testB{n: 2, k: 1, a: []int{1, 2}, b: []int{3, 4}})
	for len(tests) < 100 {
		n := rand.Intn(10) + 1
		k := rand.Intn(n + 1)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(50)
			b[i] = rand.Intn(50)
		}
		tests = append(tests, testB{n: n, k: k, a: a, b: b})
	}
	return tests[:100]
}

func expectedB(tc testB) int {
	a := append([]int(nil), tc.a...)
	b := append([]int(nil), tc.b...)
	sort.Ints(a)
	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
	for i := 0; i < tc.k && i < tc.n; i++ {
		if a[i] < b[i] {
			a[i] = b[i]
		} else {
			break
		}
	}
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTestsB()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i2, v := range tc.a {
			if i2 > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i2, v := range tc.b {
			if i2 > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		exp := fmt.Sprintf("%d", expectedB(tc))
		got, err := runProg(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
