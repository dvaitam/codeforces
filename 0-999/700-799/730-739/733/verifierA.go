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
const testcasesRaw = `100
NLTHUGXBPQBT
LQRCFO
SHOTHNLCXVJ
FPENLWQTT
HRMUVA
WAVASL
PEWBGXI
GLHJITBOWDPZDEMD
K
NZZJKGUAZHRWWMLISH
PAOLJL
AFKOXFK
EJRLUFXRZAFRYEVCB
OQUGMIVJWCKBTEPTYP
DBZUNXKNFM
NJCEENQN
FNN
MGSYVQABFRQUC
PETSUMHJCJXGWR
FARWLGS
ZLFJXZQDGNOEXOS
JOCQEDPG
YLSIUDZSKHIPEXPL
WUUKXFXDTSDGYSWAAV
CIXSKLLLOPITDW
FEEZFPGPKMUEYNPQ
CYLXMVQAC
YLNCNTGKJGD
RHVZALFYGAMYHHQKYXHI
MAESXASNVSKA
HKZRAJHYHLXNPVR
CFULLSGZHANSJBSF
SLTF
KH
TXYHNE
GVHGGRLYI
JOTFWGQAWICHY
NF
XQZF
DBCHSESPPU
IBDKC
UYGKZQQHTRBGSPIXQTE
XMM
VBQQNFRV
QLRYOMQEQTJMMWG
JZVYFFVWVXIDJSODBJP
SRNFSCOGK
CJHEOZACCLMWDPVDTQ
BCNVFWSCAFHLEORRRMNN
TLSFSUZMPFNWLUPOK
VQTHRKLXISIXDIPSRYGV
UHXVW
WWRZQKKNZJQI
N
KU
PGBSS
QIGDGFOLHNROSNLLDFK
RSGIFBCMJAFB
AZPQQMDVXVX
H
GQSGCUANSNSFCGWARKT
XAQRVRJQJA
ULNDFXRPNN
YY
ABB
EXGMXJEPTE
NWMVCWRIBZERSJVUEDRT
URQRER
QRHZZGWZG
UQIFSFNW
ILMERV
WYBERHPZLTPWZLF
ZDIHVHVWM
EE
BGA
M
IOKOFRB
MYXILTMNYIEBKVCAOLIA
DYXYKYAEYJPONZULA
FU
PV
VHGKQBXVWIPOMFUH
HXX
L
UQXAZPUHAASZVCYISV
HOQADLITQA
QEDZJQMMULMC
XDDZNXXKEMZSFNOCXBRJ
CLXVT
NTDWNGFCVNNBHNCSHRR
ICDQ
DSTDPQFSVXKZAJN
OIWFPMZSNWUYVK
YUAKCNQWFHGZMBBWTX
VLYIZFCVZSKKYBOMCBB
INVCHRORDERETWRGQGM
HVTNVHTADIPPBW
FTPYDQNWDFQDZR
THLBZROJKTWOUGZK
RABLIQGZHFYQEJX`

func solveCase(s string) string {
	vowels := map[byte]bool{'A': true, 'E': true, 'I': true, 'O': true, 'U': true, 'Y': true}
	last := 0
	ans := 0
	for i := 0; i < len(s); i++ {
		if vowels[s[i]] {
			gap := i + 1 - last
			if gap > ans {
				ans = gap
			}
			last = i + 1
		}
	}
	if len(s)+1-last > ans {
		ans = len(s) + 1 - last
	}
	return fmt.Sprintf("%d\n", ans)
}

func runCase(exe string, input string, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func parseTestcases() ([]string, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	if len(lines)-1 < t {
		return nil, fmt.Errorf("expected %d cases, found %d", t, len(lines)-1)
	}
	cases := make([]string, 0, t)
	for i := 0; i < t; i++ {
		cases = append(cases, strings.TrimSpace(lines[i+1]))
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to load embedded testcases:", err)
		os.Exit(1)
	}
	for i, s := range cases {
		input := fmt.Sprintf("%s\n", s)
		exp := solveCase(s)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
