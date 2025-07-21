package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type query struct{ l, r int }
type testCaseE struct {
	n int
	s string
	t string
	q []query
}

func genTestsE() []testCaseE {
	rand.Seed(46)
	tests := make([]testCaseE, 20)
	for i := range tests {
		n := rand.Intn(3) + 3 // 3..5
		s := randomBinary(n)
		t := randomBinary(n)
		qn := rand.Intn(3) + 1
		qs := make([]query, qn)
		for j := 0; j < qn; j++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			qs[j] = query{l, r}
		}
		tests[i] = testCaseE{n, s, t, qs}
	}
	return tests
}

func randomBinary(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rand.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func bitsFromString(s string) int {
	x := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			x |= 1 << i
		}
	}
	return x
}

type state struct{ a, b int }

func maxOnes(aStart, bStart int, k int) int {
	vis := make(map[state]bool)
	q := []state{{aStart, bStart}}
	vis[state{aStart, bStart}] = true
	best := bits.OnesCount(uint(aStart))
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if c := bits.OnesCount(uint(cur.a)); c > best {
			best = c
		}
		for i := 0; i < k-2; i++ {
			if (cur.a>>i)&1 == 0 && (cur.a>>(i+2))&1 == 0 {
				nb := cur.b | (1 << (i + 1))
				st := state{cur.a, nb}
				if !vis[st] {
					vis[st] = true
					q = append(q, st)
				}
			}
			if (cur.b>>i)&1 == 1 && (cur.b>>(i+2))&1 == 1 {
				na := cur.a | (1 << (i + 1))
				st := state{na, cur.b}
				if !vis[st] {
					vis[st] = true
					q = append(q, st)
				}
			}
		}
	}
	return best
}

func solveE(tc testCaseE) []int {
	ans := make([]int, len(tc.q))
	for idx, qq := range tc.q {
		a := tc.s[qq.l-1 : qq.r]
		b := tc.t[qq.l-1 : qq.r]
		maskA := bitsFromString(a)
		maskB := bitsFromString(b)
		ans[idx] = maxOnes(maskA, maskB, len(a))
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		fmt.Fprintln(&input, tc.s)
		fmt.Fprintln(&input, tc.t)
		fmt.Fprintln(&input, len(tc.q))
		for _, qv := range tc.q {
			fmt.Fprintf(&input, "%d %d\n", qv.l, qv.r)
		}
	}
	expected := make([][]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveE(tc)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, tc := range tests {
		for j := 0; j < len(tc.q); j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			if val != expected[i][j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
