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

type queryF struct {
	typ int
	pos int
	val int64
	x   int
	k   int64
}

const refINF int64 = 1 << 60

func refEvalA(cnt []int64, n int, y int, s int64) int64 {
	if s <= 0 {
		return 0
	}
	ops := int64(0)
	need := s
	for j := y; j < n; j++ {
		c := cnt[j]
		if need <= c {
			return ops
		}
		if j == n-1 {
			return refINF
		}
		need -= c
		take := (need + 1) >> 1
		ops += take
		need = take
	}
	return refINF
}

func refMinB(cnt []int64, n int, y int, l, r int64) int64 {
	if l > r {
		return refINF
	}
	ans := refINF
	offset := int64(0)
	for {
		c := cnt[y]
		freeR := r
		if freeR > c {
			freeR = c
		}
		if l <= freeR {
			v := offset - freeR
			if v < ans {
				ans = v
			}
		}
		if y == n-1 || r <= c {
			return ans
		}
		if r > c && ((r-c)&1) == 1 {
			t := (r - c + 1) >> 1
			a := refEvalA(cnt, n, y+1, t)
			if a < refINF {
				v := offset - c + 1 + (a - t)
				if v < ans {
					ans = v
				}
			}
		}
		if l < c+1 {
			l = c + 1
		}
		l = (l - c + 1) >> 1
		if l < 1 {
			l = 1
		}
		r = (r - c) >> 1
		offset -= c
		y++
		if l > r || y >= n {
			return ans
		}
	}
}

func refSolve(cnt []int64, n int, pow2 []int64, totalCountAll, totalUnits int64, x int, k int64) int64 {
	if k > totalUnits {
		return -1
	}
	if x == n-1 {
		if k <= totalCountAll {
			return 0
		}
		return k - totalCountAll
	}
	var g, lowCap int64
	for i := 0; i <= x; i++ {
		g += cnt[i]
		lowCap += cnt[i] * pow2[i]
	}
	if k <= g {
		return 0
	}
	d := k - g
	l := lowCap - g
	y := x + 1
	u := pow2[y]
	lb := int64(0)
	if d > l {
		lb = (d - l + u - 1) / u
	}
	ans := refINF
	upper := (d - 1) >> 1
	if lb <= upper {
		v := refMinB(cnt, n, y, lb, upper)
		if v < refINF {
			cost := d + v
			if cost < ans {
				ans = cost
			}
		}
	}
	s0 := lb
	half := (d + 1) >> 1
	if s0 < half {
		s0 = half
	}
	a := refEvalA(cnt, n, y, s0)
	if a < refINF {
		cost := a + s0
		if cost < ans {
			ans = cost
		}
	}
	if ans >= refINF {
		return -1
	}
	return ans
}

func solveQueries(n int, cnt0 []int64, qs []queryF) []int64 {
	cnt := make([]int64, n)
	copy(cnt, cnt0)
	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] << 1
	}
	var totalCountAll, totalUnits int64
	for i := 0; i < n; i++ {
		totalCountAll += cnt[i]
		totalUnits += cnt[i] * pow2[i]
	}

	var res []int64
	for _, q := range qs {
		if q.typ == 1 {
			delta := q.val - cnt[q.pos]
			cnt[q.pos] = q.val
			totalCountAll += delta
			totalUnits += delta * pow2[q.pos]
		} else {
			ans := refSolve(cnt, n, pow2, totalCountAll, totalUnits, q.x, q.k)
			res = append(res, ans)
		}
	}
	return res
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("run error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const tnum = 100
	for ti := 0; ti < tnum; ti++ {
		n := rand.Intn(5) + 1
		q := rand.Intn(8) + 2
		cnt := make([]int64, n)
		for i := range cnt {
			cnt[i] = rand.Int63n(5)
		}
		queries := make([]queryF, q)
		hasType2 := false
		for i := 0; i < q; i++ {
			if rand.Intn(2) == 0 {
				queries[i].typ = 1
				queries[i].pos = rand.Intn(n)
				queries[i].val = rand.Int63n(10)
			} else {
				queries[i].typ = 2
				queries[i].x = rand.Intn(n)
				queries[i].k = rand.Int63n(10) + 1
				hasType2 = true
			}
		}
		if !hasType2 {
			queries[0].typ = 2
			queries[0].x = 0
			queries[0].k = 1
		}
		expOut := solveQueries(n, cnt, queries)
		var in strings.Builder
		fmt.Fprintf(&in, "%d %d\n", n, q)
		for i := 0; i < n; i++ {
			if i > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", cnt[i])
		}
		in.WriteByte('\n')
		for _, qu := range queries {
			if qu.typ == 1 {
				fmt.Fprintf(&in, "1 %d %d\n", qu.pos, qu.val)
			} else {
				fmt.Fprintf(&in, "2 %d %d\n", qu.x, qu.k)
			}
		}
		out, err := runBinary(binary, in.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		outs := strings.Fields(strings.TrimSpace(out))
		if len(outs) != len(expOut) {
			fmt.Printf("test %d failed: expected %d lines got %d\noutput:\n%s\n", ti+1, len(expOut), len(outs), out)
			os.Exit(1)
		}
		for i, field := range outs {
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil || v != expOut[i] {
				fmt.Printf("test %d failed at result %d: expected %d got %s\n", ti+1, i+1, expOut[i], field)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
