package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type pair struct {
	x int64
	y int64
}

func solveD(n int, xs, ys []int64) int64 {
	cnt := make(map[pair]int64)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p := pair{xs[i] + xs[j], ys[i] + ys[j]}
			cnt[p]++
		}
	}
	var res int64
	for _, v := range cnt {
		res += v * (v - 1) / 2
	}
	return res
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(4)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(10) + 2
		xs := make([]int64, n)
		ys := make([]int64, n)
		for i := 0; i < n; i++ {
			xs[i] = int64(rand.Intn(20))
			ys[i] = int64(rand.Intn(20))
		}
		var sb strings.Builder
		fmt.Fprintln(&sb, n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", xs[i], ys[i])
		}
		expected := solveD(n, xs, ys)
		var exp strings.Builder
		fmt.Fprintln(&exp, expected)
		output, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(output)
		want := strings.TrimSpace(exp.String())
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", t+1, sb.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
