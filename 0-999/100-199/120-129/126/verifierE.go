package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ---- Embedded solver for 126E (simulated annealing approach) ----

const (
	eUP    = 0
	eRIGHT = 1
	eDOWN  = 2
	eLEFT  = 3
)

var typeColors = [10][2]int{
	{0, 3}, {0, 2}, {0, 1}, {0, 0},
	{1, 3}, {1, 2}, {1, 1},
	{2, 3}, {2, 2},
	{3, 3},
}

var typeMatrix = [4][4]int{
	{3, 2, 1, 0},
	{2, 6, 5, 4},
	{1, 5, 8, 7},
	{0, 4, 7, 9},
}

var eScoreMatrix [10][10]int
var eCost [22][22]int
var eCapNet [22][22]int
var eFlow [22][22]int

var solverInited bool

func initSolver() {
	if solverInited {
		return
	}
	solverInited = true
	for t := 0; t < 10; t++ {
		for p := 0; p < 10; p++ {
			tc1, tc2 := typeColors[t][0], typeColors[t][1]
			pc1, pc2 := typeColors[p][0], typeColors[p][1]
			s1 := 0
			if tc1 == pc1 {
				s1++
			}
			if tc2 == pc2 {
				s1++
			}

			s2 := 0
			if tc1 == pc2 {
				s2++
			}
			if tc2 == pc1 {
				s2++
			}

			if s1 > s2 {
				eScoreMatrix[t][p] = s1
			} else {
				eScoreMatrix[t][p] = s2
			}
		}
	}

	for i := 1; i <= 10; i++ {
		for j := 11; j <= 20; j++ {
			eCost[i][j] = 2 - eScoreMatrix[i-1][j-11]
			eCost[j][i] = -eCost[i][j]
		}
	}
}

func eGetType(c1, c2 int) int {
	return typeMatrix[c1][c2]
}

func ePackReq(req *[10]int) uint64 {
	var res uint64
	for i := 0; i < 10; i++ {
		res = (res << 5) | uint64(req[i])
	}
	return res
}

func eGetMaxScore(req *[10]int, inv *[10]int) int {
	for i := 0; i < 22; i++ {
		for j := 0; j < 22; j++ {
			eCapNet[i][j] = 0
			eFlow[i][j] = 0
		}
	}
	for i := 1; i <= 10; i++ {
		eCapNet[0][i] = req[i-1]
	}
	for j := 11; j <= 20; j++ {
		eCapNet[j][21] = inv[j-11]
	}
	for i := 1; i <= 10; i++ {
		for j := 11; j <= 20; j++ {
			eCapNet[i][j] = 28
		}
	}

	totalCost := 0
	totalFlow := 0

	var dist [22]int
	var parent [22]int
	var inQ [22]bool
	var q [1024]int

	for totalFlow < 28 {
		for i := 0; i < 22; i++ {
			dist[i] = 1e9
			parent[i] = -1
			inQ[i] = false
		}
		dist[0] = 0
		head, tail := 0, 0
		q[tail] = 0
		tail++
		inQ[0] = true

		for head != tail {
			u := q[head]
			head = (head + 1) & 1023
			inQ[u] = false

			for v := 0; v < 22; v++ {
				if eCapNet[u][v]-eFlow[u][v] > 0 && dist[v] > dist[u]+eCost[u][v] {
					dist[v] = dist[u] + eCost[u][v]
					parent[v] = u
					if !inQ[v] {
						inQ[v] = true
						q[tail] = v
						tail = (tail + 1) & 1023
					}
				}
			}
		}

		if dist[21] == 1e9 {
			break
		}

		push := 28 - totalFlow
		curr := 21
		for curr != 0 {
			p := parent[curr]
			avail := eCapNet[p][curr] - eFlow[p][curr]
			if avail < push {
				push = avail
			}
			curr = p
		}

		totalFlow += push
		totalCost += push * dist[21]

		curr = 21
		for curr != 0 {
			p := parent[curr]
			eFlow[p][curr] += push
			eFlow[curr][p] -= push
			curr = p
		}
	}

	return 56 - totalCost
}

