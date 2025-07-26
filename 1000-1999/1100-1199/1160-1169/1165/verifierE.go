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

const numTestsE = 100
const mod = 998244353

func solveE(a, b []uint64) uint64 {
	n := len(a)
	c := make([]uint64, n)
	nn := uint64(n)
	for i := 0; i < n; i++ {
		idx := uint64(i + 1)
		w := idx * (nn - idx + 1)
		c[i] = w * a[i]
	}
	sort.Slice(c, func(i, j int) bool { return c[i] > c[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	var ans uint64
	for i := 0; i < n; i++ {
		ci := c[i] % mod
		bi := b[i] % mod
		ans = (ans + ci*bi) % mod
	}
	return ans
}

func run(binary, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(5)
	for t := 1; t <= numTestsE; t++ {
		n := rand.Intn(30) + 1
		a := make([]uint64, n)
		b := make([]uint64, n)
		for i := 0; i < n; i++ {
			a[i] = uint64(rand.Intn(1000))
		}
		for i := 0; i < n; i++ {
			b[i] = uint64(rand.Intn(1000))
		}
		input := fmt.Sprintf("%d\n", n)
		for i, v := range a {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		for i, v := range b {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		expect := solveE(append([]uint64(nil), a...), append([]uint64(nil), b...))
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d failed to run: %v\noutput:%s\n", t, err, out)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fmt.Printf("test %d: no output\n", t)
			os.Exit(1)
		}
		var got uint64
		fmt.Sscanf(fields[0], "%d", &got)
		if got != expect {
			fmt.Printf("test %d failed\ninput:%sexpected:%d got:%d\n", t, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
