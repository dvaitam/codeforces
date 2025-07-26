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
	dp := make([]int, 26)
	colors := make([]int, n)
	maxColor := 0
	for i := 0; i < n; i++ {
		ch := int(s[i] - 'a')
		best := 0
		for j := ch + 1; j < 26; j++ {
			if dp[j] > best {
				best = dp[j]
			}
		}
		col := best + 1
		colors[i] = col
		if dp[ch] < col {
			dp[ch] = col
		}
		if col > maxColor {
			maxColor = col
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", maxColor))
	for i, c := range colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c))
	}
	return sb.String()
}

type testCase struct{ input string }

var tests = []testCase{
	{"1\na\n"},
	{"3\nabc\n"},
	{"3\ncba\n"},
	{"4\nabac\n"},
	{"1\nz\n"},
	{"3\nzyx\n"},
	{"8\nabcdabcd\n"},
	{"3\nzxy\n"},
	{"5\naaaaa\n"},
	{"6\nbacdef\n"},
	{"5\ncbacd\n"},
	{"6\nazbycx\n"},
	{"26\nzyxwvutsrqponmlkjihgfedcba\n"},
	{"6\nabcdef\n"},
	{"6\nfedcba\n"},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
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
