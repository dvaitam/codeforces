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

func solve(a, b int64) (string, int64) {
	if a < b {
		return "", 0
	}
	if a == b {
		return "infinity", 0
	}
	d := a - b
	var cnt int64
	for i := int64(1); i*i <= d; i++ {
		if d%i == 0 {
			x := i
			y := d / i
			if x > b {
				cnt++
			}
			if y != x && y > b {
				cnt++
			}
		}
	}
	return "", cnt
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
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
		var a, b int64
		fmt.Sscan(line, &a, &b)
		expStr, expCnt := solve(a, b)
		gotStr, err := run(bin, fmt.Sprintf("%d %d\n", a, b))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr = strings.TrimSpace(gotStr)
		if expStr == "infinity" {
			if strings.ToLower(gotStr) != "infinity" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected infinity got %s\n", idx, gotStr)
				os.Exit(1)
			}
			continue
		}
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: expected integer output got %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expCnt {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, expCnt, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
