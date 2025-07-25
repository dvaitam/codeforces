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

func expected(a, b string) int64 {
	n := len(a)
	m := len(b)
	pref := make([]int, m+1)
	for i := 0; i < m; i++ {
		pref[i+1] = pref[i]
		if b[i] == '1' {
			pref[i+1]++
		}
	}
	segLen := m - n + 1
	var ans int64
	for i := 0; i < n; i++ {
		l := i
		r := i + segLen - 1
		ones := pref[r+1] - pref[l]
		if a[i] == '0' {
			ans += int64(ones)
		} else {
			ans += int64(segLen - ones)
		}
	}
	return ans
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "case %d: invalid line\n", idx)
			os.Exit(1)
		}
		a := parts[0]
		b := parts[1]
		expect := expected(a, b)
		input := a + "\n" + b + "\n"
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
