package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type segment struct {
	l, r int
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func laminar(a, b segment) bool {
	return a.r < b.l || b.r < a.l || (a.l <= b.l && b.r <= a.r) || (b.l <= a.l && a.r <= b.r)
}

func addIfGood(segs *[]segment, cand segment) bool {
	if cand.l > cand.r {
		cand.l, cand.r = cand.r, cand.l
	}
	for _, s := range *segs {
		if s.l == cand.l && s.r == cand.r {
			return false
		}
		if !laminar(s, cand) {
			return false
		}
	}
	*segs = append(*segs, cand)
	return true
}

func genDisjoint() (int, []segment) {
	segs := make([]segment, 0)
	pos := 1
	cnt := rand.Intn(4) + 1
	for i := 0; i < cnt; i++ {
		pos += rand.Intn(2)
		length := rand.Intn(4) + 1
		l := pos
		r := pos + length - 1
		segs = append(segs, segment{l, r})
		pos = r + rand.Intn(3) + 1
	}
	n := pos + rand.Intn(3)
	return n, segs
}

func genChain() (int, []segment) {
	n := rand.Intn(40) + 5
	k := rand.Intn(8) + 1
	segs := make([]segment, 0, k+1)
	curL, curR := 1, n
	if rand.Intn(2) == 0 {
		segs = append(segs, segment{curL, curR})
	}
	for i := 0; i < k; i++ {
		if curL < curR {
			curL += rand.Intn(2)
		}
		if curL < curR {
			curR -= rand.Intn(2)
		}
		if curL > curR {
			curR = curL
		}
		if addIfGood(&segs, segment{curL, curR}) {
			continue
		}
	}
	return n, segs
}

func genMixed() (int, []segment) {
	rootLen := rand.Intn(30) + 15
	pos := 1
	topCnt := rand.Intn(4) + 1
	segs := make([]segment, 0)
	top := make([]segment, 0, topCnt)
	for i := 0; i < topCnt; i++ {
		pos += rand.Intn(2)
		lenHere := rand.Intn(6) + 2
		l, r := pos, pos+lenHere-1
		top = append(top, segment{l, r})
		segs = append(segs, segment{l, r})
		pos = r + rand.Intn(4) + 1
	}
	n := pos + rand.Intn(5)
	if n < rootLen {
		n = rootLen
	}
	if rand.Intn(2) == 0 {
		addIfGood(&segs, segment{1, n})
	}
	for _, p := range top {
		curL, curR := p.l, p.r
		inner := rand.Intn(4)
		for j := 0; j < inner; j++ {
			if curL < curR {
				curL += rand.Intn(2)
			}
			if curL < curR {
				curR -= rand.Intn(2)
			}
			if curL > curR {
				curR = curL
			}
			addIfGood(&segs, segment{curL, curR})
		}
	}
	return n, segs
}

func genRandom() (int, []segment) {
	n := rand.Intn(180) + 20
	target := rand.Intn(30)
	segs := make([]segment, 0, target+2)
	if rand.Intn(3) == 0 {
		addIfGood(&segs, segment{1, n})
	}
	tries := 0
	for len(segs) < target && tries < target*40+50 {
		l := rand.Intn(n) + 1
		r := l + rand.Intn(n-l+1)
		addIfGood(&segs, segment{l, r})
		tries++
	}
	return n, segs
}

func genCase() (int, []segment) {
	switch rand.Intn(4) {
	case 0:
		n := rand.Intn(6) + 1
		return n, nil
	case 1:
		return genChain()
	case 2:
		return genDisjoint()
	default:
		if rand.Intn(2) == 0 {
			return genMixed()
		}
		return genRandom()
	}
}

func genTest() []byte {
	t := rand.Intn(6) + 1
	// Occasionally insert a very large n with small m to probe limits.
	includeBig := rand.Intn(5) == 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t+boolToInt(includeBig)))
	for i := 0; i < t; i++ {
		n, segs := genCase()
		sb.WriteString(fmt.Sprintf("%d %d\n", n, len(segs)))
		for _, s := range segs {
			sb.WriteString(fmt.Sprintf("%d %d\n", s.l, s.r))
		}
	}
	if includeBig {
		n := 200000
		segs := []segment{}
		if rand.Intn(2) == 0 {
			segs = append(segs, segment{1, n})
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, len(segs)))
		for _, s := range segs {
			sb.WriteString(fmt.Sprintf("%d %d\n", s.l, s.r))
		}
	}
	return []byte(sb.String())
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refE.bin"
	if err := exec.Command("go", "build", "-o", ref, "2056E.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 200; i++ {
		input := genTest()
		want, err := run(ref, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
