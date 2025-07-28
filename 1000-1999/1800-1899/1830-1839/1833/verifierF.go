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

const MOD int = 1000000007

func modPow(a, e int) int {
	res := 1
	x := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(x) % int64(MOD))
		}
		x = int(int64(x) * int64(x) % int64(MOD))
		e >>= 1
	}
	return res
}

func modInv(a int) int { return modPow(a, MOD-2) }

func solveCase(arr []int, m int) string {
	freq := make(map[int]int)
	for _, v := range arr {
		freq[v]++
	}
	unique := make([]int, 0, len(freq))
	for v := range freq {
		unique = append(unique, v)
	}
	sort.Ints(unique)
	invCache := make(map[int]int)
	getInv := func(x int) int {
		if val, ok := invCache[x]; ok {
			return val
		}
		val := modInv(x)
		invCache[x] = val
		return val
	}
	ans := 0
	prod := 1
	left := 0
	for right := 0; right < len(unique); right++ {
		cntR := freq[unique[right]]
		prod = int(int64(prod) * int64(cntR) % int64(MOD))
		for unique[right]-unique[left] >= m {
			cntL := freq[unique[left]]
			prod = int(int64(prod) * int64(getInv(cntL)) % int64(MOD))
			left++
		}
		for right-left+1 > m {
			cntL := freq[unique[left]]
			prod = int(int64(prod) * int64(getInv(cntL)) % int64(MOD))
			left++
		}
		if right-left+1 == m && unique[right]-unique[left] < m {
			ans += prod
			if ans >= MOD {
				ans -= MOD
			}
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	m := rng.Intn(n) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(20)
	}
	input := fmt.Sprintf("1\n%d %d\n", n, m)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", arr[i])
	}
	input += "\n"
	return input, solveCase(arr, m)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
