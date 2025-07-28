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

func hasQuad(arr []int) bool {
	n := len(arr)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			sum := arr[i] + arr[j]
			for k := 0; k < n; k++ {
				if k == i || k == j {
					continue
				}
				for l := k + 1; l < n; l++ {
					if l == i || l == j {
						continue
					}
					if arr[k]+arr[l] == sum {
						return true
					}
				}
			}
		}
	}
	return false
}

func validateOutput(out string, arr []int, expect bool) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	ans := strings.ToLower(fields[0])
	if ans != "yes" && ans != "no" {
		return fmt.Errorf("first word must be YES or NO")
	}
	has := hasQuad(arr)
	if has != (ans == "yes") {
		return fmt.Errorf("wrong YES/NO")
	}
	if ans == "yes" {
		if len(fields) != 5 {
			return fmt.Errorf("expected 4 indices after YES")
		}
		idx := make([]int, 4)
		for i := 0; i < 4; i++ {
			v, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return fmt.Errorf("invalid index")
			}
			if v < 1 || v > len(arr) {
				return fmt.Errorf("index out of range")
			}
			idx[i] = v - 1
		}
		if idx[0] == idx[1] || idx[0] == idx[2] || idx[0] == idx[3] || idx[1] == idx[2] || idx[1] == idx[3] || idx[2] == idx[3] {
			return fmt.Errorf("indices must be distinct")
		}
		s1 := arr[idx[0]] + arr[idx[1]]
		s2 := arr[idx[2]] + arr[idx[3]]
		if s1 != s2 {
			return fmt.Errorf("indices do not satisfy equation")
		}
	}
	return nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
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
		if len(parts) < 1 {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil || len(parts) != n+1 {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(parts[i+1])
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(parts[i+1])
		}
		input.WriteByte('\n')
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if err := validateOutput(got, arr, hasQuad(arr)); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
