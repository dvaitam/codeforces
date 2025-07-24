package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var b, pCost, mCost int64
	fmt.Fscan(reader, &n, &b, &pCost, &mCost)

	opsType := make([]byte, n)
	opsVal := make([]int64, n)
	for i := 0; i < n; i++ {
		var t string
		fmt.Fscan(reader, &t, &opsVal[i])
		if len(t) > 0 {
			opsType[i] = t[0]
		}
	}

	// Compute original result and product of multipliers >1
	orig := int64(1)
	prodAll := int64(1)
	plusSum := int64(0)
	for i := 0; i < n; i++ {
		if opsType[i] == '+' {
			orig += opsVal[i]
			plusSum += opsVal[i]
		} else {
			orig *= opsVal[i]
			if opsVal[i] > 1 {
				prodAll *= opsVal[i]
			}
		}
	}

	// Precompute suffix product of multipliers >1
	sufProd := make([]int64, n+1)
	sufProd[n] = 1
	for i := n - 1; i >= 0; i-- {
		sufProd[i] = sufProd[i+1]
		if opsType[i] == '*' && opsVal[i] > 1 {
			sufProd[i] *= opsVal[i]
		}
	}

	// Compute minimal cost to move all plus before multipliers
	prefMult := int64(0)
	minCost := int64(1 << 60)
	plusSuffix := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		plusSuffix[i] = plusSuffix[i+1]
		if opsType[i] == '+' {
			plusSuffix[i]++
		}
	}
	for i := 0; i <= n; i++ {
		// cost if split before position i
		multBefore := prefMult
		plusAfter := int64(plusSuffix[i])
		cost := multBefore*mCost + plusAfter*pCost
		if cost < minCost {
			minCost = cost
		}
		if i < n && opsType[i] == '*' && opsVal[i] > 1 {
			prefMult++
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if minCost <= b {
		// We can move all plus before all multipliers
		result := (1 + plusSum) * prodAll
		fmt.Fprintln(writer, result)
		return
	}

	// Greedy move of plus operations to front
	profits := make([]int64, 0)
	for i := 0; i < n; i++ {
		if opsType[i] == '+' {
			gain := opsVal[i] * (prodAll - sufProd[i])
			if gain > 0 {
				profits = append(profits, gain)
			}
		}
	}
	sort.Slice(profits, func(i, j int) bool { return profits[i] > profits[j] })
	moves := int64(len(profits))
	if moves*pCost > b {
		moves = b / pCost
	}
	sumGain := int64(0)
	for i := int64(0); i < moves; i++ {
		sumGain += profits[i]
	}
	result := orig + sumGain
	fmt.Fprintln(writer, result)
}
