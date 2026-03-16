package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "981A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesRaw = `ynbiqpmzjplsgqejeydtzirwz
ejdxcvkprdlnktugrpoqibzracxmwzvuatpkhxk
cgshhzezrocckqpdjrjwdrkrgztrsjoctzmkshjfgfbtvi
ccvyeebcwrvmwqiqzhgvsnsiopvuwzl
ktdpsu
ghaxidwhlzfknbdzewhbsu
tvcadugtsdmcldbtagfwdpgxzbvarntdich
ujlnf
qobt
wmgilxp
fwvgybzvffkqidtovfapvnsqjulmvierwaoxc
xbriehypltjvlsutewjmxn
catgwkfhhuomwvsnbmwsnyvwbfociwfoqprtyabpkj
bzzngrucxeamvnkagawyavqtdgdtug
iwfdpmucaiozzdieuquu
deiabbgvirklsbxwtupwuou
lrfgmsjaeeikkzlwckytbbifesjl
rejdpxhbjfqxcjmkjnddrppkzz
dpdwpnbjkxvefusmzucczc
xhbmadmrqjopz
wvgnclhisyfngldcwaqoyvgdpmigubzgtedgom
redtpesmuvnqpvkppuvgrtha
wxkkbqeitzemsjwwzpczcq
che
jay
kfzeuolqmqqbscvzzqytcxnygjrtn
zmtshzavaxfjqsikcpijynmzmbfuehj
kbbpneptwcvwezlnbtomobdpyeabtteukdwrulgmzyypdbt
otukudvwtjzemzjxvzdqzgbzmolygolzucbbpiaqvssgh
yuyqwq
qjdensncdncdnyexazonvnapkxi
lcdlwa
lfahlctegagvvxdxajlwath
fodplwieag
kpjjrukfscdrsjfmeezhkqhh
fjlnvbetamcwcenjrnxesnjulchouluqbmnanxkogljpcfzdi
rtwezwom
ynnfhokqelou
cpygjawotoagjdyujrtenwypcvpyhrymiuadivbai
qswmodxiljyvgtcbczijrkdqhy
cnjjqesqugr
nurmxyz
jolsuefdwdmmsoervj
upxngppwqkpubojexpbtgal
maqcvcvxvmalbdtaiuwjxheysjgdnow
mfknuvneoweqkegfolzmnz
mxhzgogswbmbhucflbxuvfhtjtcwqyjy
nobuwqvurxnsopiwpgkibbb
lajuaecznvh
mrhogktdtczkkrokiaqbglcgqlggivxxjjqmipl
hbjrcaopxobznpoodcchdyengotcnrymbfhpheilkndrjt
zgwjyoqtoruiihadtzwdfxnhgjxvaxrqnbd
uidxslhvwwrvjhxhcqjvkhlup
sfezarqklsuazemefqc
ygzypsywg
xehymltseupdtaqt
pojahrufvpzxprkwcietmwg
zjmbgbkxxhkovxvvhilvf
alsrbxuelapubahbahuk
blvne
oneljfukxzxnma
irrzxvwoybsdnmfaqetvqxweckh
hazfxzvrfwcn
dtuowettbikzxxmaubpcljveohqlfxymkizpmaj
jrpbyrsrivbomxdmlpbaixbivvswyjvygy
qkmigdskzhsvxrvlfekxasbselljujkpzm
nfazesboekaxpvvyixtgcrnifqfcvufsdquzrtm
nijjanywiirqrkkgwznzeayqezvwzsmlobrnutzzyhalqfvgu
uwpaxxhshifyncsohxwoqzz
dgfocnumiinyyltkcjapayigymmnyuuvmwbsolseswikam
qebcsllacgwdvrpbkakmeyuinvetemj
bfeepwuwyxbxqbrwxwumflzsccrfigzik
wiziqoeyorebusfuqbykcguzothzoqfwkuepyrbrcqkay
dntlso
mqludekafxeaktgbnubwjz
btyzwflcnbolttyivjszonfao
gmclddalafmtwuako
xwppcbrmzziauqdckldpbequzjbamkfrweff
fuvhtkapuvmbhhujkfhlhfnolsemsyafsavmwzfeaakqabbyds
eyevzmzannsvkwhelqgrmcensvldnnhpmhmhupms
iiqlr
tt
phibtkzmudrbewmaynxmndwotoffkpnfsjyqdlleulypuqybgi
xskjmubjrnbv
imxgleedtlfansmocuwwvcynrxr
fegfhaqep
txjwzkvdznyifkuqkrwemxx
jhmlmpqjnndveeastuwqdwyuwgtutqdiwxtf
czbadlxpkdvolswivpzyhfrrc
qfaufqntvgowvmiasemfosbmzcusmkhqobptdiqprumpifhrlf
teoccpmsnrciphdjelzd
vbvesrgab
rwytpwwdprlykdvahhpjihaplq
ccjsnhxlmyehjgxypvzljm
eydmlqphuwlulnilmyywjdpjdoelhxfkphdvmmoq
sthvmqjphkqvacpkmhnbsybncigxkf
fzwlahba
raedttgcxogaqtncrfhhnmpza
gmybtiaslwlkvouetqcidwdiawet
vemgsvkgnqqdrdwpdqopforkenimcsqkhohlpynaoxarmyyohn`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
