package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Base64-encoded contents of testcasesB.txt.
const testcasesB = "MyAzIDEgMiAxIDIgMgo1IDQgMiAxIDQgMSA0IDQgNSAxIDQKNCAyIDEgMyAxIDEgMSAxIDQKNyAyIDQgNiAxIDUgMiA3IDQgNCA1IDIgMyAyIDYKMyAyIDIgMSAyIDMgMwoyIDEgMiAxIDIKNyA2IDUgNCA1IDcgNiAyIDMgMyA1IDQgNyA1IDQKNiAxIDQgMiA2IDQgNCA2IDIgMyA1IDYgNgo3IDMgMSA0IDYgNSAxIDcgMiA1IDcgNCAzIDQgNgoyIDIgMSAyIDIKNyAyIDIgNSAyIDEgNyAyIDUgNyA1IDIgNCA1IDMKOCA2IDggNSAxIDcgMyA0IDcgMSA4IDYgNCA3IDggNiA3CjQgMSAzIDQgMSAyIDIgMiAxCjggNSAxIDIgMiAxIDggMSA1IDQgNSAyIDMgNiA1IDIgMwozIDIgMyAxIDMgMiAzCjcgMyA0IDYgMyA0IDQgMSAxIDMgNCAzIDQgNyAyCjQgMSAzIDIgNCAxIDIgMSA0CjMgMSAzIDEgMiAzIDMKNyA0IDUgNyAyIDYgNyA2IDUgNCAyIDUgNiAxIDQKNyA1IDcgMyA2IDYgNCAxIDYgMyAyIDIgMSAzIDEKOCAyIDUgNSAzIDcgNSAzIDEgMSA0IDggMyAxIDcgNCA2CjIgMSAyIDEgMgoyIDIgMiAyIDEKNCA0IDMgMSAyIDIgMyAyIDMKNSAyIDMgMSA0IDUgMyA1IDQgNSAyCjIgMSAxIDEgMQozIDMgMSAyIDIgMyAzCjggNSA2IDYgNiAyIDUgNCA4IDMgMiA2IDEgNyAyIDcgMwo4IDMgNiAyIDcgMiA0IDIgNSA2IDUgMiA4IDUgMiAxIDUKMiAxIDEgMiAxCjggMSA0IDQgNyAzIDIgOCAzIDQgMyAyIDcgNyA1IDUgOAo0IDEgMiAzIDEgMSAxIDMgMwo1IDQgMyA0IDEgMSAzIDUgNCAxIDMKMyAzIDMgMyAyIDMgMgo0IDIgMiAzIDIgMiAzIDEgMwoyIDIgMSAyIDEKNSAzIDEgMyAyIDMgNSAzIDIgMyAxCjYgNSA1IDUgMSAyIDIgMSAyIDQgMSAzIDUKOCAyIDIgMSAxIDUgNiA4IDggMyAyIDYgMiAzIDMgMyAzCjggNiA1IDIgNSAzIDQgMyAxIDYgNCAzIDUgNyAzIDEgNAo0IDEgNCA0IDMgNCA0IDEgNAo4IDYgMyA1IDggMSA3IDEgMSA2IDMgMyAzIDUgNSA3IDcKMyAzIDEgMSAyIDEgMQo2IDMgNSA2IDQgNiA2IDYgMiAyIDMgNCA2CjUgMiA0IDMgNSA1IDMgMiAxIDEgNQo3IDMgMiA1IDcgNyAyIDMgMyA2IDMgNyA1IDMgMgo3IDYgNiA0IDUgMSA3IDEgNSA1IDUgNCAyIDIgMwo1IDIgNSAxIDQgNCAzIDQgNSAyIDUKNyAxIDUgMSA3IDMgNiAxIDMgNiAxIDIgNyA1IDcKNyA2IDYgMSA0IDcgMiA3IDQgNyA0IDQgMiAzIDQKMyAzIDIgMSAxIDIgMwo2IDQgMSA2IDMgMyAyIDQgNiA1IDEgMiA1CjUgNSAxIDEgNSAyIDMgMiAyIDMgMgo2IDIgMyAzIDUgMyA2IDQgMiA1IDMgNCA0CjggMiA0IDcgNCA1IDIgMSAyIDEgNSAzIDIgNiA1IDcgNgo4IDYgMSAyIDggOCA2IDUgNyA2IDggMiA3IDcgNCAxIDUKNyA1IDYgNiA3IDYgNSAyIDQgNSA3IDUgNCA2IDYKNCAyIDQgMiAzIDEgNCA0IDQKNCAxIDQgMiAzIDEgNCAyIDQKOCA1IDMgMiAxIDYgNSA3IDUgMyA4IDUgOCAzIDggMSA1CjYgMSA2IDUgNCAxIDMgMSA2IDQgMSAyIDUKNyAyIDYgMSA0IDYgNiAzIDUgMyAyIDUgMiAyIDMKNCAxIDEgMyA0IDEgMiAzIDMKNCAyIDQgNCAyIDQgMyAzIDIKNCAyIDEgNCAzIDQgMiAzIDIKMiAxIDIgMSAyCjUgNSAyIDIgMiA0IDMgMyA0IDIgMQo3IDIgNiA2IDMgMSAxIDIgNCAzIDQgMSAyIDEgMQo4IDEgNCAxIDggOCA2IDUgMiAzIDIgNCA3IDQgOCA4IDcKOCAzIDQgNCA1IDggNyA0IDggNSA2IDggMiA0IDIgMSAxCjggMSA4IDYgNyA1IDQgNyAzIDMgMSAxIDcgMyAxIDcgNQozIDEgMiAzIDIgMSAxCjYgMSA1IDIgMSAzIDEgNCAxIDIgMSA0IDYKMyAzIDIgMyAxIDMgMgo1IDMgMyAzIDIgMiAxIDUgNSAyIDMKNSA1IDUgNSAxIDMgNSA0IDUgMiA1CjUgMSAzIDUgMSAzIDIgMSAyIDEgMgo4IDcgMSAxIDIgOCA2IDIgNiAxIDMgMSA4IDMgNyA4IDEKNyA1IDMgMSAzIDcgMyAxIDMgMSA3IDQgMSA2IDMKNCAyIDMgNCAxIDMgMSA0IDIKNiA1IDIgMyAzIDUgNCA1IDQgMSAyIDYgNAo2IDUgNiA1IDYgNSA1IDEgMyA2IDIgMiAzCjUgNSAzIDEgNCAzIDIgNSAxIDEgMwo4IDYgNyA1IDYgNiA1IDYgMSAyIDMgNiA2IDYgMiA4IDUKNSA0IDMgNCAxIDUgMSAyIDEgNSA0CjYgMyAyIDYgNSA2IDMgMyA2IDMgNCAzIDQKNiAzIDUgNSAyIDEgMiAzIDYgMiA1IDIgMQozIDIgMyAzIDEgMSAzCjcgMyA2IDEgMiAzIDEgNiA1IDUgNiAxIDcgMSA3CjggNCAzIDcgMSA2IDggNSA0IDQgOCA0IDcgOCA2IDQgOAo3IDEgNyA3IDMgNCAyIDEgNiA1IDcgNCA1IDQgMQo1IDUgNSA1IDUgNCAxIDMgNCAxIDIKNCAxIDEgMyAzIDMgNCA0IDMKNSAzIDIgNSA0IDUgMiA1IDIgMyAxCjUgNSAxIDMgNCA0IDMgMSAxIDEgMQo1IDMgNCAzIDMgNCAzIDQgNCAxIDQKNCAyIDQgMiAxIDIgMyAzIDIKNiAzIDQgMyA1IDMgNiA0IDYgMyA0IDMgNAozIDMgMiAyIDMgMiAxCjIgMSAxIDEgMQo="

