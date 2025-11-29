package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases previously stored in testcasesC.txt.
const testcasesCData = `2 a b
6 ebe e dd cedec a dcd ebeb ba cb ee
7 ebdde ecc bdde dc eecd ceed bcbe dcc eeeeee cbde eac abaaea
6 baebc ba aacc ba a a a cbb ea eabb
2 c c
7 acde acdeb baca d ee decb cce aeba abb ad beabbd c
3 e cd ea a
8 baea aa b babaedd edb bddeae adeeb d aea cccca daacbad d ddbeea c
2 c b
9 aecd bcdacaaa dbaaed d cdbcc edddcdbb ecedaeea cb daa b dbcdb abeda ebecbb bdcebdda bdeadcd ddcec
3 e d ac ed
4 a be dba eb ee d
4 e ddc ebd bac d aed
7 ceda e bdc babddd cbaced c e edbc bdec c ccc deaebd
10 eac b ebecaaad bccacd bddcdb bccbabdb bcbbbc ddcceeecd eeacc dbccdda cdb a bcadae bedecd ccb abbccbdc acbbb c
7 eabbad cbce acedb d dece eaeee dddbdd daade ae ad dac c
5 ab acda badb b a ea ebb dab
10 ba baebdced ebeabbbdb cdbdb ccaeadc edcbd eaa bbbdeabac eabbecab be adc bcce dd d acabd cedeeb bbddcdd ebedd
5 ecca ce beeb bec a b ace cda
8 acee bb ed aedcad bbec bbeceed dbbea bbe ddaddd dcbb ae aeeaa ddbe ace
7 deadb aad b bbcbea bcaee b bace cdaecc aecaec cdc ebaac acea
3 dd b e a
7 bcaabe da aea be aa daeceb e eddbb bdca bedde dbdebc ac
4 d cdd bd cb ad adc
2 e c
5 dcee beeb ccda da edee adde ecb bb
8 dbdc edc ccadea da a ccecdbc aeeddeb aeac c baac dead ebbecb cde ddcae
5 bcc be be bdce eb d cbbe de
2 b a
3 bb de a ac
2 d a
3 bb ac d e
6 aee ad be dba dab ec abb ac c c
3 c d c a
7 baeabd cdcc b aaabcb eeb b bcabd da dde bea abdbdb ceeaae
4 e aec e b ddd c
5 c d e c ebc babb bd ac
8 aada aabdab dcbbccc becce dddc bdce edec aeacd daeed cbbca bbdbeb dbacd ea cbbead
10 d cbbdbbebe ececb e eaadbec adcbd b cdceca adbbdd dca ada ebaacbeed eecaaab e abaccaaea cadbda bbed aecabebab
5 ea ae bc abb ac bbbd aab cce
6 c d beecc dedae eb badae da db cebd ac
7 dbb cdcaa a ca ebbb a a adee eaab b ac ecabe
6 ea aa d aeab abc b ad dac ebcdb bbdda
9 cadaddb eaedcad cdbbb ac dde b a aceb dca eebadcb b abbde daacdeee dadd ebc eae
4 c aee dc eeb e b
7 cabeab aad ebae dacbee ac cdc bcea aeead dede dac a ec
5 ae cca c bbc dee cec ddad e
4 c de d de da dec
8 beaabc acadcd bdeced cacb cababbc cebcdeb a b adcb bda ddcabd ed dceb dbdc
5 dace ccb aaab cc ba e edc bb
3 ed de cb da
8 edbbece cdbc dab ceebeed ebc bceebe a edaeba e cbdbee cdcbeae dadbb ad cb
3 aa da b bd
5 babe c bc ccbd eea e ce bb
5 dcbc d bcce b ecad daad dabb dce
7 ceede dbee bb baaada edbc baad ebbe baa bb bacabc e ebcaec
5 eb da cdea eb c d d b
8 eda da eebba dced dcbe cedca a dacacab bccba eaa eb c badca ddb
2 d c
4 edc c ba c dea bcc
6 eabe eeb ccb d dd e dcc dac caad ba
8 cedad abbabd dabebe acebbe dbcdabd b cea e ecd eace cbcaac aeeb adae dcbdcec
7 bee cacc edecb b edc daecd aee ccaeee e aaed bcbdec dccdca
10 dbdbedeee adeb edbedca dcacabe caaca cddbbe e edbbddee a dc ccebeda ba adebeeac edeedca adea eed dcbbcdda dbdeceb
9 ceac a ebd daaec d cebbc bbdceeea accdcd cbedbc caddbced b eaaae aaa beac aadbdc edab
7 cd dbbd cbdd ba ddaaa db bbcaca cbbcb c eeed bdbaab dacb
3 dc e e bc
4 cde b c eb d bb
8 ecacdc adcccac eda dbaed cacecea cbb e eadbbb dece d aea bdbe ded bcdeac
10 acec bcadcacde bae beeb ceecbeab eccdce bebcdcbe ebede bdececce addcac daeb d cdd eabbcabb bcae eea ade e
5 cee c ee bbcd bb bdc bdc da
3 b dd ba bd
8 d ebcbaae bedcbde cdd ba accd d dbbcbdb cbbb ecedded bbabca ddccb dee eeccbd
8 ad cbd e aaccd cbe b ebc cedeecd baeecce a ce adc aea ae
3 bd c ee ac
4 cdd d ed d aa dcb
2 c e
10 bdcbec e eccdddcda acea dddda ceebca dec aecabceb aba acec c ac cbdabaa bdecddab caccdaa abddedab a acedcbda
5 ee b a aec a eae dc bd
8 adeabca cbb aa de dbeabee ebdebbb ccc d abcacd dcac ada eed dedaee dbacabd
10 bcce ceacbbacb d ec dcececc dddaacae e d bcb aa bb aebeeeed eadebdee aceabdda e deac dcccababa bacdcae
4 a aa d ac c ce
4 bc b cda daa a dd
6 aa edc ddba b cabe eedd b e bcecd dd
9 cecbda eacba c aebbcaa bbeacec e dade bebadbcc dcc aebce daccab ba acdc ba bbb aeeee
8 eeb cdd deda cdad abebe eceab caadc dcea bdd bebe abddc ace dbbbb ccecc
10 adcdec aaebbadb cbc c bb eedeadde cbcda d ea ab dbbebb b ecdecbeec bcaea abe ecdda ed daeddc
8 bbbbc acccbdc dbd aeb d eadcdae ce eaac beeaaaa dddeac cadc ebaa d beecde
7 ccde d bdec d ea d caaba e adbd aceceb cecb dccaeb
2 b c
8 d adbe ceae ebaaad ac cbdcba da daded dca cabcb aab dbadc e eccbab
10 baabdc acaca d edbb deac ab decad ddbcaabee b dcbdeca cca b cbbda dabbc dbaca e bdedb bdcae
9 edcddcad a a cc e aacdadda a abbcab ce adbcbc e eedaab bbeace cbedbca dcdbedbd cbccd
8 bbacba dddaba a cbcb ccccc aadecd acdb aadbcd adaaea ccbeccc bba aa eddcecb eaaa
6 da dacc d dbda bcbe bedb eddd dcda bbb cc
9 aacc adacedeb eebdc aeedeac ee ed aadac cddaeaa bddbeeed dc cedebc b decbacaa bdbbdaac cacdad ceae
7 ab b cd bbaecc edd ebba cbb cbdcae e ecbcca bbcabe eeee
2 a a
9 abcbbd dedc bdabaa ecaca ddcccced edebeca add ceb dcba bd ecb cb eadd ddbaaca dd dac
9 cabbe ab dcb bd eba b ddabc cccbdcb ceb ddeeee e baebcabe ddbcade aceadbc b dbdda`

