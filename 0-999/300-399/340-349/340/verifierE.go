package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod = 1000000007

func modpow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func expectedAnswer(a []int) int {
	n := len(a)
	used := make([]bool, n+1)
	for _, v := range a {
		if v != -1 {
			used[v] = true
		}
	}
	m, k := 0, 0
	for i := 1; i <= n; i++ {
		if a[i-1] == -1 {
			m++
			if !used[i] {
				k++
			}
		}
	}
	fact := make([]int, n+1)
	invfact := make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * i % mod
	}
	invfact[n] = modpow(fact[n], mod-2)
	for i := n; i > 0; i-- {
		invfact[i-1] = invfact[i] * i % mod
	}
	ans := 0
	for i := 0; i <= k; i++ {
		comb := fact[k] * invfact[i] % mod * invfact[k-i] % mod
		ways := fact[m-i]
		term := comb * ways % mod
		if i%2 == 1 {
			ans = (ans - term + mod) % mod
		} else {
			ans = (ans + term) % mod
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(50) + 2
	perm := rand.Perm(n)
	for i := range perm {
		perm[i]++
	}
	arr := make([]int, n)
	for i, v := range perm {
		arr[i] = v
	}
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			arr[i] = -1
		}
	}
	return arr
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
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
	expected := strconv.Itoa(expectedAnswer(arr))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		arr := generateCase(rng)
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
