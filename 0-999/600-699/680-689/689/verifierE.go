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

const mod int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func prepare(n int) ([]int64, []int64) {
	fact := make([]int64, n+1)
	inv := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	inv[n] = modPow(fact[n], mod-2)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
	return fact, inv
}

func C(n, r int64, fact, inv []int64) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * inv[r] % mod * inv[n-r] % mod
}

func solve(n, k int, segs [][2]int64) int64 {
	events := make(map[int64]int)
	coords := make([]int64, 0, 2*n)
	for i := 0; i < n; i++ {
		l, r := segs[i][0], segs[i][1]
		events[l]++
		events[r+1]--
	}
	for x := range events {
		coords = append(coords, x)
	}
	sort.Slice(coords, func(i, j int) bool { return coords[i] < coords[j] })
	fact, inv := prepare(n)
	var ans int64
	active := 0
	var prev int64
	for i, x := range coords {
		if i > 0 {
			length := x - prev
			if length > 0 {
				ans = (ans + C(int64(active), int64(k), fact, inv)*length) % mod
			}
		}
		active += events[x]
		prev = x
	}
	return ans % mod
}

func runCase(bin string, n, k int, segs [][2]int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", segs[i][0], segs[i][1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := solve(n, k, segs)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) (int, int, [][2]int64) {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	segs := make([][2]int64, n)
	for i := 0; i < n; i++ {
		l := rng.Int63n(100) - 50
		r := l + rng.Int63n(50)
		segs[i] = [2]int64{l, r}
	}
	return n, k, segs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, segs := randomCase(rng)
		if err := runCase(bin, n, k, segs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
