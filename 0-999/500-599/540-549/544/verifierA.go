package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func possible(k int, q string) bool {
	seen := make([]bool, 26)
	cnt := 0
	for i := 0; i < len(q) && cnt < k; i++ {
		c := q[i] - 'a'
		if !seen[c] {
			seen[c] = true
			cnt++
		}
	}
	return cnt >= k
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

const testcasesARaw = `1 a
26 abcdefghijklmnopqrstuvwxyz
2 aa
13 nbiqpmzjplsgqejeydtzirwzt
5 dxcvkprdln
11 ugrpoqibzracxmwzvuat
16 khxkwcgshhzezrocckqpdjrjwdr
11 rgztrsjoctzmkshjfgfbtvipccv
25 ebcwr
22 wqiqzhgvsnsio
16 uwzlcktdpsukghaxidwhlz
6 nbdzewhbsur
20 cadugtsdmcldbtagfwdpgx
26 va
18 tdichcujlnfbqo
2 dwmgilxpsfwvgybzvffk
17 dtovfapvn
19 qjulmvierwaoxckxbriehypltjvl
19 utewjmxnucatgwkfhhuomwvsnbmws
14 vwbfociwfoqprtyabpkjobzzn
7 ucxeamvnkagawyavqt
4 dtugjiw
6 pmuc
1 ozzdieuqu
21 deiabbgvirkl
19 bxwtupwuounlrfgmsjaeeikkzlwcky
20 bi
6 sjlmr
5 dpxhbjfqxc
10 kjnddrppkzzkd
16 wpnb
10 xvefusmzucc
26 gxh
2 admrqjopzswvg
14 lhi
19 fngldcwaqoyvgdpmigubzgted
7 mlredtpesmuvnqp
22 kppuvgrthakwxkkbqeitzemsjwwzpc
26 qbc
8 bjayo
11 fzeuolqmqqbscvzzqytcxnygjrtn
16 zmtshzavaxfjqsikcpijynmzmbfu
5 jxkbbpne
16 twcvwezlnbtomobdpyeabtteukdwr
21 gmzyypdbtwot
21 kudvwtjzemzjxvzdqzgbzmolygolzu
3 bp
9 aqvssghcyuyqwqnqjdensncdncdny
5 azonvnapkxiclcdlwallfahl
3 egagvvxdxajlwathefod
16 wieaglkpjjru
11 scdrsj
6 eezhkqhhyfjln
22 et
1 cwcenjrnxesnj
21 chouluqbmnan
24 kogljpcfzdidrtwezwomfynnfhokq
5 ouucpygjawot
15 g
10 yujr
20 nwypc
22 yhrymiuadivbaimq
19 modxiljyvgtcbczijrkdqhy
6 njj
17 squgr
4 urmxyzijolsuef
4 dmmsoervjlupxngppwqkpub
15 expbtgalpm
1 qcvcvxvmalbdtaiuwjxheysjgdno
23 mfknuvneowe
17 egfolzmnzpm
24 zgogswbm
2 ucflbxuv
6 tjtcwqyj
25 lnobuwqvurxnsopiwpgkibbbflaju
1 cznvh
20 rhogktdtczkkr
15 kiaqbglcgqlggivxxjjqmiplwhbjr
3 o
16 obznpoodcchdyengotcnrymb
6 pheilknd
18 jtrzgwjyoqtoruiihadtzwdfxnhgj
24 axrqnbdmuidxslhvwwrvjh
24 hcqjvkhlupjsfezarqklsuazeme
6 ceygzypsywghxehym
12 seupdtaqtlpojahrufvp
26 prkwcietmwgkzjmbgbkxxhko
22 vvhilvfjalsrbxuelapubahb
1 ukcblvne
7 neljfukxzxnmani
18 zxvwoybsdnmfaqetvq
24 eckhfhazfxzvrfwcntdtuow
5 tbikzxxmaubpcljveohq
12 xymkiz
16 ajqjrpbyrsriv
2 mxdmlpbaixbivvs
23 jvygyqqkmigdskzhsvxrvlfek
24 asbselljujkpzmtnfazesboekax
16 vvyixtgcrnifqfcvufsdquzrtmynij
10 n
`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
			fmt.Printf("test %d malformed line\n", idx)
			os.Exit(1)
		}
		k := 0
		for _, ch := range parts[0] {
			k = k*10 + int(ch-'0')
		}
		q := parts[1]
		input := fmt.Sprintf("%d\n%s\n", k, q)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) == 0 {
			fmt.Printf("test %d: empty output\n", idx)
			os.Exit(1)
		}
		ans := strings.TrimSpace(lines[0])
		if ans == "NO" {
			if possible(k, q) {
				fmt.Printf("test %d: answer should be YES\n", idx)
				os.Exit(1)
			}
			continue
		}
		if ans != "YES" {
			fmt.Printf("test %d: first line must be YES or NO\n", idx)
			os.Exit(1)
		}
		if !possible(k, q) {
			fmt.Printf("test %d: answer should be NO\n", idx)
			os.Exit(1)
		}
		if len(lines)-1 != k {
			fmt.Printf("test %d: expected %d lines, got %d\n", idx, k, len(lines)-1)
			os.Exit(1)
		}
		firsts := make(map[byte]bool)
		var concat strings.Builder
		for i := 1; i <= k; i++ {
			s := lines[i]
			if len(s) == 0 {
				fmt.Printf("test %d: empty string on line %d\n", idx, i)
				os.Exit(1)
			}
			if firsts[s[0]] {
				fmt.Printf("test %d: first characters not distinct\n", idx)
				os.Exit(1)
			}
			firsts[s[0]] = true
			concat.WriteString(s)
		}
		if concat.String() != q {
			fmt.Printf("test %d: concatenation mismatch\n", idx)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
