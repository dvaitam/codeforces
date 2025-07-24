package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)

	spf := make([]int, n+1)
	phi := make([]int, n+1)
	primes := make([]int, 0, n/10)
	phi[1] = 1
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			phi[i] = i - 1
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p*i > n || p > spf[i] {
				break
			}
			spf[p*i] = p
			if i%p == 0 {
				phi[p*i] = phi[i] * p
				break
			} else {
				phi[p*i] = phi[i] * (p - 1)
			}
		}
	}

	var phiSum int64
	for i := 1; i <= n; i++ {
		phiSum += int64(phi[i])
	}
	totalPairs := int64(n*(n-1)) / 2
	edges := totalPairs - (phiSum - 1)

	limit := n / 2
	freq := make(map[int]int)
	primesList := make([]int, 0)
	for i := 2; i <= n; i++ {
		p := spf[i]
		if p > limit {
			continue
		}
		if _, ok := freq[p]; !ok {
			primesList = append(primesList, p)
		}
		freq[p]++
	}

	sort.Ints(primesList)
	m := len(primesList)
	prefix := make([]int64, m+1)
	for i := m - 1; i >= 0; i-- {
		prefix[i] = prefix[i+1] + int64(freq[primesList[i]])
	}

	var count3 int64
	for i, p := range primesList {
		limitP := n / p
		j := sort.SearchInts(primesList, limitP+1)
		if j <= i {
			j = i + 1
		}
		if j < m {
			count3 += int64(freq[p]) * prefix[j]
		}
	}

	var sizeS int64
	for _, c := range freq {
		sizeS += int64(c)
	}
	pairS := sizeS * (sizeS - 1) / 2

	gcd1Pairs := pairS - edges
	count2 := gcd1Pairs - count3

	ans := edges + 2*count2 + 3*count3
	fmt.Fprintln(writer, ans)
}
