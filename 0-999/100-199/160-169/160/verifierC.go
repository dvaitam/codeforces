package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseC struct {
	n   int
	k   int64
	arr []int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(44)
	var tests []testCaseC
	for i := 0; i < 100; i++ {
		n := rand.Intn(50) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(2001) - 1000
		}
		k := rand.Int63n(int64(n)*int64(n)) + 1
		tests = append(tests, testCaseC{n, k, arr})
	}
	tests = append(tests, testCaseC{1, 1, []int{5}})

	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
		for j, v := range t.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveC(strings.NewReader(input))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveC(r io.Reader) string {
	in := bufio.NewReader(r)
	var n int
	var k int64
	fmt.Fscan(in, &n, &k)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	sort.Ints(a)
	idx1 := int((k - 1) / int64(n))
	idx2 := int((k - 1) % int64(n))
	return fmt.Sprintf("%d %d\n", a[idx1], a[idx2])
}
