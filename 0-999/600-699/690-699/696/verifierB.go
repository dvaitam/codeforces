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

type treeB struct {
	children [][]int
	size     []int
}

func solveB(n int, parents []int) string {
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		children[p] = append(children[p], i)
	}
	order := []int{}
	stack := []int{1}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, c := range children[v] {
			stack = append(stack, c)
		}
	}
	size := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		s := 1
		for _, c := range children[v] {
			s += size[c]
		}
		size[v] = s
	}
	ans := make([]float64, n+1)
	ans[1] = 1.0
	for _, v := range order {
		for _, c := range children[v] {
			ans[c] = ans[v] + 1.0 + float64(size[v]-1-size[c])/2.0
		}
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%.10f", ans[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		parents := make([]int, n-1)
		var inSB strings.Builder
		inSB.WriteString(fmt.Sprintf("%d\n", n))
		if n > 1 {
			for i := 0; i < n-1; i++ {
				scan.Scan()
				parents[i], _ = strconv.Atoi(scan.Text())
			}
			for i, p := range parents {
				if i > 0 {
					inSB.WriteByte(' ')
				}
				inSB.WriteString(strconv.Itoa(p))
			}
			inSB.WriteByte('\n')
		} else {
			inSB.WriteByte('\n')
		}
		input := inSB.String()
		expected := solveB(n, parents)
		if err := runCase(exe, input, expected); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
