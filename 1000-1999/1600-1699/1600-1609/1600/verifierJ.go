package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"bufio"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	n    int
	m    int
	grid [][]int
}

// Embedded testcases from testcasesJ.txt (gzip+base64).
const encodedTestcases = "H4sIALv/KmkC/0WXWXIDIRBD/+cUHIF9uf/FQusJp+yK7QF6ldSkpnpfpaSRztdTS/Gt3A89LD3t+yT7e8nfSLzWfeXYGn86BmJH09biA0N2+t1bRli6v4tM1kEd19hbbmi96FVaZYjTG3sop9TBC/b6lNhczpPT5flclp6hoHX2weFIhqOL7oQl0C60MEkttAghhCEFoA+a7PeTxmCivaPstOqKfZ/VqemA+eJeUPlTPlFnm/pit/MHNDjCHk3brfAejcA9hrINg9AbBsgnAqeqgmh+oPrQESGcfyQV24V+Vn1vuUZfqFLPFcYWwGy1lrUbFM0tJVDXiN3OH5vEPpExMIBz3WOF1NlNeOeHvQltzFjwRYVYLvnzMLlP3ZnNSao74D0E8Xkvplk0HqfEM9V7/efXiD7zd8cKMz2zzPMy9QuQDYSqizqVaRfK/q9N1B/PEhTe/L2tLh9IU+Al4ER6qL9U9qZynP5p72a3PZgx2m7JdFqBwVBjKWEAPmG6VaV0Jycj/NpVBeqHBENcKsSLdd6SzFW08c5keHLm4KdkItdafQ2cX9dDZK7Fyds2YVQsPVLiJ0b6vIlYRNGv5ROTY6TDt3PMD4TmGx8wVdyq1SmD82ygQelsUNUjYrtE+nmGKdpMvH5zPnPY5MkSbh2wpW57LzkBvRMx4pVTgs+DhXpVGomqf+Qucd0uMHTJ4/xyTBZT8Z9/HlPA4zHqM0iOs0zjMnUTs3HXXzt+o/n81p40fG4OjmKQ7eiYzJWIfTbYz+hkXL4uXCO1XssD3u1Vxr+jh2RygLnJGF4hy3iequjXWSctc/dqCqqYykK3smahAbzyT05PHJCByPqoI47d3mX0Nk9Yuqbmj5sCeWqmOcMr5zuRlyxjg1kLx2fVz98Yk98F2eXC3Ajno48z7/lCcbUIp7lHZTaxCErV+zmI4eKccP8C36I2Anmj+OyKOSlXgaj8k1gJWlEtu+4LJYZM/Pqeyz+XPdZ7fMnePyu3wBGibFS+CmWYQrAMz3FOuykAwKp6CBiw+BAsEuV4KCpwKe4wVYBLlcO8joGMBqJIxAq5TBV64wxK04ErlcN7HkMcoBbpXB9xfzOJwBrlcNnDt+zGMATxjDXVluNQmzvlc4hTrnC6zPG4LcMsY8xl8thvksLkC6BEEHXCLbIN8lItZ+0wKpr9Yx34LoDeX6CldCGtwioYUqrRXh8eNEKHiNh2co5La7HiyBzj6HnRhzEJNs9jOMzxyOuVoHDEDag+kxrhlthVXO30vGt88KRDDG5EetpSJR0QzL0A9GzTO4sRK7EOXSinh7ZaY45RSRM9HvBhZmb0NRMaSqhXVqoXAd3PZ8/Z++H0BxUQ8/ZD+1P/+pOimYdlJCc/pywllTpLVZdWmoU3cvQxxZmWhWwou70nKWe2Ix0Sy/ZlwkIjGyEpPyeneMG3cJzo0uoEljd++XuydH1qSzFrxydP9xD77HNpM9K9b9JmScF8221AZ9wC6xLLpaUywkPn6DArCJ7PxRttH8fkMlAeXqr19bdd+B7fN9Zm0Pqud7vlk+UL7CMjbaWFB4OOpTeXvrtPPjaWhUJoVCY7mLdhy8GDQ7RWWGx+UD3j7bGmeyUP3jgBMk/pIsPmDGmrQgMx4QwZrphxV7ZXHS+ZG8iQ7ESbELFojdTRvMoa3JnJTmMTkZzQLIMmprMFHG1Ntrb+83dl74reFVC0qXyCYONMAgZ8i1m3rYrv1YZ6Hx5RVJpsN3faBVX2sPklRm/RSb78RkJNAibD8bnuC170La5PDsZZx/AICHRXz9mKAAA="

func solve(tc testCase) string {
	n, m := tc.n, tc.m
	grid := tc.grid
	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, m)
	}
	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	var sizes []int
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if visited[i][j] {
				continue
			}
			q := [][2]int{{i, j}}
			visited[i][j] = true
			size := 0
			for len(q) > 0 {
				cell := q[0]
				q = q[1:]
				r, c := cell[0], cell[1]
				size++
				val := grid[r][c]
				for d, dir := range dirs {
					if val&(1<<(3-d)) == 0 {
						nr := r + dir[0]
						nc := c + dir[1]
						if nr >= 0 && nr < n && nc >= 0 && nc < m && !visited[nr][nc] {
							visited[nr][nc] = true
							q = append(q, [2]int{nr, nc})
						}
					}
				}
			}
			sizes = append(sizes, size)
		}
	}
	sort.Slice(sizes, func(i, j int) bool { return sizes[i] > sizes[j] })
	var sb strings.Builder
	for i, v := range sizes {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func decodeTestcases() (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedTestcases)
	if err != nil {
		return "", err
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer r.Close()
	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		return "", err
	}
	return out.String(), nil
}

func parseTestcases() ([]testCase, error) {
	raw, err := decodeTestcases()
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(raw)
	in := bufio.NewReader(reader)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			return nil, fmt.Errorf("case %d header read error: %v", caseIdx+1, err)
		}
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				if _, err := fmt.Fscan(in, &grid[i][j]); err != nil {
					return nil, fmt.Errorf("case %d grid read error: %v", caseIdx+1, err)
				}
			}
		}
		cases = append(cases, testCase{n: n, m: m, grid: grid})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.grid[i][j]))
		}
		sb.WriteByte('\n')
	}
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
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
