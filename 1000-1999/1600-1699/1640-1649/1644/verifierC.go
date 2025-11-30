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

// Embedded gzipped+base64 testcases from testcasesC.txt.
const encodedTestcases = `
H4sIAFcIK2kC/01VSW7EMAy7+xX+QABr6/L/j1Ui6UzRYCbjhaIoSrVzVu3YT85T++nXWLbPjvW9sWPz4fPUzl7MWfI9P6vPPbm+Nm77Tq74yvkyoPo+K3qvf9Ucjoni8yPXL27MqVL4wRiwjnMAjoPY92HmIlPNMbctO/8jEWPYGzKpjhBCyFlz7nYU8BCTfkfGvbmGQKfWSwcpIgRSLsACv7EB62t06JdoEQHf2eI0NEFcKNh/S1veW/YKOrrUhAAecpoQKATyaagJNr9C2RaioAiPwnYyy1W+QuGMJ6l0PwWxX5kk9uRiV1cDCJGZWODarbNyIR6w12R4mNfU8wdnXeEDF/xa5pG7WILTMqDoujyPoRA4fHgrVFu/icwRpjGSftFXWBiFHCs8n/LcIiGGpoNLV6QGZYTpjSrR39IIAc6iv5+A4B3nmjMI75KwGaYqyPNrEHzTEU9QfGhLKqMVMdpwUhAVTTGjVldOE4nxgssiSfconZBQA3ugKAyDVoh/7VCMMSyMXZ5k2sTssF+vFbWNyj1TZ5er963rfECjo2o6cH1/BsuByImdS/jlWzeqvBYaJcDokLzYJ2ZZ0yLVi6BfaF328+Cf244aMNDjIYNSU/s7NALFk93QcZgXyYHzkUkdq1d6mAZRJc6yzZEHApCLarpMUZ+RxybVfrw9ryaUMjT17VXo8jZBcii4VFDH0SyukXkps05JHRK2hF1/6ayUjMrElZ3mkrFHvrbwOX3kb3+JaG6pKiY6dJ9potanodGfP+yHdx5phgdlUZzijbfWUM7u1zf+a1xKWvwD80jqDdsGAAA=
`

type testCase struct {
	n int
	x int
	a []int
}

// solve mirrors 1644C.go.
func solve(tc testCase) string {
	n, x := tc.n, tc.x
	a := tc.a
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + a[i-1]
	}
	maxSum := make([]int, n+1)
	for length := 1; length <= n; length++ {
		best := -1 << 60
		for i := 0; i+length <= n; i++ {
			s := pref[i+length] - pref[i]
			if s > best {
				best = s
			}
		}
		maxSum[length] = best
	}
	var sb strings.Builder
	for k := 0; k <= n; k++ {
		ans := 0
		for length := 0; length <= n; length++ {
			add := k
			if length < k {
				add = length
			}
			val := maxSum[length] + add*x
			if val > ans {
				ans = val
			}
		}
		if k > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ans))
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
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, err
	}
	pos++
	cases := make([]testCase, 0, t)
	for i := 0; i < t && pos+1 < len(fields); i++ {
		n, _ := strconv.Atoi(fields[pos])
		x, _ := strconv.Atoi(fields[pos+1])
		pos += 2
		if pos+n > len(fields) {
			break
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			val, _ := strconv.Atoi(fields[pos+j])
			arr[j] = val
		}
		pos += n
		cases = append(cases, testCase{n: n, x: x, a: arr})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.x)
	for i, v := range tc.a {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
