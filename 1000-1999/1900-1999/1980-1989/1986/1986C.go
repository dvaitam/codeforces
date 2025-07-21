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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var s string
		fmt.Fscan(reader, &s)
		inds := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &inds[i])
		}
		var c string
		fmt.Fscan(reader, &c)
		letters := []byte(c)
		sort.Ints(inds)
		sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })
		bs := []byte(s)
		uniq := make([]int, 0)
		for _, v := range inds {
			if len(uniq) == 0 || uniq[len(uniq)-1] != v {
				uniq = append(uniq, v)
			}
		}
		for i := 0; i < len(uniq); i++ {
			idx := uniq[i] - 1
			if i < len(letters) {
				bs[idx] = letters[i]
			}
		}
		fmt.Fprintln(writer, string(bs))
	}
}
