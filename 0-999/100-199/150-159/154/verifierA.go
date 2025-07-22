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

const ALPH = 26

func solveA(s string, pairs []string) int {
	forbidden := make([][ALPH]bool, ALPH)
	for _, p := range pairs {
		if len(p) < 2 {
			continue
		}
		a := p[0] - 'a'
		b := p[1] - 'a'
		if a >= 0 && a < ALPH && b >= 0 && b < ALPH {
			forbidden[a][b] = true
			forbidden[b][a] = true
		}
	}
	const INF = 1000000000
	const INIT = ALPH
	dpPrev := make([]int, ALPH+1)
	dpCurr := make([]int, ALPH+1)
	for i := range dpPrev {
		dpPrev[i] = INF
	}
	dpPrev[INIT] = 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 'a' || c > 'z' {
			c = 0
		}
		ci := int(c - 'a')
		for j := range dpCurr {
			dpCurr[j] = INF
		}
		for last := 0; last <= ALPH; last++ {
			if dpPrev[last]+1 < dpCurr[last] {
				dpCurr[last] = dpPrev[last] + 1
			}
		}
		for last := 0; last <= ALPH; last++ {
			if last == INIT || !forbidden[last][ci] {
				if dpPrev[last] < dpCurr[ci] {
					dpCurr[ci] = dpPrev[last]
				}
			}
		}
		dpPrev, dpCurr = dpCurr, dpPrev
	}
	res := INF
	for _, v := range dpPrev {
		if v < res {
			res = v
		}
	}
	return res
}

func genCaseA(rng *rand.Rand) (string, []string) {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(5)))
	}
	letters := []byte{'a', 'b', 'c', 'd', 'e'}
	k := rng.Intn(3)
	pairs := make([]string, 0, k)
	used := make(map[byte]bool)
	for len(pairs) < k {
		a := letters[rng.Intn(len(letters))]
		b := letters[rng.Intn(len(letters))]
		if a == b || used[a] || used[b] {
			continue
		}
		used[a] = true
		used[b] = true
		pairs = append(pairs, string([]byte{a, b}))
	}
	return sb.String(), pairs
}

func runCaseA(bin string, s string, pairs []string) error {
	var input bytes.Buffer
	fmt.Fprintln(&input, s)
	fmt.Fprintln(&input, len(pairs))
	for _, p := range pairs {
		fmt.Fprintln(&input, p)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	expected := fmt.Sprintf("%d", solveA(s, pairs))
	got := strings.TrimSpace(string(out))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		s, pairs := genCaseA(rng)
		if err := runCaseA(bin, s, pairs); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s %v\n", t+1, err, s, pairs)
			return
		}
	}
	fmt.Println("All tests passed")
}
