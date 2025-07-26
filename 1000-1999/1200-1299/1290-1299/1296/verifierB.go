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
	readInt := func() int64 {
		in.Scan()
		var v int64
		fmt.Sscan(in.Text(), &v)
		return v
	}
	t := readInt()
	var sb strings.Builder
	for i := int64(0); i < t; i++ {
		s := readInt()
		spent := int64(0)
		for s >= 10 {
			spent += (s / 10) * 10
			s = s/10 + s%10
		}
		spent += s
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d", spent))
	}
	return sb.String()
}

type testCase struct{ input string }

var tests = []testCase{
	{"1\n1\n"},
	{"1\n9\n"},
	{"1\n10\n"},
	{"1\n19\n"},
	{"1\n99\n"},
	{"1\n100\n"},
	{"1\n873\n"},
	{"1\n123456\n"},
	{"1\n1000000000\n"},
	{"1\n15\n"},
	{"1\n10000\n"},
	{"1\n999999\n"},
	{"1\n42\n"},
	{"1\n500000000\n"},
	{"1\n987654321\n"},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\ngot:%s\ninput:%s", i+1, exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
