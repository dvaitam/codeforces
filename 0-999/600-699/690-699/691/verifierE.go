package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod int64 = 1000000007

type testCase struct {
	in  string
	out string
}

func matMul(a, b [][]int64) [][]int64 {
	n := len(a)
	res := make([][]int64, n)
	for i := range res {
		res[i] = make([]int64, n)
	}
	for i := 0; i < n; i++ {
		ai := a[i]
		ri := res[i]
		for k := 0; k < n; k++ {
			if ai[k] == 0 {
				continue
			}
			aik := ai[k]
			bk := b[k]
			for j := 0; j < n; j++ {
				if bk[j] == 0 {
					continue
				}
				ri[j] = (ri[j] + aik*bk[j]) % mod
			}
		}
	}
	return res
}

func matVecMul(a [][]int64, v []int64) []int64 {
	n := len(a)
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		row := a[i]
		var sum int64
		for j := 0; j < n; j++ {
			if row[j] == 0 || v[j] == 0 {
				continue
			}
			sum = (sum + row[j]*v[j]) % mod
		}
		res[i] = sum
	}
	return res
}

func powMatVec(m [][]int64, e int64, vec []int64) []int64 {
	for e > 0 {
		if e&1 == 1 {
			vec = matVecMul(m, vec)
		}
		m = matMul(m, m)
		e >>= 1
	}
	return vec
}

func solveCase(n int, k int64, arr []int64) string {
	if k == 1 {
		return fmt.Sprintf("%d\n", n%int(mod))
	}
	adj := make([][]int64, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			if bits.OnesCount64(uint64(arr[i]^arr[j]))%3 == 0 {
				adj[i][j] = 1
			}
		}
	}
	vec := make([]int64, n)
	for i := range vec {
		vec[i] = 1
	}
	vec = powMatVec(adj, k-1, vec)
	var ans int64
	for _, v := range vec {
		ans = (ans + v) % mod
	}
	return fmt.Sprintf("%d\n", ans)
}

func buildCase(n int, k int64, arr []int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{in: sb.String(), out: solveCase(n, k, arr)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	k := int64(rng.Intn(10) + 1)
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rng.Intn(1000))
	}
	return buildCase(n, k, arr)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(tc.out)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase(1, 1, []int64{0}))
	cases = append(cases, buildCase(2, 2, []int64{1, 2}))
	cases = append(cases, buildCase(3, 3, []int64{1, 1, 1}))
	cases = append(cases, buildCase(2, 5, []int64{5, 10}))
	cases = append(cases, buildCase(4, 4, []int64{1, 2, 3, 4}))
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
