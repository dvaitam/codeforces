package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(s string) string {
	b := []byte(s)
	n := len(b)
	for i := 1; i < n; i++ {
		if b[i] == b[i-1] {
			for c := byte('a'); c <= 'z'; c++ {
				if c != b[i] && (i+1 >= n || c != b[i+1]) {
					b[i] = c
					break
				}
			}
		}
	}
	return string(b)
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

const testcasesCRaw = `100
mltvvlmvgobqjmxrpdqkapgjpujaikcrq
kvxlggclnebauyqphvribcaoyvtikxdvjsaybkbtyciocbppn
cmzmvpyuxanqgwlre
cqzaljzuejdqkpqmnspkrqgmqylmmvqbf
xopzutgmgtpsxinabcugteabgrgnmojesvisklpahmoltd
rwyoxzwforsfvxf
ipswbbl
jsoqyurjaejiiieb
golqqfnekaq
jbmlhxojxokrstbeqkogefe
vwlteggkgijxdxmzmxjuiqqfktj
zp
mgczfpwpgaruxxdggk
qiknkawswhhtrdawpzildmqw
jwujegotcokcrhcbxixpprougudaucykglgcesvyzkrp
lfjmhvyelqlpnqnftznpqkpswprtsxvdfql
wvsljkefpixpnrdlfeewikjpqrtjskbbaojostwtgtq
jh
whqfotcnxcjxtlftxkxapheldxm
fwzgafxwypitcqskszlgajoqkfnbckgrnrhuwblm
paosjpnewgdwcymghmvmnfzzevqmatorxpoixi
yohwqfxzbjusyseyghgfdpehyscwo
ewetjnkxjxyrcuiapypeibqelyw
ojshgswbcvhvorroueskuupdgjmrtxdygbrrm
crdexyuccqmuxawfttyzpcofsybgpwhvamvwndqaosbbbtcaq
xcmsgloxazxfykd
wwjiawttxqqvulzluscu
ekpgtmwrlkctcghccqfpcnnsjmrarbloo
ktfwvnzjhxdobofegsxaovxwmzqxeczd
xohbtezkqxlerxgpvatqjptxyxumkhlfr
oiytrkmrsbaoujuqjiuukssgbvexfpurwhdi
jiwmhqgyubwxabcsdgocagawxeidnkkndtovtsfib
leobcvlkmv
pjezobrcjecuelblvvdmrmjydydglxokwa
rlmqhzqrlakgxnqfdvswwbgrynosffwpibcdhoxibzetcl
fsggcyhsxoiepcsitynakpdfxhulirzzetyuhiu
gnrfsnokwohxsbewyvqlnuphajsbj
pvbhxvm
zlzonnanmqgjbtomokkmjqufxylhjkgiwajdpbqjfvi
npwhxztqeyvqgrpuabpsrbronhx
iwgsvyyenbldemhqjufhbpgmaqibsf
aa
nmjgzkjctfujwmc
okfueoageh
mwkvrtdryvfb
hhizcehaeonpninawqstqsmwrh
lyzlryusui
pqxfzzlzsiozxaqudehtalepwprquczmrzfn
ileflmkqfzjqvxaoemsvvsedowmexfeypemmggnauhpdwnoc
fenyshbzbartuqriocyrtkpwepd
clxumdpxuhnutyapq
zgahmpyntwighvtjekhcddrunoitkpdvqij
necnuuwhtovlitgvcjxkgvgqyogbcraestgzn
zejwjlxozppdhrblqdgsfnis
daroxoarguvepowifbqpcthgri
oetheraiuxcdxvonvn
udvcovxpwnaosjhdmvekblvoldlvxyximayjzkcj
mkjlatnthokaoziaucp
zoikpjlohpawssaguzcrjjkkbu
powfcnbcsrhuvafss
tfuwzmjkprjhk
rybjgjejyuhgvvvltjowlpkhzzltvdnaz
y
ppoqflbvschfednvb
chsgmrhbkfmukplxukdgrvzorfffk
khfbsmjkurwknfikqyeoxrhyfrkasupet
cwivrcaxiygdgygtzbikjbrlcz
vdkktmozx
fipgrpbcjcbxznimmbsotnkxa
qyqzrugcg
nzidmzesinbuowkodskefbtmxxofvxia
qipjlrdylawwoncwwpawyakhomsftieok
xzeyeppilzzbumboyjwjnbno
ygkanfcpmycwnnnupyyaczzyv
aatnizvzybtbdwjrzoukcgihguvxsapicqiyrae
hrsksljiblrj
jbofyxsnrck
egitbwzjxbygceb
xnjmagyklmzllbgcstudtwjprmnzz
cyogwemupy
wnkfvvsmjfgsywbouyekdtjyprffevzksz
wzaytkfgpwqzbpbskgoryabbxwdcoyfuxtznpok
vqamwaucfgdizjsgxhkhbzisnqeqpylyehiackqw
hbeipuxnklwwfaczxlebmjwpxrafffgivqpul
putzgffsvnuekuxiyczc
lrzbsuzdkbqghmkdgimfingoupiivjv
jdvzynfixeqrbvzvavmjvltykqhitjvl
wacjsnmchjuffiacacpjmkrcqnornqnedndppxrwuobsoss
ozsgklyxmhlslkkhmdvxkuqfmshaiympixjennho
vihdbwiihzaopwaxnvfeobrvikqfcvtevpivbhns
iqxrdiyhqlxdspnxidamo
sfaoreqhjap
vqnllnxyegvujmwmmyqbshyvgcjprgbkjpnkguzcwxa
wcvdyy
phsnbtxsntcuveuezjfugltyys
xarxxuffopyktsbg
hmdtlpdprqeybxxnifjcgnihasfi
al
qwpgybgywnpwhmjclrbupkavpvudyhc
ifdvovlmfx`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data := []byte(testcasesCRaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t := 0
	fmt.Sscanf(scan.Text(), "%d", &t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s := scan.Text()
		input := s + "\n"
		exp := expected(s) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