type testCase struct {
	n   int
	a   []int
	b   []int
}

// Embedded solver logic from 1830B.go.
func solve(tc testCase) string {
	n := tc.n
	a := tc.a
	b := tc.b

	freq := make(map[int]map[int]int)
	seen := make(map[int]bool)
	uniqueA := []int{}
	for i := 0; i < n; i++ {
		ai := a[i]
		bi := b[i]
		if freq[ai] == nil {
			freq[ai] = make(map[int]int)
		}
		freq[ai][bi]++
		if !seen[ai] {
			seen[ai] = true
			uniqueA = append(uniqueA, ai)
		}
	}
	sort.Ints(uniqueA)
	limit := 2 * n
	var ans int64
	for ix := 0; ix < len(uniqueA); ix++ {
		x := uniqueA[ix]
		mapX := freq[x]
		for iy := ix; iy < len(uniqueA); iy++ {
			y := uniqueA[iy]
			if x*y > limit {
				break
			}
			mapY := freq[y]
			p := x * y
			if x == y {
				for b1, c1 := range mapX {
					b2 := p - b1
					if b2 < 1 || b2 > n {
						continue
					}
					if b2 < b1 {
						continue
					}
					if c2, ok := mapX[b2]; ok {
						if b1 == b2 {
							ans += int64(c1) * int64(c1-1) / 2
						} else {
							ans += int64(c1) * int64(c2)
						}
					}
				}
			} else {
				if len(mapX) <= len(mapY) {
					for b1, c1 := range mapX {
						b2 := p - b1
						if b2 < 1 || b2 > n {
							continue
						}
						if c2 := mapY[b2]; c2 > 0 {
							ans += int64(c1) * int64(c2)
						}
					}
				} else {
					for b2, c2 := range mapY {
						b1 := p - b2
						if b1 < 1 || b1 > n {
							continue
						}
						if c1 := mapX[b1]; c1 > 0 {
							ans += int64(c1) * int64(c2)
						}
					}
				}
			}
		}
	}
	return strconv.FormatInt(ans, 10)
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesB)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Split(bufio.ScanWords)
	cases := []testCase{}
	for sc.Scan() {
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n: %v", err)
		}
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d missing a[%d]", len(cases)+1, i)
			}
			a[i], err = strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d a[%d]: %v", len(cases)+1, i, err)
			}
		}
		for i := 0; i < n; i++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d missing b[%d]", len(cases)+1, i)
			}
			b[i], err = strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d b[%d]: %v", len(cases)+1, i, err)
			}
		}
		cases = append(cases, testCase{n: n, a: a, b: b})
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
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
		fmt.Fprintf(&input, "1\n%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for i, v := range tc.b {
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
