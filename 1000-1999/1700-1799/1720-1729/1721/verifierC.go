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

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(a, b []int) (string, string) {
	n := len(a)
	mn := make([]int, n)
	j := 0
	for i := 0; i < n; i++ {
		for j < n && b[j] < a[i] {
			j++
		}
		mn[i] = b[j] - a[i]
	}
	mx := make([]int, n)
	pos := n - 1
	for i := n - 1; i >= 0; i-- {
		mx[i] = b[pos] - a[i]
		if i > 0 && b[i-1] < a[i] {
			pos = i - 1
		}
	}
	sb1 := strings.TrimSpace(strings.Join(intsToStrs(mn), " "))
	sb2 := strings.TrimSpace(strings.Join(intsToStrs(mx), " "))
	return sb1, sb2
}

func intsToStrs(a []int) []string {
	res := make([]string, len(a))
	for i, v := range a {
		res[i] = strconv.Itoa(v)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &a[i])
		}
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &b[i])
		}
		exp1, exp2 := solveCase(a, b)
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d\n", n)
		fmt.Fprintln(&input, strings.Join(intsToStrs(a), " "))
		fmt.Fprintln(&input, strings.Join(intsToStrs(b), " "))
		got, err := runExe(bin, input.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(outLines) < 2 {
			fmt.Printf("case %d: output too short\n", caseNum)
			os.Exit(1)
		}
		g1 := strings.TrimSpace(outLines[0])
		g2 := strings.TrimSpace(outLines[1])
		if g1 != exp1 || g2 != exp2 {
			fmt.Printf("case %d failed\nexpected:\n%s\n%s\n got:\n%s\n%s\n", caseNum, exp1, exp2, g1, g2)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
