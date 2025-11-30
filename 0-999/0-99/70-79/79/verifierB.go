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

type caseB struct {
	n, m, k, t int
	waste      [][2]int
	queries    [][2]int
}

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `100
2 5 0 2
1 4
2 5
1 6 3 2
1 3
1 4
1 5
1 1
1 4
6 2 3 3
5 2
6 2
4 1
4 2
1 2
5 1
2 6 2 1
1 1
1 2
1 4
2 6 3 2
2 6
1 2
2 2
2 4
1 4
1 3 2 3
1 1
1 2
1 1
1 1
1 3
5 2 3 3
4 1
2 2
1 2
5 1
4 1
4 2
5 5 1 3
1 2
1 1
4 1
3 2
3 1 2 1
1 1
3 1
1 1
2 3 1 3
1 1
2 2
1 1
2 2
3 4 1 2
1 3
1 2
2 4
5 2 3 1
4 2
1 2
5 2
2 1
1 3 0 1
1 2
5 3 1 1
5 3
3 1
2 5 3 3
1 5
2 4
1 1
2 3
1 2
1 3
5 2 2 2
1 1
4 2
5 1
1 1
1 2 0 1
1 2
5 5 2 2
5 1
2 4
5 1
3 3
3 5 0 2
1 1
2 2
6 2 1 3
6 2
6 2
1 1
1 2
6 5 2 2
3 4
2 3
2 4
3 1
3 2 2 3
3 1
2 2
1 1
1 1
2 1
3 5 0 3
1 2
1 3
2 1
6 5 2 1
3 5
5 1
1 4
3 2 2 2
2 1
1 2
3 1
3 1
2 3 2 2
2 3
1 1
1 2
1 1
5 3 3 3
1 2
1 1
3 1
1 1
5 3
3 1
5 2 2 2
2 1
4 1
2 1
3 2
2 5 0 2
2 1
2 1
2 5 0 2
1 1
2 5
5 4 0 3
3 3
5 3
4 2
5 3 3 2
5 3
3 3
1 3
5 2
2 1
5 3 2 2
4 1
4 3
3 3
5 2
1 6 3 2
1 3
1 4
1 6
1 5
1 6
6 3 1 2
1 3
6 3
3 3
1 4 1 3
1 3
1 1
1 3
1 3
2 4 2 2
1 2
2 2
2 1
2 1
6 4 0 1
3 2
4 5 3 1
4 5
1 5
3 4
4 5
6 2 2 2
3 1
3 2
6 2
1 1
2 4 2 2
2 3
1 3
2 4
2 3
1 5 1 1
1 3
1 4
2 2 1 2
1 1
1 2
2 2
4 5 0 1
4 3
2 1 1 3
1 1
1 1
1 1
1 1
1 4 0 1
1 2
6 4 3 2
4 2
5 3
1 3
6 1
6 3
6 5 0 2
1 3
3 1
3 1 1 1
1 1
2 1
3 4 0 3
1 2
3 4
3 1
3 6 1 1
2 3
2 3
3 3 0 3
2 2
2 2
2 3
4 1 0 1
3 1
6 5 2 2
5 3
4 1
1 5
5 1
1 2 0 3
1 2
1 2
1 1
5 4 1 2
2 4
1 3
4 1
2 3 0 3
2 3
2 3
2 2
3 2 3 3
2 1
1 1
2 2
2 1
2 2
2 2
6 6 0 1
4 3
4 2 3 2
2 2
2 1
3 2
1 1
3 1
4 1 3 3
4 1
3 1
1 1
4 1
2 1
3 1
6 1 3 3
5 1
2 1
4 1
2 1
6 1
2 1
4 1 3 2
3 1
1 1
2 1
4 1
2 1
5 6 1 1
5 6
4 2
2 3 1 2
1 2
2 2
2 2
6 3 2 2
6 2
6 3
6 1
3 3
4 4 2 1
3 3
2 3
3 4
3 1 0 2
1 1
2 1
3 5 3 1
1 4
2 2
3 1
1 3
5 1 2 1
5 1
3 1
2 1
6 3 2 1
3 2
1 2
5 2
3 2 2 3
3 2
2 1
1 2
2 1
2 2
5 3 3 1
2 1
4 1
1 3
4 3
2 2 0 3
2 2
2 2
1 1
4 4 3 1
4 2
2 4
1 3
4 1
4 5 3 2
2 4
2 1
1 4
3 4
2 5
2 1 1 1
2 1
1 1
4 6 2 2
1 5
2 5
3 6
1 5
1 4 1 2
1 3
1 2
1 3
6 2 3 3
6 1
1 2
2 1
3 1
2 1
4 2
1 1 0 1
1 1
6 5 2 1
2 2
4 4
4 2
1 5 1 3
1 3
1 1
1 4
1 5
1 6 3 1
1 3
1 6
1 4
1 3
6 4 1 3
5 4
6 3
3 2
3 3
3 3 3 2
2 1
2 2
3 3
3 3
3 1
5 4 1 3
2 2
2 4
2 3
2 4
3 6 1 1
2 5
2 3
5 1 0 2
1 1
5 1
6 2 3 1
1 2
1 1
6 1
6 1
2 5 3 3
2 1
2 2
2 5
2 4
2 4
1 2
2 5 0 2
1 1
2 4
2 2 1 3
1 1
1 1
2 1
2 1
4 1 1 1
4 1
1 1
1 4 2 2
1 1
1 4
1 4
1 4
3 5 2 1
3 1
1 3
1 3
3 3 0 3
1 2
2 1
1 3
2 5 0 2
1 1
1 3
4 5 2 3
2 5
3 4
2 5
1 5
3 3
3 5 1 1
1 4
1 2`

