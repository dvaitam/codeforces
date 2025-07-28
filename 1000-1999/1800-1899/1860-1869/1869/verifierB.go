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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func expected(n, k, a, b int, xs, ys []int64) int64 {
	direct := abs(xs[a]-xs[b]) + abs(ys[a]-ys[b])
	const inf int64 = 1<<62 - 1
	dA, dB := inf, inf
	for i := 1; i <= k; i++ {
		da := abs(xs[a]-xs[i]) + abs(ys[a]-ys[i])
		if da < dA {
			dA = da
		}
		db := abs(xs[b]-xs[i]) + abs(ys[b]-ys[i])
		if db < dB {
			dB = db
		}
	}
	if dA+dB < direct {
		return dA + dB
	}
	return direct
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesB.txt:", err)
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
		fields := strings.Fields(line)
		if len(fields) < 4 {
			fmt.Fprintf(os.Stderr, "case %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		a, _ := strconv.Atoi(fields[2])
		b, _ := strconv.Atoi(fields[3])
		expectedFields := 4 + 2*n
		if len(fields) != expectedFields {
			fmt.Fprintf(os.Stderr, "case %d: wrong number of coordinates\n", idx)
			os.Exit(1)
		}
		xs := make([]int64, n+1)
		ys := make([]int64, n+1)
		p := 4
		for i := 1; i <= n; i++ {
			x64, _ := strconv.ParseInt(fields[p], 10, 64)
			y64, _ := strconv.ParseInt(fields[p+1], 10, 64)
			xs[i] = x64
			ys[i] = y64
			p += 2
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, k, a, b))
		for i := 1; i <= n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
		}
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output %q\n", idx, out)
			os.Exit(1)
		}
		exp := expected(n, k, a, b, xs, ys)
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, exp, val)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
