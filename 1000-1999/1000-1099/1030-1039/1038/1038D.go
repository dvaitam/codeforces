package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader *bufio.Reader
var writer *bufio.Writer

func init() {
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
}

func readInt() int64 {
	var sign int64 = 1
	var b byte
	var err error
	// skip non-numbers
	for {
		b, err = reader.ReadByte()
		if err != nil {
			return 0
		}
		if b == '-' {
			sign = -1
			b, _ = reader.ReadByte()
			break
		}
		if b >= '0' && b <= '9' {
			break
		}
	}
	var x int64
	for {
		if b < '0' || b > '9' {
			break
		}
		x = x*10 + int64(b-'0')
		b, _ = reader.ReadByte()
	}
	return x * sign
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	defer writer.Flush()
	n64 := readInt()
	n := int(n64)
	a := make([]int64, n)
	var ans int64
	for i := 0; i < n; i++ {
		a[i] = readInt()
		ans += abs(a[i])
	}
	if n == 1 {
		fmt.Fprint(writer, a[0])
		return
	}
	allPos := true
	for i := 0; i < n; i++ {
		if a[i] < 0 {
			allPos = false
			break
		}
	}
	if allPos {
		mn := a[0]
		for i := 1; i < n; i++ {
			if a[i] < mn {
				mn = a[i]
			}
		}
		ans -= 2 * mn
	} else {
		allNeg := true
		for i := 0; i < n; i++ {
			if a[i] > 0 {
				allNeg = false
				break
			}
		}
		if allNeg {
			mn := abs(a[0])
			for i := 1; i < n; i++ {
				ai := abs(a[i])
				if ai < mn {
					mn = ai
				}
			}
			ans -= 2 * mn
		}
	}
	fmt.Fprint(writer, ans)
}
