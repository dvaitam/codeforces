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

// Embedded gzipped+base64 testcases from testcasesA.txt.
const encodedTestcases = `
H4sIAOsFK2kC/0VU17HkMAz7VxUqwaRy/40dAvfd7IzXlhgBkKfdfnt+PXu8vnu03aLH6pF99Pj6a/E1fODq9tNj4/FoC7/2Gl6Dh9l5Aju4xGnXLrCfMlgd/7Nlg+Nqo/EYaXGA1OE4jzZbMYLpTj/MzTomIk2eDGcOxtBBMJWC4ZlsBP/4TiZ5DoI3Hg1kQ9jXFosLNQS3wP1WkaNiMeOqW/YEP2R2rcc2pyAbjENvtkawZmOrC9gAx8H8LIlAJO6CwRcQmI3VIinaaMbtsauL89GJxKMPoUW5i31fg9MASghWRya+MLkGjZQhiBoXVaG3j83BUQzSb7C3QWSSufiiSn9lk9AwHZcdK8xRv/gpi7pRDSxbnaWqrBRwBGbodktB26b2xm85PlyzFMIqXe5jbY8+IR5kK9nBt+gTaGi4Cv4qvyDafxRdOL3SG3HHzWIfBEIyWiL1P2VLZaaOHoqTrYl0Prb+DDcJtw4ky01IpppLIiLcTzHhDi9Nnr5hSM0gJv5AMDyGTv2UytFydbtMpeJeptNAhbogBFlOcNmSKj1gvIS0gFgeEspA5YgNsZLWwXZrvB2lbOfyrH6i0UMNKLsQ9a1w1Rr41QVTTVUa/VU6QUmWdXLpVD46Zbemj1eM50mRZ2PGNHdHhE4WoakmpuEYAklcBu/zh0IaHRb+iAEFR/3HqYRDsmFzXiNSGcs1SZbJ5khpIXCUkJ1DGadkO4V0Sr/eQCptE+up1SgRY4+R/yg9DG9JJMx+f4CfGsutTGTfG9akow5uB8N0pLihbYc2FxdQrdvrPRXe6V4unp4tOLwfriFzGWJAe29pFf/tdrlFr0VmrlOUDH+zhnilLm9QtqR51aAODQypolqpM3IB4v4B9R0LYH0GAAA=
`

type testCase struct {
	arr []int
}

func solve(tc testCase) string {
	n := len(tc.arr)
	if n == 0 {
		return "NO"
	}
	a := tc.arr
	prefixMax := make([]int, n)
	suffixMin := make([]int, n)
	maxVal := a[0]
	for i := 0; i < n; i++ {
		if a[i] > maxVal {
			maxVal = a[i]
		}
		prefixMax[i] = maxVal
	}
	minVal := a[n-1]
	for i := n - 1; i >= 0; i-- {
		if a[i] < minVal {
			minVal = a[i]
		}
		suffixMin[i] = minVal
	}
	for i := 0; i < n-1; i++ {
		if prefixMax[i] > suffixMin[i+1] {
			return "YES"
		}
	}
	return "NO"
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	var cases []testCase
	for i := 0; i+1 < len(lines); i += 2 {
		nLine := strings.TrimSpace(lines[i])
		arrLine := strings.TrimSpace(lines[i+1])
		if nLine == "" || arrLine == "" {
			continue
		}
		n, err := strconv.Atoi(nLine)
		if err != nil || n <= 0 {
			continue
		}
		parts := strings.Fields(arrLine)
		if len(parts) < n {
			continue
		}
		arr := make([]int, n)
		ok := true
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(parts[j])
			if err != nil {
				ok = false
				break
			}
			arr[j] = val
		}
		if ok {
			cases = append(cases, testCase{arr: arr})
		}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
