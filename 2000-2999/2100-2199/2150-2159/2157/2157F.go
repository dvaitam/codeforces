package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const mm int64 = 67

	var n int64
	fmt.Fscan(in, &n)

	sqF := make([]int64, 0, 2000000)
	sqS := make([]int64, 0, 2000000)

	var zs int64 = 0

	for ml := int64(1); ml < n; ml *= mm {
		for i := mm - 1; i > 0; i-- {
			for j := ml * i; j < n; j += ml * mm {
				zs++
				sqF = append(sqF, n-j)
				sqS = append(sqS, ml)
			}
		}
	}

	fmt.Fprintln(out, zs)
	for i := int64(0); i < zs; i++ {
		fmt.Fprintln(out, sqF[i], sqS[i])
	}
}

