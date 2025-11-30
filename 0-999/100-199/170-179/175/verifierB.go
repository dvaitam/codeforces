package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded test cases from testcasesB.txt. Format: t, then for each test case n followed by n lines "name score".
const testcasesRaw = `100
2
c 120
dddb 96
4
d 443
d 272
ea 923
aaa 665
5
d 702
da 540
dd 566
cb 693
dc 948
1
eabc 123
6
ede 849
cc 601
edea 491
dd 680
ce 903
ade 110
2
cdad 44
eee 403
6
be 232
b 552
de 352
dce 623
d 802
ee 210
4
d 890
eeb 963
dcdc 1
dea 823
2
ee 185
e 816
3
a 85
d 14
bca 816
5
cc 71
bc 540
cc 465
dda 24
dcd 815
2
ace 1000
ed 836
1
ad 149
1
de 694
4
ed 228
d 691
dac 128
ac 72
1
cbd 578
3
ae 899
e 839
ed 175
6
d 205
abe 690
ebda 960
ceda 333
cabb 878
ebc 439
2
ade 352
ebaa 86
2
be 218
cee 861
3
cca 298
ed 138
c 40
4
d 886
bc 117
aeeb 579
c 373
3
d 918
aac 12
a 423
1
b 245
5
badb 697
ba 445
ecec 728
cabc 40
a 805
3
ddc 408
a 935
eda 256
2
ccbe 212
bbc 83
3
d 92
bdc 42
bce 914
3
ca 557
b 225
b 411
1
eaa 22
6
c 768
ddb 103
aeb 183
bc 312
e 854
bbb 558
6
c 840
bc 443
ab 258
d 827
eced 871
adcb 264
4
d 999
a 708
ebe 128
cc 407
5
beab 497
b 541
edb 244
ddb 729
ceec 995
6
aa 781
beb 319
cec 169
eaae 983
bbcd 222
d 697
4
deb 557
e 92
aca 982
ea 455
2
ddbc 448
ed 982
2
d 615
accb 387
6
b 541
eaae 248
bbc 151
cc 599
dbe 365
dabe 900
4
ca 925
a 582
e 303
ae 382
5
dec 776
aad 735
cced 347
addb 570
c 650
5
de 854
cbde 684
ce 3
eddc 881
d 763
2
adb 648
cbae 10
3
dec 155
cdbd 522
c 522
1
acad 20
2
ad 651
ecb 540
2
cc 70
e 674
3
eeab 304
ceb 401
bdce 337
6
ce 723
ae 412
dbc 194
b 891
ebec 470
bb 915
6
ccdb 118
ca 108
dc 504
b 46
e 23
ad 720
5
ccae 709
ab 409
dd 386
bb 839
dee 398
2
ccde 113
aa 15
1
cdec 940
2
bbaa 396
ea 578
4
bad 667
aae 62
ac 799
d 93
2
d 652
cb 678
4
cccb 251
e 957
cd 619
c 560
4
ed 941
c 761
c 181
b 60
2
aaae 480
aca 129
5
d 680
dd 25
acc 87
ada 750
cbc 813
4
c 96
beeb 338
ede 492
b 668
4
c 761
bc 398
adc 129
a 307
6
dcc 361
cee 8
b 324
cea 462
ddc 949
aeab 49
5
ecbe 764
ccd 314
ecee 171
b 256
eb 928
1
de 51
1
abc 68
6
a 813
be 882
aecd 727
bbe 505
dd 691
ebd 743
1
dba 764
5
edad 630
acda 194
aea 841
ece 660
ede 962
5
eecd 309
ed 600
eb 258
d 753
c 430
4
aaa 866
d 275
ccdc 397
adcb 425
2
b 833
cbe 805
3
cecd 707
dcd 220
ddaa 132
2
ba 105
bda 408
6
aa 437
e 223
caae 695
accb 491
b 693
d 126
6
cedd 118
abde 927
be 263
ecde 934
ec 881
aacc 57
5
cabe 281
bdb 133
bde 645
e 854
dc 286
4
cdb 510
edb 346
eb 756
ebae 333
5
bc 637
dcab 908
cb 90
e 176
b 576
2
dca 792
c 841
5
ab 286
cee 388
a 337
bac 920
ea 355
1
a 307
3
ce 50
aab 947
cbac 280
1
acb 620
3
aeee 486
edcb 647
eba 614
5
b 246
dc 559
c 551
ecd 129
acae 371
5
e 315
bbae 145
dc 373
bbd 853
daeb 276
3
e 975
b 388
d 31
4
ccdd 620
aada 661
a 851
e 142
5
ece 930
dbe 245
e 974
baa 938
dcc 673
6
e 445
dccc 451
ee 147
c 689
e 176
caea 491
2
bdba 254
cbd 760
4
dbd 451
eaed 949
bba 385
aaab 469
4
bbe 845
c 19
dbed 5
db 678
2
bae 957
bd 599
1
db 813
1
ee 669
5
b 407
aead 91
d 46
aa 978
dcd 170
5
ec 789
edeb 715
dbdc 368
ce 286
ea 748
3
bcc 258
dce 478
b 976
5
de 948
eca 10
bac 912
eaa 799
ced 236
2
cb 762
c 138
5
d 974
a 833
e 67
dd 1
ddc 794
4
ac 910
da 817
ede 215
ddd 601
4
e 1000
bcb 726
cd 44
c 179
3
ebb 404
ddb 304
dc 327
6
ee 775
e 561
cc 752
ae 923
ec 81
aaa 873
2
a 30
ee 202
6
e 345
dbec 783
dab 324
d 380
ba 24
eb 268
5
ec 198
e 510
bc 106
dcda 587
ba 626
1
d 116
2
ece 3
cc 192
4
b 97
ec 728
ec 39
ee 814
5
d 733
ab 960
aa 943
eab 772
e 363
2
aa 993
aa 182
2
e 23
bb 443
6
da 443
da 320
cdc 672
c 441
cb 387
bae 509
3
e 593
ecdc 337
ee 974
1
e 984
5
e 509
db 134
e 146
ab 827
ddd 213
1
d 135
6
bce 236
c 57
cd 254
ea 776
db 778
ba 462
6
eb 704
dc 352
ba 853
ab 81
a 685
bb 312
4
b 8
d 761
dcab 293
bcab 294
4
de 946
d 420
a 878
e 521
2
ea 510
ca 546
6
ce 962
eba 968
e 575
e 969
e 923
b 469
2
ac 40
e 648
1
e 240
1
acd 341
3
bbb 58
c 100
dbdb 147
3
d 271
e 725
cad 938
2
ece 938
ad 609
3
e 247
e 996
b 872
4
eb 850
ca 559
adca 386
bce 317
6
aad 744
bb 302
dec 650
cd 440
dd 998
bacd 920
5
e 348
ec 92
e 431
b 411
aaac 257
5
ebd 944
ca 12
bedb 658
aa 317
ed 862
4
ae 454
ce 491
edba 243
abb 780
3
ddd 186
cab 267
cb 559
2
abc 116
eb 840
6
ebe 343
ddcb 575
dc 861
de 64
ee 190
ce 537
5
b 804
b 384
d 491
deca 773
ade 50
3
ce 406
e 848
edb 430
2
bb 337
a 493
5
de 645
eab 423
be 523
de 746
d 983
5
bb 281
be 521
eb 768
bca 177
ab 812
5
deb 166
dcad 973
ceac 139
a 10
b 937
4
e 70
e 849
de 148
aa 250
4
baa 276
eae 20
ca 807
cb 311
1
cca 700
3
abb 116
bc 155
a 511
1
ae 41
4
e 500
dbbc 813
ac 663
e 16
4
b 447
d 705
ec 339
c 14
6
e 791
aa 268
be 370
ae 304
ebe 796
ac 295
6
e 61
ce 481
e 592
e 900
bc 316
e 52
4
aad 566
aebc 818
b 1000
a 21
4
d 771
ce 371
abb 699
aa 185
1
d 55
3
baad 83
e 965
ba 8`

