package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return
		}
		m := map[int64]int64{0: 1}
		var b, a int64
		for i := int64(0); i < int64(len(s)); i++ {
			if s[i] == '1' {
				b++
			} else {
				b--
			}
			a = (a + (int64(len(s))-i)*m[b]) % MOD
			m[b] = (m[b] + i + 2) % MOD
		}
		fmt.Fprintln(writer, a%MOD)
	}
}
