package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	N = 100010
	B = 350
)

var (
	n, q int
	a    [N]int
	las  [N]int
	cnt  [N]int
)

func get(i int) {
	if i+a[i] >= n || i/B != (i+a[i])/B {
		las[i] = i
		cnt[i] = 0
	} else {
		las[i] = las[i+a[i]]
		cnt[i] = cnt[i+a[i]] + 1
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	scanInt := func() int {
		if scanner.Scan() {
			val, _ := strconv.Atoi(scanner.Text())
			return val
		}
		return 0
	}

	n = scanInt()
	q = scanInt()

	for i := 0; i < n; i++ {
		a[i] = scanInt()
	}

	for i := n - 1; i >= 0; i-- {
		get(i)
	}

	for ; q > 0; q-- {
		o := scanInt()
		if o == 1 {
			p := scanInt()
			p-- // 0-based index
			dwn := 0
			s := 0
			for ; p < n; {
				dwn = las[p]
				s += cnt[p] + 1
				p = dwn + a[dwn]
			}
			fmt.Fprintf(writer, "%d %d\n", dwn+1, s)
		} else {
			p := scanInt()
			v := scanInt()
			p-- // 0-based index
			a[p] = v
			for i := p; i >= 0 && i/B == p/B; i-- {
				get(i)
			}
		}
	}
}
