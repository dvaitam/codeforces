package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (one per line: n a b).
const embeddedTestcases = `43 35 21
154 75 38
370 51 71
133 38 15
490 94 31
957 7 40
184 67 94
73 39 52
856 43 39
425 14 13
575 62 61
346 44 16
491 15 90
510 55 5
310 43 95
704 20 22
642 73 49
826 82 12
68 11 26
768 29 8
395 2 13
404 72 67
297 58 63
808 75 92
696 28 55
86 48 29
268 75 100
171 56 25
368 15 9
842 90 4
925 68 58
771 87 26
122 64 51
263 27 83
44 28 80
150 14 26
470 49 47
560 20 14
611 63 19
578 52 82
697 55 67
508 87 42
854 64 64
651 86 26
556 79 29
10 44 91
986 96 41
838 42 5
538 19 33
618 20 49
597 38 92
724 61 9
820 11 67
890 6 9
231 17 6
308 2 98
865 58 43
883 21 20
890 84 59
992 48 65
392 68 65
35 74 12
695 67 98
615 10 96
437 97 27
297 69 77
428 62 50
623 76 30
872 3 85
907 1 95
187 39 65
584 33 43
68 64 34
966 39 99
418 50 50
64 21 83
958 17 31
294 94 43
57 5 62
428 19 63
912 78 92
84 87 90
156 46 53
37 79 60
396 59 7
104 61 100
156 3 5
613 80 17
646 42 14
717 71 84
355 25 50
803 100 100
503 15 8
625 90 60
630 81 44
667 16 88
731 80 38
808 17 50
819 38 96
894 88 16`

// Embedded solver logic from 466B.go.
func solve(n, a, b int64) (int64, int64, int64) {
	m := 6 * n
	if a*b >= m {
		return a * b, a, b
	}
	swapped := false
	if a > b {
		a, b = b, a
		swapped = true
	}
	bestArea := int64(1<<62 - 1)
	var bestA, bestB int64
	lim := int64(math.Sqrt(float64(m))) + 2
	for A := a; A <= lim; A++ {
		B := (m + A - 1) / A
		if B < b {
			B = b
		}
		area := A * B
		if area < bestArea {
			bestArea = area
			bestA, bestB = A, B
		}
	}
	A := (m + b - 1) / b
	if A < a {
		A = a
	}
	area := A * b
	if area < bestArea {
		bestArea = area
		bestA, bestB = A, b
	}
	if swapped {
		bestA, bestB = bestB, bestA
	}
	return bestArea, bestA, bestB
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Fprintf(os.Stderr, "test %d: invalid line\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.ParseInt(parts[0], 10, 64)
		a, _ := strconv.ParseInt(parts[1], 10, 64)
		b, _ := strconv.ParseInt(parts[2], 10, 64)
		expS, expA, expB := solve(n, a, b)

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		outFields := strings.Fields(got)
		if len(outFields) < 3 {
			fmt.Printf("test %d: invalid output\n", idx+1)
			os.Exit(1)
		}
		gotS, _ := strconv.ParseInt(outFields[0], 10, 64)
		gotA, _ := strconv.ParseInt(outFields[1], 10, 64)
		gotB, _ := strconv.ParseInt(outFields[2], 10, 64)
		if gotS != expS || gotA != expA || gotB != expB {
			fmt.Printf("test %d failed\nexpected: %d %d %d\ngot: %s\n", idx+1, expS, expA, expB, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
