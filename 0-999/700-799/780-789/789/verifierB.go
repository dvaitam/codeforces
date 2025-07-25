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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func expected(b1, q, l int64, bad map[int64]bool) string {
	if abs(b1) > l {
		return "0"
	}
	if q == 0 {
		if !bad[b1] {
			if bad[0] {
				return "1"
			}
			return "inf"
		}
		if bad[0] {
			return "0"
		}
		return "inf"
	}
	if q == 1 {
		if bad[b1] {
			return "0"
		}
		return "inf"
	}
	if q == -1 {
		if bad[b1] && bad[-b1] {
			return "0"
		}
		return "inf"
	}
	count := 0
	cur := b1
	for abs(cur) <= l {
		if !bad[cur] {
			count++
		}
		cur *= q
	}
	return fmt.Sprintf("%d", count)
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
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "empty testcases")
		os.Exit(1)
	}
	cases, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	for idx := 1; idx <= cases; idx++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing first line\n", idx)
			os.Exit(1)
		}
		header := strings.Fields(scanner.Text())
		if len(header) != 4 {
			fmt.Fprintf(os.Stderr, "case %d malformed header\n", idx)
			os.Exit(1)
		}
		b1, _ := strconv.ParseInt(header[0], 10, 64)
		q, _ := strconv.ParseInt(header[1], 10, 64)
		l, _ := strconv.ParseInt(header[2], 10, 64)
		m, _ := strconv.Atoi(header[3])
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing bad numbers\n", idx)
			os.Exit(1)
		}
		badFields := strings.Fields(scanner.Text())
		if len(badFields) != m {
			fmt.Fprintf(os.Stderr, "case %d expected %d bad numbers got %d\n", idx, m, len(badFields))
			os.Exit(1)
		}
		bad := make(map[int64]bool)
		for _, v := range badFields {
			x, _ := strconv.ParseInt(v, 10, 64)
			bad[x] = true
		}
		input := fmt.Sprintf("%d %d %d %d\n%s\n", b1, q, l, m, strings.Join(badFields, " "))
		want := expected(b1, q, l, bad)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", cases)
}
