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

type wave struct {
	l int64
	r int64
	a int64
}

type testCaseF struct {
	k     int64
	waves []wave
}

func solveCaseF(k int64, waves []wave) int64 {
	var sumA int64
	for _, w := range waves {
		sumA += w.a
	}
	dp := map[int64]int64{k: 0}
	for i, w := range waves {
		T := w.r - w.l + 1
		newdp := make(map[int64]int64)
		var gap int64
		if i > 0 {
			gap = w.l - waves[i-1].r
		}
		try := func(b0, cost0 int64) {
			var events int64
			if b0 > 0 {
				rem := w.a - b0
				if rem < 0 {
					rem = 0
				}
				ev := (rem + k - 1) / k
				events = 1 + ev
			} else {
				events = (w.a + k - 1) / k
			}
			if events > T {
				return
			}
			var remNew int64
			if b0 >= w.a {
				remNew = b0 - w.a
			} else {
				remm := (w.a - b0) % k
				if remm == 0 {
					remNew = 0
				} else {
					remNew = k - remm
				}
			}
			if prev, ok := newdp[remNew]; !ok || cost0 < prev {
				newdp[remNew] = cost0
			}
		}
		for remPrev, costPrev := range dp {
			try(remPrev, costPrev)
			if gap > 0 {
				try(k, costPrev+remPrev)
			}
		}
		dp = newdp
		if len(dp) == 0 {
			return -1
		}
	}
	var best int64 = -1
	for _, cost := range dp {
		if best < 0 || cost < best {
			best = cost
		}
	}
	return sumA + best
}

func buildInputF(tc testCaseF) string {
	var sb strings.Builder
	n := len(tc.waves)
	fmt.Fprintf(&sb, "%d %d\n", n, tc.k)
	for _, w := range tc.waves {
		fmt.Fprintf(&sb, "%d %d %d\n", w.l, w.r, w.a)
	}
	return sb.String()
}

func runCaseF(bin string, tc testCaseF) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputF(tc))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := solveCaseF(tc.k, tc.waves)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func generateCasesF() []testCaseF {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseF, 0, 100)
	cases = append(cases, testCaseF{k: 1, waves: []wave{{1, 1, 1}}})
	for len(cases) < 100 {
		n := rng.Intn(5) + 1
		k := int64(rng.Intn(5) + 1)
		waves := make([]wave, n)
		cur := int64(1)
		for i := 0; i < n; i++ {
			lenWave := int64(rng.Intn(3) + 1)
			waves[i].l = cur
			waves[i].r = cur + lenWave - 1
			cur = waves[i].r + int64(rng.Intn(2)) + 1
			waves[i].a = int64(rng.Intn(5) + 1)
		}
		cases = append(cases, testCaseF{k: k, waves: waves})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesF()
	for i, tc := range cases {
		if err := runCaseF(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
