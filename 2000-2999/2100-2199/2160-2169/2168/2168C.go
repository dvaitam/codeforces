package main

import (
	"bufio"
	"fmt"
	"os"
)

var columnValues = [20]int{
	1, 2, 4, 8, 16,
	3, 5, 6, 7, 9,
	10, 11, 12, 13, 14,
	15, 17, 18, 19, 20,
}

var parityPositions = []int{0, 1, 2, 3, 4}
var dataPositions = []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
var columnIndex map[int]int

func init() {
	columnIndex = make(map[int]int, len(columnValues))
	for i, val := range columnValues {
		columnIndex[val] = i
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var runType string
	fmt.Fscan(in, &runType)

	if runType == "first" {
		runFirst(in, out)
	} else {
		runSecond(in, out)
	}
}

func runFirst(in *bufio.Reader, out *bufio.Writer) {
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x int
		fmt.Fscan(in, &x)
		codeword := encodeValue(x - 1)
		count := 0
		for i := 0; i < 20; i++ {
			if codeword[i] {
				count++
			}
		}
		fmt.Fprintln(out, count)
		if count > 0 {
			first := true
			for i := 0; i < 20; i++ {
				if codeword[i] {
					if !first {
						fmt.Fprint(out, " ")
					}
					fmt.Fprint(out, i+1)
					first = false
				}
			}
			fmt.Fprintln(out)
		}
	}
}

func runSecond(in *bufio.Reader, out *bufio.Writer) {
	var t int
	fmt.Fscan(in, &t)
	var bits [20]bool
	for ; t > 0; t-- {
		for i := range bits {
			bits[i] = false
		}
		var n int
		fmt.Fscan(in, &n)
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(in, &v)
			bits[v-1] = true
		}
		value := decodeValue(&bits)
		fmt.Fprintln(out, value+1)
	}
}

func encodeValue(value int) [20]bool {
	var codeword [20]bool
	syndrome := 0
	for bitIndex, pos := range dataPositions {
		if ((value >> bitIndex) & 1) == 1 {
			codeword[pos] = true
			syndrome ^= columnValues[pos]
		}
	}
	for i, pos := range parityPositions {
		codeword[pos] = ((syndrome >> i) & 1) == 1
	}
	return codeword
}

func decodeValue(bits *[20]bool) int {
	syndrome := 0
	for i := 0; i < 20; i++ {
		if bits[i] {
			syndrome ^= columnValues[i]
		}
	}
	if syndrome != 0 {
		if idx, ok := columnIndex[syndrome]; ok {
			bits[idx] = !bits[idx]
		}
	}
	value := 0
	for bitIndex, pos := range dataPositions {
		if bits[pos] {
			value |= 1 << bitIndex
		}
	}
	return value
}
