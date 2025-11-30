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
H4sIADsGK2kC/zWU2ZHDMAxD/1WFSohuqf/GAjzSs7OxLfEAQZCj3HpqK7O0n15uXWUXHdRefTHLKq+uOurU+yxD5zK85em16eLqd/B3y7XjT5by09WuQ9554Bv5H536/CnqI+qW5dNT19vGzmtvexbDWGRoZNK7Q+rboCLjtjMBfP1TaOd7BBp6E/hyiiFdnW38dqGMDoZOSH3CwODS6AltFM58ZNxBEDWqlkctS94rvU2CbuJVDEU9K1gwVwqrjx25nY4wQ2H6h4na6cL6iAsuTIvrfukdMI778AMgL10njc+VlkrR5arS4d/EPRE3SOtzU/OcoX8svsyrT1leB3Jt/AbQYGerXvmdQDAD2CZwUx+uno+C4vdGkixhUsRAcqjPLTZprnrpeZLECQNT4UCit4U2u8MtMt6g8yWmG1RGIUQ2HKtbH5sw6MMBEs0AjynXjYx7KO5HaGfeIc+D1mmdkfXU8sR1RKyNPICInCY6jAqFrX0NC20/mBogGhDv9p+yaYEnQMFls1LhOQzKRycm8+IKaFVI96HTHbwOVC1lYkuQkUN4CHji35Yb3T0SD6BFrotOoMOFNxqhOaIRjUkzoBvaXoCGKo5CLY4ppuROw1RibICYHW0Csr3U5Kofuy23xEZjj/HFoWHYmYHLzNMzU/xjXCdI1L0YxRA+GwHRR9sj2oM5J5wZzezt0BVbZdMxTHMafjHYg/oeJrFyXs6j25bDemNPmcAZY+Bu908yA5S3rFgwVNxRLZqbSefNw4nZjdazYmaJpTvYzItVCGBLgUmJ5jGeIXOFBHOnCzcWnNQLn52m76zeQjvf3nk52qcmb+hUXNzcoPb4A1yOGodYBgAA
`

type testCase struct {
	arr []int
}

// solve mirrors 1637C.go.
func solve(tc testCase) string {
	n := len(tc.arr)
	if n == 3 && tc.arr[1]%2 == 1 {
		return "-1"
	}
	allOne := true
	for i := 1; i < n-1; i++ {
		if tc.arr[i] != 1 {
			allOne = false
			break
		}
	}
	if allOne {
		return "-1"
	}
	var res int64
	for i := 1; i < n-1; i++ {
		res += int64(tc.arr[i]+1) / 2
	}
	return fmt.Sprint(res)
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
