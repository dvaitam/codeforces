package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, p int64
	if _, err := fmt.Fscan(reader, &n, &p); err != nil {
		return
	}

	for k := int64(1); k <= 60; k++ {
		s := n - p*k
		if s <= 0 {
			continue
		}
		if s < k {
			continue
		}
		if int64(bits.OnesCount64(uint64(s))) <= k {
			fmt.Fprintln(writer, k)
			return
		}
	}
	fmt.Fprintln(writer, -1)
}
