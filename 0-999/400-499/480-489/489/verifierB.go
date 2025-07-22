package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
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

func expected(boys, girls []int) int {
	sort.Ints(boys)
	sort.Ints(girls)
	i, j, ans := 0, 0, 0
	for i < len(boys) && j < len(girls) {
		if abs(boys[i]-girls[j]) <= 1 {
			ans++
			i++
			j++
		} else if boys[i] < girls[j] {
			i++
		} else {
			j++
		}
	}
	return ans
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func parseInts(fields []string) ([]int, error) {
	vals := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	return vals, nil
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
		fields := strings.Fields(line)
		values, err := parseInts(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test line %d: %v\n", idx, err)
			os.Exit(1)
		}
		if len(values) < 3 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		n := values[0]
		if len(values) < 1+n+1 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		boys := values[1 : 1+n]
		m := values[1+n]
		if len(values) != 1+n+1+m {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		girls := values[1+n+1:]
		in := fmt.Sprintf("%d\n%s\n%d\n%s\n", n, strings.Join(fields[1:1+n], " "), m, strings.Join(fields[1+n+1:], " "))
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		tokens := strings.Fields(out)
		if len(tokens) != 1 {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx)
			os.Exit(1)
		}
		got, err := strconv.Atoi(tokens[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad integer\n", idx)
			os.Exit(1)
		}
		want := expected(append([]int(nil), boys...), append([]int(nil), girls...))
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
