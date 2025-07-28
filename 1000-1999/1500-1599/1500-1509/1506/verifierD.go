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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(arr []int) int {
	freq := make(map[int]int)
	maxCnt := 0
	for _, x := range arr {
		freq[x]++
		if freq[x] > maxCnt {
			maxCnt = freq[x]
		}
	}
	n := len(arr)
	if 2*maxCnt <= n {
		if n%2 == 0 {
			return 0
		}
		return 1
	}
	return 2*maxCnt - n
}

func main() {
	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		arg = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := arg
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesD.txt: %v\n", err)
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
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts)-1 != n {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, n, len(parts)-1)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[i+1])
			arr[i] = v
		}
		expect := fmt.Sprintf("%d", expected(arr))
		input := fmt.Sprintf("1\n%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
