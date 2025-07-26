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

type segment struct{ l, r, id int }

func solveCase(segs []segment) string {
	n := len(segs)
	sort.Slice(segs, func(i, j int) bool {
		if segs[i].l != segs[j].l {
			return segs[i].l < segs[j].l
		}
		return segs[i].r < segs[j].r
	})
	clusters := 1
	curR := segs[0].r
	for i := 1; i < n; i++ {
		if segs[i].l <= curR {
			if segs[i].r > curR {
				curR = segs[i].r
			}
		} else {
			clusters++
			if segs[i].r > curR {
				curR = segs[i].r
			}
		}
	}
	if clusters < 2 {
		return "-1"
	}
	out := make([]int, n)
	group := 1
	curR = segs[0].r
	out[segs[0].id] = group
	for i := 1; i < n; i++ {
		if segs[i].l <= curR {
			out[segs[i].id] = group
			if segs[i].r > curR {
				curR = segs[i].r
			}
		} else {
			group++
			if group > 2 {
				group = 2
			}
			out[segs[i].id] = group
			if segs[i].r > curR {
				curR = segs[i].r
			}
		}
	}
	var buf bytes.Buffer
	for i, v := range out {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, v)
	}
	return buf.String()
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
	expected := make([]string, t)
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
			segs[j] = segment{l: l, r: r, id: j}
		}
		expected[i] = solveCase(segs)
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
		got := strings.TrimSpace(outScan.Text())
		if got != expected[i] {
			fmt.Printf("test %d failed:\nexpected: %s\ngot: %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
