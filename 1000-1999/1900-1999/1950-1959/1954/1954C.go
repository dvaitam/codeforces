package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	bigTen := big.NewInt(10)
	for ; t > 0; t-- {
		var xStr, yStr string
		fmt.Fscan(in, &xStr)
		fmt.Fscan(in, &yStr)
		n := len(xStr)
		diff := big.NewInt(0)
		xb := make([]byte, n)
		yb := make([]byte, n)
		for i := 0; i < n; i++ {
			da := int64(xStr[i] - '0')
			db := int64(yStr[i] - '0')

			diff1 := new(big.Int).Set(diff)
			diff1.Mul(diff1, bigTen)
			diff1.Add(diff1, big.NewInt(da-db))

			diff2 := new(big.Int).Set(diff)
			diff2.Mul(diff2, bigTen)
			diff2.Add(diff2, big.NewInt(db-da))

			abs1 := new(big.Int).Abs(new(big.Int).Set(diff1))
			abs2 := new(big.Int).Abs(new(big.Int).Set(diff2))

			if abs2.Cmp(abs1) < 0 {
				diff = diff2
				xb[i] = byte('0' + db)
				yb[i] = byte('0' + da)
			} else {
				diff = diff1
				xb[i] = byte('0' + da)
				yb[i] = byte('0' + db)
			}
		}
		fmt.Fprintln(out, string(xb))
		fmt.Fprintln(out, string(yb))
	}
}
