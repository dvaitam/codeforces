package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `BIQPMZJ PLSGQEJ PPJGJZBQQELMSI
KTUGRP OQIBZRA ARGBUITPQKZRO
ZROC CKQ RZCKCO
RGZTRSJOC TZMKSH HKJZCSMGORTRZST
WQIQZHGVS NSIOPVU OZZGUSWSQVPVQNHII
KNBDZE WHB BDZWBEHKNR
BTAGFW DP TDAGFPWB
HCUJL NF FCUNLHJ
LXPS FWVGY SXWLVFGYP
PVN S PVUSN
WAO XCKXBRIEH AWXEXBROKCI
G WKFHHUOMWV MHHKOFWGWU
PRTYABPK JOBZZNGRU KPGTOABURJZRYEBNZP
MU CAIOZZDI IIAMCZZUOD
RKLS BXWTU RUSXKBLWT
E IKK KIKEK
SJL MRE LERSJM
JMKJNDDRP PK NJPPDKJDKKRM
ZC GX CGZX
QJOPZSW VGNCLHISY HZYOSNCPVGQLIJSW
EDGO MLREDTPESM GOEMDRTSELPED
BQEITZ EMSJWW IZTEJEMWBSQ
ZEUOLQ MQQ ZLMUQOQQE
TSHZAVAX FJQSIKC FIZXJSVHTQSACKA
P N NP
NBT OMOBDP OBMBNDOT
BTWOTUKU DV UOTTBDUKW
YGOLZUCB BPIAQV UOPBLHYZVCABIQG
DNYEXAZ ON NEXZDYOA
LLFAHL C HAFLLL
WATHE FODPLW AOWEDTLPFH
SJ FMEEZHKQH HMQHJSKZEEF
NJ RNX XJNNRJ
BMNANX KOGLJPCFZ CNBZXAFJKONPLGM
ELOUUC PYGJAWOTO JTWYEUCPOGOAOLU
DIVBA I VIDIBA
YVGTCB CZIJR CYCZBRJVGTI
NURMXYZIJ OL NYRJIXMZOU
XNGPPW QKPUBOJE ONGJPPBQKWXPU
B DTAIUW UABTIWD
WKMFKNU VNEOWEQK KQWFNKWEEOUNVMK
HUCFLBX U BFUXHLC
UWQVURXN S XSVURQWNU
AJU AECZNV CVEAUJNAZ
ROKIAQ BGLCGQ LQROGQKBCGAI
AOPXOBZNP OO ONBAZPOPOXO
H PHE HHP
RJTRZGW JY GRWYRJTZJ
HGJ XVAXRQN HXXRJQAGV
CQJV KHLU VHLUKCQJ
EFQ CEYGZYP EYZFYEPGGQC
AQ TLPOJAHRUF RAJLPUTQOHA
MBGBKX XHKOV BBXHIKOMKGVX
PUBAHB A AAUHPBB
ONE LJFU ENOJLFU
RZXVW OYBSDNMFA WNBRVDYVOZXFMAS
TD TUOWETT TDOUTTTEW
OHQLF XYM XHYFQOLM
S RIVBOMXDM RMVXIDOBMS
KMIGDSKZH SVXRVLFEK VGFHRSKISKZSLVKEXDM
XPVVYI X IXYPVV
SD QUZ SQDU
NYWII R NIIYW
YQE Z SYQEZ
UTZZYHALQ FVGULUW UGUUQVHLWTAZYZFL
NUMIINYY LT YNTIUMLYNI
SOLSESW I SLIOEWSS
LACGWDVRPB KAKMEY DKRMGAAKLPWVEBC
FLZSCCRFI GZIKWWI OSCLFIRKIWGZFZIWC
ZOQFWKUEPY RBRC RZPYWFQEKUBOR
AFX EAKTGB XGKFBETAA
OLTTYIV J IYJVLNTOT
AL AF FALA
PCBRMZZIA UQDCKLDP RMKIDCCDQABZZLPU
HHUJKFH L FHLHKUJH
AVM WZFEAAKQAB FBAZWAAEQMAKV
LQGR MCE EMQRLGC
HMHUPMSC IIQLRAT HIQACLRMSMIHPU
DWOTOFF KPNFSJY NSJFOWLTDPYKFFO
R NBVNI VNRIBN
FANSMOCUWW VCYNRX ACYCWVMOUNWFRNS
FKUQKRW EMXXR XUMKFXQKERRW
UWQDWYUWGT UTQDIWXTFM GTQFQDUXWTIWUMYWTDUW
Q NTV TQNV
MFOSBMZCUS MKH OMMZSMUKBFSC
LFJT EOCCPMSNR MNPRSLCOJEFTC
GABMRWYTPW WDPRLYKDV RKYPLRBWADDWVGMPWYT
HJGXYPV ZLJ HPZLXYJVJG
ILMYYW JDPJDOE JDLYPJIEDWMOY
HVMQJPHKQV ACPKMHNBSY MKNBJSKVQHHCPYVPAHQM
GCXOGAQTNC RFHHNMPZAN MOGQNZCNAXHTRAHPFGNC
MGSVKGNQQD RDW KNMQDVGGRSDQ
HLPY NAOXARMY HRAPMOLAXNY
QVCGJQUVDB EKB KJQEQVDBGUVCB
XKFQRWTSN EUWT TNWFRTKXSWQEU
B E EB
NU UJVKPVGFL UGVNJUKFPV
MO WRG AWGMOR
UAVEQUPJM RRCC CMVCJYRAQURPEU
IDS OIHMMXQK HMMOKSXIDIQ`

type testCase struct {
	guest string
	host  string
	pile  string
}

func solveCase(tc testCase) string {
	if len(tc.pile) != len(tc.guest)+len(tc.host) {
		return "NO"
	}
	var cnt [26]int
	for _, c := range tc.guest {
		cnt[c-'A']++
	}
	for _, c := range tc.host {
		cnt[c-'A']++
	}
	for _, c := range tc.pile {
		cnt[c-'A']--
	}
	for _, v := range cnt {
		if v != 0 {
			return "NO"
		}
	}
	return "YES"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 fields got %d", idx+1, len(parts))
		}
		cases = append(cases, testCase{
			guest: parts[0],
			host:  parts[1],
			pile:  parts[2],
		})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
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
		input := fmt.Sprintf("%s\n%s\n%s\n", tc.guest, tc.host, tc.pile)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := solveCase(tc)
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