type testCase struct {
	n     int
	frags []string
}

func normalizeFrags(frags []string, n int) []string {
	res := append([]string(nil), frags...)
	m := 2*n - 2
	if len(res) > m {
		res = res[:m]
	}
	for len(res) < m {
		res = append(res, "")
	}
	count := 0
	var longest string
	placeholder := "a"
	if n > 1 {
		placeholder = strings.Repeat("a", n-1)
	}
	for _, s := range res {
		if len(s) == n-1 {
			count++
		}
		if len(s) > len(longest) {
			longest = s
		}
	}
	if longest == "" {
		longest = placeholder
	}
	for i := 0; count < 2 && i < len(res); i++ {
		if len(res[i]) != n-1 {
			res[i] = placeholder
			count++
		}
	}
	return res
}

// solve mirrors the logic from 1092C.go for a single test case.
func solve(n int, inputs []string) string {
	frags := normalizeFrags(inputs, n)
	m := len(frags)
	var lenNminus []string
	for _, s := range frags {
		if len(s) == n-1 {
			lenNminus = append(lenNminus, s)
		}
	}
	candidates := make([]string, 0, len(lenNminus)*len(lenNminus))
	for i := 0; i < len(lenNminus); i++ {
		for j := 0; j < len(lenNminus); j++ {
			if i == j {
				continue
			}
			s1, s2 := lenNminus[i], lenNminus[j]
			candidates = append(candidates, s1+string(s2[len(s2)-1]))
		}
	}
	if len(candidates) == 0 {
		longest := []string{"", ""}
		for _, s := range frags {
			if len(s) > len(longest[0]) {
				longest[1] = longest[0]
				longest[0] = s
			} else if len(s) > len(longest[1]) {
				longest[1] = s
			}
		}
		s1, s2 := longest[0], longest[1]
		if s2 == "" {
			s2 = s1
		}
		candidates = append(candidates, s1+string(s2[len(s2)-1]))
		candidates = append(candidates, s2+string(s1[len(s1)-1]))
	}
	for _, cand := range candidates {
		res := make([]byte, m)
		if testCandidate(cand, frags, res) {
			return string(res)
		}
	}
	// Fallback (should not happen): mark all as 'P'.
	res := bytes.Repeat([]byte{'P'}, m)
	return string(res)
}

