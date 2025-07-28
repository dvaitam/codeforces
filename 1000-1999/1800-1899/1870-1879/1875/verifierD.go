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

func expected(arr []int) int64 {
	n := len(arr)
	freq := make(map[int]int)
	for _, v := range arr {
		if v <= n+1 {
			freq[v]++
		}
	}
	z := freq[0]
	if z == 0 {
		return 0
	}
	mex := 0
	for {
		if freq[mex] == 0 {
			break
		}
		mex++
	}
	if mex == 1 {
		return int64(z - 1)
	}
	costA := int64(mex) * int64(z-1)
	costB := int64(freq[1]-1)*int64(mex) + int64(z)
	if costB < costA {
		costA = costB
	}
	return costA
}

func run(bin string, input string) (string, error) {
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
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
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
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+n {
			fmt.Printf("test %d: expected %d numbers, got %d\n", idx, n+1, len(parts))
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[1+i])
			arr[i] = v
		}
		exp := expected(arr)
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d\n", n)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", idx, out)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
