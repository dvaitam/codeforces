package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m, k int
	if _, err := fmt.Fscan(reader, &m, &k); err != nil {
		return
	}
	vals := make([][]int64, k)
	S := make([]int64, k)
	Q := make([]int64, k)
	for i := 0; i < k; i++ {
		vals[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &vals[i][j])
			S[i] += vals[i][j]
			Q[i] += vals[i][j] * vals[i][j]
		}
	}

	vsum := (S[k-1] - S[0]) / int64(k-1)
	y := -1
	var delta int64
	for i := 1; i < k-1; i++ {
		expected := S[0] + int64(i)*vsum
		if S[i] != expected {
			y = i
			delta = S[i] - expected
			break
		}
	}
	if y == -1 {
		// should not happen
		return
	}

	idx := make([]int, 0, 3)
	for i := 0; i < k && len(idx) < 3; i++ {
		if i != y {
			idx = append(idx, i)
		}
	}

	t1, t2, t3 := float64(idx[0]), float64(idx[1]), float64(idx[2])
	q1, q2, q3 := float64(Q[idx[0]]), float64(Q[idx[1]]), float64(Q[idx[2]])
	yf := float64(y)
	L1 := (yf - t2) * (yf - t3) / ((t1 - t2) * (t1 - t3))
	L2 := (yf - t1) * (yf - t3) / ((t2 - t1) * (t2 - t3))
	L3 := (yf - t1) * (yf - t2) / ((t3 - t1) * (t3 - t2))
	expectedQ := q1*L1 + q2*L2 + q3*L3
	expQ := int64(math.Round(expectedQ))

	diffQ := Q[y] - expQ
	p := (diffQ/int64(delta) - delta) / 2
	fmt.Fprintf(writer, "%d %d\n", y, p)
}
