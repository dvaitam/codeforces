package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// gcd and solve mirror 1367E.go.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(n, k int, s string) int {
	freq := [26]int{}
	for _, ch := range s {
		freq[ch-'a']++
	}
	best := 0
	for length := 1; length <= n; length++ {
		g := gcd(length, k)
		cycleLen := length / g
		cycles := 0
		for _, cnt := range freq {
			cycles += cnt / cycleLen
		}
		if cycles >= g && length > best {
			best = length
		}
	}
	return best
}

type testCase struct {
	n, k int
	s    string
}

const testcaseData = `
100
16 20
dxmpeccamrjzybhq
35 24
iyfdigauzizigfjjuxlctkvmqhfhpicrjaj
37 46
jyqgnntjnofhjizbcbouiqrupwkevgcnguuoi
12 23
nxskurgkdbwh
18 49
sthdkfjoablwcjxvka
21 19
keyuntvcjtgojeimtfksa
24 3
oflzzljsdogngdbbbxftvetb
35 32
shkbdqjynugpghonpbhnohungpgbbiihqgy
15 27
iekbksdsmuuzwxb
32 25
cngsfkjvpzuknqguzvzikmpciuvgbmte
50 18
vbfwuospxmmgzagfatidmzyzmhrbgfvtkryzpqoacbwtdprity
9 3
lcyqajlcc
35 30
mgzjmhypmdcdtzlqnnwyocugujpndzrflff
46 10
kpkirawfaujdrdpzwtpqcqhnjlhyfuavbtkrozsjqzotto
26 10
iytlvkenctevztfjlgszlvtccm
42 12
kulkfjataqyczlzdffspvscydfupvxhtzvjwmthpwh
20 24
hxkruqozmqmkjonsaifr
30 45
rtlmmtaeqcoultzjcyiphuptycehcj
9 4
fmzsvtrwr
17 2
qhfydgbkcdibujtuw
12 47
eunewxcrlwax
36 10
qnegjpqcmffiqmrvjmkfmbnaivajzecfdtah
15 36
apwrfoxwvwmkfqs
13 7
ptljuntxbugis
43 38
jzpukzchkdvwbtkqplcfwbpqrtwhbgvckztewjdqqac
15 40
itaxbapweglhxxl
20 25
utmbfnpcyqqjrwjzgwwn
21 12
gapnikynunkytgiiqcami
39 19
vuetykczftcwgivgupydukaosgfchqddzkjssve
48 5
vfwxayyzfvgibskukxytlffygbbfvbviwvvpdgzemptrprkl
28 48
mbybfelhjhlsgtvlmlrwdhglfrga
2 8
gj
45 47
krpbjzxzmfemtollaoweoaxzxwngdkjtrbcpfvqbfthtk
39 2
xvamnhzsetbeftjhnxkrulfxlvrvslitkrbrels
13 32
pdcutgferwetb
35 23
cwbzterfitbjpigupyjqwybejnbalqtlcgn
44 6
znzbjgvgfgwxccaajqqqnccqukvfugqwrbwbcenedyqk
33 22
viessdsutctafbyirunyjslyeqlrohtof
31 42
srasrsgisvsxtwokfjfkjkwgrptjilv
11 21
dhzkugfodjb
34 1
bauvrfhhlummhuoxnjqckzbdfbnamfnvxb
36 15
lkqlvealldymsgplsnrxcmfcomqgyacjdeeo
29 48
lynkmvjzmxdzkmqcwpnwlhstfnjdk
35 3
qyuvvhtfyphqlnonrvdpawnylpxxigeklei
35 46
svpytzepmxpltqyzbxzvnistwfurfaovsje
3 1
hyk
10 21
koshizbvcb
29 8
xxuhhpzrcrtffgagqzghaqpklisms
9 46
ykvuwdtii
7 4
kwbymcg
46 16
bpxmlcyuwcexowxnqmwdfebfuorqpolsnhuiisaimwskwj
46 8
fqygfwymvreuftzgrcrxczjacqovduqaiigyygdzismiqi
1 7
m
24 39
zfxptxkwxgphfyfkxbggxroq
21 9
tflsqnxolkwyyczopiied
46 43
ixotqupwocilxrdtjtafmqfvmpfbkhuhmrfgpwvojmikng
45 42
bxvsbyrgxbqpkiyyybelbwxmnzlhlolmsamxeepniyfwx
29 34
llzsdzyfejtqvyecpvwrbvwitxgya
28 27
ipfadenzziroycukirxwofakiiim
22 4
stfnisdhsddptinxdmnyoq
41 29
sfvpqfjlivucovqtuxqwgjvwwrzmiroyqkbkztolr
23 10
hhtltcinojunqfswmhmjzkx
9 17
oggpdupsf
5 17
taekm
22 24
xectqbolfxlxlclqrdrovo
33 39
wpfexuknzerjouzwusffhxogunerajfrz
16 39
zkeqvhkwnxdblqjf
29 7
xjxguftcqxvkcaeetcfsifbvkkuio
31 10
xdqkjtivujbqarovizaqlcwiutbxryz
36 29
zvcpnoqprxpksbqgvvsfmsjtlpwbiengrnld
14 47
niwhsxsodfdvxw
43 44
zxcxpemsdhhsrfxzffcieywigjczswzzcnmojqmgcmo
30 41
yzqpwijxlgigphxrgpknelphzihpar
10 39
uhqqbtulov
50 26
qsawtjhyecgngffdpzyjebzeyqglwusmizakmgmjenfjddhdim
26 49
xigfkxcisxuqtxajouolqyhxmx
44 19
vnpkcnmxtuuzvayeaprlkettnhjlprgzhmcelgrhxide
18 39
nxpfduaddqdzncvxeg
20 28
dwzelvfwwsoaspsjnhpf
39 45
mupqaxhziiytpsfdbzcwkhzqszyqfnyrvepouvz
27 34
jdkepvyvoewfgdmhmyxiiydblou
1 4
z
22 50
ppvsznfsxfyznevlyqqvij
39 44
tpqnymuxvnzfvfdgoelsvwlbqoaqwlqnxygyvtm
34 20
tezamqeedswkjbzpubqjpqjnrogdwafhks
14 9
pqzcjijgtbusrg
37 36
hqbwbhkhsbvfytkesfmkvjbpoiwfpiluaghqf
22 5
tzhaycbzuxyexojlwfursz
9 32
bnzgxvuuk
17 21
zbveywmdyyarfkggc
13 37
bpytakminpuan
25 47
qbjgyeeumpcqntowoenqhgicl
44 5
ofyqqedllkqfvxjuviiisjasvcvqgshnhxvesiwenfcj
9 12
sedzmbsgy
6 44
gamkei
34 24
sdibezwkcstpbwsvriwcxtpinffhslpymp
8 38
hntmmbfn
31 35
plfrefpeizaycewveqxjzpjeqoksflt
37 49
boqfsgvkkfwxxesqhtkasctzwlxukrstarhzk
22 40
nqmcrfxfbxfutknaqilrej
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	idx := 0
	next := func() (string, error) {
		if idx >= len(fields) {
			return "", fmt.Errorf("unexpected end of data")
		}
		val := fields[idx]
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
		nStr, err := next()
		if err != nil {
			return nil, fmt.Errorf("case %d: missing n", i+1)
		}
		kStr, err := next()
		if err != nil {
			return nil, fmt.Errorf("case %d: missing k", i+1)
		}
		s, err := next()
		if err != nil {
			return nil, fmt.Errorf("case %d: missing string", i+1)
		}
		n, err := strconv.Atoi(nStr)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %v", i+1, err)
		}
		k, err := strconv.Atoi(kStr)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad k: %v", i+1, err)
		}
		cases = append(cases, testCase{n: n, k: k, s: s})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, tc := range cases {
		input := fmt.Sprintf("1\n%d %d\n%s\n", tc.n, tc.k, tc.s)
		expected := solve(tc.n, tc.k, tc.s)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(expected) {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %d\ngot: %s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
