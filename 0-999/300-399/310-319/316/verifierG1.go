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

type SAM struct {
	next []map[byte]int
	link []int
	len  []int
	cnt  []int64
	last int
	sz   int
}

func NewSAM(maxLen int) *SAM {
	size := 2*maxLen + 1
	sam := &SAM{
		next: make([]map[byte]int, size),
		link: make([]int, size),
		len:  make([]int, size),
		cnt:  make([]int64, size),
		last: 0,
		sz:   1,
	}
	sam.next[0] = make(map[byte]int)
	sam.link[0] = -1
	return sam
}

func (sam *SAM) Extend(c byte) {
	p := sam.last
	cur := sam.sz
	sam.sz++
	sam.len[cur] = sam.len[p] + 1
	sam.next[cur] = make(map[byte]int)
	sam.cnt[cur] = 1
	for p >= 0 && sam.next[p][c] == 0 {
		sam.next[p][c] = cur
		p = sam.link[p]
	}
	if p == -1 {
		sam.link[cur] = 0
	} else {
		q := sam.next[p][c]
		if sam.len[p]+1 == sam.len[q] {
			sam.link[cur] = q
		} else {
			clone := sam.sz
			sam.sz++
			sam.len[clone] = sam.len[p] + 1
			sam.next[clone] = make(map[byte]int)
			for k, v := range sam.next[q] {
				sam.next[clone][k] = v
			}
			sam.link[clone] = sam.link[q]
			sam.cnt[clone] = 0
			for p >= 0 && sam.next[p][c] == q {
				sam.next[p][c] = clone
				p = sam.link[p]
			}
			sam.link[q] = clone
			sam.link[cur] = clone
		}
	}
	sam.last = cur
}

func BuildSAM(s string) *SAM {
	sam := NewSAM(len(s))
	for i := 0; i < len(s); i++ {
		sam.Extend(s[i])
	}
	maxLen := 0
	for i := 0; i < sam.sz; i++ {
		if sam.len[i] > maxLen {
			maxLen = sam.len[i]
		}
	}
	bucket := make([]int, maxLen+1)
	for i := 0; i < sam.sz; i++ {
		bucket[sam.len[i]]++
	}
	for i := 1; i <= maxLen; i++ {
		bucket[i] += bucket[i-1]
	}
	order := make([]int, sam.sz)
	for i := sam.sz - 1; i >= 0; i-- {
		l := sam.len[i]
		bucket[l]--
		order[bucket[l]] = i
	}
	for i := sam.sz - 1; i > 0; i-- {
		v := order[i]
		p := sam.link[v]
		if p >= 0 {
			sam.cnt[p] += sam.cnt[v]
		}
	}
	return sam
}

type Rule struct {
	sam *SAM
	l   int64
	r   int64
}

func solveCase(s string, rules []Rule) string {
	good := make(map[string]struct{})
	for i := 0; i < len(s); i++ {
		cur := make([]int, len(rules))
		for k := range cur {
			cur[k] = 0
		}
		for j := i; j < len(s); j++ {
			c := s[j]
			ok := true
			for k := 0; k < len(rules); k++ {
				if cur[k] >= 0 {
					nxt, has := rules[k].sam.next[cur[k]][c]
					if !has {
						cur[k] = -1
					} else {
						cur[k] = nxt
					}
				}
				var occ int64
				if cur[k] >= 0 {
					occ = rules[k].sam.cnt[cur[k]]
				} else {
					occ = 0
				}
				if occ < rules[k].l || occ > rules[k].r {
					ok = false
					break
				}
			}
			if ok {
				sub := s[i : j+1]
				good[sub] = struct{}{}
			}
		}
	}
	return fmt.Sprintf("%d\n", len(good))
}

func genCase(rng *rand.Rand) (string, string) {
	L := rng.Intn(5) + 1
	var letters = []byte("abc")
	bs := make([]byte, L)
	for i := 0; i < L; i++ {
		bs[i] = letters[rng.Intn(len(letters))]
	}
	s := string(bs)
	n := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s %d\n", s, n)
	rules := make([]Rule, n)
	for i := 0; i < n; i++ {
		lp := rng.Intn(3) + 1
		pb := make([]byte, lp)
		for j := 0; j < lp; j++ {
			pb[j] = letters[rng.Intn(len(letters))]
		}
		p := string(pb)
		l := rng.Intn(3)
		r := l + rng.Intn(3)
		fmt.Fprintf(&sb, "%s %d %d\n", p, l, r)
		rules[i] = Rule{sam: BuildSAM(p), l: int64(l), r: int64(r)}
	}
	expected := solveCase(s, rules)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG1.go /path/to/binary")
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
