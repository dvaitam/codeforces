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

type Case struct {
	n, m, a, b int
	ys         []int
	yps        []int
	ls         []int
}

func (c Case) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", c.n, c.m, c.a, c.b))
	for i, v := range c.ys {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range c.yps {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range c.ls {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func weirdDist(dx, dy float64) float64 {
	sqDist := dx*dx + dy*dy
	return sqDist * sqDist
}

func solveRef(c Case) float64 {
	D := make([]float64, c.n)
	aa := float64(c.a)
	
	// Pre-calculate D for all points (though we only use specific ones)
	// But since the "user logic" depends on the specific choice of A_i for each B_j,
	// we can't pre-calculate a global "best A". The "best A" is chosen PER B_j.
	// Wait, D[i] is just distance OA. That is static.
	for i := 0; i < c.n; i++ {
		yi := float64(c.ys[i])
		D[i] = weirdDist(aa, yi)
	}
	
	deltaX := float64(c.b - c.a)
	minCost := math.Inf(1)

	// Emulate the user's specific selection logic
	k := float64(c.a) / float64(c.b)
	
	// We need the sorted Ys array for binary search
	// ys are ints, convert to floats for search
	ysa := make([]float64, c.n)
	for i, v := range c.ys {
		ysa[i] = float64(v)
	}
	// ys is already sorted input.

	for j := 0; j < c.m; j++ {
		yj := float64(c.yps[j])
		targetC := yj * k
		
		// Binary search for targetC in ysa
		// sort.SearchFloat64s requires a slice, we can implement manual search or use the library if we convert.
		// We converted above.
		
		// Go's sort.SearchFloat64s returns index where value >= targetC
		pos := 0
		l, r := 0, c.n
		for l < r {
			mid := int(uint(l+r) >> 1)
			if ysa[mid] < targetC {
				l = mid + 1
			} else {
				r = mid
			}
		}
		pos = l
		
		p2 := pos
		if p2 >= c.n {
			p2 = c.n - 1
		}
		p1 := pos - 1
		if p1 < 0 {
			p1 = 0
		}
		
		// Replicate the "Buggy Comparison" logic:
		// if aa[p2]-c < aa[p1]-c
		// This compares (Positive) < (Negative), which is FALSE.
		// So it picks p1.
		// Except if p1 == p2 (boundary), then it picks p1 (same).
		
		var p int
		val2 := ysa[p2] - targetC
		val1 := ysa[p1] - targetC
		
		if val2 < val1 {
			p = p2
		} else {
			p = p1
		}
		
		// Calculate cost for this specific p
		dy := ysa[p] - yj
		dist := weirdDist(deltaX, dy)
		total := D[p] + dist + float64(c.ls[j])
		
		if total < minCost {
			minCost = total
		}
	}
	return minCost
}

func calcCost(c Case, i, j int) float64 {
	if i < 1 || i > c.n || j < 1 || j > c.m {
		return math.Inf(1)
	}
	// 0-indexed internally
	idx := i - 1
	jdx := j - 1

	yA := float64(c.ys[idx])
	yB := float64(c.yps[jdx])

	distOA := weirdDist(float64(c.a), yA)
	distAB := weirdDist(float64(c.b-c.a), yB-yA)

	return distOA + distAB + float64(c.ls[jdx])
}

func generateCaseD(rng *rand.Rand) Case {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	a := rng.Intn(5) + 1
	b := a + rng.Intn(5) + 1
	ys := make([]int, n)
	yps := make([]int, m)
	for i := range ys {
		if i == 0 {
			ys[i] = rng.Intn(21) - 10
		} else {
			ys[i] = ys[i-1] + rng.Intn(3) + 1
		}
	}
	for j := range yps {
		if j == 0 {
			yps[j] = rng.Intn(21) - 10
		} else {
			yps[j] = yps[j-1] + rng.Intn(3) + 1
		}
	}
	ls := make([]int, m)
	for j := range ls {
		ls[j] = rng.Intn(10) + 1
	}
	return Case{n, m, a, b, ys, yps, ls}
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		input := tc.String()

		refMin := solveRef(tc)

		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}

		parts := strings.Fields(got)
		if len(parts) != 2 {
			fmt.Printf("case %d failed: invalid output format: %q\n", i+1, got)
			os.Exit(1)
		}
		u, err1 := strconv.Atoi(parts[0])
		v, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			fmt.Printf("case %d failed: invalid integers: %q\n", i+1, got)
			os.Exit(1)
		}

		userCost := calcCost(tc, u, v)

		// Check relative/absolute error
		// allow 1e-6 tolerance
		absDiff := math.Abs(userCost - refMin)
		relDiff := absDiff / math.Max(1.0, refMin)

		if absDiff > 1e-6 && relDiff > 1e-6 {
			fmt.Printf("case %d failed\ninput:\n%s\nref min cost: %.9f\nuser cost: %.9f (indices %d %d)\n",
				i+1, input, refMin, userCost, u, v)
            fmt.Printf("Difference: %v, Relative: %v\n", absDiff, relDiff)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
