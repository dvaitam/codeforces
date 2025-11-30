package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `
8 axihhexd
35 csnbacghqtargwuwrnhosizayzfwnkiegyk
7 cmdllti
3 xor
8 mcrjutls
13 wcbvhyjchdmio
24 fllgviwvuctufrxhfomiuwrh
21 yybhbzkmicgswkgupmuoe
17 ehxrrixsnsmlheqpc
4 deuf
28 tcmmtoqiravxdvryiyukdjnfoaxx
17 qyfqdujuqtgelyfry
34 atkpadlzjhbhsccxpcyryeevprfiqtngry
13 wjmvuloqodhhc
22 asrhshacwubhcbkcqhivpg
35 exssphzpzngddvnlnnoxbvuudbmxkzdhggr
29 enfiohcozrdburacyhfnppgmbfmam
17 zzojnwxzrvwpegjgb
38 xrbxkbbspqqfbqcfctcvhmdshstbtcnvssqkig
21 himevujokycaotsdcrgqi
9 lchljforw
20 tzuqavrjvdeiddxreijt
14 wkgvuiqpibcuni
3 aky
9 uifxorwnr
1 d
5 werbl
38 renebjlzblgvhvdlyrntxehfzzfnafxkznzvxz
16 ifzwdmbphgoljzhh
2 vg
26 kicyiluqmvrkadifsibdtnlxzk
28 tqdmsgibwnaqzrvxxxvglncvktkv
8 xjqjvnkm
19 regnvmvxftsjmrajjgn
38 tukooovgqpzzxfvcjqvutkcyhvjhzgeabhptyc
30 nusgwwmpmheuwayydynhfzwqobrhdo
9 zovqrtkyo
40 xqnrofxpoiyhuiyyqpuhiocwjhikkrceehmwewgc
27 nkronbgnmyswaysmpaljymnrxxr
39 hphinpamkvvzmxfoetramssvacuneofbimkgokk
25 iynicpaxrblhucyubyahgateh
9 pvdsgowiy
24 fttxwdyfjdsajsvmmwgcswuh
7 wyjvtzd
37 zblrnvlcqukanpdnluowenfxquitzryponxsi
21 hciohyostvmkapkfpglzi
22 itwiraqgchxnpryhwpuwpo
2 cj
15 mwhjvslprqlnxrk
23 woijihdxgkdxrywfggxpixs
34 tjdgjhlfjawreibbrjweuypdasjppokfbi
31 dcmpcsuvbeezsjchdrynttzthyqmooj
38 njstbtxdygugivcfhfrcfanowtpjbhjwjwocvh
17 zzusvzgndrhueiecb
11 zjtxsjodowj
26 iqrpoctbnxktiachvssayvisby
12 pquoifsnupcp
23 nkkvdfknwpjvmyrbockikdy
26 qavronbgqltypuoybgirejowpd
2 ut
16 wfjrarnchdoduepw
19 qwinpphoremgtqxeciy
27 kzqiajxjssvpeorplkryrmokgwh
37 mhynbkxpwzmmvzuepbeqskdodqoaxenuecpzi
22 twmuckvrmkuwyprbtchuvj
15 xcndyuwdofwjabk
4 jlln
10 hqnsvzfffc
40 mtvhpsehouioivazojvrfcolsjunwiojgmpdhmsl
37 jwjavmiasvyxbtxpjyzhtzlhugtivyxyvveud
3 jzo
3 slx
9 cjkxnfgez
35 lqqifipzjxkzdoceyhvxvmzrlczmairdolv
17 smuldvhpatrkthucu
30 wjundebbjpddhremolvxwrnsxxenud
32 tnibwlgoohldvlrulbmigdocvguutabz
22 hezsgcyrgsghkyeztaieer
17 zfdvaealzzhskafib
9 xnqdxcpoy
24 qsdoqhtbxzvqjouabpmnvdpw
29 cckteceitusrwkmtqjoqtndzwduuy
36 xgnohnkomnxdknkvilevpccccndxxlzerbsr
36 kvdnlvynxbjtjldsqgevphdlrldyishznryt
40 vuratvwiafiwyjklafesvmcexuacxqgmnokfljxk
37 tcbefytbvciovnptonigyqdlndjvvspqvjbhm
39 bagjgyeyijkdapxnfemrwhqrvzlcmxbnaocksns
26 wunjdmakfztowlcndhnsmqcmjx
22 hkyfcqudqqgyllxuehdeig
12 teyyucfyupoy
37 ysovsuutkukeocpoujzisblqcjoobbljcuctt
33 mosrzxbozsugktpqebodzkwcqufbhwooq
34 tflljmnykvtbzuukckdrvmjixvtekcsvel
20 uwvmetwcjrmuzkevwxvq
6 uvnqla
24 jfgkypgheecjzdqyrxqbvkyt
9 tmeffwytz
11 xobnlvxhotj
29 hrhjzzpglvsooyjymqqnfgzteibup
24 rdwqdjcyfioqencholanbmql

`

type testCase struct {
	n int
	s string
}

func solveCase(tc testCase) int {
	b := []byte(tc.s)
	sorted := make([]byte, tc.n)
	copy(sorted, b)
	// simple insertion sort to avoid importing sort again
	for i := 1; i < tc.n; i++ {
		j := i
		for j > 0 && sorted[j] < sorted[j-1] {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
			j--
		}
	}
	diff := 0
	for i := 0; i < tc.n; i++ {
		if b[i] != sorted[i] {
			diff++
		}
	}
	return diff
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var n int
		var s string
		if _, err := fmt.Sscan(line, &n, &s); err != nil {
			return nil, fmt.Errorf("line %d: %v", idx+1, err)
		}
		cases = append(cases, testCase{n: n, s: s})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		expected := solveCase(tc)
		input := fmt.Sprintf("1\n%d\n%s\n", tc.n, tc.s)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		vals := strings.Fields(got)
		if len(vals) != 1 {
			fmt.Printf("case %d: expected single integer output, got %q\n", i+1, got)
			os.Exit(1)
		}
		gotVal, err := strconv.Atoi(vals[0])
		if err != nil {
			fmt.Printf("case %d: non-integer output %q\n", i+1, vals[0])
			os.Exit(1)
		}
		if gotVal != expected {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %d\n", i+1, expected, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
