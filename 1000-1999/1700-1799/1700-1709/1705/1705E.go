package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxVal = 200000 + 100

var digits [maxVal]int
var top int

func add(idx int, delta int) {
	digits[idx] += delta
	for i := idx; i < maxVal-1; i++ {
		if digits[i] < 0 {
			digits[i] += 2
			digits[i+1] -= 1
		} else if digits[i] >= 2 {
			carry := digits[i] / 2
			digits[i] %= 2
			digits[i+1] += carry
		} else if digits[i+1] >= 0 && digits[i+1] <= 1 {
			break
		}
	}
	if top < idx {
		top = idx
	}
	for top+1 < maxVal && digits[top+1] > 0 {
		top++
	}
	for top > 0 && digits[top] == 0 {
		top--
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		add(a[i], 1)
	}

	for i := 0; i < q; i++ {
		var k, l int
		fmt.Fscan(reader, &k, &l)
		k--
		add(a[k], -1)
		add(l, 1)
		a[k] = l
		fmt.Fprintln(writer, top)
	}
}
