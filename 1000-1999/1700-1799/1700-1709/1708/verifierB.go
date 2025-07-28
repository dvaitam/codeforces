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

func solveB(n int, l, r int64) (bool, []int64) {
	ans := make([]int64, n)
	for i := 1; i <= n; i++ {
		ii := int64(i)
		x := ((l-1)/ii + 1) * ii
		if x > r {
			return false, nil
		}
		ans[i-1] = x
	}
	return true, ans
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
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
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			fmt.Fprintf(os.Stderr, "bad test case at line %d\n", idx+1)
			os.Exit(1)
		}
		n, err1 := strconv.Atoi(fields[0])
		l, err2 := strconv.ParseInt(fields[1], 10, 64)
		r, err3 := strconv.ParseInt(fields[2], 10, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Fprintf(os.Stderr, "bad numbers at line %d\n", idx+1)
			os.Exit(1)
		}
		ok, arr := solveB(n, l, r)
		var expected string
		if !ok {
			expected = "NO"
		} else {
			var sb strings.Builder
			sb.WriteString("YES\n")
			for i, v := range arr {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", v))
			}
			expected = strings.TrimSpace(sb.String())
		}
		var inputSB strings.Builder
		inputSB.WriteString("1\n")
		inputSB.WriteString(fmt.Sprintf("%d %d %d\n", n, l, r))
		got, err := runBinary(bin, inputSB.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
