package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ---- Embedded solver logic from 126E.go ----

const (
	le = 20
	ri = 21
)

var (
	U    = make(map[int64]struct{})
	a    [7][8]int
	P    [4][4]int
	L    [10]int
	R    [10]int
	A    [10]int
	Bv   [10]int
	T    = -1
	p    [22][22]int
	v    [7][8]bool
	w    [10][10]bool
	r2   [7][8]bool
	d2   [7][8]bool
	used [22]bool
	s    [13][15]byte
	S    [13][15]byte
	z    = []byte{'B', 'R', 'W', 'Y'}
	tAdd = 0
)

func resetGlobals() {
	U = make(map[int64]struct{})
	T = -1
	tAdd = 0
	for i := range a {
		for j := range a[i] {
			a[i][j] = 0
			v[i][j] = false
			r2[i][j] = false
			d2[i][j] = false
		}
	}
	for i := range P {
		for j := range P[i] {
			P[i][j] = 0
		}
	}
	for i := range L {
		L[i], R[i], A[i], Bv[i] = 0, 0, 0, 0
	}
	for i := range p {
		for j := range p[i] {
			p[i][j] = 0
		}
	}
	for i := range w {
		for j := range w[i] {
			w[i][j] = false
		}
	}
	for i := range used {
		used[i] = false
	}
	for i := range s {
		for j := range s[i] {
			s[i][j] = '.'
			S[i][j] = '.'
		}
	}
}

func add(x, y int) {
	P[x][y] = tAdd
	P[y][x] = tAdd
	L[tAdd] = x
	R[tAdd] = y
	tAdd++
}

func f2(x int) bool {
	if x == ri {
		return true
	}
	used[x] = true
	for i := ri; i >= 0; i-- {
		if p[x][i] > 0 && !used[i] && f2(i) {
			p[x][i]--
			p[i][x]++
			return true
		}
	}
	return false
}

func gg() {
	var c [10]int
	var D [10][10]int
	t := 0
	var h int64
	for i := 0; i < 10; i++ {
		h = h*1000000007 + int64(Bv[i])
	}
	if _, ok := U[h]; ok {
		return
	}
	U[h] = struct{}{}
	copy(c[:], A[:])
	for i := 0; i < 10; i++ {
		D[i][i] = min(A[i], Bv[i])
		t += D[i][i] * 2
	}
	if 28+t/2 <= T {
		return
	}
	for i := range p {
		for j := range p[i] {
			p[i][j] = 0
		}
	}
	for i := 0; i < 10; i++ {
		if A[i] > Bv[i] {
			p[le][i] = A[i] - Bv[i]
		} else if A[i] < Bv[i] {
			p[i+10][ri] = Bv[i] - A[i]
		}
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if w[i][j] {
				u := min(p[le][i], p[j+10][ri])
				if u > 0 {
					p[le][i] -= u
					p[j+10][ri] -= u
					p[i][j+10] = 28
					t += u
					D[i][j] += u
				}
			}
		}
	}
	for {
		for i := range used {
			used[i] = false
		}
		if f2(le) {
			t++
		} else {
			break
		}
	}
	if t <= T {
		return
	}
	T = t
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			D[i][j] += p[j+10][i]
		}
	}
	for i := 0; i < 7; i++ {
		for j := 0; j < 8; j++ {
			s[i*2][j*2] = '.'
		}
	}
	for i := 0; i < 7; i++ {
		for j := 0; j < 8; j++ {
			if r2[i][j] {
				o := P[a[i][j]][a[i][j+1]]
				for k := 0; k < 10; k++ {
					if D[k][o] > 0 {
						D[k][o]--
						c[k]--
						if a[i][j] == R[k] || a[i][j+1] == L[k] {
							L[k], R[k] = R[k], L[k]
						}
						s[i*2][j*2] = z[L[k]]
						s[i*2][j*2+2] = z[R[k]]
						break
					}
				}
			}
		}
	}
	for i := 0; i < 7; i++ {
		for j := 0; j < 8; j++ {
			if d2[i][j] {
				o := P[a[i][j]][a[i+1][j]]
				for k := 0; k < 10; k++ {
					if D[k][o] > 0 {
						D[k][o]--
						c[k]--
						if a[i][j] == R[k] || a[i+1][j] == L[k] {
							L[k], R[k] = R[k], L[k]
						}
						s[i*2][j*2] = z[L[k]]
						s[i*2+2][j*2] = z[R[k]]
						break
					}
				}
			}
		}
	}
	for i := 0; i < 7; i++ {
		for j := 0; j < 8; j++ {
			if s[i*2][j*2] == '.' && r2[i][j] {
				for k := 0; k < 10; k++ {
					if c[k] > 0 {
						c[k]--
						s[i*2][j*2] = z[L[k]]
						s[i*2][j*2+2] = z[R[k]]
						break
					}
				}
			}
		}
	}
	for i := 0; i < 7; i++ {
		for j := 0; j < 8; j++ {
			if s[i*2][j*2] == '.' && d2[i][j] {
				for k := 0; k < 10; k++ {
					if c[k] > 0 {
						c[k]--
						s[i*2][j*2] = z[L[k]]
						s[i*2+2][j*2] = z[R[k]]
						break
					}
				}
			}
		}
	}
	for i := 0; i < 13; i++ {
		for j := 0; j < 15; j++ {
			S[i][j] = s[i][j]
		}
	}
}

