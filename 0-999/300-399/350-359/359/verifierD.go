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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveCase(n int, a []int) ([]int, int) {
	logs := make([]int, n+2)
	logs[1] = 0
	for i := 2; i <= n; i++ {
		logs[i] = logs[i/2] + 1
	}
	kmax := logs[n] + 1
	st := make([][]int, kmax)
	st[0] = make([]int, n+1)
	for i := 1; i <= n; i++ {
		st[0][i] = a[i-1]
	}
	for k := 1; k < kmax; k++ {
		step := 1 << (k - 1)
		st[k] = make([]int, n+1)
		for i := 1; i+(1<<k)-1 <= n; i++ {
			st[k][i] = gcd(st[k-1][i], st[k-1][i+step])
		}
	}
	query := func(l, r int) int {
		length := r - l + 1
		k := logs[length]
		j := r - (1 << k) + 1
		return gcd(st[k][l], st[k][j])
	}
	maxd := 0
	lvals := make([]int, 0, 16)
	for j := 1; j <= n; j++ {
		aj := a[j-1]
		low, high := 1, j
		for low < high {
			mid := (low + high) >> 1
			if query(mid, j) == aj {
				high = mid
			} else {
				low = mid + 1
			}
		}
		lj := low
		low, high = j, n+1
		for low < high {
			mid := (low + high) >> 1
			if mid <= n && query(j, mid) == aj {
				low = mid + 1
			} else {
				high = mid
			}
		}
		rj := low - 1
		d := rj - lj
		if d > maxd {
			maxd = d
			lvals = lvals[:0]
			lvals = append(lvals, lj)
		} else if d == maxd {
			lvals = append(lvals, lj)
		}
	}
	sort.Ints(lvals)
	uniq := make([]int, 0, len(lvals))
	for i, v := range lvals {
		if i == 0 || v != lvals[i-1] {
			uniq = append(uniq, v)
		}
	}
	return uniq, maxd
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	lvals, d := solveCase(n, a)
	var exp strings.Builder
	exp.WriteString(fmt.Sprintf("%d %d\n", len(lvals), d))
	for i, v := range lvals {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprintf("%d", v))
	}
	exp.WriteByte('\n')
	return sb.String(), exp.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
