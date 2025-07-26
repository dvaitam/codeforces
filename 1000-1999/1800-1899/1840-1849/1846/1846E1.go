package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXN = 1000000

func precompute() []bool {
	valid := make([]bool, MAXN+1)
	for k := 2; k*k <= MAXN; k++ { // if k^2 > MAXN, sum already > MAXN
		sum := 1 + k + k*k
		power := k * k
		for sum <= MAXN {
			valid[sum] = true
			power *= k
			if power > MAXN { // avoid overflow but also ensures next iteration sum>MAXN
				break
			}
			sum += power
		}
	}
	return valid
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	valid := precompute()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		if n >= 0 && n <= MAXN && valid[n] {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
