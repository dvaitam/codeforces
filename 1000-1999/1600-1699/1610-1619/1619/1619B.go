package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func intSqrt(n int64) int64 {
	x := int64(math.Sqrt(float64(n)))
	for (x+1)*(x+1) <= n {
		x++
	}
	for x*x > n {
		x--
	}
	return x
}

func intCbrt(n int64) int64 {
	x := int64(math.Cbrt(float64(n)))
	for (x+1)*(x+1)*(x+1) <= n {
		x++
	}
	for x*x*x > n {
		x--
	}
	return x
}

func intSixthRoot(n int64) int64 {
	x := int64(math.Pow(float64(n), 1.0/6.0))
	pow6 := func(v int64) int64 {
		res := int64(1)
		for i := 0; i < 6; i++ {
			res *= v
		}
		return res
	}
	for pow6(x+1) <= n {
		x++
	}
	for pow6(x) > n {
		x--
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(reader, &n)
		sq := intSqrt(n)
		cb := intCbrt(n)
		sixth := intSixthRoot(n)
		fmt.Fprintln(writer, sq+cb-sixth)
	}
}
