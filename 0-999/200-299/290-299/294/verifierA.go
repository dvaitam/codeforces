package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type op struct{ x, y int }

func simulateA(n int, arr []int, ops []op) []int {
	a := make([]int, n)
	copy(a, arr)
	for _, o := range ops {
		x := o.x - 1
		y := o.y
		left := y - 1
		right := a[x] - y
		a[x] = 0
		if x-1 >= 0 {
			a[x-1] += left
		}
		if x+1 < n {
			a[x+1] += right
		}
	}
	return a
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(10) + 1
	}
	m := rng.Intn(5) + 1
	ops := make([]op, 0, m)
	tmp := make([]int, n)
	copy(tmp, arr)
	for i := 0; i < m; i++ {
		for {
			x := rng.Intn(n)
			if tmp[x] > 0 {
				y := rng.Intn(tmp[x]) + 1
				ops = append(ops, op{x + 1, y})
				left := y - 1
				right := tmp[x] - y
				tmp[x] = 0
				if x-1 >= 0 {
					tmp[x-1] += left
				}
				if x+1 < n {
					tmp[x+1] += right
				}
				break
			}
		}
	}
	final := simulateA(n, arr, ops)
	var in bytes.Buffer
	fmt.Fprintf(&in, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", v)
	}
	in.WriteByte('\n')
	fmt.Fprintf(&in, "%d\n", len(ops))
	for _, o := range ops {
		fmt.Fprintf(&in, "%d %d\n", o.x, o.y)
	}
	var out bytes.Buffer
	for _, v := range final {
		fmt.Fprintf(&out, "%d\n", v)
	}
	return in.String(), out.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
