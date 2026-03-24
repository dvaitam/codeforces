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

// Correct solver embedded from CF-accepted solution.
func solveCorrect(n, k int, xs []int, py int) float64 {
	x := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		x[i] = float64(xs[i])
	}
	y := float64(py)

	Sx := x[n]
	Sy := y

	Aunsorted := make([]float64, n)
	copy(Aunsorted, x[:n])

	var startVal float64
	if k <= n {
		startVal = Aunsorted[k-1]
	}

	A := make([]float64, n)
	copy(A, Aunsorted)
	sort.Float64s(A)

	dist := func(px float64) float64 {
		return math.Hypot(px-Sx, Sy)
	}

	if k == n+1 {
		ans := dist(A[0])
		if n > 1 {
			ans = A[n-1] - A[0] + math.Min(dist(A[0]), dist(A[n-1]))
		}
		return ans
	}

	startIdx := 0
	for i := 0; i < n; i++ {
		if A[i] == startVal {
			startIdx = i + 1
			break
		}
	}
	_ = startIdx

	ans := math.Inf(1)
	Astart := startVal

	for i := 0; i < n; i++ {
		target := A[i]
		c1 := (Astart - A[0]) + (A[n-1] - A[0]) + (A[n-1] - target)
		c2 := (A[n-1] - Astart) + (A[n-1] - A[0]) + (target - A[0])
		cost := math.Min(c1, c2) + dist(target)
		if cost < ans {
			ans = cost
		}
	}

	for u := 1; u < n; u++ {
		if startIdx <= u {
			for _, xIdx := range []int{0, u - 1} {
				X := A[xIdx]
				for _, yIdx := range []int{u, n - 1} {
					Y := A[yIdx]
					var CostX float64
					if xIdx == u-1 {
						CostX = Astart - A[0] + A[u-1] - A[0]
					} else {
						CostX = A[u-1] - Astart + A[u-1] - A[0]
					}
					CostY := A[n-1] - A[u]
					cost := CostX + CostY + dist(X) + dist(Y)
					if cost < ans {
						ans = cost
					}
				}
			}
		} else {
			for _, xIdx := range []int{0, u - 1} {
				X := A[xIdx]
				for _, yIdx := range []int{u, n - 1} {
					Y := A[yIdx]
					CostX := A[u-1] - A[0]
					var CostY float64
					if yIdx == u {
						CostY = A[n-1] - Astart + A[n-1] - A[u]
					} else {
						CostY = Astart - A[u] + A[n-1] - A[u]
					}
					cost := CostX + CostY + dist(X) + dist(Y)
					if cost < ans {
						ans = cost
					}
				}
			}
		}
	}

	return ans
}

func generateCase(rng *rand.Rand) (string, float64) {
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
	expected := solveCorrect(n, k, xs, py)
	return sb.String(), expected
}

func runCase(bin, input string, expected float64) error {
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
	valOut, err := strconv.ParseFloat(outStr, 64)
	if err != nil {
		return fmt.Errorf("invalid output %s", outStr)
	}
	if math.Abs(valOut-expected) > 1e-4 {
		return fmt.Errorf("expected %.10f got %s", expected, outStr)
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
