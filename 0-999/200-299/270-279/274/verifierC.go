package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type pnt struct{ x, y float64 }

const ep = 1e-9

func (a pnt) sub(b pnt) pnt       { return pnt{a.x - b.x, a.y - b.y} }
func (a pnt) add(b pnt) pnt       { return pnt{a.x + b.x, a.y + b.y} }
func (a pnt) mul(s float64) pnt   { return pnt{a.x * s, a.y * s} }
func (a pnt) div(s float64) pnt   { return pnt{a.x / s, a.y / s} }
func (a pnt) dot(b pnt) float64   { return a.x*b.x + a.y*b.y }
func (a pnt) cross(b pnt) float64 { return a.x*b.y - a.y*b.x }
func (a pnt) dist() float64       { return math.Hypot(a.x, a.y) }

func outercenter(a, b, c pnt) pnt {
	c1 := (a.dot(a) - b.dot(b)) / 2
	c2 := (a.dot(a) - c.dot(c)) / 2
	d1 := a.sub(b).cross(a.sub(c))
	x0 := (c1*(a.y-c.y) - c2*(a.y-b.y)) / d1
	y0 := (c1*(a.x-c.x) - c2*(a.x-b.x)) / -d1
	return pnt{x0, y0}
}

func chk(a, b, c pnt, pts []pnt, i, j, k int) bool {
	r := b.sub(c).dist() / 2
	o := b.add(c).div(2)
	bf := false
	for t := range pts {
		if t == i || t == j || t == k {
			continue
		}
		if (pts[t].sub(c).cross(b.sub(c)) * (a.sub(c).cross(b.sub(c)))) < -ep {
			if pts[t].sub(b).dist() < 2*r && pts[t].sub(c).dist() < 2*r {
				bf = true
			}
		}
		if o.sub(pts[t]).dist() < r-ep {
			return false
		}
	}
	return bf
}

func solveC(n int, pts []pnt) float64 {
	ans := -1.0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				fijk := pts[i].sub(pts[j]).dot(pts[k].sub(pts[j]))
				fikj := pts[i].sub(pts[k]).dot(pts[j].sub(pts[k]))
				fjik := pts[j].sub(pts[i]).dot(pts[k].sub(pts[i]))
				if fijk < -ep || fikj < -ep || fjik < -ep {
					continue
				}
				if math.Abs(fijk) < ep {
					a, b, c := pts[j], pts[i], pts[k]
					if chk(a, b, c, pts, i, j, k) {
						r := b.sub(c).dist() / 2
						if r > ans {
							ans = r
						}
					}
				} else if math.Abs(fikj) < ep {
					a, b, c := pts[k], pts[i], pts[j]
					if chk(a, b, c, pts, i, j, k) {
						r := b.sub(c).dist() / 2
						if r > ans {
							ans = r
						}
					}
				} else if math.Abs(fjik) < ep {
					a, b, c := pts[i], pts[j], pts[k]
					if chk(a, b, c, pts, i, j, k) {
						r := b.sub(c).dist() / 2
						if r > ans {
							ans = r
						}
					}
				} else {
					area := math.Abs(pts[i].sub(pts[j]).cross(pts[k].sub(pts[j])))
					r := pts[i].sub(pts[j]).dist() * pts[i].sub(pts[k]).dist() * pts[j].sub(pts[k]).dist() / area / 2
					o := outercenter(pts[i], pts[j], pts[k])
					ok := true
					for t := range pts {
						if t == i || t == j || t == k {
							continue
						}
						if o.sub(pts[t]).dist() < r-ep {
							ok = false
							break
						}
					}
					if ok && r > ans {
						ans = r
					}
				}
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	pts := make([]pnt, 0, n)
	used := map[[2]int]bool{}
	for len(pts) < n {
		x := rng.Intn(11) - 5
		y := rng.Intn(11) - 5
		if !used[[2]int{x, y}] {
			used[[2]int{x, y}] = true
			pts = append(pts, pnt{float64(x), float64(y)})
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n) + "\n")
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", int(p.x), int(p.y)))
	}
	ans := solveC(n, pts)
	if ans < -ep {
		return sb.String(), "-1\n"
	}
	return sb.String(), fmt.Sprintf("%.12f\n", ans)
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
