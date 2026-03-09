package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const tolerance = 1e-6

// ---- inline oracle from 424E.go ----

var (
	oracleData []int
	oracleMemo map[uint64]float64
)

const oracleP = 71

func getState(a, b, c int) int {
	if a > c {
		return ((c<<2)+b)<<2 + a
	}
	return ((a<<2)+b)<<2 + c
}

func setStateArr(x int) [3]int {
	return [3]int{(x >> 4) & 3, (x >> 2) & 3, x & 3}
}

func myhash(n int) uint64 {
	// Sort only non-top levels; normalize dead levels (no valid moves) to 0,
	// matching the C++ solution's approach.  Keep the top (oracleData[n-1]) separate.
	b := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		s := setStateArr(oracleData[i])
		// "had": middle present AND at least one side present → alive level
		if s[1] > 0 && (s[0] > 0 || s[2] > 0) {
			b[i] = oracleData[i]
		} else {
			b[i] = 0 // dead level: collapse to 0 to reduce state space
		}
	}
	sort.Ints(b)
	var rtn uint64
	for _, v := range b {
		rtn = rtn*oracleP + uint64(v)
	}
	rtn = rtn*oracleP + uint64(oracleData[n-1])
	return rtn
}

func oracleDFS(n int) float64 {
	s := myhash(n)
	if f, ok := oracleMemo[s]; ok {
		return f
	}
	Inf := math.Inf(1)
	c := [4]float64{0, Inf, Inf, Inf}
	an := oracleData[n-1]
	for i := 0; i < n-1; i++ {
		ai := oracleData[i]
		x0 := setStateArr(ai)
		cntx := 0
		for _, v := range x0 {
			if v > 0 {
				cntx++
			}
		}
		if cntx <= 1 {
			continue
		}
		y0 := setStateArr(an)
		nn := n
		if y0[0] > 0 && y0[1] > 0 && y0[2] > 0 {
			nn = n + 1
		}
		for j, vj := range x0 {
			if vj == 0 {
				continue
			}
			if cntx == 2 && (j == 1 || x0[1] == 0) {
				continue
			}
			x := x0
			x[j] = 0
			cntx2 := 0
			for _, v := range x {
				if v > 0 {
					cntx2++
				}
			}
			var newAi int
			if cntx2 == 1 || x[1] == 0 {
				newAi = 0
			} else {
				newAi = getState(x[0], x[1], x[2])
			}
			oracleData[i] = newAi
			y0 := setStateArr(an)
			if nn > n {
				y0 = [3]int{0, 0, 0}
			}
			for k := 0; k < 3; k++ {
				if y0[k] != 0 {
					continue
				}
				y := y0
				y[k] = vj
				newAn := getState(y[0], y[1], y[2])
				if nn == n {
					oracleData[n-1] = newAn
					c[vj] = math.Min(c[vj], oracleDFS(nn))
					oracleData[n-1] = an
				} else {
					oracleData = append(oracleData, newAn)
					c[vj] = math.Min(c[vj], oracleDFS(nn))
					oracleData = oracleData[:len(oracleData)-1]
				}
			}
			oracleData[i] = ai
		}
	}
	if math.IsInf(c[1], 1) && math.IsInf(c[2], 1) && math.IsInf(c[3], 1) {
		oracleMemo[s] = 0
		return 0
	}
	p := 1.0 / 6.0
	if math.IsInf(c[1], 1) {
		p += 1.0 / 3.0
		c[1] = 0
	}
	if math.IsInf(c[2], 1) {
		p += 1.0 / 3.0
		c[2] = 0
	}
	if math.IsInf(c[3], 1) {
		p += 1.0 / 6.0
		c[3] = 0
	}
	f := (c[1]/3.0 + c[2]/3.0 + c[3]/6.0 + 1.0) / (1.0 - p)
	oracleMemo[s] = f
	return f
}

// tag matches 424E.go: G=1, B=2, R=3
var oracleTag = map[byte]int{'G': 1, 'B': 2, 'R': 3}

func oracle(n int, levels []string) float64 {
	oracleData = make([]int, n, 2*n+5)
	oracleMemo = make(map[uint64]float64)
	for i, lv := range levels {
		oracleData[i] = getState(oracleTag[lv[0]], oracleTag[lv[1]], oracleTag[lv[2]])
	}
	return oracleDFS(n)
}

// ---- test runner ----

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func formatInput(n int, levels []string) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, lv := range levels {
		sb.WriteString(lv)
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func parseOutput(out string) (float64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := fmt.Sscanf(fields[0], "%f", new(float64))
	_ = val
	var f float64
	fmt.Sscan(fields[0], &f)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q", fields[0])
	}
	return f, nil
}

func closeEnough(expected, actual float64) bool {
	diff := math.Abs(expected - actual)
	return diff <= tolerance*math.Max(1.0, math.Abs(expected))+1e-12
}

func randomLevel(rng *rand.Rand) string {
	colors := []byte{'R', 'G', 'B'}
	b := make([]byte, 3)
	for i := range b {
		b[i] = colors[rng.Intn(3)]
	}
	return string(b)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	type tc struct {
		n      int
		levels []string
	}

	tests := []tc{
		{2, []string{"RGB", "RGB"}},
		{3, []string{"RRR", "GGG", "BBB"}},
		{4, []string{"RRG", "GBB", "BRG", "GRB"}},
		{2, []string{"RRR", "GGG"}},
		{2, []string{"BBB", "RRR"}},
	}

	// exhaustive n=2
	colors := []byte{'R', 'G', 'B'}
	for m1 := 0; m1 < 27; m1++ {
		cur := m1
		lv := make([]byte, 3)
		for i := range lv {
			lv[i] = colors[cur%3]
			cur /= 3
		}
		tests = append(tests, tc{2, []string{string(lv), string(lv)}})
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 150 {
		n := rng.Intn(3) + 2 // 2..4; larger n causes oracle to be too slow
		levels := make([]string, n)
		for i := range levels {
			levels[i] = randomLevel(rng)
		}
		tests = append(tests, tc{n, levels})
	}

	for idx, t := range tests {
		input := formatInput(t.n, t.levels)
		expect := oracle(t.n, t.levels)

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\nInput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		var candVal float64
		if _, err := fmt.Sscan(candOut, &candVal); err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %q\nInput:\n%s\n", idx+1, candOut, input)
			os.Exit(1)
		}
		if !closeEnough(expect, candVal) {
			fmt.Printf("test %d failed: expected %.10f got %.10f\nInput:\n%s\n", idx+1, expect, candVal, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
