package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test data from testcasesA.txt.
const testData = `7 94 7 -90 -34 30 24 3
5 22 -9 49 -45 29
3 -28 -65 93
2 58 -36
9 80 54 -63 -21 -75 86 -82 75 -16
8 43 -75 -10 11 -20 56 63 -48
9 22 13 33 -34 -85 40 -97 -77 84
7 81 100 71 60 -100 56 26
6 -38 86 -17 80 -84 -52
10 -44 -39 -64 39 14 -77 -80 -19 30 25
2 -23 41
5 80 -69 40 -15 38
4 54 40 50 -27
8 -77 52 -2 -19 47 -39 -26 -53
4 -53 -92 56 68
5 21 -83 -78 73 93
3 -62 -91 -80
9 74 0 80 34 -30 33 -40 -45 73
10 7 48 -30 15 26 69 64 79 -9 -79
6 56 -71 24 50 61 -15
4 -38 -96 87 -31
2 80 -44
6 -57 -15 9 -85 -75 100
3 78 -44 -89
10 62 36 54 74 -82 -94 -69 62 -52 55
10 -70 0 -77 -6 -71 -91 55 -95 -51 -53
2 22 -47
1 73
1 39
7 58 -75 -34 -83 -44 -82 65
5 -11 11 -54 -85 28
8 -90 52 -75 79 0 -49 -34 -9
8 45 -57 78 72 -48 96 -86 73
3 -59 -13 35
5 -70 52 13 70 -56
1 20
7 45 30 -21 66 -9 -1 68
5 -61 43 76 -97 17
2 -15 89
1 39
5 -66 -39 95 23 -10
10 -27 72 -9 51 62 58 -67 83 -21 -1
7 66 -80 -100 52 -51 78 -15
3 -39 -43 63
8 -4 81 72 45 6 -92 2 79
10 7 97 69 81 -89 -58 14 -84 -34 79
3 14 35 24
9 54 93 -100 -91 26 -17 -21 19 -88
7 -52 40 62 -79 85 -67 -97
7 73 6 -20 -100 -46 -97 83
1 72
9 56 -75 -52 -70 55 66 -50 -23 -29
3 -75 21 1
2 -95 -30
8 -71 -35 -66 67 33 66 65 -12
2 -61 -29
1 -90
1 -48
5 42 -20 -7 45 -90
10 67 26 82 64 17 63 11 -5 37 -55
4 -4 50 -26 -98
3 -62 -31 -15
6 -6 83 -77 -14 99 58
1 -90
5 -59 -62 49 -26 -8
7 40 -67 -25 -71 22 87 -39
1 -22
3 33 86 -82
5 3 -16 -24 6 -73
2 43 23
8 -14 -13 -69 22 -71 79 27 9
1 -23
6 88 75 -61 -58 60 44
7 63 -78 -84 -79 -50 91 -44
1 -2
1 -75
7 42 32 -26 14 25 49 82
4 8 -79 -6 -44
5 49 98 -58 10 -51
6 -71 -84 79 -93 34 15
4 -70 27 1 -35
4 64 -90 -45 59
3 -74 -50 17
7 -8 39 -62 -74 52 24 -63
10 3 63 74 8 33 26 73 -18 27 27
4 38 56 -44 -98
6 80 91 -19 -18 -91 34
3 -35 54 100
3 -3 49 -25
8 -84 -79 32 -90 -84 -43 -67 -90
5 -97 94 14 -16 -59
3 67 17 -5
9 -3 35 28 -92 46 -77 73 32 94
10 -81 91 9 93 -48 -26 37 53 6 23
7 55 50 -41 -95 68 -100 89
3 -23 29 45
5 -15 -84 26 -33 -23
7 -2 -2 -85 -59 64 -68 -39
5 86 -15 -86 -91 23
7 -64 25 54 83 -80 72 78`

// Embedded solver logic from 1376A.go. It intentionally panics.
func solverBehavior() (string, bool) {
	ok := true
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	// Directly embed the behavior of 1376A.go: panic at runtime.
	panic("intentional runtime error")
	return "", ok
}

func parseTests() ([][]int, error) {
	lines := strings.Split(testData, "\n")
	tests := make([][]int, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n in line %q", line)
		}
		if len(parts) != n+1 {
			return nil, fmt.Errorf("line %q has wrong count", line)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[i+1])
			if err != nil {
				return nil, fmt.Errorf("invalid value in line %q", line)
			}
			arr[i] = v
		}
		tests = append(tests, arr)
	}
	return tests, nil
}

func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	expectedOutput, expectedOK := solverBehavior()

	for idx, arr := range tests {
		var input strings.Builder
		input.WriteString("1\n")
		fmt.Fprintf(&input, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')

		out, err := runProg(bin, input.String())
		if expectedOK {
			if err != nil {
				fmt.Printf("case %d failed: %v\n", idx+1, err)
				os.Exit(1)
			}
			if out != expectedOutput {
				fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expectedOutput, out)
				os.Exit(1)
			}
		} else {
			if err == nil {
				fmt.Printf("case %d failed: expected runtime error but got output %s\n", idx+1, out)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
