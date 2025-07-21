package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

type caseB struct {
	n, m, k, t int
	waste      [][2]int
	queries    [][2]int
}

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(scan.Text())
	var cases []caseB
	var expected [][]string
	for tt := 0; tt < T; tt++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		t, _ := strconv.Atoi(scan.Text())
		c := caseB{n: n, m: m, k: k, t: t}
		c.waste = make([][2]int, k)
		for i := 0; i < k; i++ {
			scan.Scan()
			a, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			b, _ := strconv.Atoi(scan.Text())
			c.waste[i] = [2]int{a, b}
		}
		c.queries = make([][2]int, t)
		for i := 0; i < t; i++ {
			scan.Scan()
			x, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			y, _ := strconv.Atoi(scan.Text())
			c.queries[i] = [2]int{x, y}
		}
		cases = append(cases, c)
		expected = append(expected, solveCaseB(c))
	}
	idx := 0
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
			idx++
		}
		if outScan.Scan() {
			fmt.Printf("extra output detected in case %d\n", ci+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
