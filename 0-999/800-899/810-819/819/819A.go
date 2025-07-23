package main

import (
	"bufio"
	"fmt"
	"os"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz")

func build(a, b int, ch byte) ([]byte, []byte) {
	seq := make([]byte, a)
	for i := 0; i < a; i++ {
		seq[i] = alphabet[i]
	}
	seen := map[string]int{}
	cycle := 0
	var start, period int
	for {
		for i := 0; i < b; i++ {
			seq = append(seq, ch)
		}
		suffix := seq[len(seq)-a:]
		used := make(map[byte]bool)
		for _, c := range suffix {
			used[c] = true
		}
		t := make([]byte, 0, a)
		for i := 0; i < 26 && len(t) < a; i++ {
			if !used[alphabet[i]] {
				t = append(t, alphabet[i])
			}
		}
		seq = append(seq, t...)
		suffix = seq[len(seq)-a:]
		cycle++
		key := string(suffix)
		if val, ok := seen[key]; ok {
			start = val
			period = cycle - val
			break
		}
		seen[key] = cycle
	}
	prefixLen := a + (start-1)*(a+b)
	periodLen := period * (a + b)
	prefix := append([]byte{}, seq[:prefixLen]...)
	periodStr := append([]byte{}, seq[prefixLen:prefixLen+periodLen]...)
	return prefix, periodStr
}

func uniqueCount(prefix, period []byte, l, r int64) int {
	ans := map[byte]struct{}{}
	prefixLen := int64(len(prefix))
	periodLen := int64(len(period))
	if l <= prefixLen {
		end := r
		if end > prefixLen {
			end = prefixLen
		}
		for i := l - 1; i < end; i++ {
			ans[prefix[i]] = struct{}{}
		}
		if r <= prefixLen {
			return len(ans)
		}
		l = prefixLen + 1
	}
	if periodLen == 0 {
		return len(ans)
	}
	if r-l+1 >= periodLen {
		for _, c := range period {
			ans[c] = struct{}{}
		}
	} else {
		start := (l - prefixLen - 1) % periodLen
		for i := int64(0); i < r-l+1; i++ {
			ans[period[(start+i)%periodLen]] = struct{}{}
		}
	}
	return len(ans)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b int
	var l, r int64
	fmt.Fscan(reader, &a, &b, &l, &r)

	best := 27
	for i := 0; i < 26; i++ {
		prefix, period := build(a, b, alphabet[i])
		val := uniqueCount(prefix, period, l, r)
		if val < best {
			best = val
		}
	}
	fmt.Println(best)
}
