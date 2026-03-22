package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Embedded reference solver from the accepted solution.
func refSolve(input string) string {
	r := strings.NewReader(input)
	var n, m int
	var k int64
	fmt.Fscan(r, &n, &m, &k)
	var s string
	fmt.Fscan(r, &s)

	lcp := make([][]int, n+1)
	for i := range lcp {
		lcp[i] = make([]int, n+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if s[i] == s[j] {
				lcp[i][j] = 1 + lcp[i+1][j+1]
			} else {
				lcp[i][j] = 0
			}
		}
	}

	sa := make([]int, n)
	for i := 0; i < n; i++ {
		sa[i] = i
	}
	sort.Slice(sa, func(i, j int) bool {
		a, b := sa[i], sa[j]
		l := lcp[a][b]
		if a+l == n {
			return true
		}
		if b+l == n {
			return false
		}
		return s[a+l] < s[b+l]
	})

	adjLcp := make([]int, n)
	for i := 0; i < n-1; i++ {
		adjLcp[i] = lcp[sa[i]][sa[i+1]]
	}

	type Sub struct{ p, l int }
	var subs []Sub
	for i := 0; i < n; i++ {
		start := 0
		if i > 0 {
			start = adjLcp[i-1]
		}
		for l := start + 1; l <= n-sa[i]; l++ {
			subs = append(subs, Sub{sa[i], l})
		}
	}

	check := func(mid int) bool {
		xP := subs[mid].p
		xL := subs[mid].l

		head := make([]int, n)
		for i := range head {
			head[i] = -1
		}
		next := make([]int, n)
		for p := 0; p < n; p++ {
			l := lcp[p][xP]
			var me int
			if l >= xL {
				me = p + xL - 1
			} else {
				if p+l >= n {
					me = n
				} else if s[p+l] > s[xP+l] {
					me = p + l
				} else {
					me = n
				}
			}
			if me < n {
				next[p] = head[me]
				head[me] = p
			}
		}

		dp := make([]int64, n+1)
		dp[0] = 1

		for j := 1; j <= m; j++ {
			newDp := make([]int64, n+1)
			sumValid := int64(0)
			for i := 1; i <= n; i++ {
				for p := head[i-1]; p != -1; p = next[p] {
					sumValid += dp[p]
					if sumValid > k {
						sumValid = k + 1
					}
				}
				newDp[i] = sumValid
			}
			dp = newDp
		}

		return dp[n] >= k
	}

	if len(subs) == 0 {
		return ""
	}

	left, right := 0, len(subs)-1
	var ans string
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			ans = s[subs[mid].p : subs[mid].p+subs[mid].l]
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return ans
}

func runProg(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	m := rng.Intn(n) + 1
	k := rng.Intn(5) + 1
	s := randString(rng, n)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d %d\n%s\n", n, m, k, s)
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/bin")
		os.Exit(1)
	}
	target := os.Args[1]

	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		test := genTest(rng)
		expected := refSolve(test)
		got, err := runProg(target, test)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d execution error: %v\n", i+1, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, test, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