func solveCaseB(c caseB) []string {
	waste := make([]int, c.k)
	wasteSet := make(map[int]bool)
	for i, p := range c.waste {
		pos := (p[0]-1)*c.m + p[1]
		waste[i] = pos
		wasteSet[pos] = true
	}
	sort.Ints(waste)
	res := make([]string, c.t)
	for i, q := range c.queries {
		pos := (q[0]-1)*c.m + q[1]
		if wasteSet[pos] {
			res[i] = "Waste"
		} else {
			cnt := sort.SearchInts(waste, pos)
			r := (pos - cnt) % 3
			if r == 1 {
				res[i] = "Carrots"
			} else if r == 2 {
				res[i] = "Kiwis"
			} else {
				res[i] = "Grapes"
			}
		}
	}
	return res
}

func parseTestcases() ([]caseB, [][]string, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, nil, fmt.Errorf("empty testcase data")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, nil, fmt.Errorf("parse t: %w", err)
	}
	cases := make([]caseB, 0, t)
	expected := make([][]string, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		readInt := func(label string) (int, error) {
			if !scan.Scan() {
				return 0, fmt.Errorf("case %d: missing %s", caseIdx+1, label)
			}
			val, err := strconv.Atoi(scan.Text())
			if err != nil {
				return 0, fmt.Errorf("case %d: parse %s: %w", caseIdx+1, label, err)
			}
			return val, nil
		}
		n, err := readInt("n")
		if err != nil {
			return nil, nil, err
		}
		m, err := readInt("m")
		if err != nil {
			return nil, nil, err
		}
		k, err := readInt("k")
		if err != nil {
			return nil, nil, err
		}
		tq, err := readInt("t")
		if err != nil {
			return nil, nil, err
		}
		tc := caseB{n: n, m: m, k: k, t: tq}
		tc.waste = make([][2]int, k)
		for i := 0; i < k; i++ {
			a, err := readInt("waste a")
			if err != nil {
				return nil, nil, err
			}
			b, err := readInt("waste b")
			if err != nil {
				return nil, nil, err
			}
			tc.waste[i] = [2]int{a, b}
		}
		tc.queries = make([][2]int, tq)
		for i := 0; i < tq; i++ {
			x, err := readInt("query x")
			if err != nil {
				return nil, nil, err
			}
			y, err := readInt("query y")
			if err != nil {
				return nil, nil, err
			}
			tc.queries[i] = [2]int{x, y}
		}
		cases = append(cases, tc)
		expected = append(expected, solveCaseB(tc))
	}
	return cases, expected, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cases, expected, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to load embedded testcases:", err)
		os.Exit(1)
	}
	for ci, c := range cases {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d %d %d\n", c.n, c.m, c.k, c.t)
		for _, w := range c.waste {
			fmt.Fprintf(&buf, "%d %d\n", w[0], w[1])
		}
		for _, q := range c.queries {
			fmt.Fprintf(&buf, "%d %d\n", q[0], q[1])
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("execution failed on case %d: %v\n", ci+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		for i, exp := range expected[ci] {
			if !outScan.Scan() {
				fmt.Printf("missing output for case %d query %d\n", ci+1, i+1)
				os.Exit(1)
			}
			got := outScan.Text()
			if got != exp {
				fmt.Printf("case %d query %d failed: expected %s got %s\n", ci+1, i+1, exp, got)
				os.Exit(1)
			}
		}
		if outScan.Scan() {
			fmt.Printf("extra output detected in case %d\n", ci+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
