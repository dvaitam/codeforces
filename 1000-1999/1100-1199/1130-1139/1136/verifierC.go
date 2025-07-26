package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expected(n, m int, A, B [][]int) string {
	for sum := 0; sum <= n+m-2; sum++ {
		var da, db []int
		for i := 0; i < n; i++ {
			j := sum - i
			if j >= 0 && j < m {
				da = append(da, A[i][j])
				db = append(db, B[i][j])
			}
		}
		sort.Ints(da)
		sort.Ints(db)
		if len(da) != len(db) {
			return "NO"
		}
		for i := range da {
			if da[i] != db[i] {
				return "NO"
			}
		}
	}
	return "YES"
}

func runCase(bin, input, want string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n, m int
		A, B [][]int
	}
	tests := []test{
		{n: 1, m: 1, A: [][]int{{5}}, B: [][]int{{5}}},
		{n: 2, m: 2, A: [][]int{{1, 2}, {3, 4}}, B: [][]int{{1, 3}, {2, 4}}},
		{n: 2, m: 2, A: [][]int{{1, 2}, {3, 4}}, B: [][]int{{1, 2}, {4, 3}}},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		A := make([][]int, n)
		B := make([][]int, n)
		for r := 0; r < n; r++ {
			A[r] = make([]int, m)
			B[r] = make([]int, m)
			for c := 0; c < m; c++ {
				A[r][c] = rng.Intn(1000)
				B[r][c] = rng.Intn(1000)
			}
		}
		if i%2 == 0 {
			// make B reachable by shuffling diagonals
			for sum := 0; sum <= n+m-2; sum++ {
				var idxs []int
				for r := 0; r < n; r++ {
					c := sum - r
					if c >= 0 && c < m {
						idxs = append(idxs, r)
					}
				}
				vals := make([]int, len(idxs))
				for j, r := range idxs {
					c := sum - r
					vals[j] = A[r][c]
				}
				rng.Shuffle(len(vals), func(a, b int) { vals[a], vals[b] = vals[b], vals[a] })
				for j, r := range idxs {
					c := sum - r
					B[r][c] = vals[j]
				}
			}
		}
		tests = append(tests, test{n: n, m: m, A: A, B: B})
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", tc.A[i][j]))
			}
			sb.WriteByte('\n')
		}
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", tc.B[i][j]))
			}
			sb.WriteByte('\n')
		}
		want := expected(tc.n, tc.m, tc.A, tc.B)
		if err := runCase(bin, sb.String(), want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
