package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve embeds the logic from 1367D.go.
func solve(s string, b []int) string {
	m := len(b)
	freq := make([]int, 26)
	for _, ch := range s {
		freq[int(ch-'a')]++
	}

	res := make([]byte, m)
	used := make([]bool, m)
	remaining := m
	ch := 25
	for remaining > 0 {
		zeros := make([]int, 0)
		for i := 0; i < m; i++ {
			if !used[i] && b[i] == 0 {
				zeros = append(zeros, i)
			}
		}
		if len(zeros) == 0 {
			break
		}
		for ch >= 0 && freq[ch] < len(zeros) {
			ch--
		}
		if ch < 0 {
			break
		}
		for _, pos := range zeros {
			res[pos] = byte('a' + ch)
			used[pos] = true
		}
		freq[ch] -= len(zeros)
		ch--
		remaining -= len(zeros)
		for i := 0; i < m; i++ {
			if used[i] {
				continue
			}
			sum := 0
			for _, pos := range zeros {
				if i > pos {
					sum += i - pos
				} else {
					sum += pos - i
				}
			}
			b[i] -= sum
		}
	}
	return string(res)
}

// Embedded testcases from testcasesD.txt.
const testcaseData = `
100
sreltpusctapirhg
8
34 35 30 25 40 9 14 40
qmxavycfys
1
19
aiptxmwznmxzsoeldbepgivnyujnqmslrsnshkvaitvwfwkrss
2
45 41
usijdcpupclzcn
3
1 18 27
ndbttybmwskriqhbjacdtrbgnjtiewbkklemmoqmutvrdtzqin
4
19 27 16 33
rkaznskamtsuebuukolv
6
38 45 17 47 31 1
bvaliuojstkflfkyltijzmdyasvxejqhuzihkf
7
41 44 6 6 38 20 21
hozfckxugsoihzdbqgkzsfikzucztlsenjqziolunjns
7
2 26 9 12 0 30 39
nrwhbxoyvxqjrkhcsjdzhbbzwqgnsbapx
2
10 32
hvaqrnbtdkeirpzzblhg
2
34 7
hzizeapusmb
5
15 17 39 33 33
bpkyabyebdbcpbwcqqpkfkclmums
5
23 16 12 21 27
erawxmzc
10
11 2 23 29 38 41 50 34 24 40
tnb
6
40 31 48 44 20 26
noahgriwscznhneaklrzidowdxvqzmvdxksrdzswapehy
7
2 33 5 36 6 42 24
akdadvpwjsjz
2
2 49
qqwhdrxdrbrksfchfuhotwymiltmlrncmqhnx
3
26 44 36
svqvpeumefdpxpwqosxfeiygesqkhwryjvwntssigjaipzmgf
10
23 15 20 30 49 9 26 44 30 44
gosurapxcmzxbohhuwyvcgihgyieftwvbifbkfn
2
46 5
cijblosx
6
0 1 21 21 27 24
cgusxpmerkdicvndoqidqwlvylyojvvv
5
6 48 21 43 36
qdvpqlbwjvxsxfuuxuefluoddrekuxutnrj
3
29 30 19
wcdwfyrrsxml
2
17 17
bebpqihwyqlkmorzyclpdeisd
2
36 49
dfwgsnvxmxestemzgrqfsfgilzjazonmkrsjupqvwjvpatg
1
6
yvhpfquoggzqgbquodsjveeozctbalthqcprakkklwwectybwc
6
13 4 12 27 44 48
pkdzbncgwfmppwc
9
27 13 41 31 19 1 29 29 48
mofobxilloqltmhazgizleorgfgafsmqfuaedtfopfbam
8
20 26 2 45 45 3 15 25
mpa
4
15 6 24 30
fktdldtbzxjiz
8
50 19 31 15 35 17 1 21
lkcbvncstadavcafqbpbguqkgypkzplvbmjytumcj
3
26 7 32
rkrvymfxxmrlflznohoywplif
9
46 48 38 45 44 24 31 2 9
wyaocywvvdk
4
38 41 3 39
ooxu
6
23 0 4 12 25 50
ksjdocu
4
15 44 3 9
esadhjghrqnqytkzrzgoftcbzdtadgicdomhvtdupy
6
25 38 42 28 49 7
tomgdraojxucklgpycr
6
27 50 41 4 38 33
hlbkhnocigkfxg
4
47 39 29 45
nlgtnpynpsbjafdaxejqqbupbgxgipnblox
4
47 50 18 9
ojnocge
8
48 44 18 24 40 23 10 27
opqrhljjaollzjxhzqaa
3
40 33 9
afbagyzolzlrbpfhainkbtrzdojihvzpnxi
6
2 1 27 2 40 43
yzshexyznqy
6
35 8 17 1 10 2
pu
1
29
qvytqnlqufjfcverdnylooiziojqes
6
8 33 2 26 31 14
stiakstrdpeizwyidnvclbqpyuogjl
3
41 24 25
bigbkktmrjbenincphgxc
9
7 47 40 7 40 0 18 44 4
ipoijrrbfhpfeewfwovmuaembfyu
3
19 12 41
ebqergmyd
7
24 11 1 17 6 8 7
jemltcgale
8
15 4 22 34 31 6 45 20
axlqxzosnorrjoeromgytjxxfjfkige
1
38
ntfd
10
0 10 7 50 25 36 43 23 33 44
icorsokestgkoqslukslwtyxlwkjjifdtqhxkxuhjnioepkrzf
10
33 33 28 39 3 4 26 27 35 39
jbhmmgclqgbrpdnxmwrlajlqlmoluvdspekhalctaecgknjgaa
9
46 20 34 28 47 23 46 13 28
vsdqmhpejjrgdftcnazlma
8
11 30 17 8 24 13 32 40
suvjnjduxwccoyxkcakpnudvkuns
4
20 12 41 25
buqaqq
4
36 5 11 14
rmjlovqhgsjksesrltqlkysxxqgobxn
5
40 12 31 13 11
odqmgkf
1
20
hmvxbhmlmhfjilibwkxdghjofsfhpwgg
6
38 14 11 43 48 33
trntfujolbczropasfiqnvphmq
8
22 25 39 30 48 11 7 45
kxeskdlqejpureolpbslugtcwyr
10
19 37 43 18 34 20 26 47 41 18
n
9
47 48 22 27 30 23 36 11 15
wgsmwdkrfkxpfrmognr
10
23 6 46 3 38 30 44 12 10 32
fddrrdmrqkmicpivkhidgdplneilvfisyqgvhhizlxk
2
6 45
ychokdkye
3
2 42 31
qeloxnqsvsnznijris
7
9 2 2 10 27 40 0
zyeqxgwxulibur
10
47 22 12 7 10 13 47 4 31 27
jedrkvtxfiectcvedawcuzqmworvtoifwacjvwnfabkqxkld
1
13
azrovxkuhoohzalghfnmdmksd
9
42 38 15 18 6 1 25 1 37
hgosqiuytmvxgqsnujzcfz
8
32 33 6 18 37 39 9 28
ejqfhlqjdcpcpzhxynaaaptabwqpablvjj
8
31 19 49 19 8 19 40 36
mubdvg
4
5 10 47 25
xyxqeufcvjqqapclvpcnvwkrzmnoktakjp
5
50 30 47 16 13
duixpnpygxktbvhsbtmhmfk
2
41 11
oinnsskmnblxymlxmiqhehgytaunxdtym
8
30 24 32 32 35 27 39 12
mburwfagsvzwdemyqdddtgjrennefgw
7
31 20 40 28 32 1 16
askkrmgdjgcgmhhkhxckrjygrgrxmaqujhey
5
37 4 50 16 22
jelggurzdonibr
7
16 37 18 33 12 28 13
lvcdnymwlzuyshixdtsd
1
11
ltfhkbdyktrcnkrjapoyw
8
49 10 30 0 7 50 10 23
imqgomplwlxarkpsproapcwjfhrdvff
5
20 5 23 34 0
katbmfduyqolxpz
1
29
nurxcdycvcgdicvkazgbwobdqxgkwoxtjykwagwruagy
10
27 8 44 22 28 33 33 26 42 30
txzjujchy
3
34 21 28
hqwrkjklyyyfjlcfmxjtomqng
2
40 22
t
6
28 3 41 41 29 33
titprqxpyxgzjcwcpzomwfgiwersiafguhovsdjlvgq
2
0 20
jptbzukzjfuosspjqpkhcgrblrcla
3
2 30 10
dflshpwqurbayrhvkesjr
4
22 50 22 37
eonvwmbomauwytcinjldafzzd
7
47 2 25 44 27 48 15
cnbilfnvlnxahwtlqhxmlwljtfhlxahufpigespyqsof
7
16 24 48 36 25 5 13
udgyvslzlijwvumuryneulqcj
6
38 48 17 46 12 12
wmuguxgqhlsiohuydxdqsavsnbztrsjjtlexqeifzpyyvmi
9
15 10 43 37 48 50 48 30 5
lasvpbrnndtmk
7
19 22 40 37 25 4 29
cvioyxdirchsr
6
17 48 40 15 37 43
knbcsaegthuzdglfrycfmqoadkgmznyhdbcwcrxjv
4
11 28 11 5
lcpqveehhupimvbhtjvlcuffutaodyzhdferkwgrynajc
8
15 32 29 15 7 25 36 40
tfhvqideozznatmdhrnxfccpazfhktkpwojcdzlssbalkbcf
5
7 1 18 43 45
puanwmbikerlguxckzsjrxumemicwsgpvqcmjctx
9
35 48 20 43 15 12 14 33 46
`

