package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n int
	k int
	p []int
}

func expected(tc Test) int {
	mism := 0
	for i := 0; i < tc.n; i++ {
		if tc.p[i]%tc.k != (i+1)%tc.k {
			mism++
		}
	}
	if mism == 0 {
		return 0
	} else if mism == 2 {
		return 1
	}
	return -1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(2)
	const cases = 100
	tests := make([]Test, cases)
	for i := range tests {
		n := rand.Intn(8) + 2
		k := rand.Intn(n-1) + 1
		p := rand.Perm(n)
		for i2 := range p {
			p[i2]++
		}
		tests[i] = Test{n: n, k: k, p: p}
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.p {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("error running binary:", err)
		fmt.Print(out.String())
		return
	}

	reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
	for idx, tc := range tests {
		var ans int
		if _, err := fmt.Fscan(reader, &ans); err != nil {
			fmt.Printf("test %d: failed to read output\n", idx+1)
			return
		}
		exp := expected(tc)
		if ans != exp {
			fmt.Printf("test %d: expected %d got %d\n", idx+1, exp, ans)
			return
		}
	}
	fmt.Printf("verified %d test cases\n", len(tests))
}
