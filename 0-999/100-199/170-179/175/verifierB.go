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

func expected(names []string, scores []int) []string {
	n := len(names)
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	// sort by name
	sort.Slice(idx, func(i, j int) bool {
		return names[idx[i]] < names[idx[j]]
	})
	unique := make([]int, 0, n)
	for i := 0; i < len(idx); {
		j := i
		best := idx[i]
		for j+1 < len(idx) && names[idx[j+1]] == names[idx[i]] {
			j++
			if scores[idx[j]] > scores[best] {
				best = idx[j]
			}
		}
		unique = append(unique, best)
		i = j + 1
	}
	idx = unique
	n = len(idx)
	sort.Slice(idx, func(i, j int) bool {
		if scores[idx[i]] == scores[idx[j]] {
			return names[idx[i]] < names[idx[j]]
		}
		return scores[idx[i]] < scores[idx[j]]
	})
	res := make([]string, n)
	for i := 0; i < n; i++ {
		s := scores[idx[i]]
		r := i
		for r+1 < n && scores[idx[r+1]] == s {
			r++
		}
		nwtr := float64(r+1) / float64(n)
		btr := float64(n-(r+1)) / float64(n)
		name := names[idx[i]]
		var cat string
		switch {
		case nwtr >= 0.99:
			cat = "pro"
		case nwtr >= 0.9 && btr > 0.01:
			cat = "hardcore"
		case nwtr >= 0.8 && btr > 0.1:
			cat = "average"
		case nwtr >= 0.5 && btr > 0.2:
			cat = "random"
		case btr > 0.5:
			cat = "noob"
		}
		res[i] = fmt.Sprintf("%s %s", name, cat)
	}
	return res
}

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
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
		fmt.Fprintf(os.Stderr, "cannot read testcasesB.txt: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "bad file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	tests := make([]string, t)
	expectedOut := make([]string, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "unexpected EOF")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		names := make([]string, n)
		scores := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			scan.Scan()
			name := scan.Text()
			scan.Scan()
			sc, _ := strconv.Atoi(scan.Text())
			names[i] = name
			scores[i] = sc
			sb.WriteString(fmt.Sprintf("%s %d\n", name, sc))
		}
		tests[caseNum] = sb.String()
		res := expected(names, scores)
		var out strings.Builder
		out.WriteString(fmt.Sprintf("%d\n", len(res)))
		for _, line := range res {
			out.WriteString(line)
			out.WriteByte('\n')
		}
		expectedOut[caseNum] = strings.TrimSpace(out.String())
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s", err, errBuf.String())
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	outScan.Split(bufio.ScanLines)
	for i := 0; i < t; i++ {
		var gotLines []string
		if !outScan.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for case %d\n", i+1)
			os.Exit(1)
		}
		m, err := strconv.Atoi(strings.TrimSpace(outScan.Text()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad output for case %d\n", i+1)
			os.Exit(1)
		}
		for j := 0; j < m; j++ {
			if !outScan.Scan() {
				fmt.Fprintf(os.Stderr, "missing output for case %d\n", i+1)
				os.Exit(1)
			}
			gotLines = append(gotLines, strings.TrimSpace(outScan.Text()))
		}
		got := fmt.Sprintf("%d\n%s", m, strings.Join(gotLines, "\n"))
		if got != expectedOut[i] {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, expectedOut[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
