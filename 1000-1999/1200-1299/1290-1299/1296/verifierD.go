package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expect(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	readInt := func() int {
		in.Scan()
		v := 0
		fmt.Sscan(in.Text(), &v)
		return v
	}
	n := readInt()
	a := readInt()
	b := readInt()
	k := readInt()
	h := make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = readInt()
	}
	cycle := a + b
	need := make([]int, n)
	for i, hp := range h {
		r := hp % cycle
		if r == 0 {
			r = cycle
		}
		need[i] = (r - 1) / a
	}
	sort.Ints(need)
	points := 0
	for _, v := range need {
		if k < v {
			break
		}
		k -= v
		points++
	}
	return fmt.Sprintf("%d", points)
}

type testCase struct{ input string }

var tests = []testCase{
	{"3 1 2 1\n1 2 3\n"},
	{"3 1 2 3\n3 3 3\n"},
	{"5 2 3 2\n6 6 6 6 6\n"},
	{"4 2 1 3\n5 8 3 4\n"},
	{"1 2 2 1\n5\n"},
	{"5 5 3 4\n9 11 8 7 20\n"},
	{"5 1 5 3\n6 6 6 6 6\n"},
	{"3 10 10 1\n21 31 41\n"},
	{"2 7 3 2\n17 20\n"},
	{"4 3 4 2\n8 12 18 22\n"},
	{"3 2 2 3\n9 5 7\n"},
	{"4 2 5 2\n10 9 8 7\n"},
	{"5 2 3 1\n1 2 3 4 5\n"},
	{"2 4 3 5\n8 7\n"},
	{"2 2 1 1\n2 3\n"},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range tests {
		exp := expect(tc.input)
		got, err := runProg(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\ngot:%s\ninput:\n%s", i+1, exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
