package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesD = `100
14 11 10
1 10 7
8 5 6 5 3 3
15 -5 10
4 8 6
10 5 4 1 2 8
14 5 12
9 3 4
3 1 8 7 8 8
13 -3 10
6 1 4
3 2 9 9 7 10
8 -5 12
5 4 2
8 6 1 4 10 5
-3 -4 -3
7 8 9
6 7 7 7 4 5
-5 15 15
8 7 9
5 8 3 7 7 7
12 13 15
10 10 5
3 6 1 7 8 8
7 -5 14
5 2 4
6 3 3 8 2 2
15 -5 15
3 10 1
9 7 7 3 5 7
-4 10 -2
8 5 9
2 8 5 5 5 7
5 -3 -4
1 10 5
5 4 8 5 9 4
-4 -4 -3
7 2 7
6 1 5 5 1 1
-1 14 10
1 10 1
2 10 3 1 9 3
14 -2 -2
2 2 2
8 7 2 2 8 8
-2 -1 11
9 3 9
3 2 9 1 9 4
11 12 13
2 7 6
9 7 5 4 10 2
6 11 -5
3 10 7
1 1 3 5 4 3
10 -1 -1
3 2 8
1 3 6 6 4 2
15 15 -2
9 8 3
10 4 6 7 2 5
-4 15 7
7 9 6
1 7 9 5 3 5
5 15 8
1 5 2
1 5 4 3 5 6
10 13 14
9 2 3
10 2 7 5 1 5
5 7 13
1 1 10
9 4 1 9 4 8
14 -2 -5
9 8 10
7 6 3 8 4 7
14 10 8
10 1 4
7 3 3 7 7 8
13 -5 4
4 5 2
3 2 3 5 3 7
12 10 9
5 9 7
2 5 10 6 5 5
11 -5 15
5 7 10
9 9 4 7 5 7
12 -3 -5
1 1 7
1 8 7 9 6 9
12 10 13
8 9 8
3 10 1 10 6 4
-2 14 9
5 7 2
3 8 6 8 3 8
11 -4 11
5 8 6
5 7 9 10 2 5
-1 -3 5
3 9 3
4 1 1 8 8 2
-1 15 9
10 9 1
10 3 10 5 2 5
-4 5 -2
6 3 8
2 10 2 5 5 10
10 7 7
6 3 1
4 2 7 5 8 7
15 -3 -1
8 7 3
8 6 6 5 8 6
-1 5 10
4 3 8
3 2 9 10 10 8
-2 -2 7
6 3 2
6 6 5 6 1 4
12 10 -3
5 3 3
9 5 2 1 7 2
-1 -1 -4
7 8 8
2 3 2 7 10 2
13 -2 11
5 2 10
8 5 6 3 7 5
15 11 6
3 3 2
8 2 6 9 7 5
10 2 -3
8 1 3
3 8 7 2 6 4
14 12 5
2 6 3
8 6 10 8 9 5
14 3 11
5 1 10
9 2 5 8 3 2
11 12 14
1 8 9
10 9 7 7 6 10
10 -1 11
6 4 6
1 8 3 1 2 9
-4 5 4
1 3 2
4 2 1 4 3 6
8 -4 14
7 9 10
4 10 1 8 1 7
13 15 -1
3 4 2
6 5 5 8 6 9
11 9 12
7 8 6
6 7 6 7 5 1
7 -1 10
4 8 5
4 3 10 4 2 5
13 10 -2
9 9 3
8 4 7 3 2 1
-2 13 -1
9 1 2
3 10 4 6 6 2
-1 13 11
2 7 8
4 10 4 10 10 8
-1 8 -2
10 5 2
4 3 3 5 5 6
12 -2 -3
3 8 6
7 1 1 4 6 7
13 -2 -3
4 6 7
8 10 8 5 4 7
-3 12 9
1 6 6
1 10 3 3 9 10
7 9 4
6 8 2
1 7 1 8 2 3
14 -1 10
5 6 5
8 1 10 6 3 7
13 -3 8
1 5 3
9 8 9 3 7 6
13 -3 -3
3 9 6
4 2 4 8 4 1
13 9 11
7 6 5
1 5 8 7 4 8
7 13 -2
7 13 -2
1 9 10
6 7 3 5 5 10
-3 8 -3
3 7 1
2 10 4 9 2 4
14 -2 14
6 8 9
10 10 1 8 5 4
-3 -5 15
9 1 8
1 10 4 3 8 4
13 -3 -1
9 5 10
6 9 1 10 6 5
-5 13 -3
9 4 10
6 5 6 2 5 10
-5 15 14
3 6 1
7 1 7 4 7 10
13 9 12
8 5 1
2 1 1 2 3 7
15 12 -2
7 10 3
1 1 7 7 5 3
15 8 11
7 1 6
4 9 6 10 4 1
15 -4 13
7 8 9
4 4 1 4 2 5
-1 13 8
6 6 5
5 1 7 2 10 8
8 -1 12
3 7 4
8 10 6 1 10 9
7 10 13
4 4 9
2 9 2 3 3 5
-1 14 15
1 6 2
9 10 3 9 7 5
-4 11 9
1 6 3
10 5 7 1 4 10
-2 10 -4
5 4 6
5 2 6 5 4 6
7 -2 12
4 4 4
5 3 6 7 1 9
-2 11 -3
1 9 7
2 9 7 9 6 3
4 10 12
3 2 5
10 4 2 8 4 6
-2 -1 9
7 5 6
6 9 3 2 8 7
10 -1 -4
8 6 10
7 2 8 10 3 2
4 15 -2
2 7 4
7 10 5 9 1 10
11 15 15
1 7 7
2 2 9 10 3 5
8 -1 14
6 8 9
10 3 6 3 1 6
5 14 -2
3 2 9
9 2 8 9 10 8
-3 -4 15
5 6 10
7 5 9 6 4 10
9 -1 -4
6 5 2
2 8 5 6 5 8
15 7 -5
1 6 6
7 10 9 5 7 10
5 -5 8
1 9 2
3 2 3 1 9 5
13 -3 8
4 6 4
6 5 1 2 8 5
9 -5 -5
7 4 3
3 10 1 2 6 1
13 6 15
4 4 8
9 6 8 10 9 1
-2 -1 -2
4 3 4
1 10 9 2 10 10
`

