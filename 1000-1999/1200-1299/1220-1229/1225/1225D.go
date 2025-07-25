package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type pair struct{ p, e int }

func factor(x, k int, spf []int) (string, string) {
	if x == 1 {
		return "", ""
	}
	pairs := make([]pair, 0)
	compPairs := make([]pair, 0)
	for x > 1 {
		p := spf[x]
		if p == 0 {
			p = x
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		cnt %= k
		if cnt != 0 {
			pairs = append(pairs, pair{p, cnt})
			cp := (k - cnt) % k
			if cp != 0 {
				compPairs = append(compPairs, pair{p, cp})
			}
		}
	}
	return buildStr(pairs), buildStr(compPairs)
}

func buildStr(ps []pair) string {
	if len(ps) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, pr := range ps {
		sb.WriteString(strconv.Itoa(pr.p))
		sb.WriteByte('#')
		sb.WriteString(strconv.Itoa(pr.e))
		sb.WriteByte('#')
	}
	return sb.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	const maxA = 100000
	spf := make([]int, maxA+1)
	for i := 2; i <= maxA; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxA; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	count := make(map[string]int)
	var res int64
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		norm, comp := factor(x, k, spf)
		if v, ok := count[comp]; ok {
			res += int64(v)
		}
		count[norm]++
	}
	fmt.Fprintln(writer, res)
}