// testCase represents a single query.
type testCase struct {
	s string
	m int
	b []int
}

func parseTestcases() ([]testCase, error) {
	tokens := strings.Fields(testcaseData)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	idx := 0
	next := func() (string, error) {
		if idx >= len(tokens) {
			return "", fmt.Errorf("unexpected end of test data")
		}
		val := tokens[idx]
		idx++
		return val, nil
	}

	tStr, err := next()
	if err != nil {
		return nil, err
	}
	t, err := strconv.Atoi(tStr)
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		s, err := next()
		if err != nil {
			return nil, fmt.Errorf("case %d: missing string", i+1)
		}
		mStr, err := next()
		if err != nil {
			return nil, fmt.Errorf("case %d: missing m", i+1)
		}
		m, err := strconv.Atoi(mStr)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad m: %v", i+1, err)
		}
		b := make([]int, m)
		for j := 0; j < m; j++ {
			vStr, err := next()
			if err != nil {
				return nil, fmt.Errorf("case %d: missing b[%d]", i+1, j)
			}
			v, err := strconv.Atoi(vStr)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad b[%d]: %v", i+1, j, err)
			}
			b[j] = v
		}
		cases = append(cases, testCase{s: s, m: m, b: b})
	}
	return cases, nil
}

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, tc := range cases {
		input := &strings.Builder{}
		fmt.Fprintf(input, "1\n%s\n%d\n", tc.s, tc.m)
		for j, v := range tc.b {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		expected := solve(tc.s, append([]int(nil), tc.b...))
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input.String(), expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
