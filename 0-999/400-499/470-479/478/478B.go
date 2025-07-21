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

	var n, m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	// maximum pairs: one team with n-m+1, rest size 1
	maxSize := n - m + 1
	kmax := maxSize * (maxSize - 1) / 2
	// minimum pairs: teams as equal as possible
	q := n / m
	r := n % m
	// r teams of size q+1, m-r teams of size q
	// C(x,2) = x*(x-1)/2
	cBig := (q + 1) * q / 2
	cSmall := q * (q - 1) / 2
	kmin := r*cBig + (m-r)*cSmall
	fmt.Fprintf(writer, "%d %d", kmin, kmax)
}
