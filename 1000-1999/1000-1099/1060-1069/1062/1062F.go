package main

import (
	"fmt"
	"io"
	"math/bits"
	"os"
)

const WORDS = 64
const CHUNK = WORDS * 64

func readInt(b []byte, i int) (int, int) {
	for i < len(b) && (b[i] < '0' || b[i] > '9') {
		i++
	}
	if i >= len(b) {
		return 0, i
	}
	res := 0
	for i < len(b) && b[i] >= '0' && b[i] <= '9' {
		res = res*10 + int(b[i]-'0')
		i++
	}
	return res, i
}

func main() {
	input, _ := io.ReadAll(os.Stdin)
	offset := 0

	var n, m int
	n, offset = readInt(input, offset)
	m, offset = readInt(input, offset)

	if n == 0 {
		return
	}

	head := make([]int, n+1)
	from := make([]int, m)
	to := make([]int, m)

	for i := 0; i < m; i++ {
		from[i], offset = readInt(input, offset)
		to[i], offset = readInt(input, offset)
		from[i]--
		to[i]--
		head[from[i]+1]++
	}

	for i := 1; i <= n; i++ {
		head[i] += head[i-1]
	}

	edges := make([]int, m)
	curHead := make([]int, n)
	copy(curHead, head)

	inDegree := make([]int, n)
	for i := 0; i < m; i++ {
		u := from[i]
		v := to[i]
		edges[curHead[u]] = v
		curHead[u]++
		inDegree[v]++
	}

	topo := make([]int, 0, n)
	queue := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		topo = append(topo, u)
		for e := head[u]; e < head[u+1]; e++ {
			v := edges[e]
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	totalA := make([]int, n)
	totalR := make([]int, n)

	buffer := make([]uint64, n*WORDS)

	for start := 0; start < n; start += CHUNK {
		end := start + CHUNK
		if end > n {
			end = n
		}

		for i := 0; i < n*WORDS; i++ {
			buffer[i] = 0
		}

		for i := start; i < end; i++ {
			bitIdx := i - start
			buffer[i*WORDS+bitIdx/64] |= 1 << (bitIdx % 64)
		}

		for _, u := range topo {
			uIdx := u * WORDS
			bufU := buffer[uIdx : uIdx+WORDS]
			for e := head[u]; e < head[u+1]; e++ {
				v := edges[e]
				vIdx := v * WORDS
				bufV := buffer[vIdx : vIdx+WORDS]
				for w := 0; w < WORDS; w++ {
					bufV[w] |= bufU[w]
				}
			}
		}

		for i := 0; i < n; i++ {
			idx := i * WORDS
			buf := buffer[idx : idx+WORDS]
			cnt := 0
			for w := 0; w < WORDS; w++ {
				cnt += bits.OnesCount64(buf[w])
			}
			totalA[i] += cnt
		}
	}

	for start := 0; start < n; start += CHUNK {
		end := start + CHUNK
		if end > n {
			end = n
		}

		for i := 0; i < n*WORDS; i++ {
			buffer[i] = 0
		}

		for i := start; i < end; i++ {
			bitIdx := i - start
			buffer[i*WORDS+bitIdx/64] |= 1 << (bitIdx % 64)
		}

		for i := n - 1; i >= 0; i-- {
			u := topo[i]
			uIdx := u * WORDS
			bufU := buffer[uIdx : uIdx+WORDS]
			for e := head[u]; e < head[u+1]; e++ {
				v := edges[e]
				vIdx := v * WORDS
				bufV := buffer[vIdx : vIdx+WORDS]
				for w := 0; w < WORDS; w++ {
					bufU[w] |= bufV[w]
				}
			}
		}

		for i := 0; i < n; i++ {
			idx := i * WORDS
			buf := buffer[idx : idx+WORDS]
			cnt := 0
			for w := 0; w < WORDS; w++ {
				cnt += bits.OnesCount64(buf[w])
			}
			totalR[i] += cnt
		}
	}

	ans := 0
	for i := 0; i < n; i++ {
		if totalA[i]+totalR[i] >= n {
			ans++
		}
	}

	fmt.Println(ans)
}
