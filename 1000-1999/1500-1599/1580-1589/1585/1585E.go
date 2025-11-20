package main

import (
	"bufio"
	"os"
)

const MAXN = 1000005
const MAXQ = 1000005

var (
	gHead, gNext, gTo  [MAXN]int
	ecnt               int
	qHead              [MAXN]int
	qNext, qL, qK, qID [MAXQ]int
	qCnt               int
	val                [MAXN]int
	cnt, pos, bit      [MAXN]int
	lists              [MAXN][]int
	ans                [MAXQ]int
	reader             *bufio.Reader
	writer             *bufio.Writer
	n, q               int
)

func readInt() int {
	res := 0
	c, err := reader.ReadByte()
	for err == nil && c <= ' ' {
		c, err = reader.ReadByte()
	}
	if err != nil {
		return 0
	}
	for c >= '0' && c <= '9' {
		res = res*10 + int(c-'0')
		c, err = reader.ReadByte()
	}
	return res
}

func writeInt(x int) {
	if x < 0 {
		writer.WriteByte('-')
		writer.WriteByte('1')
		return
	}
	if x == 0 {
		writer.WriteByte('0')
		return
	}
	var b [20]byte
	p := 0
	for x > 0 {
		b[p] = byte(x%10 + '0')
		x /= 10
		p++
	}
	for p > 0 {
		p--
		writer.WriteByte(b[p])
	}
}

func addEdge(u, v int) {
	ecnt++
	gTo[ecnt] = v
	gNext[ecnt] = gHead[u]
	gHead[u] = ecnt
}

func addQuery(u, l, k, id int) {
	qCnt++
	qL[qCnt] = l
	qK[qCnt] = k
	qID[qCnt] = id
	qNext[qCnt] = qHead[u]
	qHead[u] = qCnt
}

func update(idx, v int) {
	for ; idx <= n; idx += idx & -idx {
		bit[idx] += v
	}
}

func queryBIT(idx int) int {
	res := 0
	for ; idx > 0; idx -= idx & -idx {
		res += bit[idx]
	}
	return res
}

func addVal(u int) {
	x := val[u]
	c := cnt[x]
	if c > 0 {
		update(c, -1)
		idx := pos[x]
		last := lists[c][len(lists[c])-1]
		lists[c][idx] = last
		pos[last] = idx
		lists[c] = lists[c][:len(lists[c])-1]
	}
	cnt[x]++
	c = cnt[x]
	update(c, 1)
	lists[c] = append(lists[c], x)
	pos[x] = len(lists[c]) - 1
}

func removeVal(u int) {
	x := val[u]
	c := cnt[x]
	update(c, -1)
	idx := pos[x]
	last := lists[c][len(lists[c])-1]
	lists[c][idx] = last
	pos[last] = idx
	lists[c] = lists[c][:len(lists[c])-1]

	cnt[x]--
	c = cnt[x]
	if c > 0 {
		update(c, 1)
		lists[c] = append(lists[c], x)
		pos[x] = len(lists[c]) - 1
	}
}

func solveQuery(l, k int) int {
	total := queryBIT(n) - queryBIT(l-1)
	if total < k {
		return -1
	}
	target := queryBIT(l-1) + k
	idx := 0
	current := 0
	for i := 19; i >= 0; i-- {
		nextIdx := idx + (1 << i)
		if nextIdx <= n && current+bit[nextIdx] < target {
			idx = nextIdx
			current += bit[nextIdx]
		}
	}
	freq := idx + 1
	if len(lists[freq]) == 0 {
		return -1
	}
	return lists[freq][0]
}

type StackFrame struct {
	u, e int
}

func main() {
	reader = bufio.NewReaderSize(os.Stdin, 1<<20)
	writer = bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	t := readInt()
	for i := 0; i < t; i++ {
		n = readInt()
		q = readInt()

		ecnt = 0
		qCnt = 0
		for j := 0; j <= n; j++ {
			gHead[j] = 0
			qHead[j] = 0
			cnt[j] = 0
			bit[j] = 0
			lists[j] = lists[j][:0]
		}

		for j := 1; j <= n; j++ {
			val[j] = readInt()
		}

		for j := 2; j <= n; j++ {
			p := readInt()
			addEdge(p, j)
		}

		for j := 0; j < q; j++ {
			v := readInt()
			l := readInt()
			k := readInt()
			addQuery(v, l, k, j)
		}

		stk := make([]StackFrame, 0, n)
		stk = append(stk, StackFrame{1, gHead[1]})
		addVal(1)

		idx := qHead[1]
		for idx != 0 {
			ans[qID[idx]] = solveQuery(qL[idx], qK[idx])
			idx = qNext[idx]
		}

		for len(stk) > 0 {
			top := &stk[len(stk)-1]
			if top.e != 0 {
				v := gTo[top.e]
				top.e = gNext[top.e]

				addVal(v)
				qIdx := qHead[v]
				for qIdx != 0 {
					ans[qID[qIdx]] = solveQuery(qL[qIdx], qK[qIdx])
					qIdx = qNext[qIdx]
				}
				stk = append(stk, StackFrame{v, gHead[v]})
			} else {
				removeVal(top.u)
				stk = stk[:len(stk)-1]
			}
		}

		for j := 0; j < q; j++ {
			if j > 0 {
				writer.WriteByte(' ')
			}
			writeInt(ans[j])
		}
		writer.WriteByte('\n')
	}
}