type testCase struct {
	n     int
	names []string
	scores []int
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("unexpected EOF for case %d", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d parse n: %v", caseNum+1, err)
		}
		pos++
		if pos+2*n > len(fields) {
			return nil, fmt.Errorf("case %d truncated", caseNum+1)
		}
		names := make([]string, n)
		scores := make([]int, n)
		for i := 0; i < n; i++ {
			names[i] = fields[pos]
			score, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, fmt.Errorf("case %d score %d: %v", caseNum+1, i+1, err)
			}
			scores[i] = score
			pos += 2
		}
		cases = append(cases, testCase{n: n, names: names, scores: scores})
	}
	return cases, nil
}

// Embedded solver logic from 175B.go.
func expected(tc testCase) []string {
	n := tc.n
	names := tc.names
	scores := tc.scores
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		return names[idx[i]] < names[idx[j]]
	})
	unique := make([]int, 0, n)
	for i := 0; i < len(idx); {
		j := i
		best := idx[i]
		for j+1 < len(idx) && names[idx[j+1]] == names[idx[i]] {
			j++
			if scores[idx[j]] > scores[best] {
				best = idx[j]
			}
		}
		unique = append(unique, best)
		i = j + 1
	}
	idx = unique
	n = len(idx)
	sort.Slice(idx, func(i, j int) bool {
		if scores[idx[i]] == scores[idx[j]] {
			return names[idx[i]] < names[idx[j]]
		}
		return scores[idx[i]] < scores[idx[j]]
	})
	res := make([]string, n)
	for i := 0; i < n; i++ {
		s := scores[idx[i]]
		r := i
		for r+1 < n && scores[idx[r+1]] == s {
			r++
		}
		nwtr := float64(r+1) / float64(n)
		btr := float64(n-(r+1)) / float64(n)
		name := names[idx[i]]
		var cat string
		switch {
		case nwtr >= 0.99:
			cat = "pro"
		case nwtr >= 0.9 && btr > 0.01:
			cat = "hardcore"
		case nwtr >= 0.8 && btr > 0.1:
			cat = "average"
		case nwtr >= 0.5 && btr > 0.2:
			cat = "random"
		case btr > 0.5:
			cat = "noob"
		}
		res[i] = fmt.Sprintf("%s %s", name, cat)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for caseNum, tc := range cases {
		exp := expected(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 0; i < tc.n; i++ {
			sb.WriteString(fmt.Sprintf("%s %d\n", tc.names[i], tc.scores[i]))
		}
		gotStr, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum+1, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(gotStr), "\n")
		if len(outLines) == 0 {
			fmt.Fprintf(os.Stderr, "case %d missing output\n", caseNum+1)
			os.Exit(1)
		}
		m, err := strconv.Atoi(strings.TrimSpace(outLines[0]))
		if err != nil || m != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got header %q\n", caseNum+1, len(exp), outLines[0])
			os.Exit(1)
		}
		if len(outLines)-1 != m {
			fmt.Fprintf(os.Stderr, "case %d incomplete output\n", caseNum+1)
			os.Exit(1)
		}
		for i := 0; i < m; i++ {
			if strings.TrimSpace(outLines[1+i]) != exp[i] {
				fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%d\n%s\ngot:\n%d\n%s\n",
					caseNum+1, len(exp), strings.Join(exp, "\n"), m, strings.Join(outLines[1:], "\n"))
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
