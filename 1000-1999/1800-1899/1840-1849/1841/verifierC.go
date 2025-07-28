package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var val = []int64{1, 10, 100, 1000, 10000}

const negInf int64 = -1 << 60

type testCaseC struct {
	input    string
	expected string
}

func runCandidate(bin, input string) (string, error) {
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

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solveC(s string) int64 {
	n := len(s)
	dp0 := make([]int64, 5)
	dp1 := make([]int64, 5)
	for i := 0; i < 5; i++ {
		dp0[i] = 0
		dp1[i] = negInf
	}
	for idx := n - 1; idx >= 0; idx-- {
		orig := int(s[idx] - 'A')
		newdp0 := [5]int64{negInf, negInf, negInf, negInf, negInf}
		newdp1 := [5]int64{negInf, negInf, negInf, negInf, negInf}
		for pm := 0; pm < 5; pm++ {
			if dp0[pm] != negInf {
				sign := val[orig]
				if orig < pm {
					sign = -sign
				}
				nm := pm
				if orig > nm {
					nm = orig
				}
				newdp0[nm] = max64(newdp0[nm], sign+dp0[pm])
				for ni := 0; ni < 5; ni++ {
					sign2 := val[ni]
					if ni < pm {
						sign2 = -sign2
					}
					nm2 := pm
					if ni > nm2 {
						nm2 = ni
					}
					newdp1[nm2] = max64(newdp1[nm2], sign2+dp0[pm])
				}
			}
			if dp1[pm] != negInf {
				sign := val[orig]
				if orig < pm {
					sign = -sign
				}
				nm := pm
				if orig > nm {
					nm = orig
				}
				newdp1[nm] = max64(newdp1[nm], sign+dp1[pm])
			}
		}
		for j := 0; j < 5; j++ {
			dp0[j] = newdp0[j]
			dp1[j] = newdp1[j]
		}
	}
	ans := negInf
	for i := 0; i < 5; i++ {
		ans = max64(ans, dp0[i])
		ans = max64(ans, dp1[i])
	}
	return ans
}

func generateCaseC(rng *rand.Rand) testCaseC {
	t := rng.Intn(3) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for j := 0; j < t; j++ {
		n := rng.Intn(15) + 1
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = byte('A' + rng.Intn(5))
		}
		s := string(b)
		in.WriteString(fmt.Sprintf("%s\n", s))
		out.WriteString(fmt.Sprintf("%d\n", solveC(s)))
	}
	return testCaseC{input: in.String(), expected: out.String()}
}

func runCaseC(bin string, tc testCaseC) error {
	got, err := runCandidate(bin, tc.input)
	if err != nil {
		return err
	}
	got = strings.TrimSpace(got)
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseC{generateCaseC(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseC(rng))
	}
	for i, tc := range cases {
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