func solveCase(rows []string, nums []int) (string, error) {
	initSolver()

	var targetColors [7][8]int
	for r := 0; r < 7; r++ {
		if len(rows[r]) != 8 {
			return "", fmt.Errorf("row %d has length %d", r+1, len(rows[r]))
		}
		for c := 0; c < 8; c++ {
			switch rows[r][c] {
			case 'B':
				targetColors[r][c] = 0
			case 'R':
				targetColors[r][c] = 1
			case 'W':
				targetColors[r][c] = 2
			case 'Y':
				targetColors[r][c] = 3
			default:
				return "", fmt.Errorf("invalid char %q", rows[r][c])
			}
		}
	}

	var inv [10]int
	for i := 0; i < 10; i++ {
		inv[i] = nums[i]
	}

	var dir [7][8]int
	for r := 0; r < 7; r++ {
		for c := 0; c < 8; c += 2 {
			dir[r][c] = eRIGHT
			dir[r][c+1] = eLEFT
		}
	}

	var req [10]int
	for r := 0; r < 7; r++ {
		for c := 0; c < 8; c += 2 {
			t := eGetType(targetColors[r][c], targetColors[r][c+1])
			req[t]++
		}
	}

	bestScore := eGetMaxScore(&req, &inv)
	currentScore := bestScore
	bestDir := dir

	memo := make(map[uint64]int)
	memo[ePackReq(&req)] = bestScore

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	startTime := time.Now()
	timeLimit := 1900 * time.Millisecond
	T := 5.0

	if bestScore < 56 {
		for iters := 0; ; iters++ {
			if iters&1023 == 0 {
				elapsed := time.Since(startTime)
				if elapsed > timeLimit {
					break
				}
				progress := float64(elapsed) / float64(timeLimit)
				T = 5.0 * math.Exp(-5.0*progress)
			}

			r := rng.Intn(6)
			c := rng.Intn(7)

			canFlip := false
			isHoriz := false
			if dir[r][c] == eRIGHT && dir[r+1][c] == eRIGHT {
				canFlip = true
				isHoriz = true
			} else if dir[r][c] == eDOWN && dir[r][c+1] == eDOWN {
				canFlip = true
				isHoriz = false
			}

			if !canFlip {
				continue
			}

			var t1, t2, t3, t4 int
			if isHoriz {
				t1 = eGetType(targetColors[r][c], targetColors[r][c+1])
				t2 = eGetType(targetColors[r+1][c], targetColors[r+1][c+1])
				t3 = eGetType(targetColors[r][c], targetColors[r+1][c])
				t4 = eGetType(targetColors[r][c+1], targetColors[r+1][c+1])
			} else {
				t1 = eGetType(targetColors[r][c], targetColors[r+1][c])
				t2 = eGetType(targetColors[r][c+1], targetColors[r+1][c+1])
				t3 = eGetType(targetColors[r][c], targetColors[r][c+1])
				t4 = eGetType(targetColors[r+1][c], targetColors[r+1][c+1])
			}

			req[t1]--
			req[t2]--
			req[t3]++
			req[t4]++

			hash := ePackReq(&req)
			newScore, exists := memo[hash]
			if !exists {
				newScore = eGetMaxScore(&req, &inv)
				memo[hash] = newScore
			}

			delta := newScore - currentScore
			accept := false
			if delta >= 0 {
				accept = true
			} else {
				prob := math.Exp(float64(delta) / T)
				if rng.Float64() < prob {
					accept = true
				}
			}

			if accept {
				currentScore = newScore
				if isHoriz {
					dir[r][c] = eDOWN
					dir[r+1][c] = eUP
					dir[r][c+1] = eDOWN
					dir[r+1][c+1] = eUP
				} else {
					dir[r][c] = eRIGHT
					dir[r][c+1] = eLEFT
					dir[r+1][c] = eRIGHT
					dir[r+1][c+1] = eLEFT
				}
				if currentScore > bestScore {
					bestScore = currentScore
					bestDir = dir
					if bestScore == 56 {
						break
					}
				}
			} else {
				req[t1]++
				req[t2]++
				req[t3]--
				req[t4]--
			}
		}
	}

	var bestReq [10]int
	for r := 0; r < 7; r++ {
		for c := 0; c < 8; c++ {
			if bestDir[r][c] == eRIGHT {
				t := eGetType(targetColors[r][c], targetColors[r][c+1])
				bestReq[t]++
			} else if bestDir[r][c] == eDOWN {
				t := eGetType(targetColors[r][c], targetColors[r+1][c])
				bestReq[t]++
			}
		}
	}

	eGetMaxScore(&bestReq, &inv)

	assignedPills := make([][]int, 10)
	for i := 1; i <= 10; i++ {
		for j := 11; j <= 20; j++ {
			for k := 0; k < eFlow[i][j]; k++ {
				assignedPills[i-1] = append(assignedPills[i-1], j-11)
			}
		}
	}

	var out [13][15]byte
	for r := 0; r < 13; r++ {
		for c := 0; c < 15; c++ {
			out[r][c] = '.'
		}
	}

	colorChars := []byte{'B', 'R', 'W', 'Y'}

	for r := 0; r < 7; r++ {
		for c := 0; c < 8; c++ {
			if bestDir[r][c] == eRIGHT {
				tc1 := targetColors[r][c]
				tc2 := targetColors[r][c+1]
				t := eGetType(tc1, tc2)

				p := assignedPills[t][0]
				assignedPills[t] = assignedPills[t][1:]

				pc1 := typeColors[p][0]
				pc2 := typeColors[p][1]

				s1 := 0
				if tc1 == pc1 {
					s1++
				}
				if tc2 == pc2 {
					s1++
				}
				s2 := 0
				if tc1 == pc2 {
					s2++
				}
				if tc2 == pc1 {
					s2++
				}

				if s1 >= s2 {
					out[2*r][2*c] = colorChars[pc1]
					out[2*r][2*c+2] = colorChars[pc2]
				} else {
					out[2*r][2*c] = colorChars[pc2]
					out[2*r][2*c+2] = colorChars[pc1]
				}
				out[2*r][2*c+1] = '-'

			} else if bestDir[r][c] == eDOWN {
				tc1 := targetColors[r][c]
				tc2 := targetColors[r+1][c]
				t := eGetType(tc1, tc2)

				p := assignedPills[t][0]
				assignedPills[t] = assignedPills[t][1:]

				pc1 := typeColors[p][0]
				pc2 := typeColors[p][1]

				s1 := 0
				if tc1 == pc1 {
					s1++
				}
				if tc2 == pc2 {
					s1++
				}
				s2 := 0
				if tc1 == pc2 {
					s2++
				}
				if tc2 == pc1 {
					s2++
				}

				if s1 >= s2 {
					out[2*r][2*c] = colorChars[pc1]
					out[2*r+2][2*c] = colorChars[pc2]
				} else {
					out[2*r][2*c] = colorChars[pc2]
					out[2*r+2][2*c] = colorChars[pc1]
				}
				out[2*r+1][2*c] = '|'
			}
		}
	}

	var sb strings.Builder
	fmt.Fprintln(&sb, bestScore)
	for r := 0; r < 13; r++ {
		sb.Write(out[r][:15])
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String()), nil
}

