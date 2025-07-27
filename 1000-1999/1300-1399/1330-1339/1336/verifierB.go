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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func calc(x, y, z int64) int64 {
	a := x - y
	b := y - z
	c := z - x
	return a*a + b*b + c*c
}

func best(a, b, c []int64) int64 {
	res := int64(1<<63 - 1)
	for _, y := range b {
		xi := sort.Search(len(a), func(i int) bool { return a[i] > y }) - 1
		if xi < 0 {
			continue
		}
		zi := sort.Search(len(c), func(i int) bool { return c[i] >= y })
		if zi == len(c) {
			continue
		}
		v := calc(a[xi], y, c[zi])
		if v < res {
			res = v
		}
	}
	return res
}

func solve(r, g, b []int64) int64 {
	sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
	sort.Slice(g, func(i, j int) bool { return g[i] < g[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	res := int64(1<<63 - 1)
	res = min(res, best(r, g, b))
	res = min(res, best(r, b, g))
	res = min(res, best(g, r, b))
	res = min(res, best(g, b, r))
	res = min(res, best(b, r, g))
	res = min(res, best(b, g, r))
	return res
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		if len(fields) < 3 {
			fmt.Fprintf(os.Stderr, "case %d invalid line\n", idx)
			os.Exit(1)
		}
		nr, _ := strconv.Atoi(fields[0])
		ng, _ := strconv.Atoi(fields[1])
		nb, _ := strconv.Atoi(fields[2])
		expectCount := nr + ng + nb
		if len(fields) != 3+expectCount {
			fmt.Fprintf(os.Stderr, "case %d invalid number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]int64, expectCount)
		for i := 0; i < expectCount; i++ {
			v, _ := strconv.Atoi(fields[3+i])
			arr[i] = int64(v)
		}
		r := arr[:nr]
		g := arr[nr : nr+ng]
		b := arr[nr+ng:]
		expectedAns := solve(r, g, b)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d %d\n", nr, ng, nb))
		for i := 0; i < nr; i++ {
			input.WriteString(fmt.Sprintf("%d ", r[i]))
		}
		input.WriteString("\n")
		for i := 0; i < ng; i++ {
			input.WriteString(fmt.Sprintf("%d ", g[i]))
		}
		input.WriteString("\n")
		for i := 0; i < nb; i++ {
			input.WriteString(fmt.Sprintf("%d ", b[i]))
		}
		input.WriteString("\n")
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		var ans int64
		if _, err := fmt.Sscan(got, &ans); err != nil || ans != expectedAns {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx, expectedAns, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
