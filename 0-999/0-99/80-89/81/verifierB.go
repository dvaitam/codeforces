package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesBRaw = `,,
629 234 ......  18 ...318
..., , 213 434 ,
, 120  , ,......
799... , ... 57 951 ...
238... 665..., ...
...559 457  ,
,,,, ...
971  ...,826...,,
,799...
,...,
...,...  , 705
......123, 438
770 ... ...
559  ,
195716...565...... ...  ,
... ,  959,...302
969 ...
......700
... , 700866 ... ,
,636...
...  ... 560 ... 435
675 ...475 ... 824 ...
542 ... 956 ,936746
...312,860  , ,
... ,...518, ...  343,
299584
537,340,...
...166
... ...  ..., ...
, 345
526 ...
244, 734 ...  447  876,126
890 46 ..., ,
130362  ... 512 289254782
224...147 ...,, 520
,... 523
,... ... ...13  ,
181 149,491
... 64 , ,
......  ......  385...,
, ,,,  ,... ... ,
956 298860  ,..., , ...
...180..., ...862
... 929699
, ......
... ,
660..., ,
......,
...,,,
773...
... ...
... 243
, 732 ,  ,  ...5747
,  ......
378  ... ...
... 228
834235
...825, 852
, ,
774 ...239 , 851
...164 ...  , 402... ,996
300,... ,...  ...,  ...
96 139 ...,, 770 ...
, ...
, 117,,
, ... 709 615..., 903
207 449,
366 ......,
... 106 ...  817 ...
... ...106...
806 ...  886 312... ,  ,,
405......,, 345
289 ,
, ,
, ...
...... ...  , , 505,851
, ......... 629 ...... 94
951 227 186 ...723
603303 928  170
...,
,478 904 , ,,  ,
359713,654
, ... ...230 889... ,...
, , ...  ...  ,
...  ......
,, 173,...  16...,
,... 884 ...  , ,  , 709
,..., ,... ,... ...
... 894...
,, ,,
...,... ...,877 242,
, ...,, , ...
,850
527 ... ......
...423 ...  694 ... ...,
,  671... 985 881 ,
,91 ......... 671 290
22913
... ,900484,
`

const (
	tokNum = iota
	tokComma
	tokDots
)

func expectedLine(s string) string {
	s = strings.TrimRight(s, "\r\n")
	var tokTypes []int
	var tokVals []string
	n := len(s)
	for i := 0; i < n; {
		switch {
		case s[i] >= '0' && s[i] <= '9':
			j := i
			for j < n && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			tokTypes = append(tokTypes, tokNum)
			tokVals = append(tokVals, s[i:j])
			i = j
		case i+2 < n && s[i] == '.' && s[i+1] == '.' && s[i+2] == '.':
			tokTypes = append(tokTypes, tokDots)
			tokVals = append(tokVals, "...")
			i += 3
		case s[i] == ',':
			tokTypes = append(tokTypes, tokComma)
			tokVals = append(tokVals, ",")
			i++
		case s[i] == ' ':
			i++
		default:
			i++
		}
	}
	var b strings.Builder
	m := len(tokTypes)
	for i, t := range tokTypes {
		switch t {
		case tokNum:
			if i > 0 && tokTypes[i-1] == tokNum {
				b.WriteByte(' ')
			}
			b.WriteString(tokVals[i])
		case tokComma:
			b.WriteString(",")
			if i < m-1 {
				b.WriteByte(' ')
			}
		case tokDots:
			if b.Len() > 0 {
				str := b.String()
				if str[len(str)-1] != ' ' {
					b.WriteByte(' ')
				}
			}
			b.WriteString("...")
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		idx++
		expect := expectedLine(line)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(line + "\n")
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
			fmt.Printf("test %d failed\nexpected: %q\n     got: %q\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
