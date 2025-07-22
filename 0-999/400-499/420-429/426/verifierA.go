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

func expected(n, s int, arr []int) string {
	sum := 0
	maxv := 0
	for _, v := range arr {
		sum += v
		if v > maxv {
			maxv = v
		}
	}
	if sum-maxv <= s {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, s int
		fmt.Sscan(line, &n, &s)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d malformed: missing array line\n", idx)
			os.Exit(1)
		}
		arrLine := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(arrLine)
		if len(fields) != n {
			fmt.Fprintf(os.Stderr, "case %d malformed: expected %d values got %d\n", idx, n, len(fields))
			os.Exit(1)
		}
		arr := make([]int, n)
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d malformed: %v\n", idx, err)
				os.Exit(1)
			}
			arr[i] = v
		}
		input := fmt.Sprintf("%d %d\n%s\n", n, s, strings.Join(fields, " "))
		want := strings.ToUpper(expected(n, s, arr))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx, want, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
