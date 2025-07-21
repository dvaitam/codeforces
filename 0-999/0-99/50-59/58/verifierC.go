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

func expected(a []int) int {
	n := len(a)
	freq := map[int]int{}
	maxf := 0
	for i, v := range a {
		d := i
		if n-1-i < d {
			d = n - 1 - i
		}
		x := v - d
		if x > 0 {
			freq[x]++
			if freq[x] > maxf {
				maxf = freq[x]
			}
		}
	}
	return n - maxf
}

func runCase(bin string, a []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
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
	gotStr := strings.TrimSpace(out.String())
	var got int
	fmt.Sscan(gotStr, &got)
	exp := expected(a)
	if got != exp {
		return fmt.Errorf("expected %d got %s", exp, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([][]int, 0, 100)
	cases = append(cases, []int{1})
	cases = append(cases, []int{2, 2, 2})
	for len(cases) < 100 {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rng.Intn(20) + 1
		}
		cases = append(cases, arr)
	}
	for i, a := range cases {
		if err := runCase(bin, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %v\n", i+1, err, a)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
