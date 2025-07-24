package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1e9 + 7

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var x int
		var s string
		fmt.Fscan(reader, &x)
		fmt.Fscan(reader, &s)
		arr := []byte(s)
		lenMod := int64(len(arr))
		for i := 0; i < x; i++ {
			c := int(arr[i] - '0')
			diff := (lenMod - int64(i+1)) % mod
			if diff < 0 {
				diff += mod
			}
			lenMod = (lenMod + diff*int64(c-1)) % mod
			if len(arr) < x {
				suffix := make([]byte, len(arr)-i-1)
				copy(suffix, arr[i+1:])
				for j := 0; j < c-1 && len(arr) < x; j++ {
					need := x - len(arr)
					if need >= len(suffix) {
						arr = append(arr, suffix...)
					} else {
						arr = append(arr, suffix[:need]...)
					}
				}
			}
		}
		fmt.Fprintln(writer, lenMod%mod)
	}
}
