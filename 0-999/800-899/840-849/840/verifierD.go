package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(arr []int, queries [][]int) []int {
	res := make([]int, len(queries))
	for qi, q := range queries {
		l, r, k := q[0], q[1], q[2]
		freq := make(map[int]int)
		for i := l - 1; i < r; i++ {
			freq[arr[i]]++
		}
		threshold := (r - l + 1) / k
		ans := -1
		for v, c := range freq {
			if c > threshold {
				if ans == -1 || v < ans {
					ans = v
				}
			}
		}
		res[qi] = ans
	}
	return res
}

func genCase(rng *rand.Rand) (string, []int, [][]int) {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(5) + 1
	}
	queries := make([][]int, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(arr[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		k := rng.Intn(4) + 1
		queries[i] = []int{l, r, k}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, k))
	}
	return sb.String(), arr, queries
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, arr, qs := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != len(qs) {
			fmt.Fprintf(os.Stderr, "case %d expected %d lines got %d\noutput:%s", i+1, len(qs), len(lines), out)
			os.Exit(1)
		}
		exp := expected(arr, qs)
		for j, line := range lines {
			v, err := strconv.Atoi(strings.TrimSpace(line))
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d line %d invalid int: %v\n", i+1, j+1, err)
				os.Exit(1)
			}
			if v != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d query %d failed: expected %d got %d\ninput:\n%s", i+1, j+1, exp[j], v, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
