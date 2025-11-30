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
	"885441 396",
	"794773 913",
	"441002 43",
	"271494 990",
	"536111 499",
	"424605 942",
	"821873 851",
	"318047 993",
	"499749 368",
	"611721 915",
	"952226 225",
	"529203 144",
	"295529 145",
	"792519 99",
	"648407 820",
	"262675 933",
	"558434 724",
	"849575 618",
	"945990 152",
	"325214 103",
	"765285 77",
	"942501 872",
	"717210 340",
	"495078 575",
	"105593 364",
	"455263 325",
	"640562 657",
	"957362 211",
	"579364 490",
	"464198 888",
	"546679 268",
	"65305 826",
	"963081 563",
	"960490 16",
	"97803 738",
	"880900 410",
	"744755 846",
	"823183 686",
	"655639 3",
	"641621 507",
	"868288 890",
	"349318 251",
	"765753 335",
	"737823 893",
	"66044 197",
	"961565 583",
	"232474 246",
	"842369 992",
	"149417 824",
	"569367 460",
	"95647 84",
	"335602 898",
	"532615 957",
	"513055 113",
	"316090 566",
	"305231 725",
	"130874 562",
	"348915 836",
	"967049 555",
	"213073 988",
	"838261 619",
	"573813 603",
	"301631 457",
	"96084 612",
	"836696 396",
	"332448 591",
	"253868 299",
	"192801 195",
	"861371 193",
	"34575 629",
	"688558 268",
	"499679 72",
	"94188 697",
	"794406 135",
	"919361 155",
	"968236 41",
	"883384 84",
	"941803 718",
	"967923 851",
	"566861 701",
	"410304 859",
	"739544 539",
	"289024 536",
	"851055 243",
	"890751 222",
	"938517 697",
	"618452 847",
	"995901 431",
	"607855 283",
	"472450 506",
	"692318 658",
	"734240 940",
	"831862 367",
	"86375 334",
	"642550 120",
	"510074 603",
	"660758 345",
	"886129 196",
	"254842 18",
	"767023 279",
}

func referenceSolve(n, k int64) string {
	if n == 1 {
		return "0"
	}
	target := n - 1
	maxTotal := k * (k - 1) / 2
	if maxTotal < target {
		return "-1"
	}
	l, r := int64(1), k-1
	ans := k
	for l <= r {
		m := (l + r) / 2
		sum := m * (2*k - m - 1) / 2
		if sum >= target {
			ans = m
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return strconv.FormatInt(ans, 10)
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

func parseLine(line string) (int64, int64, error) {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected 2 numbers, got %d", len(parts))
	}
	n, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse n: %w", err)
	}
	k, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse k: %w", err)
	}
	return n, k, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
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
		n, k, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d %d\n", n, k)
		expected := referenceSolve(n, k)
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
