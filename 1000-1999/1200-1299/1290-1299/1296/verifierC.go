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
	t := readInt()
	var sb strings.Builder
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := readInt()
		in.Scan()
		s := in.Text()
		type pair struct{ x, y int }
		pos := pair{0, 0}
		last := map[pair]int{pos: 0}
		bestLen := n + 1
		bestL, bestR := -1, -1
		for i := 1; i <= n; i++ {
			switch s[i-1] {
			case 'L':
				pos.x--
			case 'R':
				pos.x++
			case 'U':
				pos.y++
			case 'D':
				pos.y--
			}
			if prev, ok := last[pos]; ok {
				if i-prev < bestLen {
					bestLen = i - prev
					bestL = prev + 1
					bestR = i
				}
			}
			last[pos] = i
		}
		if caseIdx > 0 {
			sb.WriteByte('\n')
		}
		if bestLen == n+1 {
			sb.WriteString("-1")
		} else {
			sb.WriteString(fmt.Sprintf("%d %d", bestL, bestR))
		}
	}
	return sb.String()
}

type testCase struct{ input string }

var tests = []testCase{
	{"1\n2\nLR\n"},
	{"1\n2\nRL\n"},
	{"1\n4\nLRUD\n"},
	{"1\n2\nUD\n"},
	{"1\n1\nL\n"},
	{"1\n4\nLLLL\n"},
	{"1\n4\nLRLR\n"},
	{"1\n5\nLLRRR\n"},
	{"1\n6\nLRLRUD\n"},
	{"1\n3\nUDU\n"},
	{"1\n5\nURDLU\n"},
	{"1\n10\nRRRRLLLLUD\n"},
	{"1\n8\nUDUDUDUD\n"},
	{"1\n4\nLURD\n"},
	{"1\n6\nLRRLLR\n"},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
