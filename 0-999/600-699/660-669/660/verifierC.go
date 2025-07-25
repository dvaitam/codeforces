package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveC(n, k int, arr []int) (int, []int) {
	l, zeroes := 0, 0
	bestLen, bestL := 0, 0
	for r := 0; r < n; r++ {
		if arr[r] == 0 {
			zeroes++
		}
		for zeroes > k {
			if arr[l] == 0 {
				zeroes--
			}
			l++
		}
		if r-l+1 > bestLen {
			bestLen = r - l + 1
			bestL = l
		}
	}
	res := append([]int(nil), arr...)
	for i := bestL; i < bestL+bestLen; i++ {
		res[i] = 1
	}
	return bestLen, res
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(20) + 1
		k := rand.Intn(n + 1)
		arr := make([]int, n)
		for i := range arr {
			if rand.Intn(2) == 0 {
				arr[i] = 0
			} else {
				arr[i] = 1
			}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		best, res := solveC(n, k, arr)
		var exp strings.Builder
		fmt.Fprintln(&exp, best)
		for i, v := range res {
			if i > 0 {
				exp.WriteByte(' ')
			}
			fmt.Fprint(&exp, v)
		}
		exp.WriteByte('\n')
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
