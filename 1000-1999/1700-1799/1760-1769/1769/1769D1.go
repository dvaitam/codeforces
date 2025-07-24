package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const (
	suits = 4
	ranks = 9
)

var rankMap [256]int
var suitMap [256]int

func init() {
	for i := range rankMap {
		rankMap[i] = -1
	}
	for i := range suitMap {
		suitMap[i] = -1
	}
	order := "6789TJQKA"
	for i := 0; i < len(order); i++ {
		rankMap[order[i]] = i
	}
	suitMap['C'] = 0
	suitMap['D'] = 1
	suitMap['S'] = 2
	suitMap['H'] = 3
	precompute()
}

// suit state encoding
// state 0: nine not played
// otherwise: 1 + low*6 + high where low in [0..3], high in [0..5]

var playedMask [25]uint16
var frontierMask [25]uint16
var nextState [25][ranks]uint8
var totalPlayed [25]int

func encodeSuit(nine bool, low, high int) int {
	if !nine {
		return 0
	}
	return 1 + low*6 + high
}

func decodeSuit(state int) (bool, int, int) {
	if state == 0 {
		return false, 0, 0
	}
	state--
	low := state / 6
	high := state % 6
	return true, low, high
}

func precompute() {
	for i := range nextState {
		for j := range nextState[i] {
			nextState[i][j] = 255
		}
	}
	for st := 0; st < 25; st++ {
		nine, low, high := decodeSuit(st)
		var mask uint16
		if nine {
			mask |= 1 << 3
			for i := 0; i < low; i++ {
				mask |= 1 << uint(2-i)
			}
			for i := 0; i < high; i++ {
				mask |= 1 << uint(4+i)
			}
		}
		playedMask[st] = mask
		totalPlayed[st] = bits.OnesCount16(mask)
		if !nine {
			frontierMask[st] = 1 << 3
			nextState[st][3] = uint8(encodeSuit(true, 0, 0))
		} else {
			if low < 3 {
				r := 2 - low
				frontierMask[st] |= 1 << uint(r)
				nextState[st][r] = uint8(encodeSuit(true, low+1, high))
			}
			if high < 5 {
				r := 4 + high
				frontierMask[st] |= 1 << uint(r)
				nextState[st][r] = uint8(encodeSuit(true, low, high+1))
			}
		}
	}
}

func encodeBoard(b [4]int) uint32 {
	return uint32(((b[0]*25+b[1])*25+b[2])*25 + b[3])
}

func decodeBoard(code uint32) [4]int {
	var b [4]int
	b[3] = int(code % 25)
	code /= 25
	b[2] = int(code % 25)
	code /= 25
	b[1] = int(code % 25)
	code /= 25
	b[0] = int(code)
	return b
}

type move struct {
	suit int
	rank int
}

var maskPlayer [2][4]uint16
var playedCount [2][4][25]int
var memo map[uint64]bool

func remaining(board [4]int, player int) int {
	rem := 18
	for s := 0; s < 4; s++ {
		rem -= playedCount[player][s][board[s]]
	}
	return rem
}

func dfs(code uint32, turn int) bool {
	key := (uint64(code) << 1) | uint64(turn)
	if v, ok := memo[key]; ok {
		return v
	}
	board := decodeBoard(code)
	// gather available moves
	var mv [8]move
	mcount := 0
	for s := 0; s < 4; s++ {
		st := board[s]
		front := frontierMask[st]
		var mask uint16
		if turn == 0 {
			mask = maskPlayer[0][s] & front
		} else {
			mask = maskPlayer[1][s] & front
		}
		for mask != 0 {
			r := bits.TrailingZeros16(mask)
			mv[mcount] = move{suit: s, rank: r}
			mcount++
			mask &= mask - 1
		}
	}
	if mcount == 0 {
		res := !dfs(code, 1-turn)
		memo[key] = res
		return res
	}
	remCur := remaining(board, turn)
	for i := 0; i < mcount; i++ {
		m := mv[i]
		ns := nextState[board[m.suit]][m.rank]
		newBoard := board
		newBoard[m.suit] = int(ns)
		newCode := encodeBoard(newBoard)
		if remCur-1 == 0 {
			memo[key] = true
			return true
		}
		if !dfs(newCode, 1-turn) {
			memo[key] = true
			return true
		}
	}
	memo[key] = false
	return false
}

func setup(owner [4][9]int) {
	for p := 0; p < 2; p++ {
		for s := 0; s < 4; s++ {
			maskPlayer[p][s] = 0
		}
	}
	for s := 0; s < 4; s++ {
		for r := 0; r < 9; r++ {
			p := owner[s][r]
			maskPlayer[p][s] |= 1 << uint(r)
		}
	}
	for p := 0; p < 2; p++ {
		for s := 0; s < 4; s++ {
			for st := 0; st < 25; st++ {
				playedCount[p][s][st] = bits.OnesCount16(maskPlayer[p][s] & playedMask[st])
			}
		}
	}
}

func solveTest(aCards, bCards []string) string {
	var owner [4][9]int
	for _, c := range aCards {
		r := rankMap[c[0]]
		s := suitMap[c[1]]
		owner[s][r] = 0
	}
	for _, c := range bCards {
		r := rankMap[c[0]]
		s := suitMap[c[1]]
		owner[s][r] = 1
	}
	setup(owner)
	memo = make(map[uint64]bool)
	code := encodeBoard([4]int{0, 0, 0, 0})
	if dfs(code, 0) {
		return "Alice"
	}
	return "Bob"
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for {
		a := make([]string, 18)
		for i := 0; i < 18; i++ {
			if _, err := fmt.Fscan(reader, &a[i]); err != nil {
				return
			}
		}
		b := make([]string, 18)
		for i := 0; i < 18; i++ {
			if _, err := fmt.Fscan(reader, &b[i]); err != nil {
				return
			}
		}
		res := solveTest(a, b)
		fmt.Fprintln(writer, res)
	}
}
