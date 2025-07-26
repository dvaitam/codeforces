package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
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
	in.Scan()
	s := in.Text()
	colors := make([]byte, n)
	last := [2]byte{'a' - 1, 'a' - 1}
	for i := 0; i < n; i++ {
		c := s[i]
		if c >= last[0] {
			colors[i] = '0'
			last[0] = c
		} else if c >= last[1] {
			colors[i] = '1'
			last[1] = c
		} else {
			return "NO"
		}
	}
	return fmt.Sprintf("YES\n%s", string(colors))
}

type testCase struct{ input string }

var tests = []testCase{
	{"1\na\n"},
	{"2\nba\n"},
	{"3\nabc\n"},
	{"3\ncba\n"},
	{"4\nabac\n"},
	{"4\naaaa\n"},
	{"4\nabcd\n"},
	{"4\ndcba\n"},
	{"3\nbca\n"},
	{"6\nacdbbd\n"},
	{"3\nzyx\n"},
	{"3\naaz\n"},
	{"6\nabcabc\n"},
	{"6\nbacdef\n"},
	{"4\ncbba\n"},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range tests {
		exp := expect(tc.input)
		got, err := runProg(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\ngot:\n%s\ninput:%s", i+1, exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
