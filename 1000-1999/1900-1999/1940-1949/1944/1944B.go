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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		k *= 2

		a := make([]int, 2*n)
		for i := 0; i < 2*n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		cnt := make([]int, n+1)
		for i := 0; i < n; i++ {
			if a[i] >= 1 && a[i] <= n {
				cnt[a[i]]++
			}
		}

		var d0, d1, d2 []int
		for i := 1; i <= n; i++ {
			switch cnt[i] {
			case 0:
				d0 = append(d0, i)
			case 1:
				d1 = append(d1, i)
			default:
				d2 = append(d2, i)
			}
		}

		temp := 0
		for _, x := range d2 {
			if temp >= k {
				break
			}
			temp += 2
			fmt.Fprintf(writer, "%d %d ", x, x)
		}
		for _, x := range d1 {
			if temp >= k {
				break
			}
			temp++
			fmt.Fprintf(writer, "%d ", x)
		}
		fmt.Fprint(writer, "\n")

		temp = 0
		for _, x := range d0 {
			if temp >= k {
				break
			}
			temp += 2
			fmt.Fprintf(writer, "%d %d ", x, x)
		}
		for _, x := range d1 {
			if temp >= k {
				break
			}
			temp++
			fmt.Fprintf(writer, "%d ", x)
		}
		fmt.Fprint(writer, "\n")
	}
}