// ---- Verifier harness ----

// Embedded copy of testcases so the verifier is self-contained.
const testcasesRaw = `BYRRWRRY RBBRYYYB RWWBRWYY BWBYWRBB WYRWWBBB YRYRYYYR RBYBYWRB 18 24 19 10 14 19 6 22 15 21
YBRWWWYB YBYWYWYB WRYYWYRY WBWRRWBR BYYWRBYW WYBBBRYB BYYRBWBR 15 19 1 2 0 8 4 14 24 5
RRRBBYYW RBBBBRYW BYWBRBWW RRWRRRWB BWYYWRBW YBRBYBWW RBBBBRWB 20 23 26 20 22 2 9 8 14 1
BYRWWBRB WWWWBBYW WWBRYWRR BYWBRRRB WYBWRYRY BYRWWWBY BWYYWRRW 12 12 22 17 28 9 22 26 22 21
RRBYRRRB WYBBBWRR WYYWWWWY RYRYRYWW BBRWWRBW WBBYRBYY BRBRRRWW 0 9 22 8 27 21 0 10 19 1
WRYBYRYY WBWBBYYB RBBYBYRY WYBBRRRY RYRBBYYY WYWBRBYY YWYRBBWY 15 7 19 18 0 17 20 6 18 24
WRRWRYBW YRBWYRYY BWYWWYWY WBRYWRBW BYYRYRYW BBRWBRWB RRRYWWRR 20 16 15 0 24 4 15 20 19 2
BWYWRRWW RRWYWWYY YRYRWYYY BWYWBWYW YBBWWYWY BWYYBYWB WBYYRBWR 19 17 28 6 21 3 16 12 3 17
BWWWBYRB YWYWRWRR RYRRWWRW RBYRRRYR YWBBWRBB WBYWWYYW RWBWYWBR 19 15 19 27 18 28 12 16 11 25
WRRBRBWB RWYWBYBR YBYYRBBW YWRYYRWY WBBWWYRW BBYRYWRR BRRYBBWW 13 11 7 3 26 13 25 23 15 23
BYWBWYBW BWBRWRWY RYWBWWYR YYWYYYRR YBRBYBRB RRBBYRBR YRBYWRYB 3 3 6 16 3 3 17 22 2 19
BBRWYRWW RWWRYRWR RRYYYBRB WBBYBRRW BRBYYYYR WBWYYWBB BYYRYBWR 1 7 26 24 10 25 23 16 8 24
BYWRRRWB BBBRWBYW RRYYYRBR BWRBWYWR RBRRWWBB RYWRBBRR RRYRYYBW 13 17 5 0 13 20 25 18 19 4
RYWRRWRY BYYBWRWR BWYBYWRR BYBRYWWW BWWBBBYB RBBWWBBB RBWWYBYY 9 27 19 0 7 11 11 26 10 5
YRYBBWRW RBYWRBBR YRRWWYBY RWWBWWWR WWWWYYBB YRRRWWYY YRYWBRYW 18 24 3 26 0 1 2 21 22 22
RYWBBYWY YBWYWYRR BYRRRBWB YWYWRYYR YRRBWYYB WWRWBRWW RBYBYWRW 16 3 0 6 7 11 17 28 27 5
WWRWRBYW WYWBBBRR BRWRRYBY WRWBRRRR WYWYWRWR RYYRRYYW RWWYRYYY 13 23 6 28 19 1 8 18 16 23
BYBYRBBW RBBWWBWY BYYWRRRY YWYRBWYB RWRWYWBW YBBRRWWR RRWWRRBR 17 15 26 11 5 12 15 15 6 26
YYWRBYYW WBBRRYRB YRYRBRRR BWWYWWWW BWYBYWWY BBBRWRBB BWYYWWYB 24 7 8 22 10 13 25 27 10 1
BBBWRBWW RWRRRYYR BWYWWWYW RBYRWBRR RYWBYYRY BRYYWRBB WYRRRYRR 21 1 23 7 3 18 0 7 18 3
BBRRRRWY YBWBWYBW BBBWWBRB WWWRYYRR BBYRYYRB BBYRRWBB WRWWYBYB 23 9 17 22 9 13 19 10 15 22
WBBYRBBR WRBRRWYW YYYRRBRW RWRBYYWR RWBRBWBR BBBBBRYW BYBRBRRB 18 14 8 27 9 28 19 19 26 5
YWBRBRYR BBWRBWWY YRWYYWYR YBYYRBWR BBBYWRRW YBBYYWWY WYRRBBRB 10 27 19 11 16 14 26 13 25 6
YBWYYYWB RWRYWRBW YWYBYRWR BWBYWYWY RWBBBYBY WRWWWWRW WBWWWWRY 18 12 7 19 3 7 28 8 14 20
RWRWBBRW WBYBWRYR RWBBWYYW YWRBBRYW BWRYYRWW WWRWRBWY YWYBBRRY 6 18 11 24 8 1 4 25 0 22
WYWWWRBW RBYYBYYB BYYRWRYY BWYWYYYW WBBBBRYY BYBRRYWR WRRYYRBW 5 14 3 22 4 10 24 1 27 22
BWWWRRYY BWYBWWYW WBYRWYBR BBBWYRWB BBWYWRYW WWYYWYYR RYWRBRRW 24 9 8 3 10 3 0 6 2 18
WBWYYBRW RRWRBBBR BBYBWRWY RRRYYRBY YRYBRRWR YWYYWRYY YRYBBYBW 10 24 24 10 15 20 22 18 24 28
WBBBWRWR BBRBWYWB BBRRBYRW BRRYBRYR YWRRBWYB YRYYRYBB RWBRBBRR 12 11 6 2 13 6 22 24 4 22
RYYWWRWW YWRRRYRY WRRBWWBW RYRBRRWY YRRRYRBB YRBRYRYB YBWYBWYB 15 4 27 22 4 22 19 11 14 10
BBRBBYWR BWRWYYRR YBYYBYRB RRWRBYYY RBBBBYBB RYWWBYWB WWBWBBYW 0 15 18 11 5 7 24 23 5 16
RWYYRRWR BRRBWYYW BRRWRBBW RBYWWBYW WYWYWWRW YRBBWWRY WYRBRWRY 28 9 19 21 8 0 23 11 3 18
WWWRRWBR WRWWRYWY WWYYWRBW BWYRYBWR RBWWBWBR WYBRYRWB BWWWBBYR 14 21 21 26 15 9 28 4 12 0
RRYBBWBW RYBBRBRR RWBRWRRW BRYRBWBW RRBYWRYR WRYBRYWY YWRBYWBB 19 20 23 1 25 11 26 17 15 27
YYBRBWYW WWRWRBWB WWYBWRYW WWRBBRYB WRWYRYRY RYRYYWBR RBWYYYYB 18 10 4 5 8 20 22 5 18 5
YBWRYWRY BYYWWYWY YBRYYYRR RBWRRYRR YYRBWWBB YWYBWYBY YWYYWBBB 7 6 6 8 20 26 20 10 20 17
RYRRBBWB RWYYBRWB WRYYYYWR RWWRYYWW YRYBWYRB YRWBBWBB RWWWBWBR 1 25 2 28 27 17 10 28 2 25
WRWRBYWY YYYBBYWB WBYWYYBB RYYWBBRW RYWYWWBW RYBYYWWR BWRBRYWR 0 13 25 25 19 11 19 14 5 18
YWRWWWYY BYRBRYWR RYYRBRRW WRBYRYYB WYWWYYBY YWRBYWYY YWWRWWYR 16 28 6 9 2 14 4 8 3 5
WYBRYRWB WBWWBRBY BYWRRBWR RYRYBWWB BYWWWBYB BRBRBRYY WBBRWBRR 17 25 7 22 15 16 19 19 2 13
RWWYYBBY BWWBBWRB RWBWRRWB BBBWYYYY WWBRYBWR YWBYYRWY YRWBRBRB 19 5 13 14 7 28 24 21 9 19
BBYBBRBY RRYRBRRB RBRBRYYB YYWBYRYR BRYBWYBY RYYWRYWB YBWBYWYY 20 15 11 25 15 16 16 7 0 28
WRRWWYYW WRBRRYRY YBYBBBRB RRYWRRWY WWRWRYRR BRWWRWRB BYWRBWYY 14 17 11 13 22 28 24 25 23 24
WWWWWBBW BYWYYWRR BYRBRRWW YWRWBRRW WRYBYWRY YBYBYBYR WBRYYYBR 4 25 7 15 24 2 27 1 2 17
BBWBWBRY BRWWBRYR RWBBYWYB BBRBRYRB RYBRBRYR YBWYYRRB WYYYRWWW 0 6 26 27 7 7 14 0 2 26
YWBYBWYY BWWRRWWB BYYBWWWY YYYYYWBR YBWRYWBW BBBRWRBR BRRWWYBY 28 22 4 28 13 19 0 11 5 18
YWBBWRRY YYRRWYWY WYWWYBBB WYYRBYBW YRWWBBRY YYWWBBWB BRWRRYYB 28 11 26 25 1 2 4 27 27 3
BYYYRWYB RWWRRYWW BYWBWWRY WBBYBRBW BWYBWBYR WRRYBYBR BYRRWRBW 13 13 26 27 1 6 17 2 27 25
YRRYYBYB WRYRBBBY RWBBWRBY BRBRBYWR RWBYWWWB RRYYRBRR WYWYBBYY 8 25 0 25 12 3 0 3 18 0
WBBYRYBB BYRRRBBB RRYWWWBY YYYYWRBR RYYRYWBY WRBYWRBR WYWWRRRY 3 12 24 25 27 18 17 20 19 7
WBRWWWWY YRWBRWWR BWBWWBBW RYWWBBRB BRBRWBBR RWRWYWRB WRWRYYWW 0 19 25 4 8 7 9 28 9 4
RWRBBYWB RWRBWBBB WWWYYRBY YWWYBYYB YBWYWRYR YRRRBRRW WWYYBWWW 13 22 17 15 18 0 9 13 0 6
BRWRBYBR BWBYYBYB BWRWYRRW WRRRYRYR RWWBRWBY RYRYBWBY WBYWRBRB 0 13 3 10 18 10 3 21 5 28
WWBRWYRR BWWRRWYB WBBRRYBR YYYYWYYB WBYRBBWB YYYRWYBW RWWBRBRW 10 12 19 8 12 22 7 6 10 15
WRBWRWRY BYYYYYBB RWBRRRYR RYBRBBRW WWYRBBBW RBWBWRBY YYYYRBBB 1 13 19 13 28 15 6 0 10 0
YRYBBWYR WRBRYRWB RYYYBYWB WRRYRWRB WYYBRBWR YBBYWYBR RYRYRYWB 25 14 5 6 27 15 2 4 12 18
BRYYWWBB YYYRRWBB WRWBRRYB WBBYBYBB BRBYBYWY BYWWYBWB WWRRWRBR 20 7 14 27 13 26 3 28 0 16
RRBRYRYY RRBWRBYW WWYYWWYB BYBRWWBY WYRWWBWR YRYBRWWB YRWBWBRR 13 10 10 0 2 4 20 21 11 22
WYBWBWWR YBYRYYYW YRYWBYYR WWRBRRYB YWWWYBBB WWYRYRBB WBYWWYRW 12 19 18 22 3 25 25 23 17 24
YWWWBRYR RBYRBRWB WRRBWYBW BBWBBYYY WWBRWRBY BBBBRWYY BYRYRWYR 2 28 27 14 26 27 0 14 12 26
BBYBBRYR BBBYWWBW YYBBRBRY BYYRBRYB BRRBBBBW RWBWRWRB YBRWYRRY 20 26 11 4 28 1 22 5 22 8
WWRWYWWB YWRWRBYW YBBWRRYR YBRWBWRR YYYYRBBY WWYYYWYY YRBBWWBW 9 10 15 2 19 18 4 26 1 19
BYRBWRWB YYWBBRWB BWRRRRWB BRYBBYRW YYWYRBRY RBWBWYYR YRBYBBWW 9 21 5 26 4 19 0 7 20 22
BWBRRYYW BYBYYBBR YYWRYWRB BYBWRRYR RBYRRRYW BBYBYWWW BRWBWBRY 24 2 19 20 2 16 2 19 9 1
YBYYBWRY RBBRRWYR RYYBWWYB YWBBRWRY RRYRRRWY WRWYRYRB WBYYRYBB 14 22 7 7 27 19 17 26 22 28
RYBRYRYW YBBYRBRY RWBBWBWR BWYYWBBY WRBBYRWY YWWYYBYW RRRRWYWY 8 9 24 1 14 10 11 28 19 7
WBRRBWYR RBYWWWRY BRBBYBRR YWWYWWRB WWRBRWWY RRBBRYYY BYWYBRRW 12 19 20 4 22 16 28 25 24 9
WYYWWBYB WYYRRBRY RWYRWWRR WRBWWYWY YWBBRWBB RWYYBYBW BYWRWYYY 0 8 11 26 3 15 10 19 9 24
RRBRWWBR YYYWRYRB WRWWWBYB YRWRRRRB RWWBBYWW BBYBRBYR RYYWYYYB 4 9 2 27 3 14 14 22 9 14
WRWRYBYB RRBYYWWY WWRBWBRY WRYYBRYR BYWYBWYB WWBRWRWY YYWBYBRR 3 1 25 7 13 1 6 12 22 6
BRWRYRBB YRWBBYYW YBBRBBRW WRBRWRBR WBBBBWWY YWRWBRRR RWRRWBRR 22 5 12 22 20 24 8 9 17 7
BRRRYBYB YBBWYRBB YRRYYWBW YWRWWYYR YWYYRRYR RRWBBYYB WBWRRRBW 20 17 28 0 24 27 2 1 4 9
YBYRYBRW RWRBWYBB BWBWYYYW WYRBYRYW WWRYBBWY WRBYRRRW WBYYRBRR 25 13 20 22 4 27 15 0 20 12
RRBYBBRY WWYBBBRW BRYBWWWW WWYRYBWW YBRRYYWW RBRWWBYR YBRBRBWY 24 23 6 16 16 2 10 25 3 21
BRRBRRRW RRBRWWYR WWRYYBRW RWYWRBYY BWRBRWRR WBWRBBBB RYRWBBWW 19 13 18 12 7 14 1 2 27 23
BBYWBWBB RRRBYBRR YRRRRRRR BRRYBWYW RRRWWWWR RRBBYRYR BWWYRBWR 5 17 25 15 2 14 18 14 21 9
BYWBWRRB RRYYBBYW WYYYRRYB BWRBYBYR YRBRRBYR RWWWYRWR YBRRRWWB 21 19 13 25 6 9 10 23 19 15
RRBBYBWW BBYYYRYW RYWWWBWR BWRYBBRR YRRBYRRY RBYBBBYR YBRWYWRW 19 27 13 16 25 2 2 23 18 16
BYWRBBYY YWWRRBRR WBWYYRBR YBWYWRRB YRRBWRYW RBRWRRBW BYYWRBWW 3 9 22 12 0 16 24 12 27 2
RRRRWYRY YWWRWWWR BBYYRRWB RWRWRYRR YBYBWBWB WBBWYYRB WYYBBYBR 4 17 25 12 28 20 6 0 23 24
WBRWWRBB YRBRYYBB YWRBBWBY YWRYWBYB YBBBRWRB WBRYRWWY YWBBBWWY 26 16 20 5 14 13 12 9 19 28
WBWWRWWY WYYRYWWR RYRYBBYR RRYYYRYB BWWWRYWR BYBBWBWY BRBYYBYY 5 12 2 13 10 18 13 4 10 7
YYRBRWBY YYYYYRYW WRRRBYRY RRWYRRBR RBBWWYRB YRYYYRYB RRYYYBWR 13 22 17 28 15 26 5 3 27 17
WYYWRWYY BWBWRBYB WBRRRYWW BRWBRWRY WBWWWBWB WYYBYYBB RYRBBRYY 20 26 23 17 0 11 24 23 25 19
YRYBRYWR YYBYYBBW RWWYWRYY RBYWBYYR BRYBWRBW BBWWBYWW YBBRWWYB 13 3 17 20 21 2 28 25 9 25
RYRYRWRB WBYYBRRW BRRYWYWY WYBBBYWW YBWRRYWY WWBYYRWR WWYWWWYY 19 22 24 10 13 28 3 27 23 25
BBWRRRRB YWYWWWYY WWBYBBWR RYYRYWRR RRRYBBWY YWRRRBWW WWBBWBBB 6 4 7 12 2 8 19 17 4 28
YBRRYBRW WYWYWWBR RYRRYBBY WYRBYYYB YRYYBYRR BBBWBBBR WRYRRRWW 22 2 8 24 20 18 16 15 3 3
BYRWRWYY RYYYBRWW BRBWRWRR WWBBYBBB BRRBWRBB WBRWRBBW YYBBYWYR 2 6 16 6 15 5 27 17 26 3
BYRBRRWW RYWWRBYB YRBBYBYW YWWRBBRW WRYWWRRR WRYWYWWW BBBYYRRR 22 9 6 22 3 19 25 0 15 14
WYWYYYWB BWWBYYWW YRBYBWYR RWWWWRRY RRBBBBBR YWBRYRBW WYRBBYBW 4 17 6 13 7 6 24 18 3 20
YRRBYYBW BYBRRBBY YYWRRBYY WYBWWBRB BWWBRBYY YRBRWWRW YBRRYBYY 0 4 0 5 17 3 21 16 20 14
WYWBRYWW WYRYYWYW RBRYRBRR YRYBBYYY RYRRRBBB WRBWBWBW RWYWBWRY 5 20 14 20 1 12 17 21 15 3
WWBYWBRB RYRRWYBR WYBYYWBW RWRYRYWW WWWWRBBY BYRYYWYR WYBYWBBY 24 17 12 22 18 10 22 8 0 12
BYWBYWBB WBYBRRWY BBBWYYYR WBRRYWYR YWBRBBYR BWRRYWWB BRYYYYBY 27 3 14 19 2 25 27 10 24 13
RRRWYRBY WYRBRBBW YRWBWRRW BRWBWWBB YYYWWBRW BRYRYRYR BYWWRBBY 13 25 10 23 8 24 20 21 21 14
WBWYRRWB WWRRWRYB YYRYYYWB RYBYRRYW YBRYWBWR BWBYWWBW BBYBYYRW 21 0 16 20 12 16 23 22 14 16
YRBRBWWY YWYYBYYW YRBWYRYB WYYBYYRB RRRYWYRB RRWRYWWR BWWWYBYY 19 22 13 28 11 18 18 16 2 20
YWBWBWBB RYRYRWYW WRYRWRBW YRBYYBRW WYYWWRYW WWRBBWRB BYRRYBWR 8 26 27 4 23 11 6 17 4 7
BRYBWYYW WYYBWBWB BBBBYWYW YWYWBYYB RYRWWRYW WYYRBYYB RWYRBYRY 19 17 25 27 7 18 8 0 22 25`

