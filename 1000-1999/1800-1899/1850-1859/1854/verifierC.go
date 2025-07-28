package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 1000000007

func modInverse(a int64) int64 {
	return new(big.Int).ModInverse(big.NewInt(a), big.NewInt(mod)).Int64()
}

type ratCache map[uint32]*big.Rat

func expected(mask uint32, m int, cache ratCache) *big.Rat {
	if mask == 0 {
		return big.NewRat(0, 1)
	}
	if v, ok := cache[mask]; ok {
		return v
	}
	size := 0
	for i := 0; i < m; i++ {
		if mask&(1<<i) != 0 {
			size++
		}
	}
	res := big.NewRat(0, 1)
	for i := 0; i < m; i++ {
		if mask&(1<<i) == 0 {
			continue
		}
		next := mask &^ (1 << i)
		if i+1 < m && (mask&(1<<uint(i+1))) == 0 {
			next |= 1 << uint(i+1)
		}
		res.Add(res, expected(next, m, cache))
	}
	res.Quo(res, big.NewRat(int64(size), 1))
	res.Add(res, big.NewRat(1, 1))
	cache[mask] = res
	return res
}

func solveC(n, m int, arr []int) int64 {
	mask := uint32(0)
	for _, v := range arr {
		mask |= 1 << uint(v-1)
	}
	r := expected(mask, m, make(ratCache))
	num := new(big.Int).Mod(r.Num(), big.NewInt(mod))
	den := new(big.Int).Mod(r.Denom(), big.NewInt(mod))
	inv := new(big.Int).ModInverse(den, big.NewInt(mod))
	if inv == nil {
		return 0
	}
	num.Mul(num, inv)
	num.Mod(num, big.NewInt(mod))
	return num.Int64()
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(3))
	const tests = 100
	for t := 0; t < tests; t++ {
		m := 2 + r.Intn(4) // up to 5
		n := 1 + r.Intn(m)
		vals := make([]int, 0, n)
		used := map[int]bool{}
		for len(vals) < n {
			x := 1 + r.Intn(m)
			if !used[x] {
				used[x] = true
				vals = append(vals, x)
			}
		}
		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				if vals[i] > vals[j] {
					vals[i], vals[j] = vals[j], vals[i]
				}
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, v := range vals {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d error: %v\n", t+1, err)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscanf(out, "%d", &got); err != nil {
			fmt.Printf("test %d invalid output\n", t+1)
			os.Exit(1)
		}
		want := solveC(n, m, vals)
		if got != want {
			fmt.Printf("test %d failed: expected %d got %d\n", t+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
