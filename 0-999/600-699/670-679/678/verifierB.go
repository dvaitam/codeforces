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
	"18611",
	"75606",
	"9271",
	"34432",
	"16455",
	"65937",
	"59915",
	"62898",
	"86405",
	"50756",
	"28519",
	"13302",
	"64944",
	"4715",
	"52093",
	"57723",
	"80618",
	"1276",
	"92204",
	"59377",
	"35908",
	"95573",
	"30984",
	"78483",
	"14399",
	"42606",
	"5009",
	"3925",
	"4335",
	"86137",
	"71964",
	"2206",
	"50965",
	"90978",
	"29390",
	"56327",
	"96138",
	"4806",
	"70157",
	"30057",
	"58394",
	"65987",
	"73464",
	"31550",
	"46311",
	"31260",
	"89715",
	"29676",
	"61241",
	"38982",
	"3816",
	"55549",
	"73935",
	"85186",
	"14107",
	"25367",
	"83490",
	"95848",
	"39848",
	"16845",
	"98405",
	"44607",
	"95566",
	"94217",
	"66640",
	"56326",
	"67547",
	"88858",
	"25883",
	"40763",
	"38245",
	"78015",
	"66452",
	"67228",
	"52557",
	"78201",
	"5525",
	"63944",
	"32816",
	"98482",
	"53990",
	"55304",
	"88129",
	"23676",
	"49119",
	"72932",
	"93148",
	"89406",
	"97759",
	"50113",
	"12333",
	"58535",
	"88000",
	"67640",
	"15146",
	"22456",
	"69280",
	"52544",
	"49565",
	"65185",
}

func isLeap(y int) bool {
	if y%400 == 0 {
		return true
	}
	if y%100 == 0 {
		return false
	}
	return y%4 == 0
}

func solve(y int) int {
	origLeap := isLeap(y)
	shift := 0
	for {
		if isLeap(y) {
			shift = (shift + 2) % 7
		} else {
			shift = (shift + 1) % 7
		}
		y++
		if shift == 0 && isLeap(y) == origLeap {
			return y
		}
	}
}

type testCase struct {
	input string
	want  int
}

func parseCases() []testCase {
	var cases []testCase
	for _, line := range rawTestcases {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, " ") {
			// first line is the count
			continue
		}
		y, _ := strconv.Atoi(line)
		cases = append(cases, testCase{
			input: line + "\n",
			want:  solve(y),
		})
	}
	return cases
}

func run(bin string, input string) (string, error) {
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
		val, err := strconv.Atoi(strings.TrimSpace(got))
		if err != nil {
			fmt.Printf("case %d: bad output %q\n", idx+1, got)
			os.Exit(1)
		}
		if val != tc.want {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, tc.want, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
