package main

import (
	"bufio"
	"fmt"
	"os"
)

const baseStr = "What are you doing at the end of the world? Are you busy? Will you save us?"
const prefix = "What are you doing while sending \""
const middle = "\"? Are you busy? Will you send \""
const suffix = "\"?"

const maxN = 100000
const inf int64 = 1e18

var length [maxN + 1]int64

func init() {
	length[0] = int64(len(baseStr))
	for i := 1; i <= maxN; i++ {
		l := int64(len(prefix)) + int64(len(middle)) + int64(len(suffix)) + 2*length[i-1]
		if l > inf {
			l = inf
		}
		length[i] = l
	}
}

func solve(n int, k int64) byte {
	for {
		if n == 0 {
			if k > int64(len(baseStr)) {
				return '.'
			}
			return baseStr[k-1]
		}
		if k > length[n] {
			return '.'
		}
		if k <= int64(len(prefix)) {
			return prefix[k-1]
		}
		k -= int64(len(prefix))
		if k <= length[n-1] {
			n--
			continue
		}
		k -= length[n-1]
		if k <= int64(len(middle)) {
			return middle[k-1]
		}
		k -= int64(len(middle))
		if k <= length[n-1] {
			n--
			continue
		}
		k -= length[n-1]
		if k <= int64(len(suffix)) {
			return suffix[k-1]
		}
		return '.'
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	res := make([]byte, q)
	for i := 0; i < q; i++ {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		res[i] = solve(n, k)
	}
	fmt.Fprintln(writer, string(res))
}
