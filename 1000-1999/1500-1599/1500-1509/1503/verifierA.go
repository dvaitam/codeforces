package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	n int
	s string
}

func parseTests(path string) ([]testCaseA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([]testCaseA, t)
	for i := 0; i < t; i++ {
		var n int
		var s string
		if _, err := fmt.Fscan(in, &n, &s); err != nil {
			return nil, err
		}
		cases[i] = testCaseA{n: n, s: s}
	}
	return cases, nil
}

func solve(tc testCaseA) (bool, string, string) {
	n := tc.n
	s := tc.s
	ones := 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			ones++
		}
	}
	a := make([]byte, n)
	b := make([]byte, n)
	half := ones / 2
	sw := true
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			if half > 0 {
				a[i], b[i] = '(', '('
			} else {
				a[i], b[i] = ')', ')'
			}
			half--
		} else {
			if sw {
				a[i], b[i] = '(', ')'
			} else {
				a[i], b[i] = ')', '('
			}
			sw = !sw
		}
	}
	balA, balB := 0, 0
	for i := 0; i < n; i++ {
		if a[i] == '(' {
			balA++
		} else {
			balA--
		}
		if b[i] == '(' {
			balB++
		} else {
			balB--
		}
		if balA < 0 || balB < 0 {
			return false, "", ""
		}
	}
	if balA != 0 || balB != 0 {
		return false, "", ""
	}
	return true, string(a), string(b)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	cases, err := parseTests("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for idx, tc := range cases {
		expectOk, expA, expB := solve(tc)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n%s\n", tc.n, tc.s))
		output, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", idx+1, err, output)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(output), "\n")
		if !expectOk {
			if strings.ToUpper(strings.TrimSpace(outLines[0])) != "NO" {
				fmt.Printf("case %d failed: expected NO got %s\n", idx+1, outLines[0])
				os.Exit(1)
			}
			continue
		}
		if len(outLines) < 3 {
			fmt.Printf("case %d: insufficient output\n", idx+1)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(outLines[0])) != "YES" {
			fmt.Printf("case %d failed: expected YES got %s\n", idx+1, outLines[0])
			os.Exit(1)
		}
		gotA := strings.TrimSpace(outLines[1])
		gotB := strings.TrimSpace(outLines[2])
		if gotA != expA || gotB != expB {
			fmt.Printf("case %d failed: expected %s %s got %s %s\n", idx+1, expA, expB, gotA, gotB)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
