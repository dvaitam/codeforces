package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveCase(line string) string {
	fields := strings.Fields(line)
	idx := 0
	n, _ := strconv.Atoi(fields[idx])
	idx++
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i], _ = strconv.Atoi(fields[idx])
		idx++
	}
	m, _ := strconv.Atoi(fields[idx])
	idx++
	ks := make([]int, m)
	for i := 0; i < m; i++ {
		ks[i], _ = strconv.Atoi(fields[idx])
		idx++
	}
	var sb strings.Builder
	for i, k := range ks {
		if i > 0 {
			sb.WriteByte(' ')
		}
		cnt := n - k + 1
		sum := 0.0
		for x := 0; x < cnt; x++ {
			mn := arr[x]
			for j := x + 1; j < x+k; j++ {
				if arr[j] < mn {
					mn = arr[j]
				}
			}
			sum += float64(mn)
		}
		sb.WriteString(fmt.Sprintf("%.10f", sum/float64(cnt)))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
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
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
