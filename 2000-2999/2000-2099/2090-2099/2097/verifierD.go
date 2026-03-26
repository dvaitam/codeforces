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

const (
	randomTests   = 120
	totalLenLimit = 200000
	maxNPerCase   = 50000
)

// Embedded correct solver for 2097D.
const embeddedSolver = `package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
)

type Bitset []uint64

func getRREF(str string, d, q int) []Bitset {
	wordLen := (q + 63) >> 6
	pool := make([]uint64, d*wordLen)
	basis := make([]Bitset, q)

	for i := 0; i < d; i++ {
		vc := Bitset(pool[i*wordLen : (i+1)*wordLen])
		for j := 0; j < q; j++ {
			if str[i*q+j] == '1' {
				vc[j>>6] |= 1 << (j & 63)
			}
		}

		c := 0
		for c < q {
			wordIdx := c >> 6
			bitIdx := c & 63
			word := vc[wordIdx] >> bitIdx
			if word == 0 {
				c = (wordIdx + 1) << 6
				continue
			}
			tz := bits.TrailingZeros64(word)
			c += tz
			if c >= q {
				break
			}
			if basis[c] != nil {
				for k := 0; k < wordLen; k++ {
					vc[k] ^= basis[c][k]
				}
				c++
			} else {
				basis[c] = vc
				break
			}
		}
	}

	pivots := make([]int, 0, d)
	for c := 0; c < q; c++ {
		if basis[c] != nil {
			pivots = append(pivots, c)
		}
	}

	for i := len(pivots) - 1; i >= 0; i-- {
		c := pivots[i]
		for j := 0; j < i; j++ {
			r_c := pivots[j]
			if (basis[r_c][c>>6] & (1 << (c & 63))) != 0 {
				for k := 0; k < wordLen; k++ {
					basis[r_c][k] ^= basis[c][k]
				}
			}
		}
	}

	return basis
}

func compareBases(b1, b2 []Bitset) bool {
	for i := 0; i < len(b1); i++ {
		if b1[i] == nil && b2[i] == nil {
			continue
		}
		if b1[i] == nil || b2[i] == nil {
			return false
		}
		for j := 0; j < len(b1[i]); j++ {
			if b1[i][j] != b2[i][j] {
				return false
			}
		}
	}
	return true
}

func solve() {
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 1024*1024*10)
	scanner.Buffer(buf, 1024*1024*10)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return
	}
	tStr := scanner.Text()
	t, _ := strconv.Atoi(tStr)

	for tc := 0; tc < t; tc++ {
		scanner.Scan()
		nStr := scanner.Text()
		n, _ := strconv.Atoi(nStr)

		scanner.Scan()
		s := scanner.Text()

		scanner.Scan()
		tgt := scanner.Text()

		q := n
		d := 1
		for q%2 == 0 {
			q /= 2
			d *= 2
		}

		b1 := getRREF(s, d, q)
		b2 := getRREF(tgt, d, q)

		if compareBases(b1, b2) {
			fmt.Println("Yes")
		} else {
			fmt.Println("No")
		}
	}
}

func main() {
	solve()
}
`

type testCase struct {
	n int
	s string
	t string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer cleanup()

	tests := generateTests()
	input := buildInput(tests)

	expectRaw, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference failed: %v\n%s", err, expectRaw)
	}
	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}

	expect := parseOutputs(expectRaw, len(tests))
	got := parseOutputs(gotRaw, len(tests))

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch at test %d: expected %s, got %s", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "2097D-ref-")
	if err != nil {
		return "", nil, err
	}
	srcPath := filepath.Join(dir, "ref2097D.go")
	if err := os.WriteFile(srcPath, []byte(embeddedSolver), 0644); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to write embedded solver: %v", err)
	}
	binPath := filepath.Join(dir, "ref2097D.bin")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("%v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+6)

	// Deterministic coverage.
	tests = append(tests,
		testCase{n: 1, s: "0", t: "0"},
		testCase{n: 1, s: "1", t: "0"},
		testCase{n: 2, s: "00", t: "11"},
		testCase{n: 2, s: "10", t: "10"},
		testCase{n: 4, s: "0000", t: "0000"},
		testCase{n: 4, s: "0101", t: "1001"},
	)

	total := 0
	for _, tc := range tests {
		total += tc.n
	}

	for i := 0; i < randomTests && total < totalLenLimit; i++ {
		remain := totalLenLimit - total
		maxN := maxNPerCase
		if maxN > remain {
			maxN = remain
		}
		if maxN < 1 {
			break
		}
		n := rng.Intn(maxN) + 1
		total += n

		var sbS, sbT strings.Builder
		sbS.Grow(n)
		sbT.Grow(n)
		ensureOneS := rng.Intn(2) == 0
		ensureOneT := rng.Intn(2) == 0
		for j := 0; j < n; j++ {
			bs := byte('0' + rng.Intn(2))
			bt := byte('0' + rng.Intn(2))
			if ensureOneS && j == n-1 && !strings.Contains(sbS.String(), "1") {
				bs = '1'
			}
			if ensureOneT && j == n-1 && !strings.Contains(sbT.String(), "1") {
				bt = '1'
			}
			sbS.WriteByte(bs)
			sbT.WriteByte(bt)
		}
		tests = append(tests, testCase{n: n, s: sbS.String(), t: sbT.String()})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(tc.t)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("stderr not empty")
	}
	return out.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutputs(out string, t int) []string {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		fail("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]string, t)
	for i, tok := range tokens {
		res[i] = strings.ToUpper(tok)
	}
	return res
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
