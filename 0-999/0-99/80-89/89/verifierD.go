package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type P struct{ x, y, z float64 }

func (p P) add(q P) P       { return P{p.x + q.x, p.y + q.y, p.z + q.z} }
func (p P) sub(q P) P       { return P{p.x - q.x, p.y - q.y, p.z - q.z} }
func (p P) mul(k float64) P { return P{p.x * k, p.y * k, p.z * k} }
func (p P) dot(q P) float64 { return p.x*q.x + p.y*q.y + p.z*q.z }
func (p P) cross(q P) P {
	return P{
		p.y*q.z - p.z*q.y,
		p.z*q.x - p.x*q.z,
		p.x*q.y - p.y*q.x,
	}
}
func (p P) mag2() float64 { return p.dot(p) }
func (p P) mag() float64  { return math.Sqrt(p.mag2()) }

const inf = 1e20
const eps = 1e-9

func f(a, o P, r1, r2 float64, v P) (bool, float64) {
	A := v.mag2()
	B := 2 * a.sub(o).dot(v)
	C := a.sub(o).mag2() - (r1+r2)*(r1+r2)
	D := B*B - 4*A*C
	if D < 0 {
		return false, 0
	}
	if D > 0 {
		D = math.Sqrt(D)
	}
	t := (-B - D) / (2 * A)
	return t > 0, t
}

type mine struct {
	cen    P
	r      float64
	spikes []P
}

type caseD struct {
	A, V  P
	R     float64
	mines []mine
}

func solveD(c caseD) float64 {
	ans := inf
	upd := func(t float64) {
		if t < ans {
			ans = t
		}
	}
	for _, m := range c.mines {
		if ok, t := f(c.A, m.cen, c.R, m.r, c.V); ok {
			upd(t)
		}
		for _, p := range m.spikes {
			if ok, t := f(c.A, m.cen, c.R, 0, c.V); ok {
				upd(t)
			}
			if ok, t := f(c.A, m.cen.add(p), c.R, 0, c.V); ok {
				upd(t)
			}
			c1 := m.cen
			d := m.cen.add(p)
			cd := d.sub(c1)
			ca := c1.sub(c.A)
			crossVc := c.V.cross(cd)
			A1 := crossVc.mag2()
			if math.Abs(A1) < eps {
				continue
			}
			B := crossVc.mag() * ca.cross(cd).mag()
			C := ca.cross(cd).mag2()
			D := B*B - 4*A1*C
			if D >= 0 {
				if D > 0 {
					D = math.Sqrt(D)
				}
				t1 := (-B + D) / (2 * A1)
				if t1 > 0 {
					Q := c.A.add(c.V.mul(t1))
					if Q.sub(c1).dot(cd) >= 0 && Q.sub(d).dot(c1.sub(d)) >= 0 {
						upd(t1)
					}
				}
			}
		}
	}
	if ans == inf {
		return -1
	}
	return ans
}

func genCaseD(rng *rand.Rand) caseD {
	A := P{float64(rng.Intn(41) - 20), float64(rng.Intn(41) - 20), float64(rng.Intn(41) - 20)}
	V := P{}
	for V.x == 0 && V.y == 0 && V.z == 0 {
		V = P{float64(rng.Intn(11) - 5), float64(rng.Intn(11) - 5), float64(rng.Intn(11) - 5)}
	}
	R := float64(rng.Intn(10) + 10)
	m := mine{}
	m.cen = P{float64(rng.Intn(41) - 20), float64(rng.Intn(41) - 20), float64(rng.Intn(41) - 20)}
	m.r = float64(rng.Intn(5) + 1)
	spikeCount := rng.Intn(3)
	m.spikes = make([]P, spikeCount)
	for i := 0; i < spikeCount; i++ {
		for {
			p := P{float64(rng.Intn(11) - 5), float64(rng.Intn(11) - 5), float64(rng.Intn(11) - 5)}
			if d := p.mag(); d > m.r && d <= 1.5*m.r {
				m.spikes[i] = p
				break
			}
		}
	}
	// ensure initial distance > R + r
	for math.Sqrt(m.cen.sub(A).mag2()) <= R+m.r {
		m.cen.x += 10
	}
	return caseD{A, V, R, []mine{m}}
}

func runCaseD(bin string, c caseD) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d %d %d %d\n", int(c.A.x), int(c.A.y), int(c.A.z), int(c.V.x), int(c.V.y), int(c.V.z), int(c.R)))
	sb.WriteString("1\n")
	m := c.mines[0]
	sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", int(m.cen.x), int(m.cen.y), int(m.cen.z), int(m.r), len(m.spikes)))
	for _, p := range m.spikes {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", int(p.x), int(p.y), int(p.z)))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveD(c)
	if exp < 0 && got != -1 {
		return fmt.Errorf("expected -1 got %.6f", got)
	}
	if exp >= 0 && (got < 0 || math.Abs(got-exp) > 1e-6) {
		return fmt.Errorf("expected %.6f got %.6f", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := genCaseD(rng)
		if err := runCaseD(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
