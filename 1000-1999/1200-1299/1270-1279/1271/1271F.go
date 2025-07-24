package main

import (
	"bufio"
	"fmt"
	"os"
)

type piece struct {
	m, p, c int16
	typ     int
	size    int16
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveGroup(a1, b1, c1, a2, b2, c2 int, d []int) ([]int, bool) {
	totalM := d[0] + d[1] + d[2] + d[3]
	totalP := d[0] + d[1] + d[4] + d[5]
	totalC := d[0] + d[2] + d[4] + d[6]

	maxM := min(a1, totalM)
	maxP := min(b1, totalP)

	stride := maxP + 1
	size := (maxM + 1) * stride

	const INF int16 = 32000
	dp := make([]int16, size)
	parPiece := make([]int16, size)
	parPrev := make([]int32, size)
	for i := range dp {
		dp[i] = INF
		parPiece[i] = -1
		parPrev[i] = -1
	}
	dp[0] = 0

	// subject weights for each type
	wM := []int{1, 1, 1, 1, 0, 0, 0}
	wP := []int{1, 1, 0, 0, 1, 1, 0}
	wC := []int{1, 0, 1, 0, 1, 0, 1}

	pieces := make([]piece, 0)
	for typ := 0; typ < 7; typ++ {
		cnt := d[typ]
		mul := 1
		for cnt > 0 {
			take := mul
			if take > cnt {
				take = cnt
			}
			pieces = append(pieces, piece{
				m:    int16(wM[typ] * take),
				p:    int16(wP[typ] * take),
				c:    int16(wC[typ] * take),
				typ:  typ,
				size: int16(take),
			})
			cnt -= take
			mul <<= 1
		}
	}

	for idx, pc := range pieces {
		wm := int(pc.m)
		wp := int(pc.p)
		wc := pc.c
		for m := maxM; m >= wm; m-- {
			base := m * stride
			prevBase := (m - wm) * stride
			for p := maxP; p >= wp; p-- {
				prevIdx := prevBase + (p - wp)
				if dp[prevIdx] == INF {
					continue
				}
				newC := dp[prevIdx] + wc
				curIdx := base + p
				if newC < dp[curIdx] {
					dp[curIdx] = newC
					parPiece[curIdx] = int16(idx)
					parPrev[curIdx] = int32(prevIdx)
				}
			}
		}
	}

	finalIdx := -1
	for m := 0; m <= maxM && finalIdx == -1; m++ {
		base := m * stride
		for p := 0; p <= maxP; p++ {
			idx := base + p
			cUsed := int(dp[idx])
			if cUsed == int(INF) || cUsed > c1 {
				continue
			}
			if m > a1 || p > b1 {
				continue
			}
			m2 := totalM - m
			p2 := totalP - p
			c2Left := totalC - cUsed
			if m2 <= a2 && p2 <= b2 && c2Left <= c2 {
				finalIdx = idx
				break
			}
		}
	}

	if finalIdx == -1 {
		return nil, false
	}

	res := make([]int, 7)
	idx := finalIdx
	for idx != 0 {
		pcIdx := parPiece[idx]
		if pcIdx == -1 {
			break
		}
		pc := pieces[pcIdx]
		res[pc.typ] += int(pc.size)
		idx = int(parPrev[idx])
	}

	return res, true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a1, b1, c1 int
		var a2, b2, c2 int
		if _, err := fmt.Fscan(reader, &a1, &b1, &c1); err != nil {
			return
		}
		fmt.Fscan(reader, &a2, &b2, &c2)
		d := make([]int, 7)
		for i := 0; i < 7; i++ {
			fmt.Fscan(reader, &d[i])
		}
		ans, ok := solveGroup(a1, b1, c1, a2, b2, c2, d)
		if !ok {
			fmt.Fprintln(writer, -1)
		} else {
			for i, v := range ans {
				if i > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, v)
			}
			writer.WriteByte('\n')
		}
	}
}
