package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type test struct {
	input    string
	expected string
}

type lineKey struct {
	a, b, c int64
}

type lineGroup struct {
	starts []int64
	ends   []int64
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func gcd3(a, b, c int64) int64 {
	return gcd(gcd(a, b), c)
}

func solveD(n, m int, segments [][4]int64, circles [][3]int64) string {
	groups := make(map[lineKey]*lineGroup, n)
	for _, s := range segments {
		x1, y1, x2, y2 := s[0], s[1], s[2], s[3]
		dx := x2 - x1
		dy := y2 - y1
		g := gcd(dx, dy)
		dx /= g
		dy /= g
		if dx < 0 || (dx == 0 && dy < 0) {
			dx = -dx
			dy = -dy
		}
		a := -dy
		b := dx
		c2 := 2 * (a*x1 + b*y1)
		g2 := gcd3(a, b, c2)
		a /= g2
		b /= g2
		c2 /= g2
		if a < 0 || (a == 0 && b < 0) {
			a = -a
			b = -b
			c2 = -c2
		}
		key := lineKey{a, b, c2}
		grp, ok := groups[key]
		if !ok {
			grp = &lineGroup{}
			groups[key] = grp
		}
		t1 := 2 * (dx*x1 + dy*y1)
		t2 := 2 * (dx*x2 + dy*y2)
		if t1 <= t2 {
			grp.starts = append(grp.starts, t1)
			grp.ends = append(grp.ends, t2)
		} else {
			grp.starts = append(grp.starts, t2)
			grp.ends = append(grp.ends, t1)
		}
	}
	for _, grp := range groups {
		sort.Slice(grp.starts, func(i, j int) bool { return grp.starts[i] < grp.starts[j] })
		sort.Slice(grp.ends, func(i, j int) bool { return grp.ends[i] < grp.ends[j] })
	}
	type circle struct{ x, y, r int64 }
	byRadius := make(map[int64][]circle)
	for _, c := range circles {
		x, y, r := c[0], c[1], c[2]
		byRadius[r] = append(byRadius[r], circle{x, y, r})
	}
	var ans int64
	for r, cs := range byRadius {
		k := len(cs)
		if k < 2 {
			continue
		}
		lim := 4 * r * r
		for i := 0; i < k; i++ {
			xi, yi := cs[i].x, cs[i].y
			for j := i + 1; j < k; j++ {
				xj, yj := cs[j].x, cs[j].y
				dx := xj - xi
				dy := yj - yi
				if dx*dx+dy*dy <= lim {
					continue
				}
				dxd := -dy
				dyd := dx
				g := gcd(dxd, dyd)
				dxd /= g
				dyd /= g
				if dxd < 0 || (dxd == 0 && dyd < 0) {
					dxd = -dxd
					dyd = -dyd
				}
				a := -dyd
				b := dxd
				c2 := a*(xi+xj) + b*(yi+yj)
				g2 := gcd3(a, b, c2)
				a /= g2
				b /= g2
				c2 /= g2
				if a < 0 || (a == 0 && b < 0) {
					a = -a
					b = -b
					c2 = -c2
				}
				key := lineKey{a, b, c2}
				grp, ok := groups[key]
				if !ok {
					continue
				}
				tmid := dxd*(xi+xj) + dyd*(yi+yj)
				sCnt := sort.Search(len(grp.starts), func(i int) bool { return grp.starts[i] > tmid })
				eCnt := sort.Search(len(grp.ends), func(i int) bool { return grp.ends[i] >= tmid })
				ans += int64(sCnt - eCnt)
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(3) + 1
		m := rng.Intn(3) + 2
		segs := make([][4]int64, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			x1 := int64(rng.Intn(11) - 5)
			y1 := int64(rng.Intn(11) - 5)
			x2 := int64(rng.Intn(11) - 5)
			y2 := int64(rng.Intn(11) - 5)
			for x1 == x2 && y1 == y2 {
				x2 = int64(rng.Intn(11) - 5)
				y2 = int64(rng.Intn(11) - 5)
			}
			segs[i] = [4]int64{x1, y1, x2, y2}
			fmt.Fprintf(&sb, "%d %d %d %d\n", x1, y1, x2, y2)
		}
		cirs := make([][3]int64, m)
		for i := 0; i < m; i++ {
			x := int64(rng.Intn(11) - 5)
			y := int64(rng.Intn(11) - 5)
			r := int64(rng.Intn(5) + 1)
			cirs[i] = [3]int64{x, y, r}
			fmt.Fprintf(&sb, "%d %d %d\n", x, y, r)
		}
		expected := solveD(n, m, segs, cirs)
		tests = append(tests, test{sb.String(), expected})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
