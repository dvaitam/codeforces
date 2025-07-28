package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveB(a, b []int) (int, int) {
	n := len(a)
	l := 0
	for l < n && a[l] == b[l] {
		l++
	}
	if l == n {
		return 1, n
	}
	r := n - 1
	for r >= 0 && a[r] == b[r] {
		r--
	}
	for l > 0 && a[l-1] <= b[l] {
		l--
	}
	for r < n-1 && a[r+1] >= b[r] {
		r++
	}
	return l + 1, r + 1
}

func genCaseB(rng *rand.Rand) (int, []int, []int) {
	n := rng.Intn(10) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(n) + 1
	}
	l := rng.Intn(n)
	r := l + rng.Intn(n-l)
	b := append([]int(nil), a...)
	sort.Ints(b[l : r+1])
	if equalInts(a, b) {
		if r < n-1 {
			r++
		} else if l > 0 {
			l--
		}
		b = append([]int(nil), a...)
		sort.Ints(b[l : r+1])
	}
	return n, a, b
}

func equalInts(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func runCandidate(bin, input string) (string, error) {
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, a, b := genCaseB(rng)
		input := fmt.Sprintf("1\n%d\n", n)
		for j, v := range a {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		for j, v := range b {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		l, r := solveB(a, b)
		expect := fmt.Sprintf("%d %d", l, r)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
