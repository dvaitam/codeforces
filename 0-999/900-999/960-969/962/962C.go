package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func isSubsequence(target, source string) bool {
	j := 0
	for i := 0; i < len(source) && j < len(target); i++ {
		if source[i] == target[j] {
			j++
		}
	}
	return j == len(target)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var nStr string
	if _, err := fmt.Fscan(reader, &nStr); err != nil {
		return
	}

	n, _ := strconv.Atoi(nStr)
	limit := int(math.Sqrt(float64(n)))
	best := len(nStr) + 1

	for i := 1; i <= limit; i++ {
		sq := i * i
		sqStr := strconv.Itoa(sq)
		if isSubsequence(sqStr, nStr) {
			ops := len(nStr) - len(sqStr)
			if ops < best {
				best = ops
			}
		}
	}

	if best == len(nStr)+1 {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, best)
	}
}
