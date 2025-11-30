package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const modE = 51123987

const testcaseData = `100
10 jfbisldnik
23 trquidrrzbakxntpzkvszhd
14 nvcsiyrlxzebae
13 lzalzkqszsmrw
4 dtrg
7 kurmebc
30 bulgnswudngkcdxyzeihvxfgmnggaq
9 nwzgqdzyt
12 ukiqhpsgsbsv
4 inud
26 bdkkdtjdhzjzscdmpekhehebws
2 op
18 rhnhgtdkfwjdoogupi
26 pgzovftgsmlnimfextrpdzhtlv
28 siajogmhjskalbnxyaszurzaxbbn
20 jokuarmdnizxbzkbbluz
22 yvdtwxauqudsyrbgssnidd
6 derhhp
30 timgbqqpchemwvochhydenvpnitzgr
2 rw
22 mewpxszebjmspvotxznkyu
28 rhcjrqhkxjhqnijzmkxkgdjjgdhn
18 fkgnzohiqlmcspovdf
21 dzqziiammgatfzbnnknzx
13 bipvssacbetrh
6 xkneqj
30 yphucietgheawaoijxqamoztyufbfu
28 aqmlckavoxqhimbnopfgudrqdqgu
25 wxrmkihxtbkvicwztszfdoyap
12 infbwqffnxww
30 yxgjtrsbbhasyvylziozxtjiqklecu
15 rdxdefkngrixlul
9 qxaqzcvnv
30 dswnozcxdxouemxfnrblnijecusork
30 ztgrvzvvbhoxxuqnerrbxgjrgibqlw
24 jmebgtjgojabnnrwawfvdlnh
28 yiymqlsexckgvbmymikgqwyhsssb
26 pzhgiffohbchmmnuyaetjwzjub
19 gjpmdtvzwlcvolhtsbz
4 lvog
10 jxdyacvtng
16 gprvxmqpfeaqocnu
29 wkxfcsytkotgwpvznbfijrjvnysau
16 cpjxovrlptcujgyn
15 loueywplciusgse
22 qjucfxarsxyvwrozrpifau
29 nljnbwbbhnyucbmptgdiafptbcaie
15 yfghuhrvabzopke
21 cbvtbbuhmyyizapizdrec
12 xkuqugrvtztz
7 rqjqhff
12 cdjnnqiqedlm
4 fbnr
1 s
9 vxsxkhyab
30 bahzkduxuwcjiobydnhiiystxvxder
4 llzi
30 lbcnjzliaepkhgcnlrvcffuthrcjwq
25 ymbulgossgoqdogltkkwwttza
27 nyuxdkmnlhelmmwrjwwvgzgczmg
6 gbkzmu
3 dzt
18 tbleflonkjyjknxfov
5 voenr
21 lkdawewlkvgdwjltbdwou
18 uthaopchcgguztglla
10 wivaktblfn
3 xvn
30 rwevoprvukalacopwahqbvdmdnepss
30 jsmarbsgftfofogaqybqppmiuyoktt
18 uudyfsvuxuboppinhy
26 bbxxkmphtsarugsykgwezjrrwf
14 tdwimgctlmeyzv
28 nzrpdinlzlymlynviagolgckamng
19 nhuuxwuvnkdagujcgjc
7 fhswynu
11 jfhuqpayepu
20 cajoxigqreuvkjefykmw
10 lopnftogiy
14 unpbkzpkubkplu
21 occxrjujtnkuuodlvounr
2 mr
12 ualyaspqtnhv
23 bjqzrgtrgvdqmdrbikkbovx
25 eqsatuzttnipitxrvhjfhgpxh
8 yilgswke
3 zmg
25 gttgmfpwjtcxdywjgablcojjp
13 ihvrlqtdukmld
18 yetptsmqlzirevezce
3 ibw
29 uhuqlxmrlxdwxzqwsxpaygrznanvn
19 rrhdcknrjugpofuuvss
10 qrpzldxaky
12 oeribcofzpuj
21 qefscvhqvwfobdzqlyknl
8 dnzxpxts
17 cstwsojakrpdvvzrx
10 dtzlbysyfz
10 vfkphmkuwr`

type testCase struct {
	input    string
	expected string
}

func minE(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveCaseE(s string) int64 {
	n := len(s)
	d1 := make([]int, n)
	l, r := 0, -1
	for i := 0; i < n; i++ {
		k := 1
		if i <= r {
			k = minE(d1[l+r-i], r-i+1)
		}
		for i-k >= 0 && i+k < n && s[i-k] == s[i+k] {
			k++
		}
		d1[i] = k
		if i+k-1 > r {
			l = i - k + 1
			r = i + k - 1
		}
	}
	d2 := make([]int, n)
	l, r = 0, -1
	for i := 0; i < n; i++ {
		k := 0
		if i <= r {
			k = minE(d2[l+r-i+1], r-i+1)
		}
		for i-k-1 >= 0 && i+k < n && s[i-k-1] == s[i+k] {
			k++
		}
		d2[i] = k
		if i+k-1 > r {
			l = i - k
			r = i + k - 1
		}
	}
	ds := make([]int64, n+2)
	de := make([]int64, n+2)
	var total int64
	for i := 0; i < n; i++ {
		if d1[i] > 0 {
			total += int64(d1[i])
			L := i - (d1[i] - 1)
			R := i
			ds[L]++
			ds[R+1]--
			de[i]++
			de[i+(d1[i]-1)+1]--
		}
		if d2[i] > 0 {
			total += int64(d2[i])
			L := i - d2[i]
			R := i - 1
			if L <= R {
				ds[L]++
				ds[R+1]--
			}
			de[i]++
			de[i+d2[i]]--
		}
	}
	cntStart := make([]int64, n)
	cntEnd := make([]int64, n)
	var cur int64
	for i := 0; i < n; i++ {
		cur += ds[i]
		cntStart[i] = cur
	}
	cur = 0
	for i := 0; i < n; i++ {
		cur += de[i]
		cntEnd[i] = cur
	}
	prefEnd := make([]int64, n)
	var cum int64
	for i := 0; i < n; i++ {
		cum = (cum + cntEnd[i]) % modE
		prefEnd[i] = cum
	}
	var disjoint int64
	for i := 0; i < n; i++ {
		cs := cntStart[i] % modE
		if cs != 0 && i > 0 {
			disjoint = (disjoint + cs*prefEnd[i-1]) % modE
		}
	}
	totalMod := total % modE
	inv2 := (modE + 1) / 2
	t := totalMod * ((totalMod - 1 + modE) % modE) % modE
	totalPairs := t * int64(inv2) % modE
	ans := (totalPairs - disjoint + modE) % modE
	return ans
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing string", caseNum+1)
		}
		s := fields[pos]
		pos++
		if len(s) != n {
			return nil, fmt.Errorf("case %d: length mismatch", caseNum+1)
		}
		input := fmt.Sprintf("%d\n%s\n", n, s)
		cases = append(cases, testCase{
			input:    input,
			expected: strconv.FormatInt(solveCaseE(s), 10),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
