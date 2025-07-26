package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseD struct {
	strs []string
}

type strInfo struct {
	s     string
	cntS  int64
	cntH  int64
	noise int64
}

func calcInfo(t string) strInfo {
	var cntS, cntH, noise int64
	for _, ch := range t {
		if ch == 's' {
			cntS++
		} else {
			noise += cntS
			cntH++
		}
	}
	return strInfo{s: t, cntS: cntS, cntH: cntH, noise: noise}
}

func expectedD(arr []string) int64 {
	infos := make([]strInfo, len(arr))
	for i, t := range arr {
		infos[i] = calcInfo(t)
	}
	sort.Slice(infos, func(i, j int) bool {
		a := infos[i]
		b := infos[j]
		return a.cntS*b.cntH > b.cntS*a.cntH
	})
	var res, sCount int64
	for _, it := range infos {
		res += it.noise
	}
	for _, it := range infos {
		res += sCount * it.cntH
		sCount += it.cntS
	}
	return res
}

func genTestsD() []testCaseD {
	rand.Seed(4)
	tests := make([]testCaseD, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		strs := make([]string, n)
		for i := 0; i < n; i++ {
			l := rand.Intn(10) + 1
			var sb strings.Builder
			for j := 0; j < l; j++ {
				if rand.Intn(2) == 0 {
					sb.WriteByte('s')
				} else {
					sb.WriteByte('h')
				}
			}
			strs[i] = sb.String()
		}
		tests = append(tests, testCaseD{strs: strs})
	}
	return tests
}

func runCase(bin string, tc testCaseD) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(tc.strs)))
	for _, s := range tc.strs {
		input.WriteString(s)
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	expect := fmt.Sprint(expectedD(tc.strs))
	if gotStr != expect {
		return fmt.Errorf("expected %s got %s", expect, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
