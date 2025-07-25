package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func expected(a, b []int) []int {
	m := len(a)
	sort.Ints(a)
	type pair struct{ val, idx int }
	arr := make([]pair, m)
	for i, v := range b {
		arr[i] = pair{v, i}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
	res := make([]int, m)
	for i := 0; i < m; i++ {
		res[arr[i].idx] = a[m-1-i]
	}
	return res
}

func genCase(rng *rand.Rand) (string, []int) {
	m := rng.Intn(5) + 1
	a := make([]int, m)
	b := make([]int, m)
	maxB := 0
	for i := 0; i < m; i++ {
		a[i] = rng.Intn(100) + 1
	}
	for i := 0; i < m; i++ {
		b[i] = rng.Intn(100) + 1
		if b[i] > maxB {
			maxB = b[i]
		}
	}
	for i := 0; i < m; i++ {
		if a[i] < maxB {
			a[i] = maxB + rng.Intn(5)
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(m))
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(b[i]))
	}
	sb.WriteByte('\n')
	exp := expected(append([]int(nil), a...), append([]int(nil), b...))
	return sb.String(), exp
}

func parseOutput(out string, m int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != m {
		return nil, fmt.Errorf("expected %d numbers, got %d", m, len(fields))
	}
	res := make([]int, m)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "input:\n%s", input)
			os.Exit(1)
		}
		got, err := parseOutput(out, len(exp))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid output: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "output: %s\n", out)
			os.Exit(1)
		}
		for j := range exp {
			if got[j] != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\ninput:\n%s", i+1, exp, got, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
