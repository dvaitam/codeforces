package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in  := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	for ; T > 0; T-- {
		// 1) ask for p with the string "aa"
		fmt.Fprintln(out, "? aa")
		out.Flush()

		var p int64
		fmt.Fscan(in, &p)
		p-- // shift to 0‑based like the C++ (--p)

		// 2) ask for x with ten 'z'
		fmt.Fprintln(out, "? zzzzzzzzzz")
		out.Flush()

		var x int64
		fmt.Fscan(in, &x)

		// —— replicate the arithmetic ——
		hs := int64(0)
		y  := x + 1
		var o int64
		an := make([]int64, 11) // extra slot for possible borrow

		for i := 0; i < 10; i++ {
			hs = hs*p + 26
			an[i] = 26 - y%p
			y /= p
		}

		// build the 10‑letter string with carry fix
		sBytes := make([]byte, 10)
		for i := 0; i < 10; i++ {
			if an[i] < 1 {
				an[i] = 26
				an[i+1]-- // borrow
			}
			sBytes[i] = byte('a' + an[i] - 1)
		}
		s := string(sBytes)

		// 3) convert s back to number o (reverse order)
		for i := 9; i >= 0; i-- {
			o = o*p + int64(sBytes[i]-'a'+1)
		}

		// 4) ask for m with that constructed string
		fmt.Fprintf(out, "? %s\n", s)
		out.Flush()

		var m int64
		fmt.Fscan(in, &m)

		// 5) final answer
		ans := hs - x - o + m
		fmt.Fprintf(out, "! %d %d\n", p, ans)
		out.Flush()
	}
}

