package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (same format as original file).
const embeddedTestcases = `100
13 16 yuinvetemjqbf
5 17 pwuwy
2 93 qb
18 88 xwumflzsccrfigzikw
9 100 iqoeyoreb
19 22 uqbykcguzothzoqfwku
5 60 yrbrc
17 43 aycdntlsokmqludek
1 23 x
5 2 ktgbn
2 89 jz
13 6 tyzwflcnboltt
9 86 jszonfaoz
9 24 mclddalaf
13 78 wuakozrxwppcb
18 51 zziauqdckldpbequzj
2 0 mk
6 70 weffyf
8 79 kapuvmbh
1 75 g
6 61 zotiso
19 46 hiznwgaoymhkxwpmpdq
12 18 juhooztxdwpl
4 72 onmb
14 59 nkslzstkxifvzb
13 80 pbxkyrsnrupyt
5 15 peehk
20 45 bvztpvctxlvzuysmuodl
2 87 fj
12 12 koqocoavuwtd
13 83 odsonopnlkbez
13 28 swboibiwxioov
17 18 ntbrlhczjabzcvwom
14 74 rqsjhkczeygzwx
16 15 ezifpspgevpuyjmn
19 58 qnlguovjjcifnretnkn
5 64 adssy
19 95 ilplecjbhzdpwkaoeni
19 97 hmtbjffenxkbwypzese
13 99 sjaoqiowghxsj
20 68 acwjipwledbitouyyfmd
6 39 aonpdz
11 82 slsnanlxwqa
2 70 cz
2 74 ei
10 19 nugrmxfvhi
1 77 u
19 43 olpoezrmnpbuxvleapb
18 50 lvxgefkvgzypzfejxj
10 5 axgllodshj
20 75 xaojoudfrsszmfyepfpz
8 44 epdowvlx
1 44 v
4 20 lequ
6 15 uvvloe
3 56 ngy
12 41 audzileobgmn
4 21 zqhg
15 27 yjspckpoicojroz
10 61 sgnfxmqdct
18 45 kstvrbwcuokzhbvndk
5 69 qxazu
11 50 ohtpiaouxdd
13 97 mhfcaajnqmmxp
15 36 jbbgqvymjbqqlpd
11 81 dheebylpksb
10 30 iczyxomgqd
3 72 dhg
6 11 jevuzv
9 47 rzdcyozxz
5 77 lvkpz
5 68 iftjv
18 23 snlayebekioycjhbxzq
18 29 boihlobteoycpsfbxb
6 4 vjgbzj
16 50 wpwrqpoixsmrkzmg
3 2 aoj
1 80 o
5 25 nrzjb
2 27 ia
4 19 drfc
8 46 rnxsatmm
1 21 z
10 36 ihhkjsptpt
3 97 rfr
1 17 h
2 94 vn
12 8 obmngsprfylk
10 38 czvuczbgim
4 79 ccoe
20 57 neydlgnkkuweoteuvwro
14 69 gexiyxnvethxmy
13 12 kweqhtufyhqa
11 26 snkiyjznivc
6 98 yenxia
17 1 xuplwypmvwvxpnnkp
10 60 efoovlqqye
11 25 ujqmxghhpwq
3 30 kwf
4 47 zqqp
7 71 whptidu
5 58 alrve
16 15 rhpzxmgumqwxgtzy
10 32 xeroydgzqq
1 44 z
12 26 gxkjlrwnwwsy
5 22 nbvos
15 52 zgffkolsifrfqfl
8 50 yjnhbstu`

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func maxDistance(s string) int {
	dist := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		up := int('z' - c)
		down := int(c - 'a')
		if up > down {
			dist += up
		} else {
			dist += down
		}
	}
	return dist
}

func checkAnswer(s string, k int, ans string) bool {
	if ans == "-1" {
		return maxDistance(s) < k
	}
	if len(ans) != len(s) {
		return false
	}
	dist := 0
	for i := 0; i < len(s); i++ {
		if ans[i] < 'a' || ans[i] > 'z' {
			return false
		}
		diff := int(s[i]) - int(ans[i])
		if diff < 0 {
			diff = -diff
		}
		dist += diff
	}
	if dist != k {
		return false
	}
	if maxDistance(s) < k {
		return false
	}
	return true
}

// Embedded reference solution (from 628C.go, adjusted to trust string length).
func solve(n, k int, s string) string {
	_ = n // length may be inconsistent in test data; rely on actual string length.
	buf := []byte(s)
	n = len(buf)

	maxDist := 0
	for i := 0; i < n; i++ {
		c := buf[i]
		t1 := int('z' - c)
		t2 := int(c - 'a')
		if t1 > t2 {
			maxDist += t1
		} else {
			maxDist += t2
		}
	}
	if maxDist < k {
		return "-1"
	}

	for i := 0; i < n && k > 0; i++ {
		c := buf[i]
		t1 := int('z' - c)
		t2 := int(c - 'a')
		if t1 >= k {
			buf[i] = c + byte(k)
			k = 0
			break
		} else if t2 >= k {
			buf[i] = c - byte(k)
			k = 0
			break
		} else {
			if t1 > t2 {
				buf[i] = 'z'
				k -= t1
			} else {
				buf[i] = 'a'
				k -= t2
			}
		}
	}

	return string(buf)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scan := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 1; i <= t; i++ {
		if !scan.Scan() {
			fmt.Printf("missing n for case %d\n", i)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		kVal, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		s := scan.Text()
		realN := len(s)
		_ = n // trust actual string length to avoid panics on malformed cases.
		input := fmt.Sprintf("%d %d\n%s\n", realN, kVal, s)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		expect := solve(realN, kVal, s)
		if out != expect {
			fmt.Printf("case %d failed: expected %q got %q\n", i, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
