package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func sumLen(n int64) int64 {
	var res int64
	var pow10 int64 = 1
	var d int64 = 1
	for pow10 <= n {
		nxt := pow10*10 - 1
		if nxt > n {
			nxt = n
		}
		res += (nxt - pow10 + 1) * d
		pow10 *= 10
		d++
	}
	return res
}

func sumXLen(n int64) int64 {
	var res int64
	var pow10 int64 = 1
	var d int64 = 1
	for pow10 <= n {
		nxt := pow10*10 - 1
		if nxt > n {
			nxt = n
		}
		cnt := nxt - pow10 + 1
		res += d * (pow10 + nxt) * cnt / 2
		pow10 *= 10
		d++
	}
	return res
}

func prefix(n int64) int64 {
	if n <= 0 {
		return 0
	}
	return (n+1)*sumLen(n) - sumXLen(n)
}

func digitInRange(n, k int64) int {
	var pow10 int64 = 1
	var d int64 = 1
	for {
		nxt := pow10*10 - 1
		if nxt > n {
			nxt = n
		}
		cnt := nxt - pow10 + 1
		total := cnt * d
		if k > total {
			k -= total
		} else {
			idx := (k - 1) / d
			num := pow10 + idx
			pos := (k - 1) % d
			str := strconv.FormatInt(num, 10)
			return int(str[pos] - '0')
		}
		if nxt == n {
			break
		}
		pow10 *= 10
		d++
	}
	return 0
}

func query(k int64) int {
	l, r := int64(1), int64(1e9)
	for l < r {
		m := (l + r) / 2
		if prefix(m) >= k {
			r = m
		} else {
			l = m + 1
		}
	}
	n := l
	k -= prefix(n - 1)
	return digitInRange(n, k)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var k int64
		fmt.Fscan(reader, &k)
		ans := query(k)
		fmt.Fprintln(writer, ans)
	}
}
