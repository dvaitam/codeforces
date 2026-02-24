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

const MOD = 1000000007

type Cow struct {
	posL, posR int
}

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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func countLE(a []int, v int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] > v })
}

func solve(input string) string {
	in := strings.NewReader(input)
	var n, m int
	fmt.Fscan(in, &n, &m)
	// 1-based arrays
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	V := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var f, h int
		fmt.Fscan(in, &f, &h)
		V[f] = append(V[f], h)
	}
	for f := 1; f <= n; f++ {
		sort.Ints(V[f])
	}
	mx := 0
	var ans int64 = 0
	c1 := make([]int, n+1)
	c2 := make([]int, n+1)
	for i := 0; i <= n; i++ {
		for j := 1; j <= n; j++ {
			c1[j], c2[j] = 0, 0
		}
		for x := 1; x <= i; x++ {
			c1[a[x]]++
		}
		for x := i + 1; x <= n; x++ {
			c2[a[x]]++
		}
		ff := (i == 0)
		tt := 0
		if i != 0 {
			tt = 1
		}
		s := int64(1)
		ai := 0
		if i >= 1 {
			ai = a[i]
		}
		// check feasibility for anchor flavor
		if !ff {
			for _, need := range V[ai] {
				if need == c1[ai] {
					ff = true
					break
				}
			}
		}
		if !ff {
			continue
		}
		for j := 1; j <= n; j++ {
			if len(V[j]) == 0 {
				continue
			}
			x := countLE(V[j], c1[j])
			y := countLE(V[j], c2[j])
			if j == ai {
				x = 0
				if c2[j] >= c1[j] && y > 0 {
					y--
				}
			}
			if x == 0 && y == 0 {
				continue
			}
			if x > y {
				x, y = y, x
			}
			if x == 0 {
				s = (s * int64(y)) % MOD
				tt++
			} else if y == 1 {
				s = (s * 2) % MOD
				tt++
			} else {
				s = (s * int64(x)) % MOD
				s = (s * int64(y-1)) % MOD
				tt += 2
			}
		}
		if tt > mx {
			mx = tt
			ans = 0
		}
		if tt == mx {
			ans = (ans + s) % MOD
		}
	}
	return fmt.Sprintf("%d %d", mx, ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	m := rng.Intn(6) + 1
	if m > n*n {
		m = n * n
	}
	s := make([]int, n)
	for i := range s {
		s[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, v := range s {
		if i+1 == n {
			fmt.Fprintf(&sb, "%d\n", v)
		} else {
			fmt.Fprintf(&sb, "%d ", v)
		}
	}
	
	type pair struct{ f, h int }
	pairs := []pair{}
	for f := 1; f <= n; f++ {
		for h := 1; h <= n; h++ {
			pairs = append(pairs, pair{f, h})
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})
	
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d\n", pairs[i].f, pairs[i].h)
	}
	in := sb.String()
	return in, solve(in)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
