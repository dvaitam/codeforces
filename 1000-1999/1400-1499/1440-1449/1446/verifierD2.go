package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	a []int
}

func expected(a []int) int {
	n := len(a)
	pos := make([][]int, n+1)
	for i, v := range a {
		if v >= 1 && v <= n {
			pos[v] = append(pos[v], i)
		}
	}
	maxf := 0
	for v := 1; v <= n; v++ {
		if len(pos[v]) > maxf {
			maxf = len(pos[v])
		}
	}
	if maxf < 2 {
		return 0
	}
	cnt := make([]int, maxf+2)
	L := make([]int, maxf+2)
	R := make([]int, maxf+2)
	A := make([]int, maxf+2)
	for k := 1; k <= maxf; k++ {
		L[k] = n
		R[k] = -1
		A[k] = n
	}
	for v := 1; v <= n; v++ {
		pv := pos[v]
		m := len(pv)
		if m == 0 {
			continue
		}
		first := pv[0]
		for k := 1; k <= m; k++ {
			cnt[k]++
			if first < L[k] {
				L[k] = first
			}
			idx := pv[k-1]
			if idx > R[k] {
				R[k] = idx
			}
		}
		for k := 1; k < m; k++ {
			nxt := pv[k]
			if nxt < A[k] {
				A[k] = nxt
			}
		}
	}
	ans := 0
	for k := 1; k <= maxf; k++ {
		if cnt[k] < 2 {
			continue
		}
		if A[k] <= R[k] {
			continue
		}
		length := R[k] - L[k] + 1
		if length > ans {
			ans = length
		}
	}
	return ans
}

func runBin(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCase {
	tests := []testCase{
		{a: []int{1, 1, 2, 2}},
		{a: []int{1, 2, 3, 4}},
		{a: []int{1, 1, 1, 2, 2, 2}},
		{a: []int{5, 5, 5, 5}},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := r.Intn(20) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = r.Intn(n) + 1
		}
		tests = append(tests, testCase{a: arr})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(t.a)))
		for j, v := range t.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		want := expected(t.a)
		out, err := runBin(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse output\n", i+1)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
