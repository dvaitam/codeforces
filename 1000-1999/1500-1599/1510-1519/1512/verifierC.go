package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseC struct {
	A int
	B int
	S string
}

func parseTestcasesC(path string) ([]testCaseC, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cases []testCaseC
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("bad line: %s", line)
		}
		A, _ := strconv.Atoi(fields[0])
		B, _ := strconv.Atoi(fields[1])
		S := fields[2]
		cases = append(cases, testCaseC{A, B, S})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func solveC(A, B int, s string) string {
	n := len(s)
	arr := []byte(s)
	wenhao := 0
	for i := 0; i < n; i++ {
		switch arr[i] {
		case '?':
			wenhao++
		case '0':
			A--
		case '1':
			B--
		}
	}
	// Mirror known characters
	for i := 0; i < n; i++ {
		j := n - 1 - i
		if arr[i] != '?' && arr[j] == '?' {
			if arr[i] == '0' {
				arr[j] = '0'
				A--
			} else {
				arr[j] = '1'
				B--
			}
			wenhao--
		}
	}
	// Fill remaining placeholders
	for i := n / 2; i < n; i++ {
		if arr[i] != '?' {
			continue
		}
		j := n - 1 - i
		if i == j {
			if A%2 == 1 {
				arr[i] = '0'
				A--
			} else if B%2 == 1 {
				arr[i] = '1'
				B--
			}
			wenhao--
			continue
		}
		if A >= 2 {
			arr[i], arr[j] = '0', '0'
			A -= 2
			wenhao -= 2
		} else if B >= 2 {
			arr[i], arr[j] = '1', '1'
			B -= 2
			wenhao -= 2
		}
	}
	// Validate result
	ok := true
	for i := 0; i < n; i++ {
		if arr[i] != arr[n-1-i] {
			ok = false
			break
		}
	}
	if !ok || A != 0 || B != 0 || wenhao != 0 {
		return "-1"
	}
	return string(arr)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcasesC("testcasesC.txt")
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.A, tc.B))
		sb.WriteString(tc.S)
		sb.WriteByte('\n')
		expected := solveC(tc.A, tc.B, tc.S)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
