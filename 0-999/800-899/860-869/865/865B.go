package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Contestant struct {
	s              int64
	lossPerSlice   int64
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func ceilDiv(a, b int64) int64 {
	if a == 0 {
		return 0
	}
	return (a + b - 1) / b
}

func calculateLoss(group []Contestant, deficit int64) int64 {
	if deficit <= 0 {
		return 0
	}
	sort.Slice(group, func(i, j int) bool { return group[i].lossPerSlice < group[j].lossPerSlice })
	totalLoss := int64(0)
	for _, contestant := range group {
		switchCount := min(deficit, contestant.s)
		totalLoss += switchCount * contestant.lossPerSlice
		deficit -= switchCount
		if deficit == 0 {
			break
		}
	}
	return totalLoss
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	parts := strings.Fields(scanner.Text())
	n, _ := strconv.ParseInt(parts[0], 10, 64)
	sPizza, _ := strconv.ParseInt(parts[1], 10, 64)

	totalSlices := int64(0)
	hMax := int64(0)
	groupA := make([]Contestant, 0, n)
	groupB := make([]Contestant, 0, n)
	totalSlicesA := int64(0)
	totalSlicesB := int64(0)

	for i := int64(0); i < n; i++ {
		scanner.Scan()
		parts = strings.Fields(scanner.Text())
		si, _ := strconv.ParseInt(parts[0], 10, 64)
		ai, _ := strconv.ParseInt(parts[1], 10, 64)
		bi, _ := strconv.ParseInt(parts[2], 10, 64)
		totalSlices += si

		if ai > bi {
			hMax += si * ai
			totalSlicesA += si
			groupA = append(groupA, Contestant{s: si, lossPerSlice: ai - bi})
		} else if bi > ai {
			hMax += si * bi
			totalSlicesB += si
			groupB = append(groupB, Contestant{s: si, lossPerSlice: bi - ai})
		} else {
			hMax += si * ai
		}
	}

	p := ceilDiv(totalSlices, sPizza)
	reqP1 := ceilDiv(totalSlicesA, sPizza)
	p1Case1 := min(reqP1, p)
	p2Case1 := p - p1Case1
	deficitB := totalSlicesB - p2Case1*sPizza
	lossB := calculateLoss(groupB, deficitB)
	h1 := hMax - lossB

	reqP2 := ceilDiv(totalSlicesB, sPizza)
	p2Case2 := min(reqP2, p)
	p1Case2 := p - p2Case2
	deficitA := totalSlicesA - p1Case2*sPizza
	lossA := calculateLoss(groupA, deficitA)
	h2 := hMax - lossA

	fmt.Println(max(h1, h2))
}
