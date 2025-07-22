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

type testCaseA struct {
	n   int
	d   int64
	xs  []int64
	ans int64
}

func countTriples(xs []int64, d int64) int64 {
	n := len(xs)
	var res int64
	r := 0
	for l := 0; l < n; l++ {
		if r < l {
			r = l
		}
		for r+1 < n && xs[r+1]-xs[l] <= d {
			r++
		}
		cnt := r - l
		if cnt >= 2 {
			res += int64(cnt*(cnt-1)) / 2
		}
	}
	return res
}

func genCaseA() testCaseA {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10) + 3
	d := rand.Int63n(20)
	xs := make([]int64, n)
	cur := rand.Int63n(10)
	for i := 0; i < n; i++ {
		xs[i] = cur
		cur += rand.Int63n(5) + 1
	}
	return testCaseA{n, d, xs, countTriples(xs, d)}
}

func buildInputA(cases []testCaseA) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d %d\n", c.n, c.d)
		for i, v := range c.xs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := make([]testCaseA, 100)
	for i := range cases {
		cases[i] = genCaseA()
	}
	input := buildInputA(cases)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outputs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outputs) != len(cases) {
		fmt.Printf("expected %d lines got %d\n", len(cases), len(outputs))
		os.Exit(1)
	}
	for i, s := range outputs {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v != cases[i].ans {
			fmt.Printf("mismatch on case %d: expected %d got %s\n", i+1, cases[i].ans, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
