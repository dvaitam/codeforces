package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testD struct {
	s string
}

// ---------- correct embedded solver for 625D ----------

func validSingle(d int, needPos bool) bool {
	if d < 0 || d > 9 {
		return false
	}
	if needPos && d == 0 {
		return false
	}
	return true
}

func validGap(d int) bool {
	return d == 0
}

func validOverlapOuter(d int, topPos bool, botPos bool) bool {
	if d < 0 || d > 18 {
		return false
	}
	minTop := 0
	if topPos {
		minTop = 1
	}
	minBot := 0
	if botPos {
		minBot = 1
	}
	l := minTop
	if d-9 > l {
		l = d - 9
	}
	r := 9
	if d-minBot < r {
		r = d - minBot
	}
	return l <= r
}

func validOverlapCenter(d int, needPos bool) bool {
	if d < 0 || d > 18 || d%2 != 0 {
		return false
	}
	x := d / 2
	if x < 0 || x > 9 {
		return false
	}
	if needPos && x == 0 {
		return false
	}
	return true
}

func choosePair(d int, topPos bool, botPos bool) (int, int) {
	minTop := 0
	if topPos {
		minTop = 1
	}
	minBot := 0
	if botPos {
		minBot = 1
	}
	l := minTop
	if d-9 > l {
		l = d - 9
	}
	r := 9
	if d-minBot < r {
		r = d - minBot
	}
	if l > r {
		return -1, -1
	}
	return l, d - l
}

