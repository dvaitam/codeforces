package main

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
)

const MOD int64 = 1000000007

func modPow(a, e, mod int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + arr[i]
		}
		total := prefix[n]
		bestDiff := -1.0
		bestK := 1
		bestLarge := false

		for k := 1; k < n; k++ {
			denom := float64(int64(k) * int64(n-k))
			// first k smallest
			sumSmall := float64(prefix[k])
			numSmall := math.Abs(float64(n)*sumSmall - float64(k)*float64(total))
			diffSmall := numSmall / denom
			if diffSmall > bestDiff {
				bestDiff = diffSmall
				bestK = k
				bestLarge = false
			}
			// k largest
			sumLarge := float64(total - prefix[n-k])
			numLarge := math.Abs(float64(n)*sumLarge - float64(k)*float64(total))
			diffLarge := numLarge / denom
			if diffLarge > bestDiff {
				bestDiff = diffLarge
				bestK = k
				bestLarge = true
			}
		}

		k := bestK
		var subsetSum int64
		if bestLarge {
			subsetSum = total - prefix[n-k]
		} else {
			subsetSum = prefix[k]
		}
		nBig := big.NewInt(int64(n))
		kBig := big.NewInt(int64(k))
		sumBig := big.NewInt(subsetSum)
		totalBig := big.NewInt(total)

		num := new(big.Int).Mul(nBig, sumBig)
		temp := new(big.Int).Mul(kBig, totalBig)
		num.Sub(num, temp)
		num.Abs(num)

		denomBig := big.NewInt(int64(k * (n - k)))

		gcd := new(big.Int).GCD(nil, nil, num, denomBig)
		num.Div(num, gcd)
		denomBig.Div(denomBig, gcd)

		modBig := big.NewInt(MOD)
		numMod := new(big.Int).Mod(num, modBig).Int64()
		denomMod := new(big.Int).Mod(denomBig, modBig).Int64()

		invDenom := modPow(denomMod, MOD-2, MOD)
		ans := numMod * invDenom % MOD
		fmt.Fprintln(out, ans)
	}
}
