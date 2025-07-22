package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const (
	maxg = 303
	INF  = int(0x3f3f3f3f)
)

type pair struct {
	first  int
	second int
}

type block struct {
	ans      int64
	pans     int
	rightcnt int
	nextrcnt int
	ld, ldx  int
}

var (
	ab      []pair
	have    []bool
	blocks  []block
	userByB [][]int
)

func (b *block) getLeader(from, to int) {
	ourcnt := 0
	b.ans = 0
	b.pans = 0
	for i := to - 1; i >= from; i-- {
		if !have[i] {
			continue
		}
		ourcnt++
		cur := int64(b.rightcnt+ourcnt) * int64(ab[i].first)
		if cur > b.ans {
			b.ans = cur
			b.pans = ab[i].first
			b.ld = i
			b.ldx = b.rightcnt + ourcnt
		}
	}
}

func (b *block) recalc(from, to int) {
	b.getLeader(from, to)
	ourcnt := 0
	b.nextrcnt = INF
	for i := to - 1; i >= from; i-- {
		if !have[i] {
			continue
		}
		ourcnt++
		if i == b.ld || ab[i].first <= ab[b.ld].first {
			continue
		}
		curx := b.rightcnt + ourcnt
		lp := int64(ab[b.ld].first)*int64(b.ldx) - int64(ab[i].first)*int64(curx)
		rp := int64(ab[i].first) - int64(ab[b.ld].first)
		t := int(lp/rp) + 1
		if t < INF {
			nxt := b.rightcnt + t
			if nxt < b.nextrcnt {
				b.nextrcnt = nxt
			}
		}
	}
}

func (b *block) update(from, to int) {
	b.rightcnt++
	b.ans += int64(b.pans)
	if b.rightcnt >= b.nextrcnt {
		b.recalc(from, to)
	}
}

func expected(n int, w int64, arr []pair) string {
	ab = make([]pair, n)
	copy(ab, arr)
	sort.Slice(ab, func(i, j int) bool { return ab[i].first < ab[j].first })
	mb := 0
	for _, p := range ab {
		if p.second > mb {
			mb = p.second
		}
	}
	userByB = make([][]int, mb+1)
	for i := 0; i < n; i++ {
		b := ab[i].second
		userByB[b] = append(userByB[b], i)
	}
	have = make([]bool, n)
	blockCount := (n + maxg - 1) / maxg
	blocks = make([]block, blockCount)
	for i := range blocks {
		blocks[i].nextrcnt = INF
	}
	watchAds := n
	var sb strings.Builder
	for c := 0; c <= mb+1; c++ {
		if c > 0 {
			for _, u := range userByB[c-1] {
				watchAds--
				if ab[u].first == 0 {
					continue
				}
				have[u] = true
				gid := u / maxg
				from := gid * maxg
				to := from + maxg
				if to > n {
					to = n
				}
				blocks[gid].recalc(from, to)
				for i := 0; i < gid; i++ {
					f := i * maxg
					t := f + maxg
					if t > n {
						t = n
					}
					blocks[i].update(f, t)
				}
			}
		}
		bestAns := int64(0)
		bestPans := 0
		for i := range blocks {
			if blocks[i].ans > bestAns {
				bestAns = blocks[i].ans
				bestPans = blocks[i].pans
			}
		}
		total := bestAns + w*int64(c)*int64(watchAds)
		sb.WriteString(fmt.Sprintf("%d %d\n", total, bestPans))
	}
	return strings.TrimSpace(sb.String())
}

func runCase(bin, input, expected string) error {
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		w := int64(rng.Intn(5) + 1)
		arr := make([]pair, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, w))
		for j := 0; j < n; j++ {
			a := rng.Intn(5) + 1
			b := rng.Intn(5) + 1
			arr[j] = pair{a, b}
			sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
		}
		input := sb.String()
		exp := expected(n, w, arr)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
