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

type point struct{ x, y float64 }

func dist(a, b point) float64 {
	return math.Hypot(a.x-b.x, a.y-b.y)
}

func solve(n, k int, xs []int, py int) float64 {
	axis := make([]int, n)
	copy(axis, xs[:n])
	sort.Ints(axis)
	L := float64(axis[0])
	R := float64(axis[n-1])
	dLine := R - L
	px := float64(xs[n])
	pyf := float64(py)
	if n <= 2 {
		m := n + 1
		xsF := make([]float64, m)
		ysF := make([]float64, m)
		for i := 0; i < n; i++ {
			xsF[i] = float64(xs[i])
		}
		for i := 0; i < n; i++ {
			ysF[i] = 0
		}
		xsF[n] = px
		ysF[n] = pyf
		start := k - 1
		others := make([]int, 0, m-1)
		for i := 0; i < m; i++ {
			if i != start {
				others = append(others, i)
			}
		}
		best := math.Inf(1)
		var permute func([]int, int)
		permute = func(a []int, l int) {
			if l == len(a) {
				cur := 0.0
				prev := start
				for _, idx := range a {
					cur += dist(point{xsF[prev], ysF[prev]}, point{xsF[idx], ysF[idx]})
					prev = idx
				}
				if cur < best {
					best = cur
				}
				return
			}
			for i := l; i < len(a); i++ {
				a[l], a[i] = a[i], a[l]
				permute(a, l+1)
				a[l], a[i] = a[i], a[l]
			}
		}
		permute(others, 0)
		return best
	}
	idx := sort.Search(n, func(i int) bool { return float64(axis[i]) >= px })
	dv := math.Inf(1)
	for _, j := range []int{idx - 1, idx} {
		if j >= 0 && j < n {
			d := dist(point{px, pyf}, point{float64(axis[j]), 0})
			if d < dv {
				dv = d
			}
		}
	}
	ans := math.Inf(1)
	if k == n+1 {
		d1 := dist(point{px, pyf}, point{L, 0})
		d2 := dist(point{px, pyf}, point{R, 0})
		ans = dLine + math.Min(d1, d2)
	} else {
		sx := float64(xs[k-1])
		for dir := 0; dir < 2; dir++ {
			var e1, e2 float64
			if dir == 0 {
				e1, e2 = L, R
			} else {
				e1, e2 = R, L
			}
			base := math.Abs(sx-e1) + dLine
			ans = math.Min(ans, base+2*dv)
			ans = math.Min(ans, dist(point{sx, 0}, point{px, pyf})+dist(point{px, pyf}, point{e1, 0})+dLine)
			ans = math.Min(ans, base+dist(point{px, pyf}, point{e2, 0}))
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	xs := make([]int, n+1)
	used := map[int]bool{}
	for i := 0; i < n; i++ {
		for {
			v := rng.Intn(21) - 10
			if !used[v] {
				used[v] = true
				xs[i] = v
				break
			}
		}
	}
	for {
		v := rng.Intn(21) - 10
		if !used[v] {
			xs[n] = v
			break
		}
	}
	py := rng.Intn(21) - 10
	if py == 0 {
		py = 1
	}
	k := rng.Intn(n+1) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", xs[i]))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d\n", py))
	expected := fmt.Sprintf("%.10f", solve(n, k, xs, py))
	return sb.String(), expected
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
	outStr := strings.TrimSpace(out.String())
	valExp, _ := strconv.ParseFloat(expected, 64)
	if valOut, err := strconv.ParseFloat(outStr, 64); err == nil {
		if math.Abs(valOut-valExp) > 1e-6 {
			return fmt.Errorf("expected %.10f got %s", valExp, outStr)
		}
	} else {
		return fmt.Errorf("invalid output %s", outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
