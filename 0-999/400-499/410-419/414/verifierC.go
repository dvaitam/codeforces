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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func countInv(a []int) int64 {
	tmp := make([]int, len(a))
	var inv func(l, r int) int64
	inv = func(l, r int) int64 {
		if r-l <= 1 {
			return 0
		}
		m := (l + r) / 2
		res := inv(l, m) + inv(m, r)
		i, j, k := l, m, l
		for i < m && j < r {
			if a[i] <= a[j] {
				tmp[k] = a[i]
				i++
			} else {
				tmp[k] = a[j]
				res += int64(m - i)
				j++
			}
			k++
		}
		for i < m {
			tmp[k] = a[i]
			i++
			k++
		}
		for j < r {
			tmp[k] = a[j]
			j++
			k++
		}
		for i := l; i < r; i++ {
			a[i] = tmp[i]
		}
		return res
	}
	b := make([]int, len(a))
	copy(b, a)
	res := inv(0, len(b))
	return res
}

func process(a []int, q int) []int {
	seg := 1 << q
	n := len(a)
	res := make([]int, 0, n)
	for i := 0; i < n; i += seg {
		end := i + seg
		if end > n {
			end = n
		}
		for j := end - 1; j >= i; j-- {
			res = append(res, a[j])
		}
	}
	return res
}

func solveC(n int, arr []int, queries []int) []int64 {
	res := make([]int64, len(queries))
	cur := make([]int, len(arr))
	copy(cur, arr)
	for i, q := range queries {
		cur = process(cur, q)
		res[i] = countInv(cur)
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(5) + 1
	N := 1 << n
	arr := make([]int, N)
	for i := range arr {
		arr[i] = rng.Intn(20)
	}
	m := rng.Intn(6) + 1
	queries := make([]int, m)
	for i := range queries {
		queries[i] = rng.Intn(n + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i, q := range queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(q))
	}
	sb.WriteByte('\n')
	expected := solveC(n, arr, queries)
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		outTokens := strings.Fields(out)
		if len(outTokens) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\ninput:\n%soutput:\n%s", i+1, len(exp), len(outTokens), input, out)
			os.Exit(1)
		}
		for j, tok := range outTokens {
			got, err := strconv.ParseInt(tok, 10, 64)
			if err != nil || got != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at output %d: expected %d got %s\ninput:\n%s", i+1, j+1, exp[j], tok, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
