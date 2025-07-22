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

func runCandidate(bin string, input string) (string, error) {
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

func solve(n, k int, b int64, a []int64) int {
	can := func(x int) bool {
		drains := make([]int64, 0, n)
		if x < n {
			for i := 1; i <= n; i++ {
				if i == x || i == n {
					continue
				}
				drains = append(drains, a[i-1])
			}
		} else {
			for i := 1; i < n; i++ {
				if i == x {
					continue
				}
				drains = append(drains, a[i-1])
			}
		}
		need := k - 1
		if need <= 0 {
			return b < a[x-1]
		}
		if len(drains) == 0 {
			return b < a[x-1]
		}
		sort.Slice(drains, func(i, j int) bool { return drains[i] > drains[j] })
		var s int64
		for i := 0; i < need && i < len(drains); i++ {
			s += drains[i]
			if s > b {
				break
			}
		}
		return s > b-a[x-1]
	}
	l, r := 1, n
	ans := n
	for l <= r {
		mid := (l + r) / 2
		if can(mid) {
			ans = mid
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 2 // n>=2
	k := rng.Intn(n-1) + 1
	b := rng.Int63n(1000)
	a := make([]int64, n)
	for i := range a {
		a[i] = rng.Int63n(1000) + 1
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d\n", n, k)
	fmt.Fprintf(&in, "%d\n", b)
	for i, v := range a {
		if i+1 == n {
			fmt.Fprintf(&in, "%d\n", v)
		} else {
			fmt.Fprintf(&in, "%d ", v)
		}
	}
	exp := fmt.Sprintf("%d", solve(n, k, b, a))
	return in.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
