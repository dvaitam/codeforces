package main

import (
	"bufio"
	"fmt"
	"os"
)

func countOnes(n, bit int) int {
	cycle := 1 << (bit + 1)
	full := n / cycle
	cnt := full * (1 << bit)
	rem := n % cycle
	if rem > (1 << bit) {
		cnt += rem - (1 << bit)
	}
	return cnt
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	pref := make([]int, n)
	cur := 0
	for i := 1; i < n; i++ {
		var a int
		fmt.Fscan(reader, &a)
		cur ^= a
		pref[i] = cur
	}

	bits := 20
	cntPref := make([]int, bits)
	for _, v := range pref {
		for j := 0; j < bits; j++ {
			if (v>>j)&1 == 1 {
				cntPref[j]++
			}
		}
	}

	cntRange := make([]int, bits)
	for j := 0; j < bits; j++ {
		cntRange[j] = countOnes(n, j)
	}

	key := 0
	for j := 0; j < bits; j++ {
		if cntPref[j] != cntRange[j] {
			key |= 1 << j
		}
	}

	for i, v := range pref {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v^key)
	}
	writer.WriteByte('\n')
}
