package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

const eps = 1e-8

type Line struct {
	a, b, c int
	id      int
}
type Ans struct {
	x, y int
}

func dcmp(x float64) int {
	if math.Abs(x) < eps {
		return 0
	}
	if x < 0 {
		return -1
	}
	return 1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	lines := make([]Line, n)
	for i := 0; i < n; i++ {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		lines[i] = Line{a: a, b: b, c: c, id: i + 1}
	}
	cnt := n
	ans := make([]Ans, 0)
	tot := 0
	exis := true
	rand.Seed(time.Now().UnixNano())
	for cnt > 0 {
		if k >= cnt {
			// remove remaining lines one by one
			ans = append(ans, Ans{lines[cnt-1].id, -1})
			tot++
			cnt--
			continue
		}
		var interX, interY float64
		found := false
		for pos := 0; pos < 100; pos++ {
			l := rand.Intn(cnt)
			r := rand.Intn(cnt)
			// check parallel
			la, lb := lines[l].a, lines[l].b
			ra, rb := lines[r].a, lines[r].b
			det := float64(la*rb - ra*lb)
			if dcmp(det) == 0 {
				continue
			}
			// intersection point
			interX = (float64(lb*lines[r].c - lines[l].c*rb)) / det
			interY = (float64(la*lines[r].c - lines[l].c*ra)) / (float64(lines[l].b*ra - la*rb))
			// count through
			sz := 0
			for i := 0; i < cnt; i++ {
				if dcmp(float64(lines[i].a)*interX+float64(lines[i].b)*interY+float64(lines[i].c)) == 0 {
					sz++
				}
			}
			if sz*k >= cnt {
				ans = append(ans, Ans{lines[l].id, lines[r].id})
				tot++
				found = true
				break
			}
		}
		if !found {
			exis = false
			break
		}
		// filter out lines passing through intersection
		w := 0
		for i := 0; i < cnt; i++ {
			if dcmp(float64(lines[i].a)*interX+float64(lines[i].b)*interY+float64(lines[i].c)) != 0 {
				lines[w] = lines[i]
				w++
			}
		}
		cnt = w
		k--
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	if !exis {
		writer.WriteString("NO\n")
		return
	}
	writer.WriteString("YES\n")
	writer.WriteString(fmt.Sprintf("%d\n", tot))
	for i := 0; i < tot; i++ {
		writer.WriteString(fmt.Sprintf("%d %d\n", ans[i].x, ans[i].y))
	}
}
