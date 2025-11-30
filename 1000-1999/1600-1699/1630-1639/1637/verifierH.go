package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded gzipped+base64 testcases from testcasesH.txt.
const encodedTestcases = `
H4sIAFoHK2kC/21TQZIEIQi78wq+oAL/f9pCAqNdtTUHaZQQQsZl6VZT1yMrfyGmoUtPZl08by1j153xRryQ94yr8sjO/BLLul11iRCVB+YCpgnxDF9HKp4bw1u+QQzsA2xrRtEMV3Nkx2G7k693ZzBVTIA7J0cNTOLgv5uPIQIvLVbkUedqFsgkC5c7D3tHT+vZ6bQanMrRlfHh3D/sqlq4DfQ26pyYC6+miufwj+aMbLI5Qt2pJdkaTnu4FGLIxi6xQXydRrRPBbTAXt8us5feOhQkgn3esddVknvhxqNUe/Qe9Uq3UXc6XWXJsWYcBem91XME/dp9yco+ehCLO6fzxqV0O70ScIc/HF5fE++NZy/RrrKPYuxwd8kNOCeHI2fCdne5519v/ZRrNS93/kP+AKoXLVy0AwAA
`

type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(i, val int) {
	for i <= f.n {
		f.bit[i] += val
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for i > 0 {
		s += f.bit[i]
		i -= i & -i
	}
	return s
}

type testCase struct {
	arr []int
}

// solve mirrors 1637H.go.
func solve(tc testCase) string {
	n := len(tc.arr)
	p := tc.arr
	pos := make([]int, n+1)
	for i, v := range p {
		pos[v] = i + 1
	}
	fw := NewFenwick(n)
	invA := make([]int, n+1)
	for val := 1; val <= n; val++ {
		idx := pos[val]
		greater := (val - 1) - fw.Sum(idx)
		invA[val] = invA[val-1] + greater
		fw.Add(idx, 1)
	}
	fw = NewFenwick(n)
	invB := make([]int, n+1)
	for val := n; val >= 1; val-- {
		idx := pos[val]
		smaller := fw.Sum(idx - 1)
		invB[val-1] = invB[val] + smaller
		fw.Add(idx, 1)
	}

	var sb strings.Builder
	for k := 0; k <= n; k++ {
		if k > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(invA[k] + invB[k]))
	}
	return sb.String()
}

func decodeTestcases() ([]testCase, error) {
	data, err := base64.StdEncoding.DecodeString(encodedTestcases)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var out bytes.Buffer
	if _, err := out.ReadFrom(r); err != nil {
		return nil, err
	}
	fields := strings.Fields(out.String())
	pos := 0
	var cases []testCase
	for pos < len(fields) {
		if pos >= len(fields) {
			break
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, err
		}
		pos++
		if pos+n > len(fields) {
			break
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[pos+i])
			if err != nil {
				return nil, err
			}
			arr[i] = val
		}
		pos += n
		cases = append(cases, testCase{arr: arr})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", len(tc.arr))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
