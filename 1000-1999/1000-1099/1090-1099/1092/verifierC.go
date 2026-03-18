package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Embedded testcases previously stored in testcasesC.txt.
const testcasesCData = `2 a a
2 b b
2 e a
2 d e
2 a a
2 e b
2 e b
2 b d
2 a b
2 c c
3 bc ca a b
3 ce e c cc
3 ad d ea e
3 ce e ec e
3 c ab bc a
3 cd d d dc
3 c bc cb c
3 e eb b bb
3 c eb bc e
3 cd c dc c
3 bd d b dd
3 b e be bb
3 ed dc e c
3 aa a aa a
3 d de a ea
4 ece a ea ec e cea
4 b acd cdb ac a db
4 ea e ce ace e eac
4 eae e eea e ae ee
4 ab b cb bab c cba
4 a ae aeb b bb ebb
4 b c beb ebc be bc
4 ba abb bba a a ab
4 b ab aaa aa aab a
4 cd dbe c be cdb e
4 d dba a db aa baa
4 a daa aaa d da aa
4 edb b bed be b db
4 a ad ea dea ade a
4 bd ddd bdd b dd d
5 c ad a dc adc dcdc cdc adcd
5 ee e cbae cba baee cb aee c
5 bae ae ee eeb e ebae eeba e
5 daeb ebe e da be aebe dae d
5 cdb cbcd b bcdb c cb db cbc
5 ec bec aeb aebe c a ebec ae
5 dece e ee ecee de dec cee d
5 a bc aaeb ebc aa c aebc aae
5 caa aa a edc e ed dcaa edca
5 cb ed cbd ded d cbde bded c
5 cc dbac dba bacc acc c db d
5 bbb bbbb ebb bb eb b e ebbb
5 cba c ada cb da cbad a bada
5 bdc c dc abdc ba b babd bab
5 ac c aac ecaa e caac eca ec
6 a deea cc cdeea ccdee eea c ccd ea ccde
6 c dacec acec cdac cd ec cec cdace cda c
6 db b ddb bddb eb e ebbdd ebbd ebb bbddb
6 ecddd ee eec e d ddd eecd cddd eecdd dd
6 ab bbb b abcb abcbb abc cbbb bb bcbbb a
6 dddbb dd d dbba ddd bba ddbba ba dddb a
6 a e ad ee adbd adbde adb bdee dee dbdee
6 cbc c bdd bddc b ddcbc bd bc dcbc bddcb
6 cce eabb cceab ceabb c b abb cc ccea bb
6 abd dabd d edd ed eddab bd ddabd e edda
7 ee ebdbc eee c bc eeebd eeebdb dbc eebdbc e bdbc eeeb
7 eea eea adeea ea eeade e eadeea deea ee eead eeadee a
7 cdc dccd dc dcc dccdcd ccdcdc dc cdcdc dccdc d dcdc c
7 baebb b bae aebbbd baeb ba bbd bbbd bd ebbbd d baebbb
7 bcaea aeace bc b eace ace bca e ce caeace bcae bcaeac
7 eae ea eaced eaea aced ed e aeaced eaeac d ced eaeace
7 ce dbece ece bdbece d e db dbd dbdb dbdbec bece dbdbe
7 d d deedc dee edcad deedca deed eedcad ad dcad de cad
7 eaeb d eae ea abd e ebabd eaebab aebabd bd eaeba babd
7 bcecde b e bce bc ee cdee ecdee cecdee bcecd bcec dee
8 abcaeb a bb bbb ab abcaebb abca b bcaebbb abcae aebbb abc ebbb caebbb
8 bcaae aae caebc caeb ae caebcaa aebcaae c e caae caebca ca ebcaae cae
8 acdaadd add dd bacdaad aadd bac bacd bacda ba b cdaadd bacdaa daadd d
8 ed edeeeb ed d edeeebe bed edeee e eeebed edee eebed ede deeebed ebed
8 cabbbea cabbb cabbbe ab eab ca b abbbeab cab bbbeab bbeab beab cabb c
8 ebb bdaebb aebb daebb cebd cebdae bb cebda ebdaebb cebdaeb c b ce ceb
8 ea eeddaea ce ceedd ceedda ceeddae daea a c cee ddaea aea eddaea ceed
8 eab abdedcb eabdedc e eabde ea cb dcb bdedcb b edcb eabded eabd dedcb
8 ddcdea dd eada deada da ddcdead cdeada d ddcde ddcd a dcdeada ddc ada
8 eceddab ddabc bc eced ceddabc eddabc abc e c ecedda dabc ece ec ecedd
9 abadabd e ed abd edabadab d edabada dabadabd bd edabad dabd edab eda badabd edaba adabd
9 a acdcec cecacc acd c acdcecac cc acdc ecacc acc acdceca acdce cacc ac dcecacc cdcecacc
9 bda bdac bdacdddb bd acdddbd cdddbd b bdacddd bdacdd ddbd dacdddbd dbd dddbd bdacd bd d
9 dccedace edace adc ccedace ad adcc adccedac ce e a adcce dace cedace adcced ace adcceda
9 bcaac daadbc d caac aadbcaac daadbca daadbcaa dbcaac daadb ac da daad adbcaac aac daa c
9 edeb d bbd debbd edebdeb bd bdebbd ed debdebbd ebdebbd edebdebb e edebde ede ebbd edebd
9 cd cdcbddab d bd dabd abd cdcbd cdcbdd cdcbdda cdcb c cbddabd dcbddabd cdc ddabd bddabd
9 bbecbaa ecbb ecbbecba ecbbe aa ecbbec becbaa ecb a ec baa cbbecbaa e cbaa ecbbecb ecbaa
10 c ece cdcc eceecbcdc bcdcc ecbcdcc e dcc eceec ecee cbcdcc eceecb ceecbcdcc eceecbc eecbcdcc ec eceecbcd cc
10 dcced dccedca bda da dcabda cabda d dc edcabda dcce dcc cedcabda dccedcab a dccedcabd abda dccedc ccedcabda
10 bebddce ebddcedeb bebddced bebddc dcedeb be b bebdd edeb bebddcede deb b beb bebd cedeb bddcedeb eb ddcedeb
10 cbadabea a adc bea adcbadab adcbadabe adcb badabea adabea adcba dcbadabea ea dabea adcbada ad abea adcbad a
10 becbaaba baaba cbbecbaab cbbecba ba bbecbaaba aba a cbbec cbbe ecbaaba cbbecbaa cbb cb c cbbecb aaba cbaaba
10 debe de aa debeaecd d debea deb daa ecdaa ebeaecdaa debeaecda aecdaa cdaa a debeae eaecdaa debeaec beaecdaa
10 deedaa d deecde dee deecd ecdeedaa de a eedaa deecdee deecdeed daa deecdeeda aa edaa cdeedaa deec eecdeedaa`

