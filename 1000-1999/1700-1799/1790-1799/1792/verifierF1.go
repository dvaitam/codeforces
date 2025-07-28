package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

const mod int64 = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func solveF(n int) int64 {
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = modPow(fact[n], mod-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	comb := func(a, b int) int64 {
		if b < 0 || b > a {
			return 0
		}
		return fact[a] * invFact[b] % mod * invFact[a-b] % mod
	}
	h := make([]int64, n+1)
	bArr := make([]int64, n+1)
	h[1] = 1
	bArr[1] = 1
	for i := 2; i <= n; i++ {
		var val int64
		for s := 1; s < i; s++ {
			k := i - s
			bk := h[k]
			if k != 1 {
				bk = bk * 2 % mod
			}
			val = (val + comb(i-1, s-1)*h[s]%mod*bk) % mod
		}
		h[i] = val % mod
		bArr[i] = 2 * h[i] % mod
	}
	ans := bArr[n] - 2
	if ans < 0 {
		ans += mod
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/binary\n", os.Args[0])
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(47)
	const t = 100
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 3
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		expected := solveF(n)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		outBytes, err := cmd.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, "binary execution failed on test", i+1, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(bytes.NewReader(outBytes))
		scanner.Split(bufio.ScanWords)
		if !scanner.Scan() {
			fmt.Printf("no output on test %d\n", i+1)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			fmt.Printf("invalid output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("mismatch on test %d: expected %d got %d\n", i+1, expected, got)
			os.Exit(1)
		}
		if scanner.Scan() {
			fmt.Printf("extra output on test %d\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