type testCase struct {
	x, y, z    int
	x1, y1, z1 int
	a          [6]int
	expected   int
}

func calcSum(x, y, z, x1, y1, z1 int, a [6]int) int {
	sum := 0
	if y < 0 {
		sum += a[0]
	}
	if y > y1 {
		sum += a[1]
	}
	if z < 0 {
		sum += a[2]
	}
	if z > z1 {
		sum += a[3]
	}
	if x < 0 {
		sum += a[4]
	}
	if x > x1 {
		sum += a[5]
	}
	return sum
}

func parseCases(data []byte) ([]testCase, error) {
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, err
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		readInt := func() (int, error) {
			if !scan.Scan() {
				return 0, fmt.Errorf("unexpected EOF")
			}
			return strconv.Atoi(scan.Text())
		}
		x, err := readInt()
		if err != nil {
			return nil, err
		}
		y, err := readInt()
		if err != nil {
			return nil, err
		}
		z, err := readInt()
		if err != nil {
			return nil, err
		}
		x1, err := readInt()
		if err != nil {
			return nil, err
		}
		y1, err := readInt()
		if err != nil {
			return nil, err
		}
		z1, err := readInt()
		if err != nil {
			return nil, err
		}
		var a [6]int
		for j := 0; j < 6; j++ {
			v, err := readInt()
			if err != nil {
				return nil, err
			}
			a[j] = v
		}
		cases = append(cases, testCase{
			x: x, y: y, z: z,
			x1: x1, y1: y1, z1: z1,
			a:        a,
			expected: calcSum(x, y, z, x1, y1, z1, a),
		})
	}
	return cases, nil
}

func runCase(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.x, tc.y, tc.z)
	fmt.Fprintf(&sb, "%d %d %d\n", tc.x1, tc.y1, tc.z1)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases([]byte(testcasesD))
	if err != nil {
		fmt.Println("could not parse embedded testcases:", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		expected := strconv.Itoa(tc.expected)
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
