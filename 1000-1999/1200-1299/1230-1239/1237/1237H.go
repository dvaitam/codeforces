package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var TC int
	if _, err := fmt.Fscan(reader, &TC); err != nil {
		return
	}
	for tc := 0; tc < TC; tc++ {
		var s, tstr string
		fmt.Fscan(reader, &s, &tstr)
		n2 := len(s)
		N := n2 / 2
		cnt := [4]int{}
		A0 := make([]int, N)
		B0 := make([]int, N)
		for i := 0; i < N; i++ {
			A0[i] = int(s[2*i]-'0')*2 + int(s[2*i+1]-'0')
			B0[i] = int(tstr[2*i]-'0')*2 + int(tstr[2*i+1]-'0')
			cnt[A0[i]]++
			cnt[B0[i]]--
		}
		if cnt[0] != 0 || cnt[3] != 0 || cnt[1]+cnt[2] != 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		var ops []int
		var A, B []int
		for {
			ops = ops[:0]
			// reset arrays
			A = make([]int, N)
			B = make([]int, N)
			copy(A, A0)
			copy(B, B0)
			r := N
			// define doit inline
			doit := func(k int) {
				ops = append(ops, k)
				// reverse A[0:k]
				for i, j := 0, k-1; i < j; i, j = i+1, j-1 {
					A[i], A[j] = A[j], A[i]
				}
				// complement 1<->2
				for i := 0; i < k; i++ {
					if A[i] == 1 {
						A[i] = 2
					} else if A[i] == 2 {
						A[i] = 1
					}
				}
			}
			for r > 0 {
				if A[r-1] == B[r-1] {
					r--
					continue
				}
				if A[0] == B[r-1] {
					if A[0] == 1 || A[0] == 2 {
						doit(1)
					}
					doit(r)
					continue
				}
				if (A[0] == 1 && B[r-1] == 2) || (A[0] == 2 && B[r-1] == 1) {
					doit(r)
					continue
				}
				// search patterns
				var fn []int
				if r >= 2 {
					for i := 1; i < r; i++ {
						if A[i-1] == B[r-2] && A[i] == B[r-1] {
							fn = append(fn, i)
						}
					}
				}
				if len(fn) > 0 {
					i := fn[rand.Intn(len(fn))]
					doit(i + 1)
					doit(r)
					continue
				}
				fn = fn[:0]
				for i := 0; i < r; i++ {
					if A[i] == B[r-1] {
						fn = append(fn, i)
					}
				}
				if len(fn) > 0 {
					i := fn[rand.Intn(len(fn))]
					doit(i + 1)
					doit(r)
					continue
				}
				fn = fn[:0]
				for i := 0; i < r; i++ {
					if A[i] == 3-B[r-1] {
						fn = append(fn, i)
					}
				}
				if len(fn) > 0 {
					i := fn[rand.Intn(len(fn))]
					doit(i + 1)
					doit(1)
					doit(r)
					continue
				}
			}
			if len(ops) <= 2*N+1 {
				break
			}
		}
		// output
		m := len(ops)
		fmt.Fprintln(writer, m)
		if m > 0 {
			for i, v := range ops {
				if i > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, v*2)
			}
			writer.WriteByte('\n')
		} else {
			writer.WriteByte('\n')
		}
	}
}
