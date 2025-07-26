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
		val := 0
		fmt.Sscan(in.Text(), &val)
		return val
	}
	t := readInt()
	var sb strings.Builder
	for i := 0; i < t; i++ {
		n := readInt()
		sum := 0
		odd, even := 0, 0
		for j := 0; j < n; j++ {
			x := readInt()
			sum += x
			if x%2 == 0 {
				even++
			} else {
				odd++
			}
		}
		ans := "NO"
		if sum%2 == 1 || (odd > 0 && even > 0) {
			ans = "YES"
		}
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(ans)
	}
	return sb.String()
}

type testCase struct{ input string }

var tests = []testCase{
	{"1\n1\n1\n"},
	{"1\n1\n2\n"},
	{"1\n2\n1 3\n"},
	{"1\n2\n2 2\n"},
	{"1\n2\n1 2\n"},
	{"1\n3\n2 2 1\n"},
	{"1\n4\n1 1 1 1\n"},
	{"1\n4\n2 2 2 2\n"},
	{"1\n2\n3 3\n"},
	{"1\n3\n3 2 5\n"},
	{"1\n4\n2 2 2 1\n"},
	{"1\n5\n2 4 6 1 1\n"},
	{"1\n1\n999\n"},
	{"1\n1\n1000\n"},
	{"3\n2\n1 2\n3\n1 1 1\n2\n2 4\n"},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
