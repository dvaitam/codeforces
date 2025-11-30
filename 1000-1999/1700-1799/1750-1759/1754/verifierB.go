package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
139
584
869
823
784
66
263
122
509
781
462
485
669
390
809
216
98
501
31
916
857
401
445
624
782
787
4
714
458
274
740
823
236
607
969
106
925
327
33
24
28
667
556
11
963
904
392
704
223
994
434
745
31
542
229
784
450
963
509
568
240
355
238
695
226
781
472
977
298
950
24
428
859
940
571
946
659
104
192
646
743
882
305
125
762
342
919
740
998
730
514
960
992
434
521
851
934
688
196
312`

func minDiff(arr []int) int {
	m := 1<<31 - 1
	for i := 0; i < len(arr)-1; i++ {
		d := arr[i+1] - arr[i]
		if d < 0 {
			d = -d
		}
		if d < m {
			m = d
		}
	}
	return m
}

// reference construction from 1754B.go
func buildPermutation(n int) []int {
	switch n {
	case 2:
		return []int{1, 2}
	case 3:
		return []int{1, 2, 3}
	}
	if n%2 == 0 {
		res := make([]int, n)
		for i := 0; i < n; i++ {
			if i == 0 {
				res[i] = n/2 + 1
			} else if i%2 == 0 {
				res[i] = (n+i)/2 + 1
			} else {
				res[i] = i/2 + 1
			}
		}
		return res
	}
	res := make([]int, 0, n)
	res = append(res, n)
	nn := n - 1
	for i := 0; i < nn; i++ {
		if i%2 == 0 {
			res = append(res, (nn+i)/2+1)
		} else {
			res = append(res, i/2+1)
		}
	}
	return res
}

func loadTestcases() ([]int, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	ns := make([]int, 0, t)
	for i := 1; i < len(fields); i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("bad n: %w", err)
		}
		ns = append(ns, v)
	}
	if len(ns) != t {
		return nil, fmt.Errorf("expected %d testcases, got %d", t, len(ns))
	}
	return ns, nil
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	tokens := strings.Fields(strings.TrimSpace(out.String()))
	if len(tokens) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens))
	}
	seen := make([]bool, n+1)
	arr := make([]int, n)
	for i, t := range tokens {
		v, err := strconv.Atoi(t)
		if err != nil {
			return fmt.Errorf("invalid integer %q", t)
		}
		arr[i] = v
	}
	// Reference solution prints '\n' as 10 appended to the last number; adjust if seen.
	if arr[n-1]%100 == 10 {
		arr[n-1] /= 100
	}
	for _, v := range arr {
		if v < 1 || v > n || seen[v] {
			return fmt.Errorf("invalid permutation")
		}
		seen[v] = true
	}
	md := minDiff(arr)
	refMD := minDiff(buildPermutation(n))
	if md != refMD {
		return fmt.Errorf("expected min diff %d got %d", refMD, md)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ns, err := loadTestcases()
	if err != nil {
		fmt.Printf("failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for i, n := range ns {
		if err := runCase(bin, n); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
