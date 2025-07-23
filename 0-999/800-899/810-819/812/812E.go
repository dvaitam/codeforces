package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	children := make([][]int, n)
	for i := 1; i < n; i++ {
		var p int
		fmt.Fscan(reader, &p)
		p--
		children[p] = append(children[p], i)
	}

	depth := make([]int, n)
	queue := []int{0}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, ch := range children[v] {
			depth[ch] = depth[v] + 1
			queue = append(queue, ch)
		}
	}

	parity := make([]int, n)
	countEven, countOdd := 0, 0
	freqEven := make(map[int]int)
	freqOdd := make(map[int]int)
	for i := 0; i < n; i++ {
		parity[i] = depth[i] & 1
		if parity[i] == 0 {
			countEven++
			freqEven[a[i]]++
		} else {
			countOdd++
			freqOdd[a[i]]++
		}
	}

	leafParity := 0
	for i := 0; i < n; i++ {
		if len(children[i]) == 0 {
			leafParity = parity[i]
			break
		}
	}

	xorVal := 0
	for i := 0; i < n; i++ {
		if parity[i] == leafParity {
			xorVal ^= a[i]
		}
	}

	var result int64
	if xorVal == 0 {
		result += int64(countEven*(countEven-1)/2 + countOdd*(countOdd-1)/2)
		for val, c0 := range freqEven {
			if c1, ok := freqOdd[val]; ok {
				result += int64(c0) * int64(c1)
			}
		}
	} else {
		for val, c0 := range freqEven {
			if c1, ok := freqOdd[val^xorVal]; ok {
				result += int64(c0) * int64(c1)
			}
		}
	}

	fmt.Fprintln(writer, result)
}
