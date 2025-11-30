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
	"1 100",
	"2 99",
	"3 98",
	"4 97",
	"5 96",
	"6 95",
	"7 94",
	"8 93",
	"9 92",
	"10 91",
	"11 90",
	"12 89",
	"13 88",
	"14 87",
	"15 86",
	"16 85",
	"17 84",
	"18 83",
	"19 82",
	"20 81",
	"21 80",
	"22 79",
	"23 78",
	"24 77",
	"25 76",
	"26 75",
	"27 74",
	"28 73",
	"29 72",
	"30 71",
	"31 70",
	"32 69",
	"33 68",
	"34 67",
	"35 66",
	"36 65",
	"37 64",
	"38 63",
	"39 62",
	"40 61",
	"41 60",
	"42 59",
	"43 58",
	"44 57",
	"45 56",
	"46 55",
	"47 54",
	"48 53",
	"49 52",
	"50 51",
	"51 50",
	"52 49",
	"53 48",
	"54 47",
	"55 46",
	"56 45",
	"57 44",
	"58 43",
	"59 42",
	"60 41",
	"61 40",
	"62 39",
	"63 38",
	"64 37",
	"65 36",
	"66 35",
	"67 34",
	"68 33",
	"69 32",
	"70 31",
	"71 30",
	"72 29",
	"73 28",
	"74 27",
	"75 26",
	"76 25",
	"77 24",
	"78 23",
	"79 22",
	"80 21",
	"81 20",
	"82 19",
	"83 18",
	"84 17",
	"85 16",
	"86 15",
	"87 14",
	"88 13",
	"89 12",
	"90 11",
	"91 10",
	"92 9",
	"93 8",
	"94 7",
	"95 6",
	"96 5",
	"97 4",
	"98 3",
	"99 2",
	"100 1",
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

type testCase struct {
	input string
	want  string
}

func parseCases() []testCase {
	cases := make([]testCase, 0, len(rawTestcases))
	for _, line := range rawTestcases {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		a, _ := strconv.Atoi(fields[0])
		b, _ := strconv.Atoi(fields[1])
		want := fmt.Sprintf("%d", gcd(a, b))
		cases = append(cases, testCase{
			input: line + "\n",
			want:  want,
		})
	}
	return cases
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseCases()
	for idx, tc := range cases {
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.want {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx+1, tc.want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
