package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func isVowel(c byte) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

func validOutput(input, output string) bool {
	words := strings.Fields(output)
	if strings.Join(words, "") != input {
		return false
	}
	for _, w := range words {
		b := []byte(w)
		for i := 2; i < len(b); i++ {
			if !isVowel(b[i]) && !isVowel(b[i-1]) && !isVowel(b[i-2]) {
				if b[i] != b[i-1] || b[i-1] != b[i-2] {
					return false
				}
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `lfahlctegagvvxdxajlwath
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
vemgsvkgnqqdrdwpdqopforkenimcsqkhohlpynaoxarmyyohn
ippehojlupteqvcg
quvdbekbktzfyomvgnpg
mbzzykyvqgx
fqrwtsneuwtpvsgtobhpt
rxgabbehovthvvdyqnuuj
kpvgflrlwomovkocehydyemowrglwbamgcnssznrguav
qupjmrrcc
xpyagntyljcxayzhtasqleoidsoihmmxqkhsicb
kusspcqpnk
xadntveyeadwnhbuiimdljltz
gxrrbkqrafpauikeldtrnzycisdqcsepbqzhldrwg
wvhmgtwxolfi
lrggokeirvyyacvkjwzkbctajzdfwigqrlocnqrndcsokfj
unvwjbaauicvkkru
bzwhwqeitbzmnishikz
powdnlfqbhwzpmahcjdfeoaylqndqhpjimn
wbfwqlghehzjzefrirafhxnwfhyresv
xhqyjzzwxax
zgbnfmqmgmrqpdecracaxhxmovinlv
ixcyfshxgbybqzcwohowyyvetu
jatelvlljquwieapzxgtmkpjoepjzwklu
mmsmwqjvyshgngiqfqeztqyryqupkinjouqbzxg
owcagtdizagftfvlt
avvaatfjhmpxllcsnz
wayfatri
xgmgkfzfpvetsqnsuhabntb
uueqcnjkuutgwbkekvxipptkiucozyfdqirwmtonvf
nyevvyzdsydfxpjvxdvyniizdulqrzcabtayedn
zpakuecwmusxhbgnabtzmerxw
nmrbxorereugsfcspjii
mvaubqyycaujey
yqedkaxkek`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if !validOutput(line, got) {
			fmt.Printf("test %d failed\nexpected a valid segmentation of %s\n got: %s\n", idx, line, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
