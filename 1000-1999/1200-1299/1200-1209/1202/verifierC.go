package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
100
DASSWWWAS
S
A
DAD
ADASWWSSDSSWWASWAS
SSSSD
S
WSADDSSA
W
DDDDWDD
W
DW
WA
DDDDWW
WDSADWSS
SSWS
DS
WDD
DDSSDWA
AS
SD
AWS
WWWSSSWWS
W
SSDDDSDWSW
AW
WSWA
WWADDWDDW
DDWDDSSS
DSSD
WDSAWDDSA
SW
DDDWDW
S
WWAWA
WWDDA
SDWAWAW
DDWW
WDWW
SDDW
ADSS
WWDD
WSASW
WA
WSSW
SSSSW
DD
SW
DDWSS
SW
SD
S
DS
S
WDDD
DSWW
W
DDDS
DDSS
WW
SDW
SW
WDW
DD
SDWS
AD
DDW
WD
AW
WWD
AWD
D
DSDDW
S
W
D
DDWD
S
SWDD
D
WS
WWWWSWAS
DS
S
SS
DDWWDWS
DS
SDDD
DS
SW
WD
DW
SSDASASW
AW
SSWA
SSSSWD
SSAS
DS
S
DDDSS
SADW
DD
W
SA
S
`

// solve is copied from 1202C.go.
func solve(s string) int64 {
	n := len(s)
	x := make([]int, n+1)
	y := make([]int, n+1)
	for i, c := range s {
		x[i+1] = x[i]
		y[i+1] = y[i]
		switch c {
		case 'W':
			y[i+1]++
		case 'S':
			y[i+1]--
		case 'A':
			x[i+1]--
		case 'D':
			x[i+1]++
		}
	}
	prefMinX := make([]int, n+1)
	prefMaxX := make([]int, n+1)
	prefMinY := make([]int, n+1)
	prefMaxY := make([]int, n+1)
	prefMinX[0], prefMaxX[0] = x[0], x[0]
	prefMinY[0], prefMaxY[0] = y[0], y[0]
	for i := 1; i <= n; i++ {
		prefMinX[i] = min(prefMinX[i-1], x[i])
		prefMaxX[i] = max(prefMaxX[i-1], x[i])
		prefMinY[i] = min(prefMinY[i-1], y[i])
		prefMaxY[i] = max(prefMaxY[i-1], y[i])
	}
	sufMinX := make([]int, n+1)
	sufMaxX := make([]int, n+1)
	sufMinY := make([]int, n+1)
	sufMaxY := make([]int, n+1)
	sufMinX[n], sufMaxX[n] = x[n], x[n]
	sufMinY[n], sufMaxY[n] = y[n], y[n]
	for i := n - 1; i >= 0; i-- {
		sufMinX[i] = min(x[i], sufMinX[i+1])
		sufMaxX[i] = max(x[i], sufMaxX[i+1])
		sufMinY[i] = min(y[i], sufMinY[i+1])
		sufMaxY[i] = max(y[i], sufMaxY[i+1])
	}

	widthOrig := prefMaxX[n] - prefMinX[n] + 1
	heightOrig := prefMaxY[n] - prefMinY[n] + 1
	ans := int64(widthOrig) * int64(heightOrig)

	minWidth := widthOrig
	for _, dx := range []int{-1, 1} {
		for i := 0; i <= n; i++ {
			minVal := prefMinX[i]
			maxVal := prefMaxX[i]
			v := x[i] + dx
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
			if i != n {
				v = sufMinX[i+1] + dx
				if v < minVal {
					minVal = v
				}
				v = sufMaxX[i+1] + dx
				if v > maxVal {
					maxVal = v
				}
			}
			width := maxVal - minVal + 1
			if width < minWidth {
				minWidth = width
			}
		}
	}
	area := int64(minWidth) * int64(heightOrig)
	if area < ans {
		ans = area
	}

	minHeight := heightOrig
	for _, dy := range []int{-1, 1} {
		for i := 0; i <= n; i++ {
			minVal := prefMinY[i]
			maxVal := prefMaxY[i]
			v := y[i] + dy
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
			if i != n {
				v = sufMinY[i+1] + dy
				if v < minVal {
					minVal = v
				}
				v = sufMaxY[i+1] + dy
				if v > maxVal {
					maxVal = v
				}
			}
			height := maxVal - minVal + 1
			if height < minHeight {
				minHeight = height
			}
		}
	}
	area2 := int64(minHeight) * int64(widthOrig)
	if area2 < ans {
		ans = area2
	}
	return ans
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseTestcases() ([]string, error) {
	lines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, err
	}
	if len(lines)-1 < t {
		return nil, fmt.Errorf("not enough lines for %d tests", t)
	}
	cases := make([]string, 0, t)
	for i := 1; i <= t; i++ {
		cases = append(cases, strings.TrimSpace(lines[i]))
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		return
	}

	for idx, s := range testcases {
		input := fmt.Sprintf("1\n%s\n", s)
		expect := strconv.FormatInt(solve(s), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
