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

type laptop struct{ a, b, c, d int }

func expected(n int, arr []laptop) int {
	outdated := make([]bool, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if arr[i].a < arr[j].a && arr[i].b < arr[j].b && arr[i].c < arr[j].c {
				outdated[i] = true
			}
		}
	}
	bestCost := 1005
	idx := 0
	for i := 0; i < n; i++ {
		if !outdated[i] && arr[i].d < bestCost {
			bestCost = arr[i].d
			idx = i
		}
	}
	return idx + 1
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
	f, err := os.Open("testcasesB.txt")
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
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+4*n {
			fmt.Fprintf(os.Stderr, "case %d: bad line\n", idx)
			os.Exit(1)
		}
		arr := make([]laptop, n)
		p := 1
		for i := 0; i < n; i++ {
			a, _ := strconv.Atoi(parts[p])
			b, _ := strconv.Atoi(parts[p+1])
			c, _ := strconv.Atoi(parts[p+2])
			d, _ := strconv.Atoi(parts[p+3])
			arr[i] = laptop{a, b, c, d}
			p += 4
		}
		expect := expected(n, arr)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", arr[i].a, arr[i].b, arr[i].c, arr[i].d))
		}
		gotStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got, _ := strconv.Atoi(strings.Fields(gotStr)[0])
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
