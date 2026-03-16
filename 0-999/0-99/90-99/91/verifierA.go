package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesARaw = `nbiqpmz jplsgqejeydtzirwztejdxcvk
dlnktugr poqibzracxmwzvuatp
xkwcgs hhzezroc
qp djrjwdrkrgz
sjoctzmksh jfgfbtvipccvyeebcw
mwqiqzhgv snsiopvuwzlcktdpsukgha
whlzf knbd
whb surtvcadugtsdmcldbtagfwdpgxz
v arntdichcujlnfbqobtdwmgilxpsfw
bzvf fkqidtovfapvnsqjulmvierwa
ckxbrieh ypltjvlsutewjmxnucatgwkf
uomw vsnbmwsn
o ciwfoq
rtyabpkj obzzngrucxeamvnkagawyavqtdgdtu
jiwf dpmucaiozzdieuquuldeiabbgvir
sbxwtu pwuounlrfgms
eeikk z
ckytbb ifesjlmrejdpxhbjfqxcjmk
ddrpp kzzkdpdwpnbjkx
fus mzucczcgxhbmadmrqjopzswvgnclhi
fngldcwaqo yvgdpmigubzgtedgomlredtpe
uvnqpvkppu vgrthakwxkkbq
itz emsjwwzpczcqbchebjayokfzeuol
qqbscvzzq ytcxnygjrtnpz
shzavax fjqsikcpijynmzmbfueh
kbbpn eptwcvwezlnbtomobdpyeabt
ukdwrulgmz yypdb
otukudvwtj zemzjxvzdqzgbzmolygolzu
bp ia
ssghcyuyq wqnqjdensncdncdnyexazo
napkxic lcdlwallfahlctegagvvxd
l wathefodpl
aglkp jjruk
cdr sjfmeezhkqhhyfjlnvb
amc wcenjrnxesnjulchoulu
mnanxkogl jp
zd idrtwe
mfynnfho kqelouucpygjawotoagjdyujrtenwy
vpyhrymi uad
vbaim qswmodxiljyvgtcbczijrkdqhyfcn
qesqu grdnurmxyz
olsue fdwdmmsoer
upxng ppwqkpubojex
tgalpmaq cv
xv malbdtaiuwjxheysjgdnow
fknuvn eoweqkegfolzm
pmxhzgo gswbmbhucflbxuvfhtjtcwqyjy
obuwqv urxnsopiwpgkib
f la
aeczn vhtmrhogktdtczkkrokia
glcgqlggi vx
qmipl whbjrcaopx
znpoodcc hd
got cnrymbfhpheilk
rjtrzgw jyoq
ruiihadtzw dfxnhgjxvaxrqnb
ui dxslhvwwrvjhx
qjvk hlu
sfezarqk lsuazemefq
yg zypsy
xehy mltseupd
aqtlpojahr ufvpzxprkwcietmwgkzjmbgbkxx
ovxv vhilvfjalsr
u elapubahbahukcblvnegonel
ukxzx nmanir
xvwoybsdn mfaqetvqxweckhfhazfxzvrfwc
tdtuowe ttbikzxxmaubpcljveohqlfxymkizp
jqjrpby r
ivbomxdmlp baixbivvswyjvygyqq
igdskz hsvxrvlfekxas
e lljujkpzmtnfazesboe
axpvvy ixtgcrnifqfcvufsdquzrtmynijjan
rqrkk gwznzeayq
vwz smlobrnutzzyhalqfvguluwpax
hify ncsohxwoqzzwdgfocnu
inyyltk cjapayigy
nyuuvmw bsolseswikamp
bcsllacgw dvrpb
akmeyu invetemjqbfeepwuwyxbxqbrwxwumf
sccrfi gzikwwiziqoeyorebusfuqbykc
zoth zoqfwkuepyrbrcqkaycdn
sokmqludek afxeaktgbnub
mbtyz wflcnbolttyivjszonfaozigmc
dalafm twua
zrxwpp cbrmzziauqdckld
equzjbam kf
effyfuvht kapuvmbhhujkfhlhfnolsem
afsavmwzfe aakqabbydsteyevzmzannsvkw
elqg rmcensvldnnhpmhmhupmsciiqlratt
hibtkzmu drbewmaynxmndwotoffkpnfsjyqdl
ulypuq ybgif
jmubjrnbvn imxgleedtlf
s mocuwwvcynrxre
gfh aqepl
jwzkvdznyi fkuqkrwemxxrjhmlmpqjnndv
ast uwqdw
utqd iwxtfmczbadlxpkdvols
pzyhf rrcyqfaufqntvgowvmiase
osbmzcu smkhqo
t diqprumpifhrlfjt
`

func expected(s1, s2 string) int {
	n := len(s1)
	const K = 26
	next := make([][K]int, n+1)
	for c := 0; c < K; c++ {
		next[n][c] = -1
	}
	for i := n - 1; i >= 0; i-- {
		next[i] = next[i+1]
		next[i][s1[i]-'a'] = i
	}
	res := 1
	pos := 0
	for i := 0; i < len(s2); i++ {
		c := s2[i] - 'a'
		if pos <= n && next[pos][c] != -1 {
			pos = next[pos][c] + 1
		} else {
			if next[0][c] == -1 {
				return -1
			}
			res++
			pos = next[0][c] + 1
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		s1, s2 := parts[0], parts[1]
		exp := expected(s1, s2)
		input := fmt.Sprintf("%s\n%s\n", s1, s2)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errb bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errb
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errb.String())
			os.Exit(1)
		}
		var got int
		outStr := strings.TrimSpace(out.String())
		fmt.Sscan(outStr, &got)
		if got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
