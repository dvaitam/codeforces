package main

import (
	"bufio"
	"fmt"
	"math/bits"
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
	const mod = 32768
	for i := 0; i < n; i++ {
		var a int
		fmt.Fscan(reader, &a)
		ans := 15
		for add := 0; add <= 15; add++ {
			val := (a + add) % mod
			if val == 0 {
				if add < ans {
					ans = add
				}
				break
			}
			tz := bits.TrailingZeros(uint(val))
			if tz > 15 {
				tz = 15
			}
			ops := add + 15 - tz
			if ops < ans {
				ans = ops
			}
		}
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans)
	}
	writer.WriteByte('\n')
}
