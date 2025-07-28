package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func nextPermutation(a []int) bool {
	n := len(a)
	i := n - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func expected(n, s int) []int {
	mod := 998244353
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	ans := make([]int, n)
	for {
		cntS := 0
		for i := 0; i < n-1; i++ {
			if p[i]+1 == p[i+1] {
				cntS++
			}
		}
		if cntS == s {
			cntK := 0
			for i := 0; i < n-1; i++ {
				if p[i] < p[i+1] {
					cntK++
				}
			}
			ans[cntK] = (ans[cntK] + 1) % mod
		}
		if !nextPermutation(p) {
			break
		}
	}
	return ans
}

func genTest(rng *rand.Rand) (int, int) {
	n := rng.Intn(7) + 1
	s := rng.Intn(n)
	return n, s
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(4))
	for i := 0; i < 100; i++ {
		n, s := genTest(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, s))
		expectArr := expected(n, s)
		expectStr := make([]string, n)
		for j, v := range expectArr {
			expectStr[j] = fmt.Sprint(v)
		}
		expect := strings.Join(expectStr, " ")
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
