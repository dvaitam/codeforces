package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solvePerm(a []int) []int {
	n := len(a)
	used := make([]bool, n-1)
	for x := 1; x <= n; x++ {
		pos := 0
		for pos < n && a[pos] != x {
			pos++
		}
		for pos > 0 && !used[pos-1] && a[pos-1] > a[pos] {
			a[pos], a[pos-1] = a[pos-1], a[pos]
			used[pos-1] = true
			pos--
		}
	}
	return a
}

func generate() (string, string) {
	const T = 100
	rand.Seed(2)
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", T)
	for i := 0; i < T; i++ {
		n := rand.Intn(10) + 2
		p := rand.Perm(n)
		for j := 0; j < n; j++ {
			p[j]++
		}
		tmp := append([]int(nil), p...)
		res := solvePerm(tmp)
		fmt.Fprintf(&in, "%d\n", n)
		for j := 0; j < n; j++ {
			if j+1 == n {
				fmt.Fprintf(&in, "%d\n", p[j])
			} else {
				fmt.Fprintf(&in, "%d ", p[j])
			}
		}
		for j := 0; j < n; j++ {
			if j+1 == n {
				fmt.Fprintf(&out, "%d\n", res[j])
			} else {
				fmt.Fprintf(&out, "%d ", res[j])
			}
		}
	}
	return in.String(), out.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	in, exp := generate()
	out, err := runCandidate(bin, in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Fprintln(os.Stderr, "wrong answer")
		fmt.Fprintln(os.Stderr, "expected:\n"+exp)
		fmt.Fprintln(os.Stderr, "got:\n"+out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
