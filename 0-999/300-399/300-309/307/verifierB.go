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
	"865 395",
	"777 912",
	"431 42",
	"266 989",
	"524 498",
	"415 941",
	"803 850",
	"311 992",
	"489 367",
	"598 914",
	"930 224",
	"517 143",
	"289 144",
	"774 98",
	"634 819",
	"257 932",
	"546 723",
	"830 617",
	"924 151",
	"318 102",
	"748 76",
	"921 871",
	"701 339",
	"484 574",
	"104 363",
	"445 324",
	"626 656",
	"935 210",
	"990 566",
	"489 454",
	"887 534",
	"267 64",
	"825 941",
	"562 938",
	"15 96",
	"737 861",
	"409 728",
	"845 804",
	"685 641",
	"2 627",
	"506 848",
	"889 342",
	"250 748",
	"334 721",
	"892 65",
	"196 940",
	"582 228",
	"245 823",
	"991 146",
	"823 557",
	"459 94",
	"83 328",
	"897 521",
	"956 502",
	"112 309",
	"565 299",
	"724 128",
	"561 341",
	"835 945",
	"554 209",
	"987 819",
	"618 561",
	"602 295",
	"456 94",
	"611 818",
	"395 325",
	"590 248",
	"298 189",
	"194 842",
	"192 34",
	"628 673",
	"267 488",
	"71 92",
	"696 776",
	"134 898",
	"154 946",
	"40 863",
	"83 920",
	"717 946",
	"850 554",
	"700 401",
	"858 723",
	"538 283",
	"535 832",
	"242 870",
	"221 917",
	"696 604",
	"846 973",
	"430 594",
	"282 462",
	"505 677",
	"657 718",
	"939 813",
	"366 85",
	"333 628",
	"119 499",
	"602 646",
	"344 866",
	"195 249",
	"17 750",
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveCase(a, b int64) string {
	return strconv.FormatInt(gcd(a, b), 10)
}

func parseCase(line string) (int64, int64, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("invalid line")
	}
	a, err1 := strconv.ParseInt(fields[0], 10, 64)
	b, err2 := strconv.ParseInt(fields[1], 10, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("invalid numbers")
	}
	return a, b, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range rawTestcases {
		a, b, err := parseCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := solveCase(a, b)
		input := fmt.Sprintf("%d %d\n", a, b)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
