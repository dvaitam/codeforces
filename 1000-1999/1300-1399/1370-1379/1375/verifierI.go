package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesI.txt to remove external dependency.
const testcasesRaw = `3 -2 0 2 -3 -3 3 1 -3 -1
1 1 -2 -3
1 0 0 -3
2 -3 1 0 -3 3 1
1 -2 2 2
1 1 1 0
1 -2 -3 1
2 -1 0 -2 1 -3 1
3 1 3 2 -2 -3 1 1 2 -2
3 -3 1 2 -3 1 -3 1 -2 0
4 3 -1 0 1 0 -1 -1 -2 3 -2 2 3
2 -3 1 -1 1 0 -1
4 -1 1 -3 -3 1 0 -2 3 -1 -2 0 0
1 2 -3 3
3 -1 2 -1 1 0 1 3 0 -3
1 -1 0 2
1 -3 2 2
3 2 1 2 3 0 -1 2 0 2
3 -3 0 -1 -2 1 -3 0 -3 -2
3 -2 2 -2 0 0 3 0 -3 -2
4 0 1 -1 -2 3 0 3 1 -1 2 0 -1
4 -2 -2 -3 -2 -2 -2 2 -2 -3 0 3 1
2 -1 -1 -3 -2 0 1
3 1 1 -1 -2 2 3 1 1 2
1 0 3 3
4 0 0 0 -3 0 2 0 -3 -2 -3 -2 0
2 -3 -1 1 -3 -3 -3
2 1 -3 -1 1 -3 -3
2 1 0 -2 2 -1 -1
3 0 -3 -3 3 0 0 0 0 -1
1 -2 -3 2
3 2 -1 0 3 2 -2 1 -3 -2
3 -2 2 1 -3 3 1 -1 2 3
1 2 3 -1
3 -2 -1 3 -2 1 1 3 1 -1
2 1 3 3 3 3 -2
2 3 0 2 3 -2 -2
4 -1 2 -3 -3 3 -1 0 -1 -2 2 1 -1
4 3 2 -1 -1 -3 -2 -3 -2 0 -2 -1 -2
4 1 1 3 -3 0 2 -1 3 2 -3 3 2
1 0 3 2
2 0 -2 0 3 2 -1
1 3 2 0
4 0 2 -3 2 -2 -2 -2 -3 -2 1 0 3
2 1 3 1 0 2 -1
2 1 1 -2 -3 -3 3
1 1 2 -2
4 3 -2 3 3 -2 -3 -1 -2 -1 1 -2 3
3 -1 1 0 3 -2 -3 2 -1 0
4 3 1 -2 1 -2 1 1 -3 3 0 3 -2
1 3 3 -2
2 -2 0 1 2 -3 1
1 -1 2 1
4 3 3 -3 1 -3 -2 -2 -1 -3 3 -3 1
4 1 -3 3 -3 0 -1 1 1 1 1 -2 2
3 0 1 1 3 0 1 -2 2 1
3 1 -2 3 0 -2 0 -3 0 0
3 -3 2 -2 0 -3 -2 2 -1 3
1 3 -2 2
3 -2 -1 -2 0 -2 2 -3 0 0
2 2 3 -2 -2 2 0
4 -1 0 -2 -1 -1 -3 2 -1 -3 -1 1 0
4 2 -3 0 -1 1 1 -1 1 -3 -3 3 -2
1 -3 -1 -1
1 3 -2 -1
2 3 0 3 2 3 -1
4 -2 1 1 1 0 2 -1 -3 -1 -3 3 2
2 0 -3 -1 -3 2 -3
3 -3 1 3 -2 -3 -1 3 -3 0
1 -1 1 0
3 1 -2 -3 1 2 -2 -3 -2 -1
1 -2 -2 -1
3 1 3 -2 -1 0 1 2 -2 -1
3 3 -3 -1 -3 -3 -3 2 1 1
2 1 0 -2 0 -3 2
4 2 0 1 3 0 1 -1 2 -2 -2 -1 -2
2 0 -1 -3 3 -2 -3
1 2 2 -1
4 -2 -3 -3 2 3 0 3 1 2 -1 1 -2
3 -3 0 -2 -2 -1 0 -3 -1 -1
3 1 -1 -2 -3 -1 -2 -1 -2 -3
3 0 -3 0 -1 1 2 -2 -2 1
1 -3 -1 3
1 -2 0 1
1 0 -3 -1
3 2 -2 -3 1 1 3 3 -2 2
4 3 -1 2 0 -2 -1 2 1 2 -2 -3 3
4 2 2 3 1 -2 1 3 1 1 3 3 3
1 3 2 1
2 -3 -3 -3 -2 2 -1
1 0 3 0
1 2 -3 2
2 0 -1 -3 0 3 -3
1 2 1 -3
4 -1 3 -3 3 -1 -2 2 3 -2 -2 2 2
4 0 3 0 -3 0 2 -1 3 -3 1 2 2
2 -3 1 -2 -1 -1 2
3 1 1 -2 -3 0 -3 0 -1 2
1 2 -2 2
4 -1 2 1 -1 0 0 0 3 -3 1 -2 -1`

type qt struct {
	s, x, y, z int64
}

func newQt(v int64) qt { return qt{2 * v, 0, 0, 0} }
func newQtVals(s, x, y, z int64) qt {
	return qt{2 * s, 2 * x, 2 * y, 2 * z}
}
func newQtDbl(s, x, y, z int64) qt { return qt{s, x, y, z} }

