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

// solution logic from 1913F.go
func countPal(s []byte) int64 {
	n := len(s)
	rad1 := make([]int, n)
	l, r := 0, -1
	for i := 0; i < n; i++ {
		k := 1
		if i <= r {
			if rad1[l+r-i] < r-i+1 {
				k = rad1[l+r-i]
			} else {
				k = r - i + 1
			}
		}
		for i-k >= 0 && i+k < n && s[i-k] == s[i+k] {
			k++
		}
		rad1[i] = k
		if i+k-1 > r {
			l = i - k + 1
			r = i + k - 1
		}
	}
	rad2 := make([]int, n)
	l, r = 0, -1
	for i := 0; i < n; i++ {
		k := 0
		if i <= r {
			if rad2[l+r-i+1] < r-i+1 {
				k = rad2[l+r-i+1]
			} else {
				k = r - i + 1
			}
		}
		for i-k-1 >= 0 && i+k < n && s[i-k-1] == s[i+k] {
			k++
		}
		rad2[i] = k
		if i+k-1 > r {
			l = i - k
			r = i + k - 1
		}
	}
	var count int64
	for i := 0; i < n; i++ {
		count += int64(rad1[i])
		count += int64(rad2[i])
	}
	return count
}

func solveCase(n int, s string) (string, error) {
	if len(s) != n {
		return "", fmt.Errorf("string length mismatch")
	}
	bs := []byte(s)

	bestCount := countPal(bs)
	bestStr := s

	for i := 0; i < n; i++ {
		orig := bs[i]
		for c := byte('a'); c <= byte('z'); c++ {
			if c == orig {
				continue
			}
			bs[i] = c
			cnt := countPal(bs)
			str := string(bs)
			if cnt > bestCount || (cnt == bestCount && str < bestStr) {
				bestCount = cnt
				bestStr = str
			}
		}
		bs[i] = orig
	}

	return fmt.Sprintf("%d\n%s", bestCount, bestStr), nil
}

const testcasesData = `
17 nfejlrweuwpkluzui
19 ezxouhkwjzgkjldtiuh
7 ppgqaul
16 myobymbugmkdyyyy
1 z
4 ztiq
10 dioxlipvyp
6 xdzbjc
17 erprmhwudbhlkvtpf
16 ufvikiitzjwsatst
3 hmn
2 va
3 rzz
20 ysiyqpvlahbzpbqhghcl
17 iygurnnzmkeekbhol
16 njauistdnbcagaoo
7 flpyajl
1 q
12 rckcrokdragz
7 srtgbku
6 lmaqoi
15 gnaieinspqizaum
16 gtlejbdevgernzbb
12 zypixcjbatoy
7 rmwmmig
6 sclahf
10 scsbzldcpq
2 ur
11 mxsxidbkkna
18 dclfezxtjzxhrlynlb
12 mywuvyhfxlqe
4 rxyv
19 pgcsjgtnrojlckvyxbd
8 uklmxavo
2 tu
20 rvughxwdbpsfajwygjfv
2 ju
19 eykghqjzkeogmlxgviz
5 oazyi
11 siwsyqyesub
16 tooxdstbvgxaqdzl
2 ty
16 fszrsukpqhxulisp
12 qakbllrtorkd
11 uyuabeyywkx
11 tzycbbvywlw
3 yba
17 nbngryvsqgiikmxka
8 grhdcaab
16 zaonzfebfiehboft
6 uedipg
2 ky
20 nekelhnezqewgjbeqonl
1 k
18 olpmcofzzvlqmygdsh
6 gwbrfx
8 stkllwxj
8 zcazaiwd
16 gdpsvqwxfxmvoxup
11 ukgwadudibr
11 mcjewsibqat
19 pvakhldmrfekbhqysjg
11 tukfhgrfloi
14 nwgjkweosksakd
1 a
14 nvfsngrjrvjedo
10 wisdlayfqb
14 sqsspiprtfvqyv
14 zdpgsigbjgmzdj
13 akuvpvhiixdha
4 eufq
18 xjgcfryvrwyjhydubh
17 dyozcfhpkfawrqlht
12 fhcqziqblrhe
10 ymxnecrkgw
5 qzjhu
19 dfyzcjbccyhojkohgop
1 a
9 eizuedkoh
7 diqmfgw
4 omke
5 hhjxq
20 yshxbssyqkwlthpnyflr
20 qvzdaqgsjdjoyjcgrcel
14 pivqgkyczshirq
12 txbawgwrbika
3 hhd
4 ozax
18 vxjtzhycuxoiksxykz
14 lbalmtxisunjwm
17 vtntwwuqxeomrhycc
12 yhfsblyvwull
2 js
15 vgjpcgsemrrtjcq
5 fbuxe
18 swmgbifeawjfkxtitu
3 ior
2 sp
14 dpjpkusewzgrbg
8 nofcfljn
`

func lineToInput(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return "", fmt.Errorf("bad test line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	s := fields[1]
	if len(s) != n {
		return "", fmt.Errorf("string length mismatch")
	}
	return fmt.Sprintf("%d\n%s\n", n, s), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input, err := lineToInput(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test %d: %v\n", idx, err)
			os.Exit(1)
		}
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		want, err := solveCase(n, parts[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad solve on test %d: %v\n", idx, err)
			os.Exit(1)
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != want {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx, input, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
