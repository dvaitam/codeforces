package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(n uint64) (res []uint64) {
	if bits.OnesCount64(n) == 1 {
		return []uint64{n}
	}
	for i := uint(0); i < 64; i++ {
		if n&(1<<i) != 0 {
			res = append(res, n^(1<<i))
		}
	}
	res = append(res, n)
	return
}

func generate() (string, string) {
	const T = 100
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", T)
	rand.Seed(3)
	for i := 0; i < T; i++ {
		n := rand.Uint64()
		fmt.Fprintf(&in, "%d\n", n)
		ans := solve(n)
		fmt.Fprintf(&out, "%d\n", len(ans))
		for j, v := range ans {
			if j+1 == len(ans) {
				fmt.Fprintf(&out, "%d\n", v)
			} else {
				fmt.Fprintf(&out, "%d ", v)
			}
		}
	}
	return in.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	in, exp := generate()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}
	got := buf.String()
	if strings.TrimSpace(got) != strings.TrimSpace(exp) {
		fmt.Println("wrong answer")
		fmt.Println("expected:\n" + exp)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
