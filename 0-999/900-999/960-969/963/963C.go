package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct{ w, h int64 }

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func divisors(n int64) []int64 {
	divs := []int64{}
	for i := int64(1); i*i <= n; i++ {
		if n%i == 0 {
			divs = append(divs, i)
			if i*i != n {
				divs = append(divs, n/i)
			}
		}
	}
	sort.Slice(divs, func(i, j int) bool { return divs[i] < divs[j] })
	return divs
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)

	rect := make(map[pair]int64, n)
	wSum := make(map[int64]int64)
	hSum := make(map[int64]int64)
	widths := make([]int64, 0)
	heights := make([]int64, 0)
	widthSeen := make(map[int64]bool)
	heightSeen := make(map[int64]bool)
	var total int64

	for i := 0; i < n; i++ {
		var w, h, c int64
		fmt.Fscan(reader, &w, &h, &c)
		rect[pair{w, h}] = c
		wSum[w] += c
		hSum[h] += c
		total += c
		if !widthSeen[w] {
			widthSeen[w] = true
			widths = append(widths, w)
		}
		if !heightSeen[h] {
			heightSeen[h] = true
			heights = append(heights, h)
		}
	}

	if int64(len(widths))*int64(len(heights)) != int64(len(rect)) {
		fmt.Println(0)
		return
	}

	// choose base width and height
	w0 := widths[0]
	h0 := heights[0]
	cross, ok := rect[pair{w0, h0}]
	if !ok {
		fmt.Println(0)
		return
	}

	var rowGCD int64
	for _, w := range widths {
		val, ok := rect[pair{w, h0}]
		if !ok {
			fmt.Println(0)
			return
		}
		rowGCD = gcd(rowGCD, val)
	}
	var colGCD int64
	for _, h := range heights {
		val, ok := rect[pair{w0, h}]
		if !ok {
			fmt.Println(0)
			return
		}
		colGCD = gcd(colGCD, val)
	}

	divs := divisors(cross)
	var ans int
	for _, d := range divs {
		if rowGCD%d != 0 {
			continue
		}
		c0 := cross / d
		if colGCD%c0 != 0 {
			continue
		}

		rowCount := make(map[int64]int64)
		colCount := make(map[int64]int64)
		valid := true
		for _, h := range heights {
			val := rect[pair{w0, h}]
			if val%c0 != 0 {
				valid = false
				break
			}
			rowCount[h] = val / c0
		}
		if !valid {
			continue
		}
		for _, w := range widths {
			val := rect[pair{w, h0}]
			if val%d != 0 {
				valid = false
				break
			}
			colCount[w] = val / d
		}
		if !valid {
			continue
		}
		for _, w := range widths {
			for _, h := range heights {
				val := rect[pair{w, h}]
				if val != colCount[w]*rowCount[h] {
					valid = false
					break
				}
			}
			if !valid {
				break
			}
		}
		if valid {
			ans++
		}
	}

	fmt.Println(ans)
}