func ff(x, y int) {
	if x == 7 {
		gg()
		return
	}
	if y == 8 {
		ff(x+1, 0)
		return
	}
	if v[x][y] {
		ff(x, y+1)
		return
	}
	if y < 7 && !v[x][y+1] {
		v[x][y], v[x][y+1] = true, true
		idx := P[a[x][y]][a[x][y+1]]
		Bv[idx]++
		s[x*2][y*2+1] = '-'
		r2[x][y] = true
		ff(x, y+1)
		v[x][y], v[x][y+1] = false, false
		Bv[idx]--
		s[x*2][y*2+1] = '.'
		r2[x][y] = false
	}
	if x < 6 && !v[x+1][y] {
		v[x][y], v[x+1][y] = true, true
		idx := P[a[x][y]][a[x+1][y]]
		Bv[idx]++
		s[x*2+1][y*2] = '|'
		d2[x][y] = true
		ff(x, y+1)
		v[x][y], v[x+1][y] = false, false
		Bv[idx]--
		s[x*2+1][y*2] = '.'
		d2[x][y] = false
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveCase(rows []string, nums []int) (string, error) {
	if len(rows) != 7 || len(nums) != 10 {
		return "", fmt.Errorf("invalid testcase parts")
	}
	resetGlobals()
	for i := 0; i < 7; i++ {
		if len(rows[i]) != 8 {
			return "", fmt.Errorf("row %d has length %d", i+1, len(rows[i]))
		}
		for j := 0; j < 8; j++ {
			switch rows[i][j] {
			case 'B':
				a[i][j] = 0
			case 'R':
				a[i][j] = 1
			case 'W':
				a[i][j] = 2
			case 'Y':
				a[i][j] = 3
			default:
				return "", fmt.Errorf("invalid char %q", rows[i][j])
			}
		}
	}
	add(0, 3)
	add(0, 2)
	add(0, 1)
	add(0, 0)
	add(1, 3)
	add(1, 2)
	add(1, 1)
	add(2, 3)
	add(2, 2)
	add(3, 3)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if L[i] == L[j] || L[i] == R[j] || R[i] == L[j] || R[i] == R[j] {
				w[i][j] = true
			}
		}
	}
	for i := 0; i < 10; i++ {
		A[i] = nums[i]
	}
	for i := 0; i < 13; i++ {
		for j := 0; j < 15; j++ {
			s[i][j] = '.'
		}
	}
	ff(0, 0)
	var out strings.Builder
	fmt.Fprintln(&out, T)
	for i := 0; i < 13; i++ {
		out.Write(s[i][:15])
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String()), nil
}

// ---- Verifier harness ----

// Embedded copy of testcasesE.txt so the verifier is self-contained.
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
		for i, row := range tc.rows {
			sb.WriteString(row)
			sb.WriteByte('\n')
			_ = i
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
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
