package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const maxBit = 30

func dfsC(a []int, bit int, cnt01, cnt10 []int64) {
	if bit < 0 || len(a) < 2 {
		return
	}
	v0 := make([]int, 0, len(a))
	v1 := make([]int, 0, len(a))
	var zeros, ones int64
	for _, v := range a {
		if (v>>bit)&1 == 1 {
			cnt01[bit] += zeros
			ones++
			v1 = append(v1, v)
		} else {
			cnt10[bit] += ones
			zeros++
			v0 = append(v0, v)
		}
	}
	dfsC(v0, bit-1, cnt01, cnt10)
	dfsC(v1, bit-1, cnt01, cnt10)
}

func solveC(a []int) (int64, int64) {
	cnt01 := make([]int64, maxBit+1)
	cnt10 := make([]int64, maxBit+1)
	dfsC(a, maxBit, cnt01, cnt10)
	var x, inv int64
	for b := 0; b <= maxBit; b++ {
		if cnt10[b] <= cnt01[b] {
			inv += cnt10[b]
		} else {
			inv += cnt01[b]
			x |= 1 << b
		}
	}
	return inv, x
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(3)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(20) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(1000)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		got, err := runBinary(binary, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", t, err)
			os.Exit(1)
		}
		inv, x := solveC(a)
		expected := fmt.Sprintf("%d %d", inv, x)
		if got != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\nexpected: %s\ngot: %s\n", t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
