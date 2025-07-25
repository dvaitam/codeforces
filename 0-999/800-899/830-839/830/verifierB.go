package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Add(i, v int) {
	for i <= f.n {
		f.tree[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for i > 0 {
		s += f.tree[i]
		i -= i & -i
	}
	return s
}

func (f *Fenwick) RangeSum(l, r int) int {
	if l <= r {
		return f.Sum(r) - f.Sum(l-1)
	}
	return f.Sum(f.n) - f.Sum(l-1) + f.Sum(r)
}

func solveCase(a []int) int64 {
	n := len(a)
	posMap := make(map[int][]int)
	for i := 0; i < n; i++ {
		posMap[a[i]] = append(posMap[a[i]], i+1)
	}
	vals := make([]int, 0, len(posMap))
	for v := range posMap {
		vals = append(vals, v)
	}
	sort.Ints(vals)

	fw := NewFenwick(n)
	for i := 1; i <= n; i++ {
		fw.Add(i, 1)
	}
	cur := 1
	var ops int64
	for _, v := range vals {
		idxs := posMap[v]
		start := sort.Search(len(idxs), func(i int) bool { return idxs[i] >= cur })
		for i := start; i < len(idxs); i++ {
			idx := idxs[i]
			ops += int64(fw.RangeSum(cur, idx))
			fw.Add(idx, -1)
			cur = idx
		}
		for i := 0; i < start; i++ {
			idx := idxs[i]
			ops += int64(fw.RangeSum(cur, idx))
			fw.Add(idx, -1)
			cur = idx
		}
	}
	return ops
}

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read testcasesB.txt:", err)
		os.Exit(1)
	}

	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())

	cases := make([][]int, t)
	expected := make([]int64, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			val, _ := strconv.Atoi(scan.Text())
			arr[j] = val
		}
		cases[i] = arr
		expected[i] = solveCase(arr)
	}

	for i, arr := range cases {
		var input bytes.Buffer
		fmt.Fprintln(&input, len(arr))
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')

		outStr, err := run(bin, input.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(outStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output\n", i+1)
			os.Exit(1)
		}
		if got != expected[i] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", t)
}
