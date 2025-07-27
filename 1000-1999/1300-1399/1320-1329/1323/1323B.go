package main

import (
	"bufio"
	"fmt"
	"os"
)

func getFreq(arr []int) []int {
	n := len(arr)
	freq := make([]int, n+1)
	count := 0
	for _, v := range arr {
		if v == 1 {
			count++
		} else if count > 0 {
			freq[count]++
			count = 0
		}
	}
	if count > 0 {
		freq[count]++
	}
	return freq
}

func buildPrefix(freq []int) ([]int64, []int64) {
	n := len(freq) - 1
	cnt := make([]int64, n+2)
	sum := make([]int64, n+2)
	for i := n; i >= 1; i-- {
		cnt[i] = cnt[i+1] + int64(freq[i])
		sum[i] = sum[i+1] + int64(freq[i]*i)
	}
	return cnt, sum
}

func val(cnt, sum []int64, length int) int64 {
	if length >= len(cnt) || length <= 0 {
		return 0
	}
	return sum[length] - int64(length-1)*cnt[length]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &b[i])
	}

	freqA := getFreq(a)
	freqB := getFreq(b)
	cntA, sumA := buildPrefix(freqA)
	cntB, sumB := buildPrefix(freqB)

	var ans int64
	for p := int64(1); p*p <= k; p++ {
		if k%p != 0 {
			continue
		}
		q := k / p
		if p <= int64(n) && q <= int64(m) {
			ans += val(cntA, sumA, int(p)) * val(cntB, sumB, int(q))
		}
		if p != q {
			if q <= int64(n) && p <= int64(m) {
				ans += val(cntA, sumA, int(q)) * val(cntB, sumB, int(p))
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
