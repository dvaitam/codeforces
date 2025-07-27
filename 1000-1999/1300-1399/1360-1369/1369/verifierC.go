package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solveC(n, k int, a []int, w []int) int64 {
	sort.Ints(a)
	sort.Ints(w)
	l := 0
	r := n - 1
	var res int64
	for i := 0; i < k; i++ {
		res += int64(a[r])
		if w[i] == 1 {
			res += int64(a[r])
		}
		r--
	}
	for i := k - 1; i >= 0; i-- {
		if w[i] == 1 {
			continue
		}
		res += int64(a[l])
		l += w[i] - 1
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(15) + 1
		k := rand.Intn(n) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(100) + 1
		}
		w := make([]int, k)
		for i := range w {
			w[i] = 1
		}
		remaining := n - k
		for j := 0; j < remaining; j++ {
			idx := rand.Intn(k)
			w[idx]++
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		for i, v := range w {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		expected := fmt.Sprintf("%d\n", solveC(n, k, append([]int(nil), a...), append([]int(nil), w...)))
		output, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", t+1, sb.String(), expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
