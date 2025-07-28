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

type testCaseF struct {
	n int
	k int64
	a []int64
	b []int
}

func generateCaseF(rng *rand.Rand) (string, testCaseF) {
	n := rng.Intn(5) + 3
	k := int64(rng.Intn(20) + 1)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(int(k)) + 1)
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 1
		} else {
			b[i] = 2
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d %d\n", n, k))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), testCaseF{n, k, a, b}
}

func solveCaseF(n int, k int64, a []int64, b []int) []int64 {
	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + a[i]
	}
	L := pref[n]
	cost1 := []int{}
	for i := 0; i < n; i++ {
		if b[i] == 1 {
			cost1 = append(cost1, i)
		}
	}
	if len(cost1) == 0 {
		ans := make([]int64, n)
		for i := range ans {
			ans[i] = 2 * L
		}
		return ans
	}
	dist := func(i, j int) int64 {
		if j >= i {
			return pref[j] - pref[i]
		}
		return L - (pref[i] - pref[j])
	}
	base := int64(0)
	m := len(cost1)
	for idx := 0; idx < m; idx++ {
		i := cost1[idx]
		j := cost1[(idx+1)%m]
		var D int64
		if m == 1 {
			D = L
		} else {
			D = dist(i, j)
		}
		if D > k {
			base += D - k
		}
	}
	cost1Cost := L + base
	ans := make([]int64, n)
	for s := 0; s < n; s++ {
		if b[s] == 1 {
			ans[s] = cost1Cost
			continue
		}
		pos := sort.SearchInts(cost1, s)
		prev := cost1[(pos-1+len(cost1))%len(cost1)]
		nxt := cost1[0]
		if pos < len(cost1) {
			nxt = cost1[pos]
		}
		var D int64
		if prev == nxt {
			D = L
		} else {
			D = dist(prev, nxt)
		}
		x := dist(prev, s)
		var add int64
		if D <= k {
			add = D - x
		} else if k > x {
			add = k - x
		}
		ans[s] = cost1Cost + add
	}
	return ans
}

func expectedF(tc testCaseF) string {
	arr := solveCaseF(tc.n, tc.k, tc.a, tc.b)
	var sb strings.Builder
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCaseF(rng)
		expect := expectedF(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
