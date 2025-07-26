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

func run(bin string, in []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func mex(arr []int) int {
	used := make(map[int]bool)
	for _, v := range arr {
		used[v] = true
	}
	for i := 0; ; i++ {
		if !used[i] {
			return i
		}
	}
}

func genCase() (string, [][]int) {
	n := rand.Intn(5) + 1
	ops := make([][]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		a := rand.Intn(n + 1)
		l := rand.Intn(n + 1)
		r := rand.Intn(n + 1)
		k := rand.Intn(n + 1)
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a, l, r, k))
		ops[i] = []int{a, l, r, k}
	}
	return sb.String(), ops
}

func solveNaive(ops [][]int) []int {
	n := len(ops)
	arr := make([]int, 0, n)
	last := 0
	res := make([]int, n)
	for i := 0; i < n; i++ {
		a := (ops[i][0] + last) % (n + 1)
		l := (ops[i][1]+last)%(i+1) + 1
		r := (ops[i][2]+last)%(i+1) + 1
		if l > r {
			l, r = r, l
		}
		k := (ops[i][3] + last) % (n + 1)
		arr = append(arr, a)
		ans := 0
		for x := l - 1; x < r; x++ {
			for y := x; y < r; y++ {
				if mex(arr[x:y+1]) == k {
					ans++
				}
			}
		}
		res[i] = ans
		last = ans
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		input, ops := genCase()
		out, err := run(bin, []byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) != len(ops) {
			fmt.Fprintf(os.Stderr, "wrong number of answers on test %d\n", t+1)
			os.Exit(1)
		}
		cand := make([]int, len(ops))
		for i := range fields {
			v, err := strconv.Atoi(fields[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "bad output on test %d\n", t+1)
				os.Exit(1)
			}
			cand[i] = v
		}
		exp := solveNaive(ops)
		for i := range exp {
			if cand[i] != exp[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%soutput:\n%s", t+1, input, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