func (q qt) conj() qt { return newQtDbl(q.s, -q.x, -q.y, -q.z) }

func norm(q qt) int64 { return (q.s*q.s + q.x*q.x + q.y*q.y + q.z*q.z) >> 2 }

func (a qt) add(b qt) qt { return newQtDbl(a.s+b.s, a.x+b.x, a.y+b.y, a.z+b.z) }
func (a qt) sub(b qt) qt { return newQtDbl(a.s-b.s, a.x-b.x, a.y-b.y, a.z-b.z) }

func (q qt) mulScalar(c int64) qt { return newQtDbl(q.s*c, q.x*c, q.y*c, q.z*c) }

func (a qt) mul(b qt) qt {
	return newQtDbl(
		(a.s*b.s-a.x*b.x-a.y*b.y-a.z*b.z)>>1,
		(a.s*b.x+a.x*b.s+a.y*b.z-a.z*b.y)>>1,
		(a.s*b.y+a.y*b.s+a.z*b.x-a.x*b.z)>>1,
		(a.s*b.z+a.z*b.s+a.x*b.y-a.y*b.x)>>1,
	)
}

func floorDiv(u, v int64) int64 {
	d := u / v
	r := u % v
	if (u^v) < 0 && r != 0 {
		d--
	}
	return d
}

func rightDiv(a, b qt) (qt, qt) {
	numer := b.conj().mul(a)
	den := norm(b)
	s := floorDiv(numer.s, den)
	x := floorDiv(numer.x, den)
	y := floorDiv(numer.y, den)
	z := floorDiv(numer.z, den)
	q1 := newQtDbl(s|1, x|1, y|1, z|1)
	r1 := a.sub(b.mul(q1))
	q2 := newQtDbl((s+1)&^1, (x+1)&^1, (y+1)&^1, (z+1)&^1)
	r2 := a.sub(b.mul(q2))
	if norm(r1) < norm(r2) {
		return q1, r1
	}
	return q2, r2
}

func rightGcd(a, b qt) qt {
	for a.s != 0 || a.x != 0 || a.y != 0 || a.z != 0 {
		_, r := rightDiv(b, a)
		b, a = a, r
	}
	return b
}

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveCase(n int, coords []int64) string {
	A := make([]qt, n)
	var k int64
	for i := 0; i < n; i++ {
		x := coords[3*i]
		y := coords[3*i+1]
		z := coords[3*i+2]
		A[i] = newQtVals(0, x, y, z)
		k = gcd(k, x)
		k = gcd(k, y)
		k = gcd(k, z)
	}
	for i := range A {
		A[i].x /= k
		A[i].y /= k
		A[i].z /= k
	}
	var G qt
	for _, a := range A {
		G = rightGcd(G, a)
	}
	normG := norm(G)
	var r2 int64
	for _, a0 := range A {
		a := a0
		normA := norm(a)
		a = G.conj().mul(a)
		a.s /= normG
		a.x /= normG
		a.y /= normG
		a.z /= normG
		g1 := gcd(normA, normG)
		normA /= g1
		curR2 := g1
		a = a.mul(G)
		v1 := a.x >> 1
		if v1 < 0 {
			v1 = -v1
		}
		v2 := a.y >> 1
		if v2 < 0 {
			v2 = -v2
		}
		v3 := a.z >> 1
		if v3 < 0 {
			v3 = -v3
		}
		g2 := gcd(normA, gcd(gcd(v1, v2), v3))
		curR2 *= g2
		r2 = gcd(r2, curR2)
	}
	var r int64 = 1
	nTmp := r2
	for p := int64(2); p*p <= nTmp; p++ {
		if nTmp%p == 0 {
			cnt := int64(0)
			for nTmp%p == 0 {
				cnt++
				nTmp /= p
			}
			for i := int64(0); i < cnt/2; i++ {
				r *= p
			}
		}
	}
	Q := rightGcd(G, newQt(r))
	area := k * r * k * r

	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(area))
	e1 := newQtVals(0, 1, 0, 0)
	e2 := newQtVals(0, 0, 1, 0)
	e3 := newQtVals(0, 0, 0, 1)
	for _, e := range []qt{e1, e2, e3} {
		t := Q.mul(e).mul(Q.conj()).mulScalar(k)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", t.x/2, t.y/2, t.z/2))
	}
	return strings.TrimRight(sb.String(), "\n")
}

type testCase struct {
	n      int
	coords []int64
}

func parseTests() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	var cases []testCase
	lineNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineNum++
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", lineNum, err)
		}
		expect := 1 + 3*n
		if len(fields) != expect {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", lineNum, expect, len(fields))
		}
		coords := make([]int64, 0, 3*n)
		for i := 1; i < expect; i++ {
			v, err := strconv.ParseInt(fields[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse int: %v", lineNum, err)
			}
			coords = append(coords, v)
		}
		cases = append(cases, testCase{n: n, coords: coords})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Println("failed to parse embedded tests:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 0; i < tc.n; i++ {
			x := tc.coords[3*i]
			y := tc.coords[3*i+1]
			z := tc.coords[3*i+2]
			input.WriteString(fmt.Sprintf("%d %d %d\n", x, y, z))
		}
		expected := solveCase(tc.n, tc.coords)
		got, err := runCandidate(bin, []byte(input.String()))
		if err != nil {
			fmt.Printf("Test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("Test %d failed\nexpected:\n%s\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