const embeddedRefGo = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return
	}
	var n int
	fmt.Sscanf(scanner.Text(), "%d", &n)

	strs := make([]string, 2*n-2)
	var l1, l2 string
	for i := 0; i < 2*n-2; i++ {
		scanner.Scan()
		strs[i] = scanner.Text()
		if len(strs[i]) == n-1 {
			if l1 == "" {
				l1 = strs[i]
			} else {
				l2 = strs[i]
			}
		}
	}

	cand1 := l1 + string(l2[n-2])
	cand2 := l2 + string(l1[n-2])

	if check(cand1, n, strs) {
		return
	}
	check(cand2, n, strs)
}

func check(s string, n int, strs []string) bool {
	usedP := make([]bool, n)
	ans := make([]byte, 2*n-2)
	for i, str := range strs {
		l := len(str)
		isP := s[:l] == str
		isS := s[n-l:] == str

		if isP && !usedP[l] {
			ans[i] = 'P'
			usedP[l] = true
		} else if isS {
			ans[i] = 'S'
		} else {
			return false
		}
	}
	fmt.Println(string(ans))
	return true
}
`

type testCase struct {
	n     int
	frags []string
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
		if len(parts) < 2*n-1 {
			return nil, fmt.Errorf("line %d: expected %d fragments, got %d", idx+1, 2*n-2, len(parts)-1)
		}
		cases = append(cases, testCase{n: n, frags: parts[1 : 2*n-1]})
	}
	return cases, nil
}

func buildRef() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %v", err)
	}
	ref := filepath.Join(wd, "refC.bin")
	goPath := filepath.Join(wd, "refC.go")
	if err := os.WriteFile(goPath, []byte(embeddedRefGo), 0644); err != nil {
		return "", fmt.Errorf("write go: %v", err)
	}
	cmd := exec.Command("go", "build", "-o", ref, goPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference go: %v: %s", err, string(out))
	}
	return ref, nil
}

func runBin(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin, ref string, tc testCase) error {
	if tc.n <= 1 {
		return nil
	}
	frags := tc.frags
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

	inStr := input.String()

	exp, err := runBin(ref, inStr)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}

	got, err := runBin(bin, inStr)
	if err != nil {
		return fmt.Errorf("candidate error: %v", err)
	}

	if len(got) != len(exp) {
		return fmt.Errorf("invalid answer length %d, expected %d", len(got), len(exp))
	}

	// Both must be valid P/S assignments of the same length.
	// We validate structurally: for each position, it must be P or S,
	// and the implied string must be consistent.
	if got != exp {
		// Different assignments can both be valid. Validate got independently.
		if err := validateAnswer(tc.n, frags, got); err != nil {
			return fmt.Errorf("invalid answer %q: %v", got, err)
		}
	}
	return nil
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

func validateAnswer(n int, frags []string, answer string) error {
	m := len(frags)
	if len(answer) != m {
		return fmt.Errorf("answer length %d != expected %d", len(answer), m)
	}
	for _, ch := range answer {
		if ch != 'P' && ch != 'S' {
			return fmt.Errorf("invalid character %q in answer", ch)
		}
	}

	var longestP, longestS string
	for i, ch := range answer {
		if len(frags[i]) == n-1 {
			if ch == 'P' {
				longestP = frags[i]
			} else {
				longestS = frags[i]
			}
		}
	}
	if longestP == "" || longestS == "" {
		return fmt.Errorf("no length-%d fragment assigned as P or S", n-1)
	}
	if n >= 3 && longestP[1:] != longestS[:n-2] {
		return fmt.Errorf("longest P %q and S %q overlap mismatch", longestP, longestS)
	}
	s := longestP + string(longestS[n-2])

	for i, ch := range answer {
		frag := frags[i]
		if ch == 'P' {
			if !hasPrefix(s, frag) {
				return fmt.Errorf("fragment %d %q marked P but is not prefix of %q", i, frag, s)
			}
		} else {
			if !hasSuffix(s, frag) {
				return fmt.Errorf("fragment %d %q marked S but is not suffix of %q", i, frag, s)
			}
		}
	}

	pCount := make(map[int]int)
	sCount := make(map[int]int)
	for i, ch := range answer {
		l := len(frags[i])
		if ch == 'P' {
			pCount[l]++
		} else {
			sCount[l]++
		}
	}
	for l := 1; l < n; l++ {
		if pCount[l] != 1 || sCount[l] != 1 {
			return fmt.Errorf("length %d: P=%d S=%d, expected 1 each", l, pCount[l], sCount[l])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer func() {
		os.Remove(ref)
		wd, _ := os.Getwd()
		os.Remove(filepath.Join(wd, "refC.go"))
	}()

	cases, err := parseTestCases(testcasesCData)
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for i, tc := range cases {
		if err := runCase(bin, ref, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
