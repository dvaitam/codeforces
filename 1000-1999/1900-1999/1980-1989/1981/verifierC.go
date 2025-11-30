package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesC = `100
2 4 -1
2 4 -1
3 5 -1 -1
3 -1 2 3
3 -1 -1 5
4 -1 3 -1 -1
3 -1 -1 4
4 -1 -1 -1 1
2 4 -1
2 1 3
4 3 5 4 -1
3 -1 4 -1
3 4 -1 2
4 -1 4 -1 4
4 -1 -1 -1 -1
4 3 5 -1 -1
3 1 2 5
3 3 -1 2
2 -1 -1
4 -1 -1 2 -1
3 -1 5 3
3 3 -1 -1
2 4 3
3 5 1 -1
2 -1 3
3 -1 1 -1
2 5 3
3 5 -1 4
3 1 -1 2
3 1 4 4
4 -1 1 -1 1
2 4 -1
2 -1 3
2 5 -1
3 -1 -1 -1
2 -1 -1
3 -1 -1 -1
2 1 2
4 -1 4 -1 2
4 -1 -1 1 5
2 2 5
4 -1 -1 5 -1
3 2 -1 3
2 -1 4
4 2 3 2 1
4 -1 -1 -1 -1
2 -1 3
2 5 -1
3 2 4 -1
3 -1 -1 -1
4 -1 2 5 -1
3 -1 1 -1
4 -1 5 1 -1
2 -1 -1
3 -1 5 2
3 3 4 -1
2 5 5
4 3 3 -1 -1
4 -1 2 -1 2
3 4 1 3
4 3 5 3 3
2 -1 3
2 -1 1
4 2 1 -1 -1
4 -1 -1 4 5
2 4 -1
4 5 -1 4 -1
4 -1 1 5 2
3 1 5 -1
2 -1 5
2 -1 2
2 1 2
2 -1 -1
2 -1 5
3 -1 3 2
4 -1 4 -1 4
2 -1 3
4 4 -1 -1 -1
3 5 1 -1
2 -1 3
2 -1 4
3 2 -1 2
4 3 4 -1 1
3 4 -1 4
4 3 1 -1 4
4 -1 -1 3 1
2 -1 3
2 5 -1
4 2 4 2 -1
2 -1 -1
4 -1 -1 4 5
2 -1 -1
4 -1 -1 -1 -1
2 2 2
2 3 -1
4 -1 1 4 -1
2 -1 -1
4 2 5 5 4
4 3 5 -1 -1
4 2 1 3 -1`

func get(x int64) int {
	for i := 30; i >= 0; i-- {
		if (x>>i)&1 == 1 {
			return i
		}
	}
	return -1
}

func solveCase(a []int64) string {
	n := len(a)
	const inf = int64(1e9)
	b := make([]int64, n+2)
	flag := true
	l := 1
	for l <= n {
		if a[l-1] != -1 {
			b[l] = a[l-1]
			l++
		} else {
			r := l
			for r <= n && a[r-1] == -1 {
				r++
			}
			if l == 1 {
				if r <= n {
					b[r] = a[r-1]
				} else {
					b[r] = 1
				}
				for i := r - 1; i >= 1; i-- {
					if 2*b[i+1] <= inf {
						b[i] = 2 * b[i+1]
					} else {
						b[i] = b[i+1] / 2
					}
				}
			} else if r == n+1 {
				for i := l; i < r; i++ {
					if b[i-1]*2 <= inf {
						b[i] = b[i-1] * 2
					} else {
						b[i] = b[i-1] / 2
					}
				}
			} else {
				L := a[l-2]
				R := a[r-1]
				can := false
				length := r - l + 1
				for suff := 0; suff <= 30 && suff <= length; suff++ {
					cur := L >> suff
					if cur == 0 {
						continue
					}
					tmp := (r - l) - suff + 1
					v1 := get(cur)
					v2 := get(R)
					if v1 > v2 {
						continue
					}
					need := v2 - v1
					value := R >> need
					if cur != value {
						continue
					}
					if need <= tmp && (tmp-need)%2 == 0 {
						can = true
						for i := l; i < l+suff; i++ {
							b[i] = b[i-1] / 2
						}
						x := need
						for i := l + suff; i < l+suff+need; i++ {
							x--
							bit := (R >> x) & 1
							b[i] = b[i-1]*2 + bit
						}
						start := l + suff + need
						for i := start; i < r; i++ {
							if i%2 == start%2 {
								b[i] = b[i-1] * 2
							} else {
								b[i] = b[i-1] / 2
							}
						}
						break
					}
				}
				if !can {
					flag = false
				}
			}
			l = r
		}
	}
	for i := 1; i < n; i++ {
		if !(b[i] == b[i+1]/2 || b[i+1] == b[i]/2) {
			flag = false
			break
		}
	}
	if !flag {
		return "-1"
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(b[i], 10))
	}
	return sb.String()
}

type testCase struct {
	arr []int64
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesC)
	pos := 0
	readInt := func() (int64, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.ParseInt(fields[pos], 10, 64)
		pos++
		return v, err
	}
	t64, err := readInt()
	if err != nil {
		return nil, err
	}
	t := int(t64)
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n64, err := readInt()
		if err != nil {
			return nil, err
		}
		n := int(n64)
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			v, err := readInt()
			if err != nil {
				return nil, err
			}
			arr[j] = v
		}
		tests[i] = testCase{arr: arr}
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.arr)))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}

	allInput := buildAllInput(tests)
	allOutput, err := runCandidate(bin, allInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime error:", err)
		os.Exit(1)
	}
	outLines := strings.Split(strings.TrimSpace(allOutput), "\n")
	if len(outLines) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outLines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solveCase(tc.arr)
		if strings.TrimSpace(outLines[i]) != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput: %v\nexpected: %s\ngot: %s\n", i+1, tc.arr, want, outLines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
