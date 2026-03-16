package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

var mapping = map[rune]string{
	'0': "",
	'1': "",
	'2': "2",
	'3': "3",
	'4': "322",
	'5': "5",
	'6': "53",
	'7': "7",
	'8': "7222",
	'9': "7332",
}

func solveDigits(s string) string {
	var res []byte
	for _, ch := range s {
		if rep, ok := mapping[ch]; ok {
			for i := 0; i < len(rep); i++ {
				res = append(res, rep[i])
			}
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i] > res[j] })
	return string(res)
}

const testcasesCRaw = `100
11 85225090925
6 445769
7 2029072
6 074931
9 642821291
9 896644406
13 4488855362082
13 9657086930582
11 35703934261
10 7378132716
11 64465951407
1 4
4 6666
12 097592945067
15 820195501318705
15 056246292648022
15 270808625911824
4 4544
9 727820928
1 5
14 13793782527808
2 85
1 3
10 5975685125
15 022059306046092
6 160759
15 944976207436151
2 05
1 2
7 9057877
2 08
14 64081155170284
1 8
3 822
3 239
6 076033
11 45235326759
3 690
3 906
12 220058000199
3 260
7 6953258
4 8612
7 9516637
7 3637691
14 44858099734095
14 61802606661797
3 257
7 2948155
3 578
13 0342954604860
11 64635221952
1 6
10 6711688222
4 2308
3 759
12 451642585826
12 843656784661
11 22099813998
2 49
3 610
1 7
6 175947
13 3288182028693
8 64092627
10 0619653870
8 91487867
5 23852
5 92711
8 69681473
14 14251202983006
9 971785510
4 3744
4 0758
15 511496356234375
5 69216
6 873556
5 56417
15 417253257316678
8 99368473
6 801756
15 360906143521250
4 0068
1 2
2 99
15 317308961823366
8 06360940
10 5557298141
12 140292639536
9 818171877
3 785
3 646
2 98
6 373576
1 7
12 086736347723
8 45792813
5 81025
1 5
14 92763673260583
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data := []byte(testcasesCRaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		digits := scan.Text()
		expected[i] = solveDigits(digits)
		_ = n // n not needed
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		if outScan.Text() != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], outScan.Text())
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
