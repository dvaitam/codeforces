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

type pair struct {
	cnt uint64
	sum uint64
}

var memo = map[uint64]pair{0: {0, 0}}

func icbrt(n uint64) uint64 {
	lo, hi := uint64(0), uint64(1_000_001)
	for lo+1 < hi {
		mid := (lo + hi) / 2
		mid3 := mid * mid * mid
		if mid3 <= n {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo
}

func solve(m uint64) pair {
	if res, ok := memo[m]; ok {
		return res
	}
	t := icbrt(m)
	t3 := t * t * t
	a := m - t3
	var b uint64
	if t > 0 {
		b = t3 - 1
	}
	alt1 := solve(a)
	alt1.cnt++
	alt1.sum += t3
	alt2 := solve(b)
	var res pair
	if alt1.cnt > alt2.cnt || (alt1.cnt == alt2.cnt && alt1.sum >= alt2.sum) {
		res = alt1
	} else {
		res = alt2
	}
	memo[m] = res
	return res
}

func runBinary(bin, input string) (string, string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), errb.String(), err
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
		input := line + "\n"
		m, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse number on line %d: %v\n", idx, err)
			os.Exit(1)
		}
		res := solve(m)
		exp := fmt.Sprintf("%d %d", res.cnt, res.sum)

		got, errStr2, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errStr2)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
