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

func solveCase(q []int) (string, string) {
	n := len(q)
	minAns := make([]int, n)
	maxAns := make([]int, n)
	usedMin := make([]bool, n+1)
	usedMax := make([]bool, n+1)

	nextMin := 1
	for i := 0; i < n; i++ {
		if i == 0 || q[i] != q[i-1] {
			minAns[i] = q[i]
			usedMin[q[i]] = true
		} else {
			for nextMin <= n && usedMin[nextMin] {
				nextMin++
			}
			minAns[i] = nextMin
			usedMin[nextMin] = true
		}
	}

	cur := 0
	stack := []int{}
	for i := 0; i < n; i++ {
		if i == 0 || q[i] != q[i-1] {
			maxAns[i] = q[i]
			for x := cur + 1; x < q[i]; x++ {
				if !usedMax[x] {
					stack = append(stack, x)
					usedMax[x] = true
				}
			}
			usedMax[q[i]] = true
			cur = q[i]
		} else {
			idx := len(stack) - 1
			maxAns[i] = stack[idx]
			stack = stack[:idx]
		}
	}

	var sb1 strings.Builder
	var sb2 strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb1.WriteByte(' ')
			sb2.WriteByte(' ')
		}
		sb1.WriteString(strconv.Itoa(minAns[i]))
		sb2.WriteString(strconv.Itoa(maxAns[i]))
	}
	return sb1.String(), sb2.String()
}

func main() {
	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		arg = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := arg
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesE.txt: %v\n", err)
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
		q := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[i+1])
			q[i] = v
		}
		e1, e2 := solveCase(q)
		expect := e1 + "\n" + e2
		input := fmt.Sprintf("1\n%d\n", n)
		for i, v := range q {
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
			fmt.Printf("case %d failed:\nexpected:\n%s\n\ngot:\n%s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
