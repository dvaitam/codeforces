package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type pairDP struct {
	l, r uint64
	mem  [61][2][2][2][2]int8
}

func b2u(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func (p *pairDP) dfs(pos int, ul, ur, vl, vr uint8) bool {
	if pos < 0 {
		return true
	}
	if p.mem[pos][ul][ur][vl][vr] != 0 {
		return p.mem[pos][ul][ur][vl][vr] == 1
	}
	lb := (p.l >> pos) & 1
	rb := (p.r >> pos) & 1
	for bu := uint64(0); bu <= 1; bu++ {
		for bv := uint64(0); bv <= 1; bv++ {
			if bu == 1 && bv == 1 {
				continue
			}
			if ul == 1 && bu < lb {
				continue
			}
			if ur == 1 && bu > rb {
				continue
			}
			if vl == 1 && bv < lb {
				continue
			}
			if vr == 1 && bv > rb {
				continue
			}
			nul := ul == 1 && bu == lb
			nur := ur == 1 && bu == rb
			nvl := vl == 1 && bv == lb
			nvr := vr == 1 && bv == rb
			if p.dfs(pos-1, b2u(nul), b2u(nur), b2u(nvl), b2u(nvr)) {
				p.mem[pos][ul][ur][vl][vr] = 1
				return true
			}
		}
	}
	p.mem[pos][ul][ur][vl][vr] = -1
	return false
}

func (p *pairDP) buildU() (uint64, bool) {
	if !p.dfs(60, 1, 1, 1, 1) {
		return 0, false
	}
	var u uint64
	ul, ur, vl, vr := uint8(1), uint8(1), uint8(1), uint8(1)
	for pos := 60; pos >= 0; pos-- {
		lb := (p.l >> pos) & 1
		rb := (p.r >> pos) & 1
		chosen := uint64(2)
		var nUL, nUR, nVL, nVR uint8
		for bu := uint64(0); bu <= 1; bu++ {
			if ul == 1 && bu < lb {
				continue
			}
			if ur == 1 && bu > rb {
				continue
			}
			for bv := uint64(0); bv <= 1; bv++ {
				if bu == 1 && bv == 1 {
					continue
				}
				if vl == 1 && bv < lb {
					continue
				}
				if vr == 1 && bv > rb {
					continue
				}
				nul := ul == 1 && bu == lb
				nur := ur == 1 && bu == rb
				nvl := vl == 1 && bv == lb
				nvr := vr == 1 && bv == rb
				if p.dfs(pos-1, b2u(nul), b2u(nur), b2u(nvl), b2u(nvr)) {
					chosen = bu
					nUL, nUR, nVL, nVR = b2u(nul), b2u(nur), b2u(nvl), b2u(nvr)
					break
				}
			}
			if chosen != 2 {
				break
			}
		}
		if chosen == 2 {
			return 0, false
		}
		u |= chosen << pos
		ul, ur, vl, vr = nUL, nUR, nVL, nVR
	}
	return u, true
}

type singleDP struct {
	l, r, u uint64
	mem     [61][2][2]int8
}

func (s *singleDP) dfs(pos int, tl, tr uint8) bool {
	if pos < 0 {
		return true
	}
	if s.mem[pos][tl][tr] != 0 {
		return s.mem[pos][tl][tr] == 1
	}
	lb := (s.l >> pos) & 1
	rb := (s.r >> pos) & 1
	allow1 := ((s.u >> pos) & 1) == 0
	for bv := uint64(0); bv <= 1; bv++ {
		if bv == 1 && !allow1 {
			continue
		}
		if tl == 1 && bv < lb {
			continue
		}
		if tr == 1 && bv > rb {
			continue
		}
		ntl := tl == 1 && bv == lb
		ntr := tr == 1 && bv == rb
		if s.dfs(pos-1, b2u(ntl), b2u(ntr)) {
			s.mem[pos][tl][tr] = 1
			return true
		}
	}
	s.mem[pos][tl][tr] = -1
	return false
}

func (s *singleDP) buildV() (uint64, bool) {
	if !s.dfs(60, 1, 1) {
		return 0, false
	}
	var v uint64
	tl, tr := uint8(1), uint8(1)
	for pos := 60; pos >= 0; pos-- {
		lb := (s.l >> pos) & 1
		rb := (s.r >> pos) & 1
		allow1 := ((s.u >> pos) & 1) == 0
		chosen := uint64(2)
		var nTL, nTR uint8
		for bv := uint64(0); bv <= 1; bv++ {
			if bv == 1 && !allow1 {
				continue
			}
			if tl == 1 && bv < lb {
				continue
			}
			if tr == 1 && bv > rb {
				continue
			}
			ntl := tl == 1 && bv == lb
			ntr := tr == 1 && bv == rb
			if s.dfs(pos-1, b2u(ntl), b2u(ntr)) {
				chosen = bv
				nTL, nTR = b2u(ntl), b2u(ntr)
				break
			}
		}
		if chosen == 2 {
			return 0, false
		}
		v |= chosen << pos
		tl, tr = nTL, nTR
	}
	return v, true
}

func expectedAnswer(n, l, r, k uint64) (uint64, bool) {
	if n%2 == 1 {
		return l, true
	}
	if n == 2 {
		return 0, false
	}
	pdp := pairDP{l: l, r: r}
	u, ok := pdp.buildU()
	if !ok {
		return 0, false
	}
	sdp := singleDP{l: l, r: r, u: u}
	v, ok := sdp.buildV()
	if !ok {
		return 0, false
	}
	if k <= n-2 {
		return u, true
	}
	return v, true
}

type test struct {
	raw string
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(2119))
	var tests []test

	// Include the statement sample.
	sample := "9\n1 4 4 1\n3 1 3 3\n4 6 9 2\n4 6 9 3\n4 6 7 4\n2 5 5 1\n2 3 6 2\n999999999999999999 1000000000000000000 1000000000000000000 999999999999999999\n1000000000000000000 1 999999999999999999 1000000000000000000\n"
	tests = append(tests, test{raw: sample})

	for len(tests) < 120 {
		t := rng.Intn(4) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for i := 0; i < t; i++ {
			n := uint64(rng.Intn(10) + 1)
			l := uint64(rng.Intn(20) + 1)
			r := l + uint64(rng.Intn(20))
			k := uint64(rng.Intn(int(n)) + 1)
			// occasionally stress large values
			if rng.Intn(10) == 0 {
				n = uint64(rng.Intn(5)+1) * 1_000_000_000_000_000_00 // up to 5e17
				if n == 0 {
					n = 1
				}
				l = uint64(rng.Intn(1_000)) + 1
				r = l + uint64(rng.Intn(1_000))
				k = n
			}
			fmt.Fprintf(&sb, "%d %d %d %d\n", n, l, r, k)
		}
		tests = append(tests, test{raw: sb.String()})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseAndJudge(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return fmt.Errorf("failed to read t")
	}
	expected := make([]uint64, t)
	exists := make([]bool, t)
	for i := 0; i < t; i++ {
		var n, l, r, k uint64
		fmt.Fscan(in, &n, &l, &r, &k)
		val, ok := expectedAnswer(n, l, r, k)
		expected[i] = val
		exists[i] = ok
	}

	out := bufio.NewReader(strings.NewReader(output))
	for i := 0; i < t; i++ {
		var gotStr string
		if _, err := fmt.Fscan(out, &gotStr); err != nil {
			return fmt.Errorf("missing output for test case %d", i+1)
		}
		if !exists[i] {
			if gotStr != "-1" {
				return fmt.Errorf("case %d: expected -1 got %s", i+1, gotStr)
			}
			continue
		}
		var got uint64
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			return fmt.Errorf("case %d: invalid number %s", i+1, gotStr)
		}
		if got != expected[i] {
			return fmt.Errorf("case %d: expected %d got %d", i+1, expected[i], got)
		}
	}
	var extra string
	if _, err := fmt.Fscan(out, &extra); err == nil {
		return fmt.Errorf("extra output detected")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(bin, tc.raw)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := parseAndJudge(tc.raw, got); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, tc.raw, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
