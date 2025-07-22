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

const mod = 1 << 30
const maxN = 1000000

var divisors []int

func init() {
	divisors = make([]int, maxN+1)
	for i := 1; i <= maxN; i++ {
		for j := i; j <= maxN; j += i {
			divisors[j]++
		}
	}
}

func solve(a, b, c int) int {
	var sum int64
	for i := 1; i <= a; i++ {
		for j := 1; j <= b; j++ {
			ij := i * j
			for k := 1; k <= c; k++ {
				sum += int64(divisors[ij*k])
			}
		}
	}
	return int(sum % mod)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesB.txt: %v\n", err)
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
		if len(fields) != 3 {
			fmt.Printf("bad test %d: %s\n", idx, line)
			os.Exit(1)
		}
		a, _ := strconv.Atoi(fields[0])
		b, _ := strconv.Atoi(fields[1])
		c, _ := strconv.Atoi(fields[2])
		want := solve(a, b, c)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(fmt.Sprintf("%d %d %d\n", a, b, c))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", idx, err, out.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
