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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var q int
		fmt.Fscan(reader, &q)
		// counts for s and t
		cntS := make([]int64, 26)
		cntT := make([]int64, 26)
		cntS[0] = 1
		cntT[0] = 1
		lenS := int64(1)
		lenT := int64(1)
		for ; q > 0; q-- {
			var d, k int
			var x string
			fmt.Fscan(reader, &d, &k, &x)
			freq := make([]int64, 26)
			for i := 0; i < len(x); i++ {
				freq[x[i]-'a']++
			}
			if d == 1 {
				for i := 0; i < 26; i++ {
					if freq[i] > 0 {
						cntS[i] += int64(k) * freq[i]
					}
				}
				lenS += int64(k) * int64(len(x))
			} else {
				for i := 0; i < 26; i++ {
					if freq[i] > 0 {
						cntT[i] += int64(k) * freq[i]
					}
				}
				lenT += int64(k) * int64(len(x))
			}

			// check
			hasTBig := false
			for i := 1; i < 26; i++ {
				if cntT[i] > 0 {
					hasTBig = true
					break
				}
			}
			if hasTBig {
				fmt.Fprintln(writer, "YES")
				continue
			}
			// t has only 'a'
			hasSBig := false
			for i := 1; i < 26; i++ {
				if cntS[i] > 0 {
					hasSBig = true
					break
				}
			}
			if hasSBig {
				fmt.Fprintln(writer, "NO")
				continue
			}
			if lenS < lenT {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
