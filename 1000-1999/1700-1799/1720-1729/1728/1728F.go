package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	ones := []int{}
	others := []int{}
	for _, v := range arr {
		if v == 1 {
			ones = append(ones, v)
		} else {
			others = append(others, v)
		}
	}
	sort.Ints(others)

	// count occurrences of other values
	cnt := make(map[int]int)
	for _, v := range others {
		cnt[v]++
	}
	uniq := []int{}
	for v := range cnt {
		uniq = append(uniq, v)
	}
	sort.Ints(uniq)

	sequence := make([]int, 0, n)
	sequence = append(sequence, ones...)
	for _, v := range uniq {
		sequence = append(sequence, v)
	}
	extras := []int{}
	for _, v := range uniq {
		for i := 1; i < cnt[v]; i++ {
			extras = append(extras, v)
		}
	}
	sort.Ints(extras)
	sequence = append(sequence, extras...)

	prev := 0
	total := 0
	for _, v := range sequence {
		next := prev + v - (prev % v)
		if next <= prev {
			next += v
		}
		total += next
		prev = next
	}
	fmt.Fprintln(out, total)
}
