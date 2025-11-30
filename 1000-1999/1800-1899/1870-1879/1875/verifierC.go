package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `138 583
868 822
783 65
262 121
508 780
461 484
668 389
808 215
97 500
30 915
856 400
444 623
781 786
3 713
457 273
739 822
235 606
968 105
924 326
32 23
27 666
555 10
962 903
391 703
222 993
433 744
30 541
228 783
449 962
508 567
239 354
237 694
225 780
471 976
297 949
23 427
858 939
570 945
658 103
191 645
742 881
304 124
761 341
918 739
997 729
513 959
991 433
520 850
933 687
195 311
291 602
997 904
512 867
964 518
403 604
874 36
492 249
762 817
414 425
681 178
376 562
904 720
795 691
756 384
89 450
680 521
111 798
168 534
861 403
380 502
751 31
481 45
316 721
869 630
608 593
404 663
175 173
515 233
13 790
205 553
943 881
562 238
415 527
353 976
868 592
362 471
932 276
676 562
624 981
747 6
393 803
878 841
978 908
961 759
525 829
133 532
797 575
211 437
973 58
493 891`

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expected(n, m int64) int64 {
	g := gcd(n, m)
	n1 := n / g
	m1 := m / g
	if m1&(m1-1) != 0 {
		return -1
	}
	q := n1 / m1
	r := n1 % m1
	needed := m*q + m*int64(bits.OnesCount64(uint64(r)))
	ops := needed - n
	if ops < 0 {
		ops = 0
	}
	return ops
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func loadCases() ([][2]int64, error) {
	tokens := strings.Fields(testcasesRaw)
	if len(tokens)%2 != 0 {
		return nil, fmt.Errorf("invalid embedded testcases")
	}
	cases := make([][2]int64, 0, len(tokens)/2)
	for i := 0; i < len(tokens); i += 2 {
		n, err1 := strconv.ParseInt(tokens[i], 10, 64)
		m, err2 := strconv.ParseInt(tokens[i+1], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid embedded value")
		}
		cases = append(cases, [2]int64{n, m})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		n64, m64 := tc[0], tc[1]
		exp := expected(n64, m64)
		input := fmt.Sprintf("1\n%d %d\n", n64, m64)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", idx+1, out)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
