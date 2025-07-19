package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for {
		var N int
		if _, err := fmt.Fscan(reader, &N); err != nil {
			return
		}
		A := make([]int, N)
		for i := 0; i < N; i++ {
			fmt.Fscan(reader, &A[i])
			A[i]--
		}
		// Prepare
		vis := make([]bool, N)
		small := make([][][]int, 8)
		ans := make([][]int, 0, N)
		sep := make([]bool, 0, N)
		ansLen := make([]int, 0, N)
		// Decompose cycles
		for u := 0; u < N; u++ {
			if vis[u] {
				continue
			}
			cyc := make([]int, 0)
			for v := u; !vis[v]; v = A[v] {
				vis[v] = true
				cyc = append(cyc, v)
			}
			clen := len(cyc)
			// Split long cycles
			i := 0
			for ; i+5 <= clen; i += 4 {
				ansLen = append(ansLen, 5)
				sep = append(sep, false)
				tmp := make([]int, 5)
				for j := 0; j < 4; j++ {
					tmp[j] = cyc[i+j]
				}
				tmp[4] = cyc[clen-1]
				ans = append(ans, tmp)
			}
			// Remainder
			x := clen - i
			if x == 1 {
				continue
			}
			row := make([]int, x)
			for j := 0; j < x; j++ {
				row[j] = cyc[clen-x+j]
			}
			small[x] = append(small[x], row)
		}
		// Combine small cycles
		m2, m3, m4 := len(small[2]), len(small[3]), len(small[4])
		sho := m2 + m3
		km := 0
		for k := 0; k <= m3; k += 3 {
			m2t := m2
			m3t := m3 - k
			tmpCnt := (k / 3) * 2
			tmp23 := m2t
			if m3t < tmp23 {
				tmp23 = m3t
			}
			tmpCnt += tmp23
			m2t -= tmp23
			m3t -= tmp23
			tmpCnt += (m2t + 1) / 2
			tmpCnt += m3t
			if sho > tmpCnt {
				sho = tmpCnt
				km = k
			}
		}
		// Use selection
		ii := 0
		i := 0
		// process groups of three 3-cycles
		for ; i < km; i += 3 {
			// first op
			ansLen = append(ansLen, 5)
			sep = append(sep, true)
			t1 := small[3][i]
			t2 := small[3][i+1]
			t3 := small[3][i+2]
			ans = append(ans, []int{t1[0], t1[2], t2[0], t2[1], t2[2]})
			// second op
			ansLen = append(ansLen, 5)
			sep = append(sep, true)
			ans = append(ans, []int{t1[1], t1[2], t3[0], t3[1], t3[2]})
		}
		// pair 2-cycles with 3-cycles
		for ; ii < m2 && i < m3; ii, i = ii+1, i+1 {
			ansLen = append(ansLen, 5)
			sep = append(sep, true)
			a2 := small[2][ii]
			a3 := small[3][i]
			ans = append(ans, []int{a2[0], a2[1], a3[0], a3[1], a3[2]})
		}
		// remaining 2-cycles
		for ; ii < m2; ii += 2 {
			// one or two 2-cycles
			length := 2
			if ii+1 < m2 {
				length = 4
			}
			ansLen = append(ansLen, length)
			sep = append(sep, true)
			if length == 4 {
				a := small[2][ii]
				b := small[2][ii+1]
				ans = append(ans, []int{a[0], a[1], b[0], b[1]})
			} else {
				a := small[2][ii]
				ans = append(ans, []int{a[0], a[1]})
			}
		}
		// remaining 3-cycles
		for ; i < m3; i++ {
			ansLen = append(ansLen, 3)
			sep = append(sep, false)
			c := small[3][i]
			ans = append(ans, []int{c[0], c[1], c[2]})
		}
		// 4-cycles
		for idx := 0; idx < m4; idx++ {
			ansLen = append(ansLen, 4)
			sep = append(sep, false)
			c := small[4][idx]
			ans = append(ans, []int{c[0], c[1], c[2], c[3]})
		}
		// Output
		M := len(ansLen)
		fmt.Fprintln(writer, M)
		for t := 0; t < M; t++ {
			L := ansLen[t]
			fmt.Fprintln(writer, L)
			for j := 0; j < L; j++ {
				if j > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, ans[t][j]+1)
			}
			fmt.Fprintln(writer)
			for j := 0; j < L; j++ {
				if j > 0 {
					writer.WriteByte(' ')
				}
				jj := j + 1
				if sep[t] {
					if jj == 2 {
						jj = 0
					} else if jj == L {
						jj = 2
					}
				} else {
					if jj == L {
						jj = 0
					}
				}
				fmt.Fprint(writer, ans[t][jj]+1)
			}
			fmt.Fprintln(writer)
		}
	}
}
