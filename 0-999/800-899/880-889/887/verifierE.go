package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = `package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/bits"
	"os"
	"sort"
)

type U128 struct {
	hi uint64
	lo uint64
}

type Interval struct {
	l float64
	r float64
}

const two64 = 18446744073709551616.0

func mulU(a, b uint64) U128 {
	hi, lo := bits.Mul64(a, b)
	return U128{hi: hi, lo: lo}
}

func cmpU(a, b U128) int {
	if a.hi < b.hi {
		return -1
	}
	if a.hi > b.hi {
		return 1
	}
	if a.lo < b.lo {
		return -1
	}
	if a.lo > b.lo {
		return 1
	}
	return 0
}

func subU(a, b U128) U128 {
	lo, borrow := bits.Sub64(a.lo, b.lo, 0)
	hi, _ := bits.Sub64(a.hi, b.hi, borrow)
	return U128{hi: hi, lo: lo}
}

func shlU(a U128, k uint) U128 {
	if k == 0 {
		return a
	}
	if k < 64 {
		return U128{
			hi: (a.hi << k) | (a.lo >> (64 - k)),
			lo: a.lo << k,
		}
	}
	if k < 128 {
		return U128{
			hi: a.lo << (k - 64),
			lo: 0,
		}
	}
	return U128{}
}

func toFloatU(a U128) float64 {
	return float64(a.hi)*two64 + float64(a.lo)
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	idx := 0
	nextInt := func() int64 {
		for idx < len(data) {
			c := data[idx]
			if c == ' ' || c == '\n' || c == '\r' || c == '\t' {
				idx++
			} else {
				break
			}
		}
		sign := int64(1)
		if data[idx] == '-' {
			sign = -1
			idx++
		}
		var val int64
		for idx < len(data) {
			c := data[idx]
			if c < '0' || c > '9' {
				break
			}
			val = val*10 + int64(c-'0')
			idx++
		}
		return sign * val
	}

	x1 := nextInt()
	y1 := nextInt()
	x2 := nextInt()
	y2 := nextInt()

	n := int(nextInt())

	dx := x2 - x1
	dy := y2 - y1
	D := dx*dx + dy*dy
	Df := float64(D)
	sqrtD := math.Sqrt(Df)
	a2 := Df / 4.0

	intervals := make([]Interval, 0, n)

	for i := 0; i < n; i++ {
		xi := nextInt()
		yi := nextInt()
		ri := nextInt()

		s := dx*(yi-y1) - dy*(xi-x1)
		sAbs := uint64(abs64(s))
		r2 := uint64(ri) * uint64(ri)

		s2 := mulU(sAbs, sAbs)
		r2D := mulU(r2, uint64(D))
		Nu := subU(s2, r2D)
		Nf := toFloatU(Nu)

		tx := 2*xi - x1 - x2
		ty := 2*yi - y1 - y2
		K4 := tx*tx + ty*ty - D - 4*ri*ri
		K := float64(K4) / 4.0

		kAbs := uint64(abs64(K4))
		K4sq := mulU(kAbs, kAbs)
		sixteenR2D := shlU(r2D, 4)

		var Cf float64
		if cmpU(K4sq, sixteenR2D) >= 0 {
			Cf = toFloatU(subU(K4sq, sixteenR2D)) / 64.0
		} else {
			Cf = -toFloatU(subU(sixteenR2D, K4sq)) / 64.0
		}

		A := Nf / Df
		B := -K * float64(s) / sqrtD
		sqrtDisc := float64(ri) * math.Sqrt(K*K+Nf)

		var q float64
		if B >= 0 {
			q = -0.5 * (B + sqrtDisc)
		} else {
			q = -0.5 * (B - sqrtDisc)
		}

		r1 := q / A
		r2f := Cf / q
		if r1 > r2f {
			r1, r2f = r2f, r1
		}
		intervals = append(intervals, Interval{l: r1, r: r2f})
	}

	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].l == intervals[j].l {
			return intervals[i].r < intervals[j].r
		}
		return intervals[i].l < intervals[j].l
	})

	bestT := 0.0
	covered := false

	for i := 0; i < len(intervals); {
		L := intervals[i].l
		R := intervals[i].r
		i++
		for i < len(intervals) && intervals[i].l < R {
			if intervals[i].r > R {
				R = intervals[i].r
			}
			i++
		}
		if L < 0 && 0 < R {
			covered = true
			if -L < R {
				bestT = -L
			} else {
				bestT = R
			}
			break
		}
		if L > 0 {
			break
		}
	}

	ans := sqrtD / 2.0
	if covered {
		ans = math.Sqrt(a2 + bestT*bestT)
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(out, "%.15f\n", ans)
	out.Flush()
}
`

func buildEmbeddedRef(dir string) (string, error) {
	src := filepath.Join(dir, "ref_embedded_887E.go")
	if err := os.WriteFile(src, []byte(refSource), 0644); err != nil {
		return "", err
	}
	bin := filepath.Join(dir, "ref_887E_bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build ref: %v\n%s", err, out)
	}
	os.Remove(src)
	return bin, nil
}

func runBin(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tmpDir, err := os.MkdirTemp("", "v887E")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	refBin, err := buildEmbeddedRef(tmpDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rand.Seed(5)
	for i := 0; i < 100; i++ {
		x1 := rand.Float64()*200 - 100
		y1 := rand.Float64()*200 - 100
		x2 := rand.Float64()*200 - 100
		y2 := rand.Float64()*200 - 100
		n := rand.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", int(x1), int(y1), int(x2), int(y2))
		fmt.Fprintf(&sb, "%d\n", n)
		circles := make([][3]int, n)
		for j := 0; j < n; j++ {
			cx := int(rand.Float64()*200 - 100)
			cy := int(rand.Float64()*200 - 100)
			cr := int(rand.Float64()*10) + 1
			circles[j] = [3]int{cx, cy, cr}
			fmt.Fprintf(&sb, "%d %d %d\n", cx, cy, cr)
		}
		input := sb.String()

		exp, err := runBin(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		var expVal, gotVal float64
		fmt.Sscan(exp, &expVal)
		fmt.Sscan(got, &gotVal)
		if math.Abs(gotVal-expVal) > 1e-4*math.Max(1, math.Abs(expVal)) {
			fmt.Fprintf(os.Stderr, "case %d: expected %.6f got %.6f\ninput:\n%s", i+1, expVal, gotVal, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
