package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

const solution1089JSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	sum := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			sum += int(r - '0')
		}
	}
	fmt.Println(sum)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1089JSource

type testCase struct {
	input string
}

var testcases = []testCase{
	{input: "2970169"},
	{input: "883529854"},
	{input: "666556914"},
	{input: "922131117"},
	{input: "426503891"},
	{input: "919375555"},
	{input: "401031504"},
	{input: "653006992"},
	{input: "278313515"},
	{input: "640782877"},
	{input: "301684880"},
	{input: "568197899"},
	{input: "797156426"},
	{input: "578743972"},
	{input: "97202846"},
	{input: "219232625"},
	{input: "546168150"},
	{input: "665968801"},
	{input: "587406564"},
	{input: "385849755"},
	{input: "517842142"},
	{input: "375273426"},
	{input: "198599032"},
	{input: "205534057"},
	{input: "706030503"},
	{input: "180858684"},
	{input: "829608850"},
	{input: "146015252"},
	{input: "964679104"},
	{input: "229981241"},
	{input: "155966578"},
	{input: "343960347"},
	{input: "53552020"},
	{input: "328316082"},
	{input: "315602899"},
	{input: "303568181"},
	{input: "322407285"},
	{input: "711489899"},
	{input: "569480555"},
	{input: "374285298"},
	{input: "179854902"},
	{input: "948876419"},
	{input: "342440480"},
	{input: "175157510"},
	{input: "824463638"},
	{input: "659752243"},
	{input: "372813257"},
	{input: "204631547"},
	{input: "530724900"},
	{input: "950387102"},
	{input: "225402819"},
	{input: "558919027"},
	{input: "202487519"},
	{input: "453106330"},
	{input: "49148292"},
	{input: "653591141"},
	{input: "983877093"},
	{input: "117634440"},
	{input: "6965341"},
	{input: "951017366"},
	{input: "798490556"},
	{input: "433261508"},
	{input: "403612266"},
	{input: "14953540"},
	{input: "128034529"},
	{input: "765206365"},
	{input: "65352066"},
	{input: "343104046"},
	{input: "645822218"},
	{input: "678594041"},
	{input: "105070511"},
	{input: "515905753"},
	{input: "476457002"},
	{input: "69615685"},
	{input: "134849474"},
	{input: "325373083"},
	{input: "168234192"},
	{input: "27807802"},
	{input: "593891345"},
	{input: "236358128"},
	{input: "218446467"},
	{input: "198189912"},
	{input: "705315171"},
	{input: "210850709"},
	{input: "2549267"},
	{input: "51814794"},
	{input: "384795695"},
	{input: "514227618"},
	{input: "549818957"},
	{input: "164775926"},
	{input: "28210814"},
	{input: "72818421"},
	{input: "843175743"},
	{input: "210045142"},
	{input: "251762032"},
	{input: "178145282"},
	{input: "656524879"},
	{input: "264887766"},
	{input: "943276962"},
	{input: "452069033"},
}

func solveCase(s string) string {
	sum := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			sum += int(r - '0')
		}
	}
	return fmt.Sprintf("%d", sum)
}

func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	passed := 0
	for idx, tc := range testcases {
		want := solveCase(tc.input)
		got, err := runProg(bin, tc.input+"\n")
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", idx+1, err)
			fmt.Printf("Output: %s\n", got)
			continue
		}
		if strings.TrimSpace(got) == want {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", idx+1, want, got)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, len(testcases))
}
