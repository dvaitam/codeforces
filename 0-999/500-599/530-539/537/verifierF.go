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

func expectedAnswerF(a []int64) string {
	n := len(a) - 1
	ans := make([]int, n)
	const K = 350
	smallK := K
	if smallK > n-1 {
		smallK = n - 1
	}
	for k := 1; k <= smallK; k++ {
		cnt := 0
		for j := 1; ; j++ {
			l := k*(j-1) + 2
			if l > n {
				break
			}
			r := k*j + 1
			if r > n {
				r = n
			}
			pj := a[j]
			for v := l; v <= r; v++ {
				if a[v] < pj {
					cnt++
				}
			}
		}
		ans[k] = cnt
	}
	diff := make([]int, n+2)
	for u := 2; u <= n; u++ {
		x := u - 2
		maxP := x/(smallK+1) + 1
		if maxP > u-1 {
			maxP = u - 1
		}
		for p := 1; p <= maxP; p++ {
			if a[p] <= a[u] {
				continue
			}
			var L, R int
			if p == 1 {
				L = x + 1
				R = n - 1
			} else {
				L = x/p + 1
				R = x / (p - 1)
			}
			if L <= smallK {
				L = smallK + 1
			}
			if L > R || L > n-1 {
				continue
			}
			if R > n-1 {
				R = n - 1
			}
			diff[L]++
			diff[R+1]--
		}
	}
	cur := 0
	for k := smallK + 1; k <= n-1; k++ {
		cur += diff[k]
		ans[k] = cur
	}
	var sb strings.Builder
	for k := 1; k <= n-1; k++ {
		if k > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(ans[k]))
	}
	return sb.String()
}

func generateCaseF(rng *rand.Rand) []int64 {
	n := rng.Intn(30) + 2
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		a[i] = int64(rng.Intn(100) - 50)
	}
	return a
}

func runCaseF(bin string, a []int64) error {
	var sb strings.Builder
	n := len(a) - 1
	sb.WriteString(fmt.Sprint(n, "\n"))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedAnswerF(a)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a := generateCaseF(rng)
		if err := runCaseF(bin, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
