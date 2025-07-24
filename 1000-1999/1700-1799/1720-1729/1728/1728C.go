package main

import (
	"bufio"
	"fmt"
	"os"
)

func digits(x int) int {
	count := 0
	for x > 0 {
		count++
		x /= 10
	}
	return count
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		countA := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			countA[x]++
		}
		countB := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			countB[x]++
		}

		// cancel identical numbers
		for val, cntA := range countA {
			if cntB, ok := countB[val]; ok {
				if cntA < cntB {
					countB[val] -= cntA
					delete(countA, val)
				} else if cntB < cntA {
					countA[val] -= cntB
					delete(countB, val)
				} else { // equal
					delete(countA, val)
					delete(countB, val)
				}
			}
		}

		ops := 0
		// convert numbers >=10
		newA := make(map[int]int)
		for val, cnt := range countA {
			if val >= 10 {
				d := digits(val)
				newA[d] += cnt
				ops += cnt
			} else {
				newA[val] += cnt
			}
		}
		countA = newA

		newB := make(map[int]int)
		for val, cnt := range countB {
			if val >= 10 {
				d := digits(val)
				newB[d] += cnt
				ops += cnt
			} else {
				newB[val] += cnt
			}
		}
		countB = newB

		// cancel again after first conversion
		for val, cntA := range countA {
			if cntB, ok := countB[val]; ok {
				if cntA < cntB {
					countB[val] -= cntA
					delete(countA, val)
				} else if cntB < cntA {
					countA[val] -= cntB
					delete(countB, val)
				} else {
					delete(countA, val)
					delete(countB, val)
				}
			}
		}

		// convert remaining numbers >1 to 1
		for val, cnt := range countA {
			if val > 1 {
				ops += cnt
				countA[val] -= cnt
				countA[1] += cnt
				delete(countA, val)
			}
		}
		for val, cnt := range countB {
			if val > 1 {
				ops += cnt
				countB[val] -= cnt
				countB[1] += cnt
				delete(countB, val)
			}
		}

		fmt.Fprintln(writer, ops)
	}
}
