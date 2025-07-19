package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Precompute smallest prime factors up to 1000
	const maxN = 1000
	spf := make([]int, maxN+1)
	for i := 0; i <= maxN; i++ {
		spf[i] = i
	}
	for i := 2; i*i <= maxN; i++ {
		for j := i * 2; j <= maxN; j += i {
			if spf[j] > i {
				spf[j] = i
			}
		}
	}

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		v := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &v[i])
		}
		// Assign colors based on smallest prime factor
		mp := make(map[int]int)
		id := 1
		for _, x := range v {
			p := spf[x]
			if _, ok := mp[p]; !ok {
				mp[p] = id
				id++
			}
		}
		// Output number of colors
		writer.WriteString(strconv.Itoa(len(mp)))
		writer.WriteByte('\n')
		// Output colors for each element
		for i, x := range v {
			writer.WriteString(strconv.Itoa(mp[spf[x]]))
			if i+1 < n {
				writer.WriteByte(' ')
			}
		}
		writer.WriteByte('\n')
	}
}
