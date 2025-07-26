package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func fValue(arr []int) int {
	tmp := make([]int, len(arr))
	copy(tmp, arr)
	for len(tmp) > 1 {
		next := make([]int, len(tmp)-1)
		for i := 0; i < len(tmp)-1; i++ {
			next[i] = tmp[i] ^ tmp[i+1]
		}
		tmp = next
	}
	return tmp[0]
}

func maxF(a []int, l, r int) int {
	maxv := 0
	for i := l - 1; i < r; i++ {
		for j := i; j < r; j++ {
			v := fValue(a[i : j+1])
			if v > maxv {
				maxv = v
			}
		}
	}
	return maxv
}

func generateB(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(16)
	}
	l := rng.Intn(n) + 1
	r := rng.Intn(n-l+1) + l
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteString("\n1\n")
	fmt.Fprintf(&sb, "%d %d\n", l, r)
	res := maxF(a, l, r)
	return sb.String(), fmt.Sprint(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(43))
	for i := 0; i < 100; i++ {
		in, exp := generateB(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i+1, err, out.String())
			return
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
