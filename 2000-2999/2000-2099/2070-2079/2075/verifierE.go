package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input string
	t     int
}

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for i := 0; i < t; i++ {
		var n, m, A, B int64
		fmt.Fscan(reader, &n, &m, &A, &B)

		A_mod := A % MOD
		B_mod := B % MOD

		dp := make([]int64, 16)
		nextDp := make([]int64, 16)
		dp[15] = 1

		for k := 29; k >= 0; k-- {
			bitA := (A >> k) & 1
			bitB := (B >> k) & 1
			for j := 0; j < 16; j++ {
				nextDp[j] = 0
			}

			for mask := 0; mask < 16; mask++ {
				if dp[mask] == 0 {
					continue
				}
				tx1 := (mask >> 3) & 1
				tx2 := (mask >> 2) & 1
				ty1 := (mask >> 1) & 1
				ty2 := mask & 1

				for bx1 := int64(0); bx1 <= 1; bx1++ {
					if tx1 == 1 && bx1 > bitA {
						continue
					}
					ntx1 := tx1
					if bx1 < bitA {
						ntx1 = 0
					}

					for bx2 := int64(0); bx2 <= 1; bx2++ {
						if tx2 == 1 && bx2 > bitA {
							continue
						}
						ntx2 := tx2
						if bx2 < bitA {
							ntx2 = 0
						}

						for by1 := int64(0); by1 <= 1; by1++ {
							if ty1 == 1 && by1 > bitB {
								continue
							}
							nty1 := ty1
							if by1 < bitB {
								nty1 = 0
							}

							by2 := bx1 ^ bx2 ^ by1
							if ty2 == 1 && by2 > bitB {
								continue
							}
							nty2 := ty2
							if by2 < bitB {
								nty2 = 0
							}

							nmask := (ntx1 << 3) | (ntx2 << 2) | (nty1 << 1) | nty2
							nextDp[nmask] = (nextDp[nmask] + dp[mask]) % MOD
						}
					}
				}
			}
			dp, nextDp = nextDp, dp
		}

		var C int64 = 0
		for mask := 0; mask < 16; mask++ {
			C = (C + dp[mask]) % MOD
		}

		ways0 := (A_mod + 1) * (B_mod + 1) % MOD

		C = (C - ways0 + MOD) % MOD
		C = (C * 748683265) % MOD

		powN := (power(2, n) - 2 + MOD) % MOD
		powM := (power(2, m) - 2 + MOD) % MOD

		ans := ways0

		ans = (ans + (A_mod+1)*(B_mod+1)%MOD*B_mod%MOD*499122177%MOD*powM) % MOD
		ans = (ans + (B_mod+1)*(A_mod+1)%MOD*A_mod%MOD*499122177%MOD*powN) % MOD
		ans = (ans + powN*powM%MOD*C) % MOD

		fmt.Fprintln(writer, ans)
	}
}

func power(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % MOD
		}
		base = (base * base) % MOD
		exp /= 2
	}
	return res
}
`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup := buildRef()
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		expect, err := parseOutputs(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutputs(gotOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}

		if len(expect) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d: output count mismatch, expected %d got %d\nInput:\n%sOutput:\n%s", idx+1, len(expect), len(got), tc.input, gotOut)
			os.Exit(1)
		}
		for i := range expect {
			if expect[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch at case %d, expected %d got %d\nInput:\n%sOutput:\n%s", idx+1, i+1, expect[i], got[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildRef() (string, func()) {
	tmpDir, err := os.MkdirTemp("", "ref2075E")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	srcPath := filepath.Join(tmpDir, "ref.go")
	os.WriteFile(srcPath, []byte(refSource), 0644)
	binPath := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build ref: %v\n%s\n", err, string(out))
		os.Exit(1)
	}
	return binPath, func() { os.RemoveAll(tmpDir) }
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() []testCase {
	tests := make([]testCase, 0, 50)
	tests = append(tests, sampleTest())
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func sampleTest() testCase {
	input := "6\n2 2 2 2\n2 3 4 5\n5 7 4 3\n1337 42 1337 42\n4 2 13 3\n753687090 2 536370902 536390912\n"
	return testCase{input: input, t: 6}
}

func randomTest(rng *rand.Rand) testCase {
	t := rng.Intn(5) + 1 // 1..5 cases
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := randRange(rng, 2, 200)
		m := randRange(rng, 2, 200)
		A := randRange(rng, 2, 1<<20)
		B := randRange(rng, 2, 1<<20)
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, A, B))
	}
	return testCase{input: sb.String(), t: t}
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}