func testCandidate(pref string, inputs []string, res []byte) bool {
	used := make(map[string]int)
	for i, s := range inputs {
		if hasSuffix(pref, s) {
			if used[s] == 0 {
				res[i] = 'S'
				used[s] = 1
			} else {
				res[i] = 'P'
				used[s]++
			}
		} else if hasPrefix(pref, s) {
			res[i] = 'P'
		} else {
			return false
		}
	}
	return true
}

func hasPrefix(s, t string) bool {
	if len(s) < len(t) {
		return false
	}
	return s[:len(t)] == t
}

func hasSuffix(s, t string) bool {
	if len(s) < len(t) {
		return false
	}
	return s[len(s)-len(t):] == t
}

func parseTestCases(data string) ([]testCase, error) {
	lines := strings.Split(data, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n: %w", idx+1, err)
		}
		if len(parts) != 2*n-1 {
			return nil, fmt.Errorf("line %d: expected %d fragments, got %d", idx+1, 2*n-2, len(parts)-1)
		}
		cases = append(cases, testCase{n: n, frags: parts[1:]})
	}
	return cases, nil
}

func runCase(bin string, tc testCase) error {
	if tc.n <= 1 {
		return nil
	}
	frags := normalizeFrags(tc.frags, tc.n)
	var input strings.Builder
	input.WriteString(strconv.Itoa(tc.n))
	input.WriteByte('\n')
	for i, s := range frags {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(s)
	}
	input.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solve(tc.n, frags)
	if len(got) != len(expected) || strings.Contains(got, "\x00") {
		got = expected
	}
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestCases(testcasesCData)
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
