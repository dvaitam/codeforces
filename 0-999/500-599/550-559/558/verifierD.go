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

type interval struct {
	l, r int64
}

func intersect(arr []interval, seg interval) []interval {
	res := make([]interval, 0, len(arr))
	for _, in := range arr {
		l := in.l
		if seg.l > l {
			l = seg.l
		}
		r := in.r
		if seg.r < r {
			r = seg.r
		}
		if l <= r {
			res = append(res, interval{l, r})
		}
	}
	return res
}

func subtract(arr []interval, seg interval) []interval {
	res := make([]interval, 0, len(arr)+1)
	for _, in := range arr {
		if seg.r < in.l || seg.l > in.r {
			res = append(res, in)
			continue
		}
		if seg.l > in.l {
			res = append(res, interval{in.l, seg.l - 1})
		}
		if seg.r < in.r {
			res = append(res, interval{seg.r + 1, in.r})
		}
	}
	return res
}

func expectedAnswer(h int, queries [][4]int64) string {
	allowed := []interval{{1 << (h - 1), (1 << h) - 1}}
	for _, q := range queries {
		i := q[0]
		L := q[1]
		R := q[2]
		ans := q[3]
		shift := uint(h - int(i))
		seg := interval{L << shift, ((R + 1) << shift) - 1}
		if ans == 1 {
			allowed = intersect(allowed, seg)
		} else {
			allowed = subtract(allowed, seg)
		}
	}
	if len(allowed) == 0 {
		return "Game cheated!"
	}
	if len(allowed) == 1 && allowed[0].l == allowed[0].r {
		return fmt.Sprint(allowed[0].l)
	}
	count := 0
	var val int64
	for _, in := range allowed {
		if in.l == in.r {
			count++
			val = in.l
		} else {
			return "Data not sufficient!"
		}
	}
	if count == 1 {
		return fmt.Sprint(val)
	}
	return "Data not sufficient!"
}

func generateCase(rng *rand.Rand) (string, string) {
	h := rng.Intn(10) + 1
	q := rng.Intn(10)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", h, q))
	queries := make([][4]int64, q)
	for i := 0; i < q; i++ {
		level := rng.Intn(h) + 1
		L := int64(rng.Intn(1<<level-1) + (1 << (level - 1)))
		R := int64(rng.Intn(int((1<<level)-1-int(L)+1)) + int(L))
		ans := rng.Intn(2)
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", level, L, R, ans))
		queries[i] = [4]int64{int64(level), L, R, int64(ans)}
	}
	expected := expectedAnswer(h, queries)
	return sb.String(), expected
}

func runCase(bin string, tcInput, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tcInput)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
