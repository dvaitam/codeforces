package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

type edge struct {
	to     uint8
	vision []byte
}

type state struct {
	cur         uint8
	pendingType int8
	pending     string
}

func nextInt(r *bufio.Reader) int {
	sign, val := 1, 0
	c, err := r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	n := nextInt(reader)
	m := nextInt(reader)
	if n == 0 {
		return
	}

	graph := make([][]edge, n+1)
	for i := 0; i < m; i++ {
		x := nextInt(reader)
		y := nextInt(reader)
		k := nextInt(reader)
		seq := make([]byte, k)
		for j := 0; j < k; j++ {
			seq[j] = byte(nextInt(reader))
		}
		graph[x] = append(graph[x], edge{to: byte(y), vision: seq})
	}

	limit := 2 * n
	ans := make([]int, limit+1)
	curStates := make(map[state]int)
	for start := 1; start <= n; start++ {
		pending := []byte{byte(start)}
		st := state{cur: byte(start), pendingType: 1, pending: string(pending)}
		curStates[st] = (curStates[st] + 1) % mod
	}

	for length := 1; length < limit && len(curStates) > 0; length++ {
		nextStates := make(map[state]int)
		for st, cnt := range curStates {
			if cnt == 0 {
				continue
			}
			baseQueue := []byte(st.pending)
			queueLen := len(baseQueue)
			lenVision := length
			switch st.pendingType {
			case 1:
				lenVision = length - queueLen
			case -1:
				lenVision = length + queueLen
			}
			if lenVision < 0 || lenVision > limit {
				continue
			}
			for _, e := range graph[int(st.cur)] {
				if lenVision+len(e.vision) > limit {
					continue
				}
				queue := append([]byte(nil), baseQueue...)
				ptype := st.pendingType
				ok := true
				for _, val := range e.vision {
					if ptype == 1 {
						if len(queue) == 0 || queue[0] != val {
							ok = false
							break
						}
						queue = queue[1:]
						if len(queue) == 0 {
							ptype = 0
						}
					} else {
						queue = append(queue, val)
						ptype = -1
					}
				}
				if !ok {
					continue
				}
				dest := e.to
				if ptype == -1 {
					if len(queue) == 0 || queue[0] != dest {
						continue
					}
					queue = queue[1:]
					if len(queue) == 0 {
						ptype = 0
					}
				} else {
					queue = append(queue, dest)
					ptype = 1
				}
				newState := state{cur: dest, pendingType: ptype, pending: string(queue)}
				newCnt := nextStates[newState] + cnt
				if newCnt >= mod {
					newCnt -= mod
				}
				nextStates[newState] = newCnt
				if ptype == 0 {
					ans[length+1] += cnt
					if ans[length+1] >= mod {
						ans[length+1] -= mod
					}
				}
			}
		}
		curStates = nextStates
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 1; i <= limit; i++ {
		fmt.Fprintln(writer, ans[i]%mod)
	}
}
