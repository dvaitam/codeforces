package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	input    string
	expected string
}

func brute(n, m, k int, arr []int, queries [][2]int) []int64 {
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] ^ arr[i-1]
	}
	res := make([]int64, m)
	for idx, q := range queries {
		l := q[0]
		r := q[1]
		var count int64
		for i := l; i <= r; i++ {
			for j := i; j <= r; j++ {
				if pref[j]^pref[i-1] == k {
					count++
				}
			}
		}
		res[idx] = count
	}
	return res
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(6) + 1
		m := rng.Intn(5) + 1
		k := rng.Intn(8)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(8)
		}
		queries := make([][2]int, m)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			queries[i] = [2]int{l, r}
			fmt.Fprintf(&sb, "%d %d\n", l, r)
		}
		ansSlice := brute(n, m, k, arr, queries)
		var out strings.Builder
		for i, v := range ansSlice {
			if i > 0 {
				out.WriteByte('\n')
			}
			out.WriteString(fmt.Sprintf("%d", v))
		}
		tests = append(tests, test{sb.String(), out.String()})
	}
	return tests
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
