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

type node struct {
	x   int64
	id  int
	id2 int
}

func expected(n, m int, aVals, bVals []int64) string {
	t := make([]node, 2*n)
	for i := 1; i <= n; i++ {
		ai := aVals[i]
		bi := bVals[i]
		if bi-ai < ai {
			t[i-1] = node{bi >> 1, i, i << 1}
			t[i-1+n] = node{bi >> 1, i + n, i<<1 | 1}
		} else {
			t[i-1] = node{ai, i, i << 1}
			t[i-1+n] = node{bi - ai, i + n, i<<1 | 1}
		}
	}
	sort.Slice(t, func(i, j int) bool {
		if t[i].x != t[j].x {
			return t[i].x < t[j].x
		}
		if t[i].id2 != t[j].id2 {
			return t[i].id2 < t[j].id2
		}
		return t[i].id < t[j].id
	})
	cnt := make([]int, 2*n+2)
	var ans int64
	for i := 0; i < m; i++ {
		ans += t[i].x
		cnt[t[i].id]++
	}
	if t[m-1].id <= n && cnt[t[m-1].id+n] == 0 {
		id0 := t[m-1].id
		if t[m-1].x != aVals[id0] {
			tmp := ans - t[m-1].x
			ans = tmp + aVals[id0]
			cnt[id0]--
			C := []int{id0}
			Cval := []int{1}
			for i := 1; i <= n; i++ {
				if i == id0 {
					continue
				}
				if cnt[i] == 0 {
					if aVals[i]+tmp < ans {
						ans = aVals[i] + tmp
						C = []int{i}
						Cval = []int{1}
					}
				}
				if cnt[i] > 0 && cnt[i+n] == 0 {
					if tmp+(bVals[i]-aVals[i]) < ans {
						ans = tmp + (bVals[i] - aVals[i])
						C = []int{i + n}
						Cval = []int{1}
					}
					if tmp-aVals[i]+bVals[id0] < ans {
						ans = tmp - aVals[i] + bVals[id0]
						C = []int{i, id0, id0 + n}
						Cval = []int{-1, 1, 1}
					}
				}
				if cnt[i] > 0 && cnt[i+n] > 0 {
					if tmp-(bVals[i]-aVals[i])+bVals[id0] < ans {
						ans = tmp - (bVals[i] - aVals[i]) + bVals[id0]
						C = []int{i + n, id0, id0 + n}
						Cval = []int{-1, 1, 1}
					}
				}
			}
			for k, idx := range C {
				cnt[idx] += Cval[k]
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", ans>>1))
	for i := 1; i <= n; i++ {
		sel := cnt[i] + cnt[i+n]
		sb.WriteByte(byte('0' + sel))
	}
	return sb.String()
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
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		n := rng.Intn(3) + 1
		m := rng.Intn(2*n) + 1
		a := make([]int64, n+1)
		b := make([]int64, n+1)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 1; j <= n; j++ {
			ai := int64(rng.Intn(10) + 1)
			bi := ai + int64(rng.Intn(10)+1)
			a[j] = ai << 1
			b[j] = bi << 1
			sb.WriteString(fmt.Sprintf("%d %d\n", ai, bi))
		}
		input := sb.String()
		exp := expected(n, m, a, b)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