func parseTestcases() ([]struct {
	rows []string
	nums []int
}, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]struct {
		rows []string
		nums []int
	}, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 17 {
			return nil, fmt.Errorf("line %d: expected 17 fields, got %d", idx+1, len(parts))
		}
		rows := parts[:7]
		nums := make([]int, 10)
		for i := 0; i < 10; i++ {
			val, err := strconv.Atoi(parts[7+i])
			if err != nil {
				return nil, fmt.Errorf("line %d number %d: %v", idx+1, i+1, err)
			}
			nums[i] = val
		}
		cases = append(cases, struct {
			rows []string
			nums []int
		}{rows: rows, nums: nums})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		expect, err := solveCase(tc.rows, tc.nums)
		if err != nil {
			fmt.Printf("solver failed on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		var sb strings.Builder
		for _, row := range tc.rows {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		for i, num := range tc.nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(num))
		}
		sb.WriteByte('\n')
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		// For this problem, compare only the score (first line).
		// The grid layout may differ due to randomization but score must match or exceed.
		expectLines := strings.SplitN(strings.TrimSpace(expect), "\n", 2)
		gotLines := strings.SplitN(strings.TrimSpace(got), "\n", 2)
		expectScore, _ := strconv.Atoi(strings.TrimSpace(expectLines[0]))
		gotScore, err2 := strconv.Atoi(strings.TrimSpace(gotLines[0]))
		if err2 != nil {
			fmt.Printf("case %d: cannot parse candidate score: %v\ngot:\n%s\n", idx+1, err2, got)
			os.Exit(1)
		}
		if gotScore < expectScore {
			fmt.Printf("case %d failed\nexpected score: %d\ngot score: %d\n", idx+1, expectScore, gotScore)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
