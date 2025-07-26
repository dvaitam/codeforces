package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type pair struct{ x, y float64 }
type interval struct{ l, r float64 }

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(pts []pair, m int) float64 {
	n := len(pts)
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x != pts[j].x {
			return pts[i].x < pts[j].x
		}
		return pts[i].y < pts[j].y
	})
	uniq := pts[:1]
	for i := 1; i < n; i++ {
		if pts[i].x != pts[i-1].x || pts[i].y != pts[i-1].y {
			uniq = append(uniq, pts[i])
		}
	}
	pts = uniq
	n = len(pts)
	dist := make([]float64, n)
	th := make([]float64, n)
	rb := math.Inf(1)
	for i, p := range pts {
		d := math.Hypot(p.x, p.y)
		dist[i] = d
		if d < rb {
			rb = d
		}
		th[i] = math.Atan2(p.y, p.x)
	}
	lb := 0.0
	b := make([]interval, n)
	var c []interval
	twoPI := 2 * math.Pi
	chk := func(r float64) bool {
		for i := 0; i < n; i++ {
			ac := math.Acos(r / dist[i])
			l := th[i] - ac
			rr := th[i] + ac
			if l < 0 {
				l += twoPI
				rr += twoPI
			}
			b[i].l = l
			b[i].r = rr
		}
		sort.Slice(b, func(i, j int) bool { return b[i].l < b[j].l })
		c = c[:0]
		for i := 0; i < n; i++ {
			for len(c) > 0 && c[len(c)-1].r >= b[i].r {
				c = c[:len(c)-1]
			}
			if len(c) == 0 || c[0].r > b[i].r-twoPI {
				c = append(c, b[i])
			}
		}
		aa := len(c)
		if aa == 0 {
			return true
		}
		orig := make([]interval, aa)
		copy(orig, c)
		for i := 0; i < aa; i++ {
			c = append(c, interval{orig[i].l + twoPI, orig[i].r + twoPI})
		}
		st := make([][17]int, aa)
		j := 0
		for i := 0; i < aa; i++ {
			for j < i+aa && c[j].l <= c[i].r {
				j++
			}
			st[i][0] = j - i
		}
		for k := 1; k < 17; k++ {
			for i := 0; i < aa; i++ {
				nxt := (i + st[i][k-1]) % aa
				st[i][k] = st[i][k-1] + st[nxt][k-1]
			}
		}
		for s := 0; s < aa; s++ {
			used, covered, idx := 0, 0, s
			for k := 16; k >= 0; k-- {
				if used+(1<<k) <= m {
					covered += st[idx][k]
					idx = (idx + st[idx][k]) % aa
					used += 1 << k
				}
			}
			if covered >= aa {
				return true
			}
		}
		return false
	}
	for rb-lb > 1e-7 {
		mb := (lb + rb) / 2
		if chk(mb) {
			lb = mb
		} else {
			rb = mb
		}
	}
	return (lb + rb) / 2
}

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(n) + 1
	pts := make([]pair, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		x := rng.Float64()*5 - 2.5
		y := rng.Float64()*5 - 2.5
		pts[i] = pair{x, y}
		sb.WriteString(fmt.Sprintf("%.3f %.3f\n", x, y))
	}
	ans := solve(pts, m)
	return sb.String(), fmt.Sprintf("%.9f", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		ev, _ := strconv.ParseFloat(expected, 64)
		ov, err2 := strconv.ParseFloat(out, 64)
		if err2 != nil || math.Abs(ev-ov) > 1e-3 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
