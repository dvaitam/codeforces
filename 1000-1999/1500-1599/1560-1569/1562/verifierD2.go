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

func solveCase(n, q int, s string, queries [][2]int) string {
	p := make([]int, n+1)
	v1 := make([][]int, n+1)
	v2 := make([][]int, n+1)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			if s[i] == '+' {
				p[i+1] = p[i] + 1
			} else {
				p[i+1] = p[i] - 1
			}
		} else {
			if s[i] == '+' {
				p[i+1] = p[i] - 1
			} else {
				p[i+1] = p[i] + 1
			}
		}
		if p[i+1] >= 0 {
			v1[p[i+1]] = append(v1[p[i+1]], i+1)
		} else {
			v2[-p[i+1]] = append(v2[-p[i+1]], i+1)
		}
	}
	var out strings.Builder
	for _, qr := range queries {
		l, r := qr[0], qr[1]
		c := p[r] - p[l-1]
		if c == 0 {
			out.WriteString("0\n")
			continue
		}
		if c%2 == 0 {
			r0 := r
			r = r0 - 1
			var f int
			if p[r] > p[l-1] {
				f = p[l-1] + (p[r]-p[l-1]+1)/2
			} else {
				f = p[l-1] - (p[l-1]-p[r]+1)/2
			}
			var pos int
			if f >= 0 {
				arr := v1[f]
				idx := sort.Search(len(arr), func(i int) bool { return arr[i] > l-1 })
				pos = arr[idx]
			} else {
				arr := v2[-f]
				idx := sort.Search(len(arr), func(i int) bool { return arr[i] > l-1 })
				pos = arr[idx]
			}
			out.WriteString("2\n")
			out.WriteString(fmt.Sprintf("%d %d\n", r0, pos))
		} else {
			var f int
			if p[r] > p[l-1] {
				f = p[l-1] + (p[r]-p[l-1]+1)/2
			} else {
				f = p[l-1] - (p[l-1]-p[r]+1)/2
			}
			var pos int
			if f >= 0 {
				arr := v1[f]
				idx := sort.Search(len(arr), func(i int) bool { return arr[i] > l-1 })
				pos = arr[idx]
			} else {
				arr := v2[-f]
				idx := sort.Search(len(arr), func(i int) bool { return arr[i] > l-1 })
				pos = arr[idx]
			}
			out.WriteString("1\n")
			out.WriteString(fmt.Sprintf("%d\n", pos))
		}
	}
	return out.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	q := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	str := make([]byte, n)
	for i := range str {
		if rng.Intn(2) == 0 {
			str[i] = '+'
		} else {
			str[i] = '-'
		}
	}
	s := string(str)
	sb.WriteString(s)
	sb.WriteByte('\n')
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	input := sb.String()
	expected := solveCase(n, q, s, queries)
	return input, expected
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
