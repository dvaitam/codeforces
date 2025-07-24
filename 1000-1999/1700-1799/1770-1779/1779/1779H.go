package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Player struct {
	val int64
	idx int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	vals := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &vals[i])
		total += vals[i]
	}
	players := make([]Player, n)
	for i := 0; i < n; i++ {
		players[i] = Player{vals[i], i}
	}
	sort.Slice(players, func(i, j int) bool { return players[i].val > players[j].val })
	var pref int64
	k := 0
	half := total / 2
	for k < n {
		pref += players[k].val
		k++
		if pref >= half {
			break
		}
	}
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = '0'
	}
	for i := 0; i < k; i++ {
		res[players[i].idx] = '1'
	}
	fmt.Fprintln(writer, string(res))
}
