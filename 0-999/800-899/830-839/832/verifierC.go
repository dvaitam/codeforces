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

type Person struct {
	x int
	v int
	t int
}

type Test struct {
	n   int
	s   int
	ppl []Person
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.s))
	for _, p := range tc.ppl {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", p.x, p.v, p.t))
	}
	return sb.String()
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// --- embedded solver ---

type sLine struct {
	m, c float64
}

func (l sLine) eval(x int) float64 {
	return l.m*float64(x) + l.c
}

type sLCTNode struct {
	line  sLine
	left  int
	right int
}

var sLCT []sLCTNode

func sAddNode() int {
	sLCT = append(sLCT, sLCTNode{line: sLine{0, 1e18}, left: -1, right: -1})
	return len(sLCT) - 1
}

func sInsert(node int, l, r int, newLine sLine) {
	mid := l + (r-l)/2
	betterLeft := newLine.eval(l) < sLCT[node].line.eval(l)
	betterMid := newLine.eval(mid) < sLCT[node].line.eval(mid)

	if betterMid {
		sLCT[node].line, newLine = newLine, sLCT[node].line
	}

	if l == r {
		return
	}

	if betterLeft != betterMid {
		if sLCT[node].left == -1 {
			sLCT[node].left = sAddNode()
		}
		sInsert(sLCT[node].left, l, mid, newLine)
	} else {
		if sLCT[node].right == -1 {
			sLCT[node].right = sAddNode()
		}
		sInsert(sLCT[node].right, mid+1, r, newLine)
	}
}

func sQuery(node int, l, r int, x int) float64 {
	if node == -1 {
		return 1e18
	}
	res := sLCT[node].line.eval(x)
	if l == r {
		return res
	}
	mid := l + (r-l)/2
	if x <= mid {
		leftRes := sQuery(sLCT[node].left, l, mid, x)
		if leftRes < res {
			res = leftRes
		}
	} else {
		rightRes := sQuery(sLCT[node].right, mid+1, r, x)
		if rightRes < res {
			res = rightRes
		}
	}
	return res
}

func solve(tc Test) string {
	n := tc.n
	s := float64(tc.s)

	const X = 1000000

	CL := 1e18
	CR := 1e18

	leftByX := make([][]int, X+1)
	rightByX := make([][]int, X+1)

	mL := make([]float64, n)
	cL := make([]float64, n)
	mR := make([]float64, n)
	cR := make([]float64, n)

	for i := 0; i < n; i++ {
		x := tc.ppl[i].x
		v := float64(tc.ppl[i].v)
		dir := tc.ppl[i].t

		if dir == 1 {
			t := float64(x) / v
			if t < CL {
				CL = t
			}
			mL[i] = s / (s*s - v*v)
			cL[i] = -v * float64(x) / (s*s - v*v)
			leftByX[x] = append(leftByX[x], i)
		} else {
			t := float64(X-x) / v
			if t < CR {
				CR = t
			}
			mR[i] = -s / (s*s - v*v)
			cR[i] = (s*float64(X) - v*float64(X-x)) / (s*s - v*v)
			rightByX[x] = append(rightByX[x], i)
		}
	}

	sLCT = make([]sLCTNode, 0, 2000000)
	rootL := sAddNode()

	fL := make([]float64, X+1)
	for Y := 0; Y <= X; Y++ {
		for _, i := range leftByX[Y] {
			sInsert(rootL, 0, X, sLine{mL[i], cL[i]})
		}
		q := sQuery(rootL, 0, X, Y)
		fL[Y] = math.Min(CL, q)
	}

	sLCT = sLCT[:0]
	rootR := sAddNode()

	fR := make([]float64, X+1)
	for Y := X; Y >= 0; Y-- {
		for _, i := range rightByX[Y] {
			sInsert(rootR, 0, X, sLine{mR[i], cR[i]})
		}
		q := sQuery(rootR, 0, X, Y)
		fR[Y] = math.Min(CR, q)
	}

	ans := 1e18
	for Y := 0; Y <= X; Y++ {
		maxT := math.Max(fL[Y], fR[Y])
		if maxT < ans {
			ans = maxT
		}
	}

	return fmt.Sprintf("%.7f", ans)
}

func genTests() []Test {
	rand.Seed(time.Now().UnixNano())
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2
		s := rand.Intn(9) + 2
		ppl := make([]Person, n)
		for j := 0; j < n; j++ {
			x := rand.Intn(1_000_000-1) + 1
			v := rand.Intn(s-1) + 1
			t := rand.Intn(2) + 1
			ppl[j] = Person{x: x, v: v, t: t}
		}
		tests = append(tests, Test{n: n, s: s, ppl: ppl})
	}
	// simple small test
	tests = append(tests, Test{n: 2, s: 3, ppl: []Person{{1, 1, 1}, {999999, 1, 2}}})
	return tests
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp := solve(tc)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expVal := parseFloat(strings.TrimSpace(exp))
		gotVal := parseFloat(strings.TrimSpace(got))
		diff := math.Abs(expVal - gotVal)
		denom := math.Max(1.0, math.Abs(expVal))
		if diff/denom > 1e-6 {
			fmt.Printf("Test %d failed\nInput:%sExpected:%s\nGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
