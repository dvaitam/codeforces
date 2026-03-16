package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesARaw = `pxnupcqroattdosfmodtcotpsvzxy
nbkzfyldgxlfpgigif
ocguevegahuqsrnskbboah
lfpxvrkzjuopegnvivwnh
vdyncuwxre
wujrdugghbmbfwxp
jnxqwjbabzgf
vgnkcqdfabaefhqdrnogjrquzokin
vwbcwbwgbcqcio
fwrswnlyowrspjgx
ypbnahxpzjrcytwx
atxcnkjra
qusvnkovyjrtyogqckow
tyrghsvkgxdiapoozgosqyq
bjmfvnmwoxhikcrdbrkkllwzli
cioyv
cgeunmyawjwfwwe
jnmllbhvoizeygkdxdk
ocwmambyeciimidp
hqibafrgwzocsgufuf
mwirxhv
zfahitpihq
rw
wx
mcmismgdbysnvkmuq
hhoiftlreiknwhzdcoltefls
jydampryhwmtsyerbqrizeblsv
rmycdinfstlawrdasyajywzkv
rfzhxsgelvxkwcjycp
fmqqxxcz
ewjsryzqsj
rwdp
rohcsb
oczzwnxqhpuy
xrxinhikendzzjw
buruiwutquiezntksmbgavncyw
lbponcjwu
fzfuxfkcjlsounitzghm
y
sseb
tocpvmdkeaafuh
sdpttjhjuajdortlnlswdvgodglsqg
btzvdxoda
rsb
yvqlylznixtjj
jplo
zxq
ondenflabiunyijtyquiazcn
mhvsjcuwfootwai
guvqjfkds
abbyyay
ngqusebzrweni
nhv
hrkijnzttcnutwbkiuhmn
zwuldlbltipikgkajfzbbwqax
askdszlmfyusfcrjjcem
wmfmbanzlfjazggcxbvgl
hmdtcnjhyq
ejiapgmyhjtunu
drkcthbhlarhonqvk
yttctnihnatu
pjxqlzodbmmrwoeeusnneqeby
i
zufaevyyqekonpmxqti
alcwuqowtduexfeco
gboudbciflhmtwrfsvtgqaqlgnn
qmkfvsgllrp
kqkc
pllwhgeieavpe
fytveycainmeevwnhgsfigpjzbicq
qfhbxsou
gfcxjweqignzkmaeabaafwntg
zdzhfreoahawfwbfcbj
jlpftrl
gpptririyhk
xsrpfyerwsbpmpazxcecfkpxqvbgkz
senvkkffogbdp
btnrqvwciyshrovoywotzumr
byizedxcjfq
evxdhcoisbbqfgpnc
e
alrjtgnnbvxxxvhbfky
hozwcvjuzwcpahvuvrwuefjx
nhselpxtwmb
uwhdxqsswdgsopdtrgfamwocqs
znoflblwueialoqidshcdyvs
tnojked
rjqwbjnapkwjs
dqjrjplfbvtjsjtiau
hmwtxzizchnfdm
fopoxugvwjzaoitlicies
xxghikiun
ougchyxkebqprxilayhfe
bqbxugptwlitydtvvqexmiydafv
eeydbzbzcqjwfnhmhduxnymmrbo
xqwdjvhryztsqgwfrtusmliuxsqai
h
keeoxivdqmoadtuvpesooyilyy
yorybkofgnuvytivudy
wipqjmvbpypimmztwpnoowloyi
`

func expected(s string) string {
	stack := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		n := len(stack)
		if n > 0 && stack[n-1] == c {
			stack = stack[:n-1]
		} else {
			stack = append(stack, c)
		}
	}
	return string(stack)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		expect := expected(s)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(s + "\n")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %q got %q\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
