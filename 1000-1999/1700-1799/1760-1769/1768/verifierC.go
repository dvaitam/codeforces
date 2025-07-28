package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func solveC(n int, v []int) (bool, []int, []int) {
	c := make([]int, n+1)
	for _, x := range v {
		c[x]++
	}
	pq := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if c[i] == 0 {
			pq = append(pq, i)
		}
	}
	sort.Ints(pq)
	pair := make([]int, n+1)
	f := false
	for i := n; i >= 1; i-- {
		switch c[i] {
		case 1:
			pair[i] = i
		case 2:
			if len(pq) == 0 {
				f = true
				break
			}
			m := pq[len(pq)-1]
			if m < i {
				pq = pq[:len(pq)-1]
				pair[i] = m
				pair[m] = i
			} else {
				f = true
			}
		default:
			if c[i] > 2 {
				f = true
			}
		}
		if f {
			break
		}
	}
	if f {
		return false, nil, nil
	}
	vis := make([]bool, n+1)
	p := make([]int, n)
	q := make([]int, n)
	for i, x := range v {
		if !vis[x] {
			p[i] = x
			q[i] = pair[x]
			vis[x] = true
		} else {
			p[i] = pair[x]
			q[i] = x
		}
	}
	return true, p, q
}

func isPermutation(arr []int, n int) bool {
	seen := make([]bool, n+1)
	for _, v := range arr {
		if v < 1 || v > n || seen[v] {
			return false
		}
		seen[v] = true
	}
	for i := 1; i <= n; i++ {
		if !seen[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n int
		a []int
	}

	var cases []test
	// deterministic cases
	cases = append(cases, test{n: 1, a: []int{1}})
	cases = append(cases, test{n: 2, a: []int{2, 1}})
	cases = append(cases, test{n: 3, a: []int{3, 3, 1}})

	for len(cases) < 100 {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(n) + 1
		}
		cases = append(cases, test{n: n, a: arr})
	}

	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		ok, _, _ := solveC(tc.n, append([]int(nil), tc.a...))
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fmt.Fprintf(os.Stderr, "case %d: empty output\n", idx+1)
			os.Exit(1)
		}
		ans := strings.ToUpper(fields[0])
		if ans == "NO" {
			if ok {
				fmt.Fprintf(os.Stderr, "case %d failed: expected YES got NO\n", idx+1)
				os.Exit(1)
			}
			continue
		}
		if ans != "YES" {
			fmt.Fprintf(os.Stderr, "case %d: expected YES/NO got %q\n", idx+1, fields[0])
			os.Exit(1)
		}
		if len(fields)-1 != 2*tc.n {
			fmt.Fprintf(os.Stderr, "case %d: wrong number of integers\n", idx+1)
			os.Exit(1)
		}
		p := make([]int, tc.n)
		q := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			val, err := strconv.Atoi(fields[1+i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: cannot parse p\n", idx+1)
				os.Exit(1)
			}
			p[i] = val
		}
		for i := 0; i < tc.n; i++ {
			val, err := strconv.Atoi(fields[1+tc.n+i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: cannot parse q\n", idx+1)
				os.Exit(1)
			}
			q[i] = val
		}
		if !isPermutation(p, tc.n) || !isPermutation(q, tc.n) {
			fmt.Fprintf(os.Stderr, "case %d: output is not permutation\n", idx+1)
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			m := p[i]
			if q[i] > m {
				m = q[i]
			}
			if m != tc.a[i] {
				fmt.Fprintf(os.Stderr, "case %d: condition failed\n", idx+1)
				os.Exit(1)
			}
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "case %d failed: expected NO got YES\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
