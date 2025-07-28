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

func expected(n int, a, b []int) string {
	pos := make(map[int]int, n)
	for i := 0; i < n; i++ {
		pos[b[i]] = i
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		p, ok := pos[a[i]]
		if !ok {
			return "NO"
		}
		perm[i] = p
	}
	visited := make([]bool, n)
	cycles := 0
	for i := 0; i < n; i++ {
		if !visited[i] {
			cycles++
			j := i
			for !visited[j] {
				visited[j] = true
				j = perm[j]
			}
		}
	}
	parity := (n - cycles) % 2
	if parity == 0 {
		return "YES"
	}
	return "NO"
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("failed to open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid testcases file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Printf("missing n for case %d\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			a[i], _ = strconv.Atoi(scan.Text())
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			b[i], _ = strconv.Atoi(scan.Text())
		}
		expect := expected(n, a, b)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", b[i]))
		}
		sb.WriteByte('\n')
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %q got %q\n", caseNum, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
