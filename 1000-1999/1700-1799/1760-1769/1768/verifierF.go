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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func solveF(n int, a []int) []int64 {
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i/2] + 1
	}
	K := log[n] + 1
	st := make([][]int, K)
	st[0] = make([]int, n)
	copy(st[0], a)
	for k := 1; k < K; k++ {
		size := n - (1 << k) + 1
		st[k] = make([]int, size)
		for i := 0; i < size; i++ {
			x := st[k-1][i]
			y := st[k-1][i+(1<<(k-1))]
			if x < y {
				st[k][i] = x
			} else {
				st[k][i] = y
			}
		}
	}
	query := func(l, r int) int {
		if l > r {
			l, r = r, l
		}
		k := log[r-l+1]
		x := st[k][l]
		y := st[k][r-(1<<k)+1]
		if x < y {
			return x
		}
		return y
	}

	dp := make([]int64, n)
	const INF int64 = 1 << 60
	for i := 1; i < n; i++ {
		dp[i] = INF
	}
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			m := query(j, i)
			d := int64(i - j)
			cost := int64(m) * d * d
			if dp[j]+cost < dp[i] {
				dp[i] = dp[j] + cost
			}
		}
	}
	return dp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n int
		a []int
	}

	var cases []test
	cases = append(cases, test{n: 2, a: []int{1, 2}})
	cases = append(cases, test{n: 3, a: []int{2, 1, 3}})
	cases = append(cases, test{n: 4, a: []int{4, 4, 4, 4}})

	for len(cases) < 100 {
		n := rng.Intn(8) + 2
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(n) + 1
		}
		cases = append(cases, test{n: n, a: arr})
	}

	for i, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j, v := range tc.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		expected := solveF(tc.n, tc.a)
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) != tc.n {
			fmt.Fprintf(os.Stderr, "case %d: expected %d numbers got %q\n", i+1, tc.n, out)
			os.Exit(1)
		}
		for idx := 0; idx < tc.n; idx++ {
			val, err := strconv.ParseInt(fields[idx], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: cannot parse integer\n", i+1)
				os.Exit(1)
			}
			if val != expected[idx] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, expected[idx], val, sb.String())
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
