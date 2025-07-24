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

	var n, m, a, d int64
	if _, err := fmt.Fscan(reader, &n, &m, &a, &d); err != nil {
		return
	}
	clients := make([]int64, m)
	for i := int64(0); i < m; i++ {
		fmt.Fscan(reader, &clients[i])
	}
	empIdx := int64(1)
	cliIdx := int64(0)
	opens := int64(0)

	for empIdx <= n || cliIdx < m {
		nextEmp := int64(1<<63 - 1)
		if empIdx <= n {
			nextEmp = empIdx * a
		}
		nextCli := int64(1<<63 - 1)
		if cliIdx < m {
			nextCli = clients[cliIdx]
		}
		next := nextEmp
		if nextCli < next {
			next = nextCli
		}
		opens++
		closeTime := next + d
		if empIdx <= n {
			jump := closeTime/a + 1
			if jump > n+1 {
				jump = n + 1
			}
			if jump > empIdx {
				empIdx = jump
			}
		}
		for cliIdx < m && clients[cliIdx] <= closeTime {
			cliIdx++
		}
	}
	fmt.Fprintln(writer, opens)
}
