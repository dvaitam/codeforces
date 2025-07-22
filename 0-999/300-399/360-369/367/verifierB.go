package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func multisetEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	ma := make(map[int]int)
	mb := make(map[int]int)
	for _, v := range a {
		ma[v]++
	}
	for _, v := range b {
		mb[v]++
	}
	if len(ma) != len(mb) {
		return false
	}
	for k, v := range ma {
		if mb[k] != v {
			return false
		}
	}
	return true
}

func solveB(n, m, p int, a, b []int) []int {
	var res []int
	for q := 0; q+(m-1)*p < n; q++ {
		sub := make([]int, m)
		for i := 0; i < m; i++ {
			sub[i] = a[q+i*p]
		}
		if multisetEqual(sub, b) {
			res = append(res, q+1)
		}
	}
	sort.Ints(res)
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	m := rng.Intn(min(10, n)) + 1
	p := rng.Intn(n) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(10) + 1
	}
	b := make([]int, m)
	for i := 0; i < m; i++ {
		b[i] = rng.Intn(10) + 1
	}
	input := fmt.Sprintf("%d %d %d\n", n, m, p)
	for i, v := range a {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	for i, v := range b {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	ans := solveB(n, m, p, a, b)
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d", len(ans)))
	if len(ans) > 0 {
		sb.WriteByte('\n')
		for i, v := range ans {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
	}
	expected := sb.String()
	return input, expected
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
