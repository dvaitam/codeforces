package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `865 395
777 912
431 42
266 989
524 498
415 941
803 850
311 992
489 367
598 914
930 224
517 143
289 144
774 98
634 819
257 932
546 723
830 617
924 151
318 102
600 577
595 176
690 634
390 771
14 812
10 893
837 350
318 771
782 460
863 158
85 718
646 298
680 990
698 647
158 840
872 763
636 327
571 772
653 189
948 314
118 607
489 50
14 984
399 458
13 252
197 9
616 795
68 715
241 733
486 548
900 663
676 801
302 177
231 252
944 989
237 203
965 392
726 474
493 903
938 38
455 503
72 546
130 156
726 252
345 44
772 773
776 855
761 948
465 423
114 584
395 361
968 368
354 547
559 982
174 422
85 390
385 92
729 954
330 565
6 43
516 608
86 754
8 39
191 329
956 816
95 153
66 351
755 709
576 215
40 998
26 410
722 872
805 200
606 93
325 276
725 236
304 24
370 601
69 186
770 889
405 110`

func expected(n, m int64) string {
	if n > m {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%d", m)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func loadCases() ([]string, []string) {
	lines := strings.Fields(testcasesRaw)
	if len(lines)%2 != 0 {
		fmt.Println("invalid embedded testcases")
		os.Exit(1)
	}
	var inputs []string
	var expects []string
	for i := 0; i < len(lines); i += 2 {
		n, _ := strconv.ParseInt(lines[i], 10, 64)
		m, _ := strconv.ParseInt(lines[i+1], 10, 64)
		inputs = append(inputs, fmt.Sprintf("1\n%d %d\n", n, m))
		expects = append(expects, expected(n, m))
	}
	return inputs, expects
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expects[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expects[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
