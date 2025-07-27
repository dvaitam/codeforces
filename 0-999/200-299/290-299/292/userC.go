package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var (
	digits   []byte
	needMask int
	results  []string
	seen     = make(map[string]struct{})
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	digits = make([]byte, n)
	for i := 0; i < n; i++ {
		var d int
		fmt.Fscan(in, &d)
		digits[i] = byte('0' + d)
		needMask |= 1 << d
	}

	for totalLen := 4; totalLen <= 12; totalLen++ {
		halfLen := (totalLen + 1) / 2
		if halfLen < n { // impossible to use all digits
			continue
		}
		buf := make([]byte, halfLen)
		dfs(0, halfLen, buf, 0, totalLen)
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, len(results))
	for _, addr := range results {
		fmt.Fprintln(out, addr)
	}
	out.Flush()
}

func dfs(pos, halfLen int, buf []byte, mask, totalLen int) {
	if pos == halfLen {
		if mask != needMask {
			return
		}
		full := make([]byte, totalLen)
		copy(full, buf)
		for i := 0; i < totalLen/2; i++ {
			full[totalLen-1-i] = buf[i]
		}
		process(full)
		return
	}
	for _, d := range digits {
		buf[pos] = d
		dfs(pos+1, halfLen, buf, mask|(1<<(d-'0')), totalLen)
	}
}

func process(s []byte) {
	L := len(s)
	for l1 := 1; l1 <= 3; l1++ {
		for l2 := 1; l2 <= 3; l2++ {
			for l3 := 1; l3 <= 3; l3++ {
				l4 := L - l1 - l2 - l3
				if l4 < 1 || l4 > 3 {
					continue
				}
				idx := 0
				v1, ok := segmentVal(s, idx, l1)
				if !ok {
					continue
				}
				idx += l1
				v2, ok := segmentVal(s, idx, l2)
				if !ok {
					continue
				}
				idx += l2
				v3, ok := segmentVal(s, idx, l3)
				if !ok {
					continue
				}
				idx += l3
				v4, ok := segmentVal(s, idx, l4)
				if !ok {
					continue
				}

				addr := strconv.Itoa(v1) + "." + strconv.Itoa(v2) + "." + strconv.Itoa(v3) + "." + strconv.Itoa(v4)
				if _, exists := seen[addr]; !exists {
					seen[addr] = struct{}{}
					results = append(results, addr)
				}
			}
		}
	}
}

func segmentVal(s []byte, start, length int) (int, bool) {
	if length > 1 && s[start] == '0' {
		return 0, false
	}
	val := 0
	for i := 0; i < length; i++ {
		val = val*10 + int(s[start+i]-'0')
		if val > 255 {
			return 0, false
		}
	}
	return val, true
}