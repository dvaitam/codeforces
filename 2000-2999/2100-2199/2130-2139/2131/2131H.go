package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxV = 1000000

var spf [maxV + 1]int
var mu [maxV + 1]int
var freq = make([]int, maxV+1)
var cntDiv = make([]int, maxV+1)

func init() {
	buildSieve()
}

func buildSieve() {
	primes := make([]int, 0)
	spf[1] = 1
	mu[1] = 1
	for i := 2; i <= maxV; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if p > spf[i] || i*p > maxV {
				break
			}
			spf[i*p] = p
			if i%p == 0 {
				mu[i*p] = 0
				break
			}
			mu[i*p] = -mu[i]
		}
	}
}

type Item struct {
	val int
	idx int
}

func getDivisors(x int) []int {
	if x == 1 {
		return []int{1}
	}
	divs := []int{1}
	tmp := x
	for tmp > 1 {
		p := spf[tmp]
		cnt := 0
		for tmp%p == 0 {
			tmp /= p
			cnt++
		}
		curLen := len(divs)
		mult := 1
		for i := 0; i < cnt; i++ {
			mult *= p
			for j := 0; j < curLen; j++ {
				divs = append(divs, divs[j]*mult)
			}
		}
	}
	return divs
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func findPair(items []Item, limit int) (Item, Item, bool) {
	if len(items) < 2 || limit == 0 {
		return Item{}, Item{}, false
	}
	touched := make([]int, 0, len(items))
	for _, it := range items {
		if freq[it.val] == 0 {
			touched = append(touched, it.val)
		}
		freq[it.val]++
	}

	for d := 1; d <= limit; d++ {
		cnt := 0
		for j := d; j <= limit; j += d {
			cnt += freq[j]
		}
		cntDiv[d] = cnt
	}

	hasPartner := make(map[int]bool, len(touched))
	for _, val := range touched {
		divs := getDivisors(val)
		sum := 0
		for _, d := range divs {
			sum += mu[d] * cntDiv[d]
		}
		fmt.Println("debug", val, sum, len(items))
		hasPartner[val] = sum > 0
	}

	var first, second Item
	found := false
	for i := 0; i < len(items) && !found; i++ {
		if !hasPartner[items[i].val] {
			continue
		}
		for j := 0; j < len(items); j++ {
			if i == j {
				continue
			}
			if gcd(items[i].val, items[j].val) == 1 {
				first = items[i]
				second = items[j]
				found = true
				break
			}
		}
	}

	for _, val := range touched {
		freq[val] = 0
	}

	if !found {
		return Item{}, Item{}, false
	}
	return first, second, true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		items := make([]Item, n)
		maxVal := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &items[i].val)
			items[i].idx = i + 1
			if items[i].val > maxVal {
				maxVal = items[i].val
			}
		}

		firstA, firstB, ok := findPair(items, maxVal)
		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}

		remain := make([]Item, 0, len(items)-2)
		maxVal2 := 0
		for _, it := range items {
			if it.idx == firstA.idx || it.idx == firstB.idx {
				continue
			}
			remain = append(remain, it)
			if it.val > maxVal2 {
				maxVal2 = it.val
			}
		}

		secondA, secondB, ok := findPair(remain, maxVal2)
		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}

		fmt.Fprintf(out, "%d %d %d %d\n", firstA.idx, firstB.idx, secondA.idx, secondB.idx)
	}
}
