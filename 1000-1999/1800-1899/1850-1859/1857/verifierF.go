package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func isqrt(x int64) int64 {
	if x < 0 {
		return 0
	}
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func solveQuery(freq map[int64]int64, x, y int64) int64 {
	delta := x*x - 4*y
	if delta < 0 {
		return 0
	}
	s := isqrt(delta)
	if s*s != delta {
		return 0
	}
	if (x+s)%2 != 0 || (x-s)%2 != 0 {
		return 0
	}
	r1 := (x + s) / 2
	r2 := (x - s) / 2
	if r1 == r2 {
		c := freq[r1]
		if c >= 2 {
			return c * (c - 1) / 2
		}
		return 0
	}
	return freq[r1] * freq[r2]
}

func expected(a []int64, queries [][2]int64) string {
	freq := make(map[int64]int64)
	for _, v := range a {
		freq[v]++
	}
	ans := make([]string, len(queries))
	for i, q := range queries {
		ans[i] = fmt.Sprintf("%d", solveQuery(freq, q[0], q[1]))
	}
	return strings.Join(ans, " ")
}

func runCase(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("failed to read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Printf("case %d missing n\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
		if !scan.Scan() {
			fmt.Printf("case %d missing array\n", caseNum)
			os.Exit(1)
		}
		fields := strings.Fields(scan.Text())
		if len(fields) != n {
			fmt.Printf("case %d wrong array length\n", caseNum)
			os.Exit(1)
		}
		a := make([]int64, n)
		for i, f := range fields {
			a[i], _ = strconv.ParseInt(f, 10, 64)
		}
		if !scan.Scan() {
			fmt.Printf("case %d missing q\n", caseNum)
			os.Exit(1)
		}
		q, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
		queries := make([][2]int64, q)
		for i := 0; i < q; i++ {
			if !scan.Scan() {
				fmt.Printf("case %d missing query %d\n", caseNum, i)
				os.Exit(1)
			}
			parts := strings.Fields(scan.Text())
			if len(parts) != 2 {
				fmt.Printf("case %d malformed query\n", caseNum)
				os.Exit(1)
			}
			x, _ := strconv.ParseInt(parts[0], 10, 64)
			y, _ := strconv.ParseInt(parts[1], 10, 64)
			queries[i] = [2]int64{x, y}
		}
		inputLines := []string{"1", fmt.Sprintf("%d", n), strings.Join(fields, " "), fmt.Sprintf("%d", q)}
		for _, qv := range queries {
			inputLines = append(inputLines, fmt.Sprintf("%d %d", qv[0], qv[1]))
		}
		input := strings.Join(inputLines, "\n") + "\n"
		want := expected(a, queries)
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
