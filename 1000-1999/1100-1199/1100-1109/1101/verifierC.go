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

type segment struct{ l, r int }

type testCase struct {
	segs     []segment
	possible bool
}

func canSplit(segs []segment) bool {
	sort.Slice(segs, func(i, j int) bool {
		if segs[i].l != segs[j].l {
			return segs[i].l < segs[j].l
		}
		return segs[i].r < segs[j].r
	})
	curR := segs[0].r
	for i := 1; i < len(segs); i++ {
		if segs[i].l > curR {
			return true
		}
		if segs[i].r > curR {
			curR = segs[i].r
		}
	}
	return false
}

func validate(line string, tc testCase) error {
	line = strings.TrimSpace(line)
	if line == "-1" {
		if tc.possible {
			return fmt.Errorf("partition exists but got -1")
		}
		return nil
	}
	if !tc.possible {
		return fmt.Errorf("partition not possible but got result")
	}
	fields := strings.Fields(line)
	if len(fields) != len(tc.segs) {
		return fmt.Errorf("expected %d numbers, got %d", len(tc.segs), len(fields))
	}
	groups := make([]int, len(tc.segs))
	has1, has2 := false, false
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		if v != 1 && v != 2 {
			return fmt.Errorf("invalid group %d", v)
		}
		groups[i] = v
		if v == 1 {
			has1 = true
		} else {
			has2 = true
		}
	}
	if !has1 || !has2 {
		return fmt.Errorf("both groups must be non-empty")
	}
	segs := tc.segs
	for i := 0; i < len(segs); i++ {
		for j := i + 1; j < len(segs); j++ {
			if groups[i] != groups[j] {
				if !(segs[i].r < segs[j].l || segs[j].r < segs[i].l) {
					return fmt.Errorf("segments %d and %d intersect but have different groups", i+1, j+1)
				}
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) > 2 {
		bin = os.Args[2]
	}
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		segs := make([]segment, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			l, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			r, _ := strconv.Atoi(scan.Text())
			segs[j] = segment{l: l, r: r}
		}
		segCopy := make([]segment, n)
		copy(segCopy, segs)
		cases[i] = testCase{segs: segs, possible: canSplit(segCopy)}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanLines)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		line := outScan.Text()
		if err := validate(line, cases[i]); err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
