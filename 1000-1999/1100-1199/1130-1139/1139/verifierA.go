package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const solution1139ASource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	var res int64
	for i := 0; i < n && i < len(s); i++ {
		if (s[i]-'0')%2 == 0 {
			res += int64(i + 1)
		}
	}
	fmt.Fprint(writer, res)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1139ASource

type testCase struct {
	n int
	s string
}

var testcases = []testCase{
	{n: 13, s: "7159875864935"},
	{n: 5, s: "25935"},
	{n: 4, s: "2689"},
	{n: 4, s: "6764"},
	{n: 18, s: "889519127186462444"},
	{n: 5, s: "98226"},
	{n: 17, s: "82595296949582764"},
	{n: 10, s: "3431582233"},
	{n: 2, s: "29"},
	{n: 13, s: "9594475886262"},
	{n: 16, s: "6441524636712341"},
	{n: 19, s: "9212427262114328411"},
	{n: 18, s: "725242567319812745"},
	{n: 12, s: "834133695283"},
	{n: 1, s: "8"},
	{n: 14, s: "95675391826195"},
	{n: 5, s: "48656"},
	{n: 19, s: "3577214634487717713"},
	{n: 15, s: "253898911865817"},
	{n: 7, s: "9231776"},
	{n: 1, s: "4"},
	{n: 1, s: "1"},
	{n: 17, s: "24245532872158253"},
	{n: 17, s: "62351114596618876"},
	{n: 18, s: "347513356662611533"},
	{n: 19, s: "5679352841539257657"},
	{n: 4, s: "2988"},
	{n: 11, s: "62828715633"},
	{n: 19, s: "7222441712799588472"},
	{n: 12, s: "453746221984"},
	{n: 4, s: "8754"},
	{n: 2, s: "43"},
	{n: 4, s: "4876"},
	{n: 18, s: "328377986884941666"},
	{n: 2, s: "93"},
	{n: 9, s: "375822912"},
	{n: 8, s: "31518633"},
	{n: 15, s: "697991292745978"},
	{n: 13, s: "4113595628557"},
	{n: 13, s: "7133456118738"},
	{n: 20, s: "23671878128311362964"},
	{n: 13, s: "8218625375294"},
	{n: 2, s: "78"},
	{n: 12, s: "486211851944"},
	{n: 3, s: "997"},
	{n: 17, s: "52377227227318771"},
	{n: 16, s: "6526226166314623"},
	{n: 7, s: "1421561"},
	{n: 20, s: "43382865314668559632"},
	{n: 4, s: "9537"},
	{n: 5, s: "34694"},
	{n: 8, s: "35671317"},
	{n: 3, s: "237"},
	{n: 10, s: "9737562486"},
	{n: 17, s: "17717684658232529"},
	{n: 20, s: "38737734869368284518"},
	{n: 20, s: "81452593782849751251"},
	{n: 1, s: "5"},
	{n: 13, s: "9782565421255"},
	{n: 18, s: "629432755939492797"},
	{n: 9, s: "586332277"},
	{n: 19, s: "8395687488968185381"},
	{n: 20, s: "41687192271612155435"},
	{n: 7, s: "2786736"},
	{n: 14, s: "73839634386778"},
	{n: 13, s: "4484171423613"},
	{n: 8, s: "52956781"},
	{n: 17, s: "97885846511136151"},
	{n: 5, s: "27479"},
	{n: 8, s: "84622669"},
	{n: 15, s: "651914624964455"},
	{n: 10, s: "9758641592"},
	{n: 1, s: "8"},
	{n: 16, s: "8178882224237482"},
	{n: 14, s: "97134843566729"},
	{n: 10, s: "9458989554"},
	{n: 1, s: "2"},
	{n: 20, s: "23744519971275264954"},
	{n: 8, s: "29564685"},
	{n: 19, s: "3319966137339234844"},
	{n: 5, s: "47638"},
	{n: 4, s: "1968"},
	{n: 15, s: "514938896253746"},
	{n: 10, s: "7141646845"},
	{n: 12, s: "351691361811"},
	{n: 8, s: "11462167"},
	{n: 5, s: "48736"},
	{n: 10, s: "3677175998"},
	{n: 2, s: "27"},
	{n: 13, s: "3193932643413"},
	{n: 18, s: "327283156711826538"},
	{n: 8, s: "96376587"},
	{n: 1, s: "5"},
	{n: 17, s: "59819951872768115"},
	{n: 2, s: "55"},
	{n: 7, s: "9967542"},
	{n: 19, s: "6496336113665568773"},
	{n: 1, s: "3"},
	{n: 19, s: "1836185429753932329"},
	{n: 18, s: "775551755999664731"},
}

func solveCase(tc testCase) int64 {
	var res int64
	for i := 0; i < tc.n && i < len(tc.s); i++ {
		if (tc.s[i]-'0')%2 == 0 {
			res += int64(i + 1)
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
		want := solveCase(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d failed: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		if gotStr != fmt.Sprintf("%d", want) {
			fmt.Printf("test %d failed: expected %d got %s\n", idx+1, want, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
