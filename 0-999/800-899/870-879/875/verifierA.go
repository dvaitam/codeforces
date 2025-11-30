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
	"906691060",
	"413654000",
	"813847340",
	"955892129",
	"451585302",
	"43469774",
	"278009743",
	"548977049",
	"521760890",
	"434794719",
	"985946605",
	"841597327",
	"891047769",
	"325679555",
	"511742082",
	"384452588",
	"626401696",
	"957413343",
	"975078789",
	"234551095",
	"541903390",
	"149544007",
	"302621085",
	"150050892",
	"811538591",
	"101823754",
	"663968656",
	"858351977",
	"268979134",
	"976832603",
	"571835845",
	"757172937",
	"869964136",
	"646287426",
	"968693315",
	"157798603",
	"333018423",
	"106046332",
	"783650879",
	"79180333",
	"965120264",
	"913189318",
	"734422155",
	"354546568",
	"506959382",
	"601095368",
	"108127102",
	"379880546",
	"466188457",
	"339513622",
	"655934895",
	"687649392",
	"980338160",
	"219556307",
	"593267778",
	"512185346",
	"475338373",
	"929119464",
	"559799207",
	"279701489",
	"66872193",
	"864392047",
	"986194170",
	"589161386",
	"983541587",
	"15077163",
	"100149904",
	"772777020",
	"902041077",
	"428233517",
	"762628806",
	"885670548",
	"842938613",
	"717424033",
	"671374074",
	"1227090",
	"657019496",
	"529975200",
	"889126175",
	"931581387",
	"357701129",
	"261897307",
	"784130655",
	"349185523",
	"755530427",
	"934661371",
	"67628852",
	"205156724",
	"984641620",
	"609360020",
	"238052748",
	"256211902",
	"862585180",
	"153002189",
	"862407392",
	"583031025",
	"481003666",
	"97942385",
	"86378037",
	"343656009",
}

func sumDigits(x int) int {
	s := 0
	for x > 0 {
		s += x % 10
		x /= 10
	}
	return s
}

type testCase struct {
	n     int
	input string
	want  string
}

func parseCases() []testCase {
	cases := make([]testCase, 0, len(rawTestcases))
	for _, line := range rawTestcases {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, _ := strconv.Atoi(line)
		var sb strings.Builder
		sb.WriteString(line)
		sb.WriteByte('\n')
		start := n - 100
		if start < 1 {
			start = 1
		}
		var ans []int
		for x := start; x <= n; x++ {
			if x+sumDigits(x) == n {
				ans = append(ans, x)
			}
		}
		var want strings.Builder
		fmt.Fprintln(&want, len(ans))
		for _, v := range ans {
			fmt.Fprintln(&want, v)
		}
		cases = append(cases, testCase{
			n:     n,
			input: sb.String(),
			want:  strings.TrimSpace(want.String()),
		})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseCases()
	for idx, tc := range cases {
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != tc.want {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx+1, tc.want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
