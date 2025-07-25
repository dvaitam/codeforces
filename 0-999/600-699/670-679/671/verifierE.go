package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseE struct {
	n int
	k int64
	w []int64
	g []int64
}

func genTestsE() []testCaseE {
	rand.Seed(5)
	tests := make([]testCaseE, 100)
	for i := range tests {
		n := rand.Intn(5) + 2 // 2..6
		k := int64(rand.Intn(10))
		w := make([]int64, n)
		for j := 1; j < n; j++ {
			w[j] = int64(rand.Intn(5) + 1)
		}
		g := make([]int64, n+1)
		for j := 1; j <= n; j++ {
			g[j] = int64(rand.Intn(5))
		}
		tests[i] = testCaseE{n, k, w, g}
	}
	return tests
}

// solver copied from 671E.go
func solveE(tc testCaseE) int {
	n := tc.n
	k := tc.k
	w := tc.w
	g := tc.g
	W := make([]int64, n+2)
	G := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		W[i] = W[i-1]
		if i >= 2 {
			W[i] += w[i-1]
		}
		G[i] = G[i-1] + g[i]
	}
	A := make([]int64, n)
	B := make([]int64, n)
	for j := 1; j < n; j++ {
		A[j] = W[j+1] - G[j]
		B[j] = W[j] - G[j]
	}
	check := func(L int) bool {
		if L <= 1 {
			return true
		}
		d := L - 1
		qa := make([]int, 0, n)
		qb := make([]int, 0, n)
		for j := 1; j < n; j++ {
			start := j - d + 1
			if len(qa) > 0 && qa[0] < start {
				qa = qa[1:]
			}
			if len(qb) > 0 && qb[0] < start {
				qb = qb[1:]
			}
			for len(qa) > 0 && A[qa[len(qa)-1]] <= A[j] {
				qa = qa[:len(qa)-1]
			}
			qa = append(qa, j)
			for len(qb) > 0 && B[qb[len(qb)-1]] >= B[j] {
				qb = qb[:len(qb)-1]
			}
			qb = append(qb, j)
			if j >= d {
				l := j - d + 1
				r := l + L - 1
				ma := A[qa[0]]
				mb := B[qb[0]]
				C := W[l] - G[l-1]
				E := W[r] - G[r]
				fdef := ma - C
				if fdef < 0 {
					fdef = 0
				}
				bdef := E - mb
				if bdef < 0 {
					bdef = 0
				}
				need := fdef
				if bdef > need {
					need = bdef
				}
				if need <= k {
					return true
				}
			}
		}
		return false
	}
	low, high := 1, n
	ans := 1
	for low <= high {
		mid := (low + high) / 2
		if check(mid) {
			ans = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return ans
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for j := 1; j < tc.n; j++ {
			if j > 1 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, tc.w[j])
		}
		input.WriteByte('\n')
		for j := 1; j <= tc.n; j++ {
			if j > 1 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, tc.g[j])
		}
		input.WriteByte('\n')
		expect := solveE(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
