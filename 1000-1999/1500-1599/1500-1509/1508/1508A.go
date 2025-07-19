package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		var n int
		fmt.Fscan(reader, &n)
		var sx, sy, sz string
		fmt.Fscan(reader, &sx, &sy, &sz)

		c1x := 0
		for i := 0; i < 2*n; i++ {
			if sx[i] == '1' {
				c1x++
			}
		}
		c0x := 2*n - c1x
		c1y := 0
		for i := 0; i < 2*n; i++ {
			if sy[i] == '1' {
				c1y++
			}
		}
		c0y := 2*n - c1y
		c1z := 0
		for i := 0; i < 2*n; i++ {
			if sz[i] == '1' {
				c1z++
			}
		}
		c0z := 2*n - c1z

		var aStr, bStr string
		need := byte('0')
		if c1x >= n && c1y >= n {
			aStr, bStr, need = sx, sy, '1'
		} else if c1x >= n && c1z >= n {
			aStr, bStr, need = sx, sz, '1'
		} else if c0x >= n && c0y >= n {
			aStr, bStr, need = sx, sy, '0'
		} else if c0x >= n && c0z >= n {
			aStr, bStr, need = sx, sz, '0'
		} else if c1y >= n && c1z >= n {
			aStr, bStr, need = sy, sz, '1'
		} else {
			aStr, bStr, need = sy, sz, '0'
		}

		ans := make([]byte, 0, 3*n)
		j := 0
		for i := 0; i < 2*n; i++ {
			if aStr[i] == need {
				for j < 2*n && bStr[j] != need {
					ans = append(ans, bStr[j])
					j++
				}
				if j < 2*n {
					j++
				}
			}
			ans = append(ans, aStr[i])
		}
		for j < 2*n {
			ans = append(ans, bStr[j])
			j++
		}

		writer.Write(ans)
		writer.WriteByte('\n')
	}
}