func solveForM(nd []int, m int) (string, bool) {
	L := len(nd)
	if m < 1 || m > L || m < L-1 {
		return "", false
	}
	extra := 0
	if m == L {
		extra = 0
	} else {
		if nd[m] != 1 {
			return "", false
		}
		extra = 1
	}

	h := m / 2
	prev := make([][12]int16, h+1)
	act := make([][12]byte, h+1)
	for i := 0; i <= h; i++ {
		for j := 0; j < 12; j++ {
			prev[i][j] = -1
		}
	}

	var cur [12]bool
	cur[extra] = true

	for p := 0; p < h; p++ {
		var next [12]bool
		nk := nd[p]
		nj := nd[m-1-p]

		for id := 0; id < 12; id++ {
			if !cur[id] {
				continue
			}
			phase := id / 4
			st := id % 4
			cl := st >> 1
			cr := st & 1

			add := func(nphase, nst int, action byte) {
				nid := nphase*4 + nst
				if prev[p+1][nid] == -1 {
					prev[p+1][nid] = int16(id)
					act[p+1][nid] = action
					next[nid] = true
				}
			}

			try := func(kind int, first bool, needPos bool, nphase int, action byte) {
				for ql := 0; ql <= 1; ql++ {
					d := 10*ql + nk - cl
					if d < 0 || d > 18 {
						continue
					}
					pr := 10*cr + nj - d
					if pr < 0 || pr > 1 {
						continue
					}
					ok := false
					if kind == 0 {
						ok = validSingle(d, needPos)
					} else if kind == 1 {
						ok = validOverlapOuter(d, p == 0, first)
					} else {
						ok = validGap(d)
					}
					if ok {
						add(nphase, (ql<<1)|pr, action)
					}
				}
			}

			if phase == 0 {
				try(0, false, p == 0, 0, 0)
				try(0, false, true, 2, 1)
				try(1, true, false, 1, 2)
			} else if phase == 1 {
				try(1, false, false, 1, 3)
			} else {
				try(2, false, false, 2, 4)
			}
		}

		cur = next
	}

	finalID := -1
	centerAct := byte(0)

	if m%2 == 0 {
		for id := 0; id < 12; id++ {
			if !cur[id] {
				continue
			}
			phase := id / 4
			if phase == 0 {
				continue
			}
			st := id % 4
			if (st >> 1) == (st & 1) {
				finalID = id
				break
			}
		}
		if finalID == -1 {
			return "", false
		}
	} else {
		for id := 0; id < 12; id++ {
			if !cur[id] {
				continue
			}
			phase := id / 4
			st := id % 4
			cl := st >> 1
			cr := st & 1
			d := 10*cr + nd[h] - cl
			if d < 0 || d > 18 {
				continue
			}
			if phase == 0 {
				if validOverlapCenter(d, true) {
					finalID = id
					centerAct = 5
					break
				}
			} else if phase == 1 {
				if validOverlapCenter(d, false) {
					finalID = id
					centerAct = 6
					break
				}
			} else {
				if validGap(d) {
					finalID = id
					centerAct = 7
					break
				}
			}
		}
		if finalID == -1 {
			return "", false
		}
	}

	ids := make([]int, h+1)
	acts := make([]byte, h)
	ids[h] = finalID
	for p := h; p > 0; p-- {
		acts[p-1] = act[p][ids[p]]
		ids[p-1] = int(prev[p][ids[p]])
	}

	t := -1
	for p := 0; p < h; p++ {
		if acts[p] == 2 {
			t = p
			break
		}
	}
	if t == -1 && m%2 == 1 && centerAct == 5 {
		t = h
	}
	if t == -1 {
		for p := 0; p < h; p++ {
			if acts[p] == 1 {
				t = m - (p + 1)
				break
			}
		}
	}
	if t == -1 {
		return "", false
	}
	s := m - t
	if s < 1 {
		return "", false
	}

	b := make([]int, s)
	for i := 0; i < s; i++ {
		b[i] = -1
	}

	setDigit := func(idx, val int) bool {
		if idx < 0 || idx >= s || val < 0 || val > 9 {
			return false
		}
		if b[idx] != -1 && b[idx] != val {
			return false
		}
		b[idx] = val
		return true
	}

	for p := 0; p < h; p++ {
		stBefore := ids[p] % 4
		cl := stBefore >> 1
		stAfter := ids[p+1] % 4
		ql := stAfter >> 1
		d := 10*ql + nd[p] - cl

		switch acts[p] {
		case 0, 1:
			if !setDigit(s-1-p, d) {
				return "", false
			}
		case 2:
			u, v := choosePair(d, p == 0, true)
			if !setDigit(s-1-p, u) || !setDigit(p-t, v) {
				return "", false
			}
		case 3:
			u, v := choosePair(d, p == 0, false)
			if !setDigit(s-1-p, u) || !setDigit(p-t, v) {
				return "", false
			}
		case 4:
		default:
			return "", false
		}
	}

	if m%2 == 1 {
		st := ids[h] % 4
		cl := st >> 1
		cr := st & 1
		d := 10*cr + nd[h] - cl
		if centerAct == 5 || centerAct == 6 {
			if d%2 != 0 {
				return "", false
			}
			if !setDigit(h-t, d/2) {
				return "", false
			}
		}
	}

	for i := 0; i < s; i++ {
		if b[i] == -1 {
			return "", false
		}
	}
	if b[0] == 0 || b[s-1] == 0 {
		return "", false
	}

	out := make([]byte, s+t)
	for i := 0; i < s; i++ {
		out[i] = byte('0' + b[s-1-i])
	}
	for i := s; i < s+t; i++ {
		out[i] = '0'
	}
	return string(out), true
}

func solveD(s string) string {
	L := len(s)
	nd := make([]int, L)
	for i := 0; i < L; i++ {
		nd[i] = int(s[L-1-i] - '0')
	}

	if ans, ok := solveForM(nd, L); ok {
		return ans
	}
	if ans, ok := solveForM(nd, L-1); ok {
		return ans
	}
	return "0"
}

// ---------- test generation and verification ----------

func genTests() []testD {
	rand.Seed(4)
	tests := make([]testD, 100)
	for i := range tests {
		l := rand.Intn(20) + 1
		b := make([]byte, l)
		for j := range b {
			if j == 0 {
				b[j] = byte(rand.Intn(9)+1) + '0'
			} else {
				b[j] = byte(rand.Intn(10)) + '0'
			}
		}
		tests[i] = testD{s: string(b)}
	}
	tests = append(tests, testD{s: "4"})
	tests = append(tests, testD{s: "11"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := fmt.Sprintf("%s\n", t.s)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		expected := solveD(t.s)
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("Test %d failed\nInput:%sExpected %s Got %s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
