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

type testF struct {
	n int
	m int
	l []int
	r []int
}

func solveF(tc testF) int {
	total := 0
	for a := 1; a <= tc.m; a++ {
		for b := a; b <= tc.m; b++ {
			valid := false
			ok := true
			for i := 0; i < tc.n; i++ {
				start := a
				if tc.l[i] > start {
					start = tc.l[i]
				}
				end := b
				if tc.r[i] < end {
					end = tc.r[i]
				}
				length := end - start + 1
				if start > end {
					length = 0
				}
				if length > 0 {
					valid = true
				}
				if length%2 == 0 {
					if length != 0 {
						ok = false
						break
					}
				}
			}
			if ok && valid {
				total += b - a + 1
			}
		}
	}
	return total
}

func genF() (string, int) {
	n := rand.Intn(3) + 1
	m := rand.Intn(6) + 1
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		a := rand.Intn(m) + 1
		b := rand.Intn(m) + 1
		if a > b {
			a, b = b, a
		}
		l[i], r[i] = a, b
	}
	tc := testF{n, m, l, r}
	ans := solveF(tc)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", l[i], r[i]))
	}
	input := sb.String()
	return input, ans
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genF()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s", i+1, err, in, got)
			return
		}
		got = strings.TrimSpace(got)
		var val int
		if _, err := fmt.Sscan(got, &val); err != nil || val != exp {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%d\nGot:\n%s", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
