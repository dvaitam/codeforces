package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var testcases = []string{
	"7917",
	"4155",
	"7209",
	"8292",
	"4445",
	"1331",
	"3121",
	"8909",
	"5188",
	"4980",
	"4317",
	"8522",
	"7420",
	"7798",
	"3484",
	"8928",
	"4904",
	"3933",
	"5779",
	"8304",
	"8439",
	"2789",
	"5134",
	"2140",
	"3308",
	"2144",
	"7191",
	"1776",
	"6065",
	"7548",
	"3052",
	"8452",
	"5362",
	"6776",
	"7637",
	"5930",
	"8390",
	"2203",
	"3540",
	"1809",
	"6978",
	"1604",
	"8363",
	"7967",
	"6603",
	"3704",
	"4867",
	"5585",
	"1824",
	"3898",
	"4556",
	"3590",
	"6004",
	"6246",
	"8479",
	"2675",
	"8918",
	"5526",
	"4907",
	"4626",
	"8088",
	"5270",
	"3133",
	"1510",
	"7594",
	"8524",
	"5494",
	"8503",
	"1115",
	"1764",
	"6895",
	"7882",
	"4267",
	"6818",
	"7757",
	"7431",
	"6473",
	"6122",
	"1009",
	"6012",
	"5043",
	"7783",
	"8107",
	"3729",
	"2998",
	"6982",
	"3664",
	"6764",
	"8130",
	"1515",
	"2565",
	"8512",
	"5649",
	"2816",
	"2954",
	"7581",
	"8926",
	"2167",
	"7579",
	"5448",
}

func distinct(n int) bool {
	var seen [10]bool
	for n > 0 {
		d := n % 10
		if seen[d] {
			return false
		}
		seen[d] = true
		n /= 10
	}
	return true
}

func referenceSolve(y int) string {
	for i := y + 1; ; i++ {
		if distinct(i) {
			return strconv.Itoa(i)
		}
	}
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	idx := 0
	for _, tc := range testcases {
		line := strings.TrimSpace(tc)
		if line == "" {
			continue
		}
		idx++
		y, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse y: %v\n", idx, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n", y)
		expected := referenceSolve(y)
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\ngot: %s\n", idx, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
