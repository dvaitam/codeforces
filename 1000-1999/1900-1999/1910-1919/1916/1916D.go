package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t, n int
	fmt.Fscan(reader, &t)
	// precomputed answers for n <= 13 (odd n only)
	ans := map[int][]string{
		1:  {"9"},
		3:  {"169", "196", "961"},
		5:  {"16384", "31684", "36481", "38416", "43681"},
		7:  {"1493284", "3214849", "3912484", "4239481", "4293184", "4932841", "9132484"},
		9:  {"236759769", "297769536", "369677529", "526977936", "677925369", "769729536", "773562969", "796763529", "927567936"},
		11: {"48458977956", "54785487969", "56487979584", "59988745476", "64755998784", "68597895744", "77956548849", "85745894976", "95887457649", "95896747584", "97598758464"},
		13: {"5898894567696", "6559678848969", "6586598874969", "6599894588676", "6788958546969", "6795885899664", "6965988947856", "8579966588964", "8585696978496", "8594788665969", "8698865979456", "8995896467856", "9579866858496"},
	}

	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &n)
		if n <= 13 {
			if list, ok := ans[n]; ok {
				for _, v := range list {
					writer.WriteString(v)
					writer.WriteByte('\n')
				}
			}
		} else {
			z := n / 2
			found := 0
			// generate n numbers
			for x := 0; x < z && found < n; x++ {
				for y := x + 1; y < z && found < n; y++ {
					if x+z == y*2 {
						continue
					}
					a := make([]int, n)
					a[x*2]++
					a[y*2]++
					a[z*2]++
					a[x+y] += 2
					a[x+z] += 2
					a[y+z] += 2
					var sb strings.Builder
					for idx := n - 1; idx >= 0; idx-- {
						sb.WriteByte(byte('0' + a[idx]))
					}
					writer.WriteString(sb.String())
					writer.WriteByte('\n')
					found++
				}
			}
		}
	}
}
