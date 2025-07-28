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

func solve(a []int) []int {
	n := len(a)
	freq := make([]int, n+2)
	for _, v := range a {
		if v <= n+1 {
			freq[v]++
		}
	}
	res := make([]int, 0)
	for i := 0; i < n; {
		mex := 0
		for mex <= n+1 && freq[mex] > 0 {
			mex++
		}
		if mex == 0 {
			freq[a[i]]--
			res = append(res, 0)
			i++
			continue
		}
		seen := make([]bool, mex)
		need := mex
		j := i
		for j < n && need > 0 {
			v := a[j]
			freq[v]--
			if v < mex && !seen[v] {
				seen[v] = true
				need--
			}
			j++
		}
		res = append(res, mex)
		i = j
	}
	return res
}

func buildInput(n int, a []int) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(res []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return strings.TrimSpace(sb.String())
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rng.Intn(n + 1)
		}
		input := buildInput(n, a)
		res := solve(a)
		exp := expected(res)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Printf("case %d wrong answer\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
