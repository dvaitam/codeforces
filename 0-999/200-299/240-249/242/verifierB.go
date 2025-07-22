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

type segment struct{ l, r int64 }

func solveCase(segs []segment) string {
	minL := int64(1<<63 - 1)
	maxR := int64(-1 << 63)
	for _, s := range segs {
		if s.l < minL {
			minL = s.l
		}
		if s.r > maxR {
			maxR = s.r
		}
	}
	res := -1
	for i, s := range segs {
		if s.l == minL && s.r == maxR {
			res = i + 1
			break
		}
	}
	return fmt.Sprintf("%d\n", res)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	segs := make([]segment, n)
	cover := rng.Intn(2) == 0
	var coverL, coverR int64
	if cover {
		coverL = int64(rng.Intn(20) + 1)
		coverR = coverL + int64(rng.Intn(10)+1)
		segs[0] = segment{coverL, coverR}
		for i := 1; i < n; i++ {
			l := coverL + int64(rng.Intn(int(coverR-coverL+1)))
			r := l + int64(rng.Intn(int(coverR-l+1)))
			segs[i] = segment{l, r}
		}
	} else {
		used := make(map[[2]int64]bool)
		for i := 0; i < n; i++ {
			for {
				l := int64(rng.Intn(30) + 1)
				r := l + int64(rng.Intn(10))
				if !used[[2]int64{l, r}] {
					segs[i] = segment{l, r}
					used[[2]int64{l, r}] = true
					break
				}
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, s := range segs {
		fmt.Fprintf(&sb, "%d %d\n", s.l, s.r)
	}
	return sb.String(), solveCase(segs)
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
