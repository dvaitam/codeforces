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

type pair struct{ v, d int }

func expected(a []int) int {
	stack := make([]pair, 0, len(a))
	ans := 0
	for _, val := range a {
		death := 0
		for len(stack) > 0 && stack[len(stack)-1].v <= val {
			if stack[len(stack)-1].d > death {
				death = stack[len(stack)-1].d
			}
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			death = 0
		} else {
			death++
		}
		if death > ans {
			ans = death
		}
		stack = append(stack, pair{val, death})
	}
	return ans
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(arr)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := [][]int{
		{1},
		{2, 1},
		{1, 2},
		{3, 2, 1},
		{1, 2, 3},
	}
	for i := 0; i < 95; i++ {
		n := rng.Intn(30) + 1
		perm := rng.Perm(n)
		for j := range perm {
			perm[j]++ // make 1..n
		}
		cases = append(cases, perm)
	}
	for i, arr := range cases {
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\narray: %v\n", i+1, err, arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
