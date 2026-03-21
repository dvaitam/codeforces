package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	x int64
	y int64
	z int64
	k int64
}

// Embedded solver for 2035E.
func getCost(x, y, D, attacks int64) int64 {
	if x > 0 && D > 2000000000000000000/x {
		return 2000000000000000000
	}
	return x*D + y*attacks
}

func solveCase(x, y, z, k int64) int64 {
	ans := int64(2000000000000000000)

	for q := int64(0); ; q++ {
		Rq := z - q*(q+1)/2*k
		if Rq <= 0 {
			low := int64(1)
			high := q * k
			if high < 1 {
				high = 1
			}
			ans_D := q * k
			for low <= high {
				mid := low + (high-low)/2
				c := mid / k
				if c > q {
					c = q
				}
				dmg := c*(c+1)/2*k + (q-c)*mid
				if dmg >= z {
					ans_D = mid
					high = mid - 1
				} else {
					low = mid + 1
				}
			}
			cost := getCost(x, y, ans_D, q)
			if cost < ans {
				ans = cost
			}
			break
		}

		L := q * k
		if L < 1 {
			L = 1
		}
		R_max := (q+1)*k - 1

		pL := (Rq + L - 1) / L
		costL := getCost(x, y, L, q+pL)
		if costL < ans {
			ans = costL
		}

		pR := (Rq + R_max - 1) / R_max
		costR := getCost(x, y, R_max, q+pR)
		if costR < ans {
			ans = costR
		}

		if x >= y {
			D_unc := math.Sqrt(float64(y) * float64(Rq) / float64(x))
			Dc := int64(D_unc)
			D_center := Dc
			if D_center < L {
				D_center = L
			}
			if D_center > R_max {
				D_center = R_max
			}
			start := D_center - 500
			if start < L {
				start = L
			}
			end := D_center + 500
			if end > R_max {
				end = R_max
			}
			for D := start; D <= end; D++ {
				p := (Rq + D - 1) / D
				cost := getCost(x, y, D, q+p)
				if cost < ans {
					ans = cost
				}
			}
		} else {
			p_unc := math.Sqrt(float64(x) * float64(Rq) / float64(y))
			pc := int64(p_unc)
			P_min := (Rq + R_max - 1) / R_max
			if P_min < 1 {
				P_min = 1
			}
			P_max := (Rq + L - 1) / L

			p_center := pc
			if p_center < P_min {
				p_center = P_min
			}
			if p_center > P_max {
				p_center = P_max
			}
			start := p_center - 500
			if start < 1 {
				start = 1
			}
			end := p_center + 500
			for p := start; p <= end; p++ {
				D := (Rq + p - 1) / p
				if D < L {
					D = L
				}
				if D <= R_max {
					cost := getCost(x, y, D, q+p)
					if cost < ans {
						ans = cost
					}
				}
			}
		}
	}
	return ans
}

func solveAll(tests []testCase) []int64 {
	results := make([]int64, len(tests))
	for i, tc := range tests {
		results[i] = solveCase(tc.x, tc.y, tc.z, tc.k)
	}
	return results
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	input := buildInput(tests)

	expected := solveAll(tests)

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d, got %d\nInput: x=%d y=%d z=%d k=%d\n", i+1, expected[i], got[i], tests[i].x, tests[i].y, tests[i].z, tests[i].k)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(x, y, z, k int64) {
		tests = append(tests, testCase{x: x, y: y, z: z, k: k})
	}
	add(1, 1, 1, 1)
	add(2, 3, 5, 5)
	add(10, 20, 40, 5)
	add(1, 60, 100, 10)
	add(60, 1, 100, 10)
	add(1, 100000000, 100000000, 1)
	add(100000000, 1, 100000000, 100000000)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		tests = append(tests, randomCase(rng))
	}
	return tests
}

func randomCase(rng *rand.Rand) testCase {
	x := int64(rng.Intn(1_0000_0000) + 1)
	y := int64(rng.Intn(1_0000_0000) + 1)
	z := int64(rng.Intn(1_0000_0000) + 1)
	k := int64(rng.Intn(1_0000_0000) + 1)
	return testCase{x: x, y: y, z: z, k: k}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.x, tc.y, tc.z, tc.k))
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok, i+1)
		}
		if val < 0 {
			return nil, fmt.Errorf("negative cost %d at position %d", val, i+1)
		}
		res[i] = val
	}
	return res, nil
}
