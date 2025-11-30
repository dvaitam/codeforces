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

// Embedded gzipped+base64 testcases from testcasesE.txt.
const encodedTestcases = `
H4sIAKgIK2kC/1VU3XVGIQh7ZwpGUPx/dwL3H6YIwfu1Paf1qmASAjklGrXy3vts/0NZMtvKd87Z/kttTb6fcUZlJT7+RU3KXZ8v4NBYjbdvbg8kWcIHr1nSMbvH7Xjn0BJhS2sI/L1z86XCttL17JP9EAg0rqfsyQOkJ9D9GtgAz25XR3zOR/5mzqKsgqUuqDQLV47F+Oj/yog4ejr9ukaWwvGqIpRq7GnOwc7QhTBIrXHc3DRT5wMsfj6k/erk+WV0/oeWWlu8SWqCFgGaJC+Geq4Y9fGzY7hyg5oec4CA5mj8thFur9VxnaFwh2sPKlKhI60MXc6Dcm/Uavd/X7ITNU3nKIhy7m4FfHs1VlLYWqvFyHVVpQ7suAltPg/PVjkc9cjRrAwuVteVhN3GjYPmOV8Bc26vzIalrs57PxZuosSvbo8gDtAZLaP8EecP0JqZd2heot5erTIbg8zVz/03lvDXFGFkrX8yNfaTA73V1Xh+jXJy7gFEPZy9J7QDMyjA6V/TX/CJv3JcMCKNn602pBKRcMQTUa+WweiMuvi8XoVOOgy686G5Fp+wztfSVATQYkDZD8m0SXSuMttwD/7xDdUOQTCIBAo64V69Y/FZbaztsMV2JbylzkeI2ij8edtjsYfK60N9RYWe+a6zk4A+1eazLVIgnQ/GX3eWWWICnOhC1RGj+Xa4caJROoemd7/3xJ/NXKUpkOdr6+uJufiNCBypGoJ0VC9qfWAm3m+Mu+NKKvxmszp5yE+mkFLbGybWwbA4EJq+fwE/mwd5BgAA
`

type testCase struct {
	n int64
	s string
}

func solveCase(n int64, s string) int64 {
	cntD, cntR := 0, 0
	kD, kR := -1, -1
	for i, r := range s {
		if r == 'D' {
			cntD++
			if kD == -1 {
				kD = i
			}
		} else {
			cntR++
			if kR == -1 {
				kR = i
			}
		}
	}
	if cntD == 0 || cntR == 0 {
		return n
	}
	Sx := n - int64(cntD) - 1
	Sy := n - int64(cntR) - 1
	var ans int64
	if kD < kR {
		k := int64(kR)
		Sv := Sx
		Sh := Sy
		ans = k + (Sv+1)*(Sh+2)
		remD := int64(cntD) - k
		remR := int64(cntR) - 1
		ans += remD * (Sh + 1)
		ans += remR * (Sv + 1)
	} else {
		k := int64(kD)
		Sv := Sy
		Sh := Sx
		ans = k + (Sv+1)*(Sh+2)
		remR := int64(cntR) - k
		remD := int64(cntD) - 1
		ans += remR * (Sh + 1)
		ans += remD * (Sv + 1)
	}
	return ans
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
		n, _ := strconv.ParseInt(fields[pos], 10, 64)
		s := fields[pos+1]
		pos += 2
		cases = append(cases, testCase{n: n, s: s})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("1\n%d %s\n", tc.n, tc.s)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := fmt.Sprint(solveCase(tc.n, tc.s))
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
