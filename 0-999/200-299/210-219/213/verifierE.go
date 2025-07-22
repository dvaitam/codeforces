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

type caseE struct {
	n, m int
	a    []int
	b    []int
}

func genPerm(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func genCaseE(rng *rand.Rand) caseE {
	n := rng.Intn(5) + 1
	m := rng.Intn(10) + n
	a := genPerm(rng, n)
	b := genPerm(rng, m)
	return caseE{n, m, a, b}
}

func matchPat(plBase, glBase, plTxt, glTxt []int, i, j int) bool {
	s := i - j
	if plBase[j] == 0 {
		if plTxt[i] > s {
			return false
		}
	} else {
		if plTxt[i] != s+(j-plBase[j]) {
			return false
		}
	}
	if glBase[j] == 0 {
		if glTxt[i] > s {
			return false
		}
	} else {
		if glTxt[i] != s+(j-glBase[j]) {
			return false
		}
	}
	return true
}

func matchTxt(plPat, glPat, plTxt, glTxt []int, i, j int) bool {
	s := i - j
	if plPat[j] == 0 {
		if plTxt[i] > s {
			return false
		}
	} else {
		if plTxt[i] != s+(j-plPat[j]) {
			return false
		}
	}
	if glPat[j] == 0 {
		if glTxt[i] > s {
			return false
		}
	} else {
		if glTxt[i] != s+(j-glPat[j]) {
			return false
		}
	}
	return true
}

func solveE(tc caseE) int {
	n := tc.n
	m := tc.m
	a := append([]int{0}, tc.a...)
	b := append([]int{0}, tc.b...)
	P := make([]int, m+2)
	for i := 1; i <= m; i++ {
		P[b[i]] = i
	}
	T := make([]int, m+1)
	for v := 1; v <= m; v++ {
		T[v] = P[v]
	}
	plPat := make([]int, n+1)
	glPat := make([]int, n+1)
	stack := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			plPat[i] = 0
		} else {
			plPat[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			glPat[i] = 0
		} else {
			glPat[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	plTxt := make([]int, m+1)
	glTxt := make([]int, m+1)
	stack = stack[:0]
	for i := 1; i <= m; i++ {
		for len(stack) > 0 && T[stack[len(stack)-1]] >= T[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			plTxt[i] = 0
		} else {
			plTxt[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := 1; i <= m; i++ {
		for len(stack) > 0 && T[stack[len(stack)-1]] <= T[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			glTxt[i] = 0
		} else {
			glTxt[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	pi := make([]int, n+1)
	for i := 2; i <= n; i++ {
		j := pi[i-1]
		for j > 0 && !matchPat(plPat, glPat, plPat, glPat, i, j+1) {
			j = pi[j]
		}
		if matchPat(plPat, glPat, plPat, glPat, i, j+1) {
			j++
		}
		pi[i] = j
	}
	q := 0
	count := 0
	for i := 1; i <= m; i++ {
		for q > 0 && !matchTxt(plPat, glPat, plTxt, glTxt, i, q+1) {
			q = pi[q]
		}
		if matchTxt(plPat, glPat, plTxt, glTxt, i, q+1) {
			q++
		}
		if q == n {
			count++
			q = pi[q]
		}
	}
	return count
}

func runE(bin string, tc caseE) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveE(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCaseE(rng)
		if err := runE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
