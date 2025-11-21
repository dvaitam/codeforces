package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextString() string {
	b, err := fs.r.ReadByte()
	for err == nil && b <= ' ' {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return ""
	}
	buf := []byte{b}
	for {
		b, err = fs.r.ReadByte()
		if err != nil || b <= ' ' {
			break
		}
		buf = append(buf, b)
	}
	return string(buf)
}

type fenwick struct {
	bit []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{bit: make([]int, n+2)}
}

func (f *fenwick) add(idx, delta int) {
	idx++
	for idx < len(f.bit) {
		f.bit[idx] += delta
		idx += idx & -idx
	}
}

func (f *fenwick) sum(idx int) int {
	idx++
	if idx <= 0 {
		return 0
	}
	res := 0
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

type action struct {
	typ int
	pos int
}

func charIndex(c byte) int {
	switch c {
	case 'Y':
		return 0
	case 'D':
		return 1
	default:
		return 2
	}
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	tStr := fs.nextString()
	t, _ := strconv.Atoi(tStr)
	typeToChar := []byte{'Y', 'D', 'X'}

	for ; t > 0; t-- {
		s := []byte(fs.nextString())
		n := len(s)
		if n%3 != 0 {
			out.WriteString("NO\n")
			continue
		}
		if !isValidBase(s) {
			out.WriteString("NO\n")
			continue
		}

		left := make([]int, n)
		right := make([]int, n)
		for i := 0; i < n; i++ {
			left[i] = i - 1
			if i+1 < n {
				right[i] = i + 1
			} else {
				right[i] = -1
			}
		}

		alive := make([]bool, n)
		safe := make([]bool, n)
		for i := range alive {
			alive[i] = true
		}

		queues := make([][]int, 3)
		heads := make([]int, 3)
		for i := 0; i < n; i++ {
			val := computeSafe(i, left, right, alive, s)
			safe[i] = val
			if val {
				queues[charIndex(s[i])] = append(queues[charIndex(s[i])], i)
			}
		}

		ft := newFenwick(n)
		for i := 0; i < n; i++ {
			ft.add(i, 1)
		}

		totalOps := n / 3
		ops := make([][]action, 0, totalOps)
		valid := true

		for op := 0; op < totalOps && valid; op++ {
			step := make([]action, 0, 3)
			used := [3]bool{}
			for rem := 0; rem < 3; rem++ {
				best := int(1e9)
				selected := -1
				selectedIdx := -1
				nextHead := 0
				for tType := 0; tType < 3; tType++ {
					if used[tType] {
						continue
					}
					head := heads[tType]
					q := queues[tType]
					for head < len(q) {
						idx := q[head]
						if alive[idx] && safe[idx] {
							break
						}
						head++
					}
					heads[tType] = head
					avail := len(q) - head
					if avail == 0 {
						continue
					}
					if avail < best {
						best = avail
						selected = tType
						selectedIdx = q[head]
						nextHead = head + 1
					}
				}
				if selected == -1 {
					valid = false
					break
				}
				heads[selected] = nextHead
				pos := ft.sum(selectedIdx - 1)
				step = append(step, action{typ: selected, pos: pos})
				alive[selectedIdx] = false
				safe[selectedIdx] = false
				ft.add(selectedIdx, -1)

				L := left[selectedIdx]
				R := right[selectedIdx]
				if L != -1 {
					right[L] = R
					newSafe := computeSafe(L, left, right, alive, s)
					safe[L] = newSafe
					if newSafe {
						queues[charIndex(s[L])] = append(queues[charIndex(s[L])], L)
					}
				}
				if R != -1 {
					left[R] = L
					newSafe := computeSafe(R, left, right, alive, s)
					safe[R] = newSafe
					if newSafe {
						queues[charIndex(s[R])] = append(queues[charIndex(s[R])], R)
					}
				}
				used[selected] = true
			}
			if !valid {
				break
			}
			ops = append(ops, step)
		}

		if !valid {
			out.WriteString("NO\n")
			continue
		}

		out.WriteString("YES\n")
		out.Write(s)
		out.WriteByte('\n')
		for op := len(ops) - 1; op >= 0; op-- {
			step := ops[op]
			var sb strings.Builder
			for i := len(step) - 1; i >= 0; i-- {
				if sb.Len() > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteByte(typeToChar[step[i].typ])
				sb.WriteByte(' ')
				sb.WriteString(strconv.Itoa(step[i].pos))
			}
			sb.WriteByte('\n')
			out.WriteString(sb.String())
		}
	}
}

func computeSafe(i int, left, right []int, alive []bool, s []byte) bool {
	if i == -1 || !alive[i] {
		return false
	}
	l := left[i]
	r := right[i]
	if l == -1 || r == -1 {
		return true
	}
	return s[l] != s[r]
}

func isValidBase(s []byte) bool {
	n := len(s)
	if n == 0 {
		return true
	}
	cnt := [3]int{}
	for i := 0; i < n; i++ {
		switch s[i] {
		case 'Y':
			cnt[0]++
		case 'D':
			cnt[1]++
		case 'X':
			cnt[2]++
		default:
			return false
		}
		if i+1 < n && s[i] == s[i+1] {
			return false
		}
	}
	if cnt[0] != cnt[1] || cnt[1] != cnt[2] {
		return false
	}
	return true
}
