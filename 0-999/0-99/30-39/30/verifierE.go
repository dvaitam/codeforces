package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func isPalindrome(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

type segment struct{ pos, len int }

func bruteForce(s string) (int, []segment) {
	n := len(s)
	best := 0
	bestK := 1
	bestSeg := []segment{{1, 1}}
	// k=1
	for i := 0; i < n; i++ {
		for l := 1; i+l <= n; l++ {
			if l%2 == 1 && isPalindrome(s[i:i+l]) {
				if l > best {
					best = l
					bestK = 1
					bestSeg = []segment{{i + 1, l}}
				}
			}
		}
	}
	// k=3
	for i1 := 0; i1 < n; i1++ {
		for l1 := 1; i1+l1 <= n; l1++ {
			prefix := s[i1 : i1+l1]
			for i3 := i1 + l1; i3+l1 <= n; i3++ {
				if s[i3:i3+l1] != prefix {
					continue
				}
				for i2 := i1 + l1; i2 < i3; i2++ {
					for l2 := 1; i2+l2 <= i3; l2++ {
						if l2%2 == 1 && isPalindrome(s[i2:i2+l2]) {
							total := 2*l1 + l2
							if total > best {
								best = total
								bestK = 3
								bestSeg = []segment{{i1 + 1, l1}, {i2 + 1, l2}, {i3 + 1, l1}}
							}
						}
					}
				}
			}
		}
	}
	return bestK, bestSeg
}

func segmentsValid(s string, segs []segment) bool {
	if len(segs) == 1 {
		if segs[0].len%2 == 1 && isPalindrome(s[segs[0].pos-1:segs[0].pos-1+segs[0].len]) {
			return true
		}
		return false
	}
	if len(segs) != 3 {
		return false
	}
	a := segs[0]
	b := segs[1]
	c := segs[2]
	if !(a.pos+a.len-1 <= b.pos-1 && b.pos+b.len-1 <= c.pos-1) {
		return false
	}
	pre := s[a.pos-1 : a.pos-1+a.len]
	suf := s[c.pos-1 : c.pos-1+c.len]
	if pre != suf {
		return false
	}
	if b.len%2 == 0 || !isPalindrome(s[b.pos-1:b.pos-1+b.len]) {
		return false
	}
	return true
}

func generateCase(rng *rand.Rand) (string, string, int, []segment) {
	n := rng.Intn(8) + 1
	bytesS := make([]byte, n)
	for i := 0; i < n; i++ {
		bytesS[i] = byte('a' + rng.Intn(3))
	}
	s := string(bytesS)
	k, segs := bruteForce(s)
	var sb strings.Builder
	sb.WriteString(s + "\n")
	var exp strings.Builder
	exp.WriteString(fmt.Sprint(k, "\n"))
	for _, seg := range segs {
		exp.WriteString(fmt.Sprintf("%d %d\n", seg.pos, seg.len))
	}
	return sb.String(), strings.TrimSpace(exp.String()), k, segs
}

func runCase(bin, input, expected string, s string, segs []segment) error {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(outLines) == 0 {
		return fmt.Errorf("no output")
	}
	kOut, err := strconv.Atoi(strings.TrimSpace(outLines[0]))
	if err != nil {
		return fmt.Errorf("invalid k")
	}
	var resSegs []segment
	for i := 1; i < len(outLines); i++ {
		var x, l int
		fmt.Sscanf(outLines[i], "%d %d", &x, &l)
		resSegs = append(resSegs, segment{x, l})
	}
	if kOut != len(resSegs) {
		return fmt.Errorf("k mismatch")
	}
	if !segmentsValid(s, resSegs) {
		return fmt.Errorf("invalid segments")
	}
	// check length optimality
	_, best := bruteForce(s)
	bestLen := 0
	for _, sg := range best {
		bestLen += sg.len
	}
	outLen := 0
	for _, sg := range resSegs {
		outLen += sg.len
	}
	if outLen != bestLen {
		return fmt.Errorf("expected length %d got %d", bestLen, outLen)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp, _, segs := generateCase(rng)
		s := strings.TrimSpace(strings.Split(in, "\n")[0])
		if err := runCase(bin, in, exp, s, segs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
