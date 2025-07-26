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

func parseLine(line string) (int, []int, error) {
	parts := strings.Fields(line)
	if len(parts) < 1 {
		return 0, nil, fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, nil, err
	}
	if len(parts) != 1+n {
		return 0, nil, fmt.Errorf("wrong number of integers")
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i], err = strconv.Atoi(parts[1+i])
		if err != nil {
			return 0, nil, err
		}
	}
	return n, arr, nil
}

func expectedCost(arr []int) int64 {
	var cost int64
	negatives := 0
	zeros := 0
	for _, x := range arr {
		if x > 0 {
			cost += int64(x - 1)
		} else if x < 0 {
			cost += int64(-1 - x)
			negatives++
		} else {
			zeros++
		}
	}
	if negatives%2 == 0 {
		cost += int64(zeros)
	} else {
		if zeros > 0 {
			cost += int64(zeros)
		} else {
			cost += 2
		}
	}
	return cost
}

func runCase(bin string, n int, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
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
	expected := expectedCost(arr)
	got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	f, err := os.Open(filepath.Join(dir, "testcasesB.txt"))
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
		n, arr, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if err := runCase(bin, n, arr); err != nil {
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
