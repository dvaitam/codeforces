package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

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

func parseLine(line string) (int, []int, int, []int, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return 0, nil, 0, nil, fmt.Errorf("invalid line")
	}
	idx := 0
	n, err := strconv.Atoi(parts[idx])
	if err != nil {
		return 0, nil, 0, nil, err
	}
	idx++
	if len(parts) < 1+n+1 {
		return 0, nil, 0, nil, fmt.Errorf("not enough numbers for A")
	}
	A := make([]int, n)
	for i := 0; i < n; i++ {
		A[i], err = strconv.Atoi(parts[idx])
		if err != nil {
			return 0, nil, 0, nil, err
		}
		idx++
	}
	m, err := strconv.Atoi(parts[idx])
	if err != nil {
		return 0, nil, 0, nil, err
	}
	idx++
	if len(parts) != idx+m {
		return 0, nil, 0, nil, fmt.Errorf("wrong number of B elements")
	}
	B := make([]int, m)
	for i := 0; i < m; i++ {
		B[i], err = strconv.Atoi(parts[idx+i])
		if err != nil {
			return 0, nil, 0, nil, err
		}
	}
	return n, A, m, B, nil
}

func runCase(bin string, n int, A []int, m int, B []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range A {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i, v := range B {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return fmt.Errorf("expected 2 numbers got %d", len(fields))
	}
	x, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid integer %q", fields[0])
	}
	y, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("invalid integer %q", fields[1])
	}
	ma := make(map[int]struct{}, len(A))
	for _, v := range A {
		ma[v] = struct{}{}
	}
	mb := make(map[int]struct{}, len(B))
	for _, v := range B {
		mb[v] = struct{}{}
	}
	if _, ok := ma[x]; !ok {
		return fmt.Errorf("%d is not in A", x)
	}
	if _, ok := mb[y]; !ok {
		return fmt.Errorf("%d is not in B", y)
	}
	if _, ok := ma[x+y]; ok {
		return fmt.Errorf("sum %d exists in A", x+y)
	}
	if _, ok := mb[x+y]; ok {
		return fmt.Errorf("sum %d exists in B", x+y)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	f, err := os.Open(filepath.Join(dir, "testcasesA.txt"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, A, m, B, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if err := runCase(bin, n, A, m, B); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
