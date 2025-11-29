package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `100
<b><a><b><c></c></b></a></b>
<c><a><c></c></a><a></a></c><c><c></c></c>
<a><c><a></a></c></a><a><b></b></a>
<c><b></b></c>
<c></c>
<c></c>
<c></c>
<a><c><b></b></c><c></c></a>
<c><b></b></c>
<a><b><c><b></b></c></b></a>
<b><a><c></c></a><a></a></b><a><a></a></a>
<a><a><c><a></a></c></a></a>
<c><c><b><c></c></b></c></c>
<c><b><b><b></b></b></b></c>
<b></b>
<c><c><b><a></a></b></c><a><b></b></a><a></a></c><c><a><a></a></a></c>
<a></a>
<c><c></c></c><a></a>
<a></a>
<c><a></a></c>
<c></c>
<c><a><b><a></a></b></a><a><a></a></a><a></a></c><b><b><c></c></b></b>
<a><c></c></a><c></c>
<c><a><a><b></b></a></a></c>
<b><c><b></b></c></b><a><b></b></a>
<a><b></b></a><b></b>
<c><b></b></c><c></c>
<c><b></b></c><b></b>
<a><b><c><c></c></c></b><b><a></a></b><c></c></a>
<a></a>
<b><a><b></b></a></b>
<c></c>
<b></b>
<b><a><a><a></a></a><c></c></a><a><c></c></a></b><a><a><b></b></a><a></a></a>
<b></b>
<a><a><b></b></a><a></a></a>
<c><a><c></c></a></c><c><b></b></c><c></c>
<c><b><a><a></a></a><b></b></b><c><a></a></c><c></c></c><c><a><b></b></a></c>
<a><c><c></c></c><b></b></a><c><b></b></c>
<a><b><a></a></b><b></b></a>
<a><c></c></a>
<a></a>
<c><a></a></c><b></b>
<b><b><c></c></b><a></a></b>
<c><a><b></b></a><c></c></c>
<a><b></b></a><a></a>
<c><a></a></c><b></b>
<a><c></c></a><b></b>
<c><b><c><b></b></c></b></c>
<a><b></b></a>
<b><a><c></c></a><a></a></b>
<a><a><c><a></a></c><b></b></a></a><a><a><c></c></a><b></b></a>
<c></c>
<c></c>
<b><c></c></b><b></b>
<c><c><a><a></a></a><a></a></c></c><b><b><b></b></b><b></b></b>
<a><b></b></a><a></a>
<b><c></c></b>
<b><b></b></b><b></b>
<a><a><a><c></c></a></a><c><c></c></c><b></b></a><c><b><c></c></b><a></a></c>
<a><b><b></b></b></a>
<c></c>
<a><b></b></a>
<b><c><a><a></a></a><b></b></c></b>
<a><a></a></a>
<c><b><a><a></a></a></b></c>
<a><a><b><a></a></b><b></b></a><b><b></b></b><a></a></a>
<b><a><a></a></a><b></b></b>
<a><a></a></a>
<c></c>
<c><a><c></c></a><b></b></c><b><a></a></b><b></b>
<b><c><c></c></c></b><c><c></c></c><a></a>
<b><c></c></b>
<b><b><c></c></b><a></a></b>
<b><b></b></b>
<c><b></b></c><a></a>
<c><c><a></a></c><c></c></c><b><b></b></b><a></a>
<a><c><c></c></c></a><b><a></a></b><a></a>
<c><a><b></b></a></c><a><b></b></a>
<a><a><b><a></a></b></a></a><c><b><a></a></b></c><b><b></b></b><b></b>
<a></a>
<b><a><c><b></b></c></a></b><a><a><a></a></a><b></b></a>
<a><b></b></a><a></a>
<c><a></a></c>
<c><b><b><b></b></b></b><c><a></a></c></c><c><b><b></b></b></c><b><c></c></b><a></a>
<c><b></b></c>
<b></b>
<c></c>
<c></c>
<a><b><a><a></a></a></b></a>
<a><c></c></a><b></b>
<a><b><b><c></c></b><b></b></b><b><a></a></b><b></b></a><b><b><c></c></b><a></a></b><b><a></a></b><a></a>
<a><c></c></a><c></c>
<b><b><a></a></b></b>
<c><b><b><b></b></b><b></b></b><b><a></a></b><a></a></c>
<c><a><b><a></a></b><a></a></a><b><b></b></b></c><a><a><c></c></a></a><c><c></c></c>
<b><b><b><c></c></b><c></c></b><c><b></b></c></b>
<b><a><a><a></a></a></a><a><b></b></a></b>
<a><a><a><b></b></a></a><b><b></b></b><b></b></a>
<b><c><c></c></c></b><a><a></a></a>`

func parseTestcases() ([]string, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %v", err)
	}
	if len(lines)-1 < t {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(lines)-1)
	}
	cases := make([]string, 0, t)
	for i := 0; i < t; i++ {
		cases = append(cases, strings.TrimSpace(lines[i+1]))
	}
	return cases, nil
}

// solve implements the logic from 125B.go for a single XML string.
func solve(s string) string {
	indent := 0
	var out strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '<' {
			j := i + 1
			for j < len(s) && s[j] != '>' {
				j++
			}
			if j >= len(s) {
				break
			}
			token := s[i : j+1]
			if len(token) >= 2 && token[1] == '/' {
				indent--
				if indent < 0 {
					indent = 0
				}
			}
			out.WriteString(strings.Repeat(" ", indent*2))
			out.WriteString(token)
			out.WriteByte('\n')
			if len(token) >= 2 && token[1] != '/' {
				indent++
			}
			i = j + 1
		} else {
			i++
		}
	}
	return strings.TrimSpace(out.String())
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, s := range cases {
		expect := solve(s)
		got, err := runCandidate(bin, s+"\n")
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, s, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
