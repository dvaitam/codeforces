package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Base64-encoded contents of testcasesE.txt.
const testcasesE = "MTAwCjIgMyAyIDIgMgoyIDEgMiAxIDIKMyAxIDMgMyAzIDIKMyAyIDIgMiAxIDEKMyAzIDMgMSAyIDIKMyAxIDYgMyA2IDQKMiAyIDMgMiAyCjEgMyAyIDIKMyAyIDUgNCA1IDMKNCA0IDMgMSAyIDIgMQoxIDEgNSAzCjQgMyAzIDMgMSAxIDIKMiA0IDQgMiAzCjQgMyA2IDIgMyAxIDEKMiAzIDYgNSAyCjEgMyAzIDIKNCAxIDIgMiAxIDIgMgoxIDMgNCAzCjIgNCA2IDYgMQozIDIgNSAzIDIgMwo0IDIgNCAxIDMgMSA0CjIgMyA0IDMgMQo0IDIgNSAyIDEgMSAxCjEgMiA2IDYKMiAxIDYgNCA1CjIgMyAyIDEgMgo0IDIgNSAyIDIgNCA0CjQgMSAzIDIgMiAxIDMKNCAyIDUgMiAxIDEgMwozIDIgNiAyIDIgNAozIDIgNCAxIDMgMQo0IDEgNSA0IDEgNCAyCjIgMyA0IDQgMwo0IDIgNCAzIDQgNCAxCjMgMiAyIDIgMSAyCjEgMiA1IDUKNCA0IDUgMiAxIDIgMgoxIDMgMiAyCjQgMiA2IDEgMiAyIDYKMyA0IDYgNCAxIDEKMSAxIDUgNQozIDIgMiAyIDEgMQozIDMgMiAxIDIgMgoyIDMgNSAyIDQKNCAxIDIgMSAyIDIgMgo0IDEgMyAzIDIgMiAyCjEgMiA0IDIKMiAyIDQgNCAzCjMgMSAzIDEgMyAyCjEgMSA1IDUKNCAxIDYgMiA0IDMgMwoyIDIgMiAxIDIKNCAzIDYgNCA1IDUgNAo0IDIgNCAzIDMgMiA0CjEgMiA2IDIKMyAzIDMgMyAyIDMKMSAxIDUgMgozIDMgNCAyIDMgMQoxIDEgNCAxCjIgMiA2IDQgNgoxIDEgMyAzCjQgMiA2IDYgMyA2IDQKMiA0IDMgMiAyCjIgMyA2IDYgNQo0IDQgNiA0IDMgMyA0CjQgMSA0IDIgNCAzIDQKNCAxIDMgMyAxIDIgMwozIDMgMiAyIDIgMQo0IDEgMyAxIDEgMiAxCjEgMiA1IDUKMyAxIDYgMiAyIDEKMiAxIDQgMSAxCjMgMSA0IDIgMiA0CjIgMSA2IDMgNgoxIDIgNiA0CjIgMiA0IDQgMQo0IDIgMyAyIDMgMiAzCjMgNCA0IDIgNCAxCjQgMSA0IDEgMyAyIDEKMiAxIDYgMSAyCjIgMSA1IDUgMgo0IDQgNCAyIDIgMSA0CjMgMyA1IDUgMSAyCjMgMyA1IDMgMSAyCjMgMSAyIDIgMiAyCjEgMiAyIDIKMiAxIDMgMyAxCjMgMiA0IDEgMSAxCjIgMyA2IDEgNgoxIDEgNSAyCjIgMyAzIDMgMwozIDMgNSA1IDQgMQoyIDQgNSAxIDUKMyAzIDMgMyAzIDIKMyAyIDMgMSAyIDIKMyAzIDUgNCAzIDUKMiAzIDQgMSAxCjQgMyA2IDMgNiA2IDIKMyAxIDMgMyAxIDMKMiAzIDMgMyAyCg=="

type testCase struct {
	n, k int
	arr  []int
}

// Embedded correct solver logic from the ACCEPTED CF solution.
func solve(tc testCase) string {
	n := tc.n
	k := tc.k
	a := tc.arr

	const maxDiff = 45
	const offset = maxDiff
	const numDiffs = maxDiff*2 + 1

	dp := make([]int, numDiffs*(k+1))
	nextDp := make([]int, numDiffs*(k+1))

	dp[offset*(k+1)+0] = 1

	for j := 0; j < n; j++ {
		x := a[j]

		for i := 0; i < len(nextDp); i++ {
			nextDp[i] = 0
		}

		for diff := -maxDiff; diff <= maxDiff; diff++ {
			diffOffset := diff + offset
			baseIdx := diffOffset * (k + 1)

			for dist := 0; dist <= k; dist++ {
				val := dp[baseIdx+dist]
				if val == 0 {
					continue
				}

				for y := 0; y <= 1; y++ {
					newDiff := diff + y - x
					if newDiff < -maxDiff || newDiff > maxDiff {
						continue
					}

					absDiff := newDiff
					if absDiff < 0 {
						absDiff = -absDiff
					}

					newDist := dist + absDiff
					if newDist <= k {
						nIdx := (newDiff+offset)*(k+1) + newDist
						nextDp[nIdx] += val
						if nextDp[nIdx] >= 1000000007 {
							nextDp[nIdx] -= 1000000007
						}
					}
				}
			}
		}

		dp, nextDp = nextDp, dp
	}

	ans := 0
	baseIdx := offset * (k + 1)
	for d := 0; d <= k; d++ {
		if (k-d)%2 == 0 {
			ans += dp[baseIdx+d]
			if ans >= 1000000007 {
				ans -= 1000000007
			}
		}
	}

	return strconv.Itoa(ans)
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesE)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("invalid test data")
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d n: %v", i+1, err)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing k", i+1)
		}
		k, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d k: %v", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d missing a[%d]", i+1, j)
			}
			arr[j], err = strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d a[%d]: %v", i+1, j, err)
			}
		}
		cases = append(cases, testCase{n: n, k: k, arr: arr})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		want := solve(tc)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
