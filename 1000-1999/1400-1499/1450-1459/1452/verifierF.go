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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveQueries(n int, cnt []int64, qs []queryF) []int64 {
	c := make([]int64, n)
	copy(c, cnt)
	res := []int64{}
	for _, q := range qs {
		if q.typ == 1 {
			c[q.pos] = q.val
		} else {
			x := q.x
			k := q.k
			m := n - x - 1
			arr := make([]int64, maxInt(1, m+1))
			var small int64
			for i := 0; i <= x && i < n; i++ {
				small += c[i]
			}
			arr[0] = small
			for j := 1; j <= m; j++ {
				arr[j] = c[x+j]
			}
			need := k
			if arr[0] >= need {
				res = append(res, 0)
				continue
			}
			var totUnits int64
			for j, v := range arr {
				if v == 0 {
					continue
				}
				totUnits += v << uint(j)
			}
			if totUnits < need {
				res = append(res, -1)
				continue
			}
			var ans int64
			need2 := need - arr[0]
			for need2 > 0 {
				r1 := (need2 + 1) / 2
				var avail1 int64
				if len(arr) > 1 {
					avail1 = arr[1]
				}
				need1 := r1 - avail1
				if need1 <= 0 {
					ans += r1
					break
				}
				jj := 2
				for jj <= m && (jj >= len(arr) || arr[jj] == 0) {
					jj++
				}
				if jj > m {
					ans = -1
					break
				}
				sj := int64(1) << (jj - 1)
				needCj := (need1 + sj - 1) / sj
				if needCj > arr[jj] {
					needCj = arr[jj]
				}
				ans += needCj
				arr[jj] -= needCj
				arr[jj-1] += needCj * 2
			}
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
