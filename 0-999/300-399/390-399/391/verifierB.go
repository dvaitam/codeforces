package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"M", "Y", "N", "B", "I", "Q", "P", "M", "Z", "J",
	"P", "L", "S", "G", "Q", "EJ", "EY", "DT", "ZI", "RW",
	"ZT", "EJ", "DX", "CV", "KP", "RD", "LN", "KT", "UG", "RP",
	"OQI", "BZR", "ACX", "MWZ", "VUA", "TPK", "HXK", "WCG", "SHH", "ZEZ",
	"ROC", "CKQ", "PDJ", "RJW", "DRK", "RGZT", "RSJO", "CTZM", "KSHJ", "FGFB",
	"TVIP", "CCVY", "EEBC", "WRVM", "WQIQ", "ZHGV", "SNSI", "OPVU", "WZLC", "KTDP",
	"SUKGH", "AXIDW", "HLZFK", "NBDZE", "WHBSU", "RTVCA", "DUGTS", "DMCLD", "BTAGF", "WDPGX",
	"ZBVAR", "NTDIC", "HCUJL", "NFBQO", "BTDWM", "GILXPS", "FWVGYB", "ZVFFKQ", "IDTOVF", "APVNSQ",
	"JULMVI", "ERWAOX", "CKXBRI", "EHYPLT", "JVLSUT", "EWJMXN", "UCATGW", "KFHHUO", "MWVSNB", "MWSNYV",
	"WBFOCIW", "FOQPRTY", "ABPKJOB", "ZZNGRUC", "XEAMVNK", "AGAWYAV", "QTDGDTU", "GJIWFDP", "MUCAIOZ", "ZDIEUQU",
}

func solveCase(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	ans := 1
	for p0 := 0; p0 < n; p0++ {
		memo := make(map[int]int)
		var dfs func(pTip, w int) int
		dfs = func(pTip, w int) int {
			key := pTip*(n+1) + w
			if v, ok := memo[key]; ok {
				return v
			}
			best := 1
			maxF := (w + pTip - 2) / 2
			for f := pTip; f <= maxF; f++ {
				pNew := 2*f - pTip + 1
				if pNew < 0 || pNew >= w {
					continue
				}
				if s[pNew] != s[p0] {
					continue
				}
				h := 1 + dfs(pNew, f+1)
				if h > best {
					best = h
				}
			}
			memo[key] = best
			return best
		}
		h := dfs(p0, n)
		if h > ans {
			ans = h
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, s := range rawTestcases {
		expected := strconv.Itoa(solveCase(s))
		got, err := run(bin, s+"\n")
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
