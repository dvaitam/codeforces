package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"100",
	"5 4 38 46 8 25 8",
	"2 30 11 6",
	"1 27 40",
	"3 24 39 16 23",
	"7 17 18 18 31 49 4 7 49",
	"5 25 5 20 20 42 7",
	"7 42 41 6 21 32 17 30 40",
	"4 5 49 2 14 12",
	"1 20 9",
	"1 21 20",
	"8 19 38 37 8 5 21 44 18 49",
	"4 17 16 47 5 39",
	"1 45 34",
	"6 45 44 26 4 33 38 36",
	"6 45 46 16 44 35 17 36",
	"5 50 18 5 40 29 16",
	"6 12 38 25 10 8 27 10",
	"3 3 42 20 19",
	"5 46 31 49 9 44 20",
	"2 15 30 38",
	"6 44 2 1 50 31 11 14",
	"1 30 15",
	"3 16 49 13 8",
	"3 17 24 21 11",
	"4 14 8 8 14 1",
	"8 20 39 25 34 38 47 27 38 21",
	"3 50 15 40 14",
	"6 47 22 8 14 18 49 38",
	"7 3 6 9 14 4 34 24 38",
	"5 22 45 36 42 32 6",
	"8 6 25 37 28 42 17 8 42 39",
	"2 16 23 36",
	"6 3 6 7 49 23 35 18",
	"4 29 4 18 1 40",
	"8 17 22 11 11 16 42 14 40 17",
	"3 43 35 30 27",
	"2 28 45 12",
	"3 46 6 22 32",
	"2 18 12 22",
	"1 33 27",
	"3 27 37 47 33",
	"1 50 16",
	"6 40 37 29 15 13 50 30",
	"7 31 12 11 21 24 30 40 22",
	"3 42 8 9 33",
	"7 50 2 8 47 48 42 30 49",
	"7 11 9 29 4 3 5 34 32",
	"2 32 34 35",
	"3 24 21 12 50",
	"8 20 22 15 42 32 45 20 34 34",
	"2 20 14 14",
	"6 37 21 4 7 24 32 10",
	"1 4 30",
	"3 48 14 34 30",
	"4 33 14 35 41 7",
	"4 36 34 33 33 43",
	"7 5 43 8 7 20 25 3 14",
	"7 48 34 4 36 46 36 15 4",
	"2 37 39 14",
	"7 44 34 30 22 45 20 31 17",
	"7 35 14 32 28 15 16 32 18",
	"1 39 40",
	"1 46 31",
	"7 43 48 34 21 47 23 41 23",
	"5 30 28 22 23 3 50",
	"4 23 7 18 31 38",
	"3 18 18 22 32",
	"1 26 45",
	"7 39 32 5 44 10 32 15 17",
	"5 17 24 21 3 38 45",
	"7 38 48 5 28 41 41 42 9",
	"6 41 45 42 8 16 48 27",
	"3 4 5 30 41",
	"1 7 29",
	"6 32 1 36 44 27 15 28",
	"3 15 44 36 21",
	"8 20 25 21 27 18 9 45 26 8",
	"4 16 21 12 17 11",
	"5 28 28 2 26 44 32",
	"4 47 29 41 5 24",
	"4 1 18 41 46 47",
	"7 21 43 14 35 46 46 1 42",
	"6 50 4 35 49 36 5 50",
	"6 11 19 18 40 47 3 23",
	"5 23 31 36 15 14 29",
	"2 43 23 7",
	"2 49 2 8",
	"2 26 42 15",
	"5 21 16 42 37 50 39",
	"6 20 29 7 24 28 29 4",
	"8 4 19 22 7 40 41 17 38 18",
	"3 40 31 31 34",
	"1 49 35",
	"3 12 29 27 3",
	"8 26 39 3 19 12 17 9 4 38",
	"5 48 6 23 2 15 4",
	"1 49 43",
	"1 37 43",
	"3 40 18 22 31",
	"3 2 18 44 44",
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

type testCase struct {
	input string
	want  string
}

func parseCases() []testCase {
	var cases []testCase
	for idx, line := range rawTestcases {
		if idx == 0 {
			// count
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.ParseInt(fields[1], 10, 64)
		nums := make([]int64, n)
		for i := 0; i < n; i++ {
			val, _ := strconv.ParseInt(fields[2+i], 10, 64)
			nums[i] = val
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		g := int64(1)
		for _, c := range nums {
			g = lcm(g, gcd(c, k))
		}
		want := "Yes"
		if g != k {
			want = "No"
		}
		cases = append(cases, testCase{input: sb.String(), want: want})
	}
	return cases
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: verifierB <solution-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseCases()
	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.want {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
