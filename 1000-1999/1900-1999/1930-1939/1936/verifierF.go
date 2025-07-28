package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const eps = 1e-6

type vect struct{ x, y float64 }

type circle struct {
	v vect
	r float64
}

func (v vect) add(o vect) vect    { return vect{v.x + o.x, v.y + o.y} }
func (v vect) sub(o vect) vect    { return vect{v.x - o.x, v.y - o.y} }
func (v vect) mul(s float64) vect { return vect{v.x * s, v.y * s} }
func (v vect) len() float64       { return math.Hypot(v.x, v.y) }
func (v vect) lensq() float64     { return v.x*v.x + v.y*v.y }

func isin(a, b circle) bool { return a.v.sub(b.v).len() <= b.r-a.r+1e-10 }

func f2(a, b circle) circle {
	v := b.v.sub(a.v)
	d := v.len()
	va := v.mul(a.r / d)
	vb := v.mul((d - b.r) / d)
	ctr := a.v.add(va.add(vb).mul(0.5))
	rad := va.sub(vb).len() * 0.5
	return circle{ctr, rad}
}

func f3(a, b, c circle) circle {
	x1, y1, r1 := a.v.x, a.v.y, a.r
	x2, y2, r2 := b.v.x, b.v.y, b.r
	x3, y3, r3 := c.v.x, c.v.y, c.r
	a2 := x1 - x2
	a3 := x1 - x3
	b2 := y1 - y2
	b3 := y1 - y3
	c2 := r2 - r1
	c3 := r3 - r1
	d1 := x1*x1 + y1*y1 - r1*r1
	d2 := d1 - x2*x2 - y2*y2 + r2*r2
	d3 := d1 - x3*x3 - y3*y3 + r3*r3
	ab := a3*b2 - a2*b3
	xa := (b2*d3-b3*d2)/(2*ab) - x1
	xb := (b3*c2 - b2*c3) / ab
	ya := (a3*d2-a2*d3)/(2*ab) - y1
	yb := (a2*c3 - a3*c2) / ab
	A := xb*xb + yb*yb - 1
	B := 2 * (r1 + xa*xb + ya*yb)
	C := xa*xa + ya*ya - r1*r1
	var r float64
	if math.Abs(A) > 1e-10 {
		disc := B*B - 4*A*C
		r = -(B - math.Sqrt(disc)) / (2 * A)
	} else {
		r = -C / B
	}
	cx := x1 + xa + xb*r
	cy := y1 + ya + yb*r
	return circle{vect{cx, cy}, r}
}

func minimalCircle(circles []circle) circle {
	rand.Seed(1)
	rand.Shuffle(len(circles), func(i, j int) { circles[i], circles[j] = circles[j], circles[i] })
	ans := circle{vect{0, 0}, 1e18}
	for i := 0; i < len(circles); i++ {
		if !isin(ans, circles[i]) {
			ans = circles[i]
			for j := 0; j < i; j++ {
				if !isin(ans, circles[j]) {
					ans = f2(circles[i], circles[j])
					for k := 0; k < j; k++ {
						if !isin(ans, circles[k]) {
							ans = f3(circles[i], circles[j], circles[k])
						}
					}
				}
			}
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func refSolveF(input string) (circle, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return circle{}, err
	}
	v := make([]circle, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &v[i].v.x, &v[i].v.y, &v[i].r)
	}
	return minimalCircle(v), nil
}

func genCaseF(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := rng.Float64()*10 - 5
		y := rng.Float64()*10 - 5
		r := rng.Float64()*3 + 0.1
		sb.WriteString(fmt.Sprintf("%.3f %.3f %.3f\n", x, y, r))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCaseF(rng)
		expect, err := refSolveF(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse generated case: %v", err)
			os.Exit(1)
		}
		outStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var x, y, r float64
		if _, err := fmt.Sscan(outStr, &x, &y, &r); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output\n", i+1)
			os.Exit(1)
		}
		if math.Abs(x-expect.v.x) > eps || math.Abs(y-expect.v.y) > eps || math.Abs(r-expect.r) > eps {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f %.6f %.6f got %s\ninput:\n%s", i+1, expect.v.x, expect.v.y, expect.r, outStr, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
