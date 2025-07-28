package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var ans = []string{
	`AH QH TH 9H 8H 7H QS JS 7S AD KD 7D 6D AC JC 9C 8C 6C`,
	`KH JH 6H AS KS TS 9S 8S 6S QD JD TD 9D 8D KC QC TC 7C`,
	`AH TH 8H 7H 6H AS 9S 8S 6S AD KD QD JD TD 9D AC KC 6C`,
	`KH QH JH 9H KS QS JS TS 7S 8D 7D 6D QC JC TC 9C 8C 7C`,
	`AH QH JH TH 8H 7H KS QS 7S QD 9D 8D 6D AC TC 9C 8C 7C`,
	`KH 9H 6H AS JS TS 9S 8S 6S AD KD JD TD 7D KC QC JC 6C`,
	`KH JH 8H 7H 6H KS TS 7S AD KD JD 9D 8D 7D KC JC 9C 7C`,
	`AH QH TH 9H AS QS JS 9S 8S 6S QD TD 6D AC QC TC 8C 6C`,
	`AH JH 7H 6H KS 7S KD QD TD 9D 8D AC QC TC 9C 8C 7C 6C`,
	`KH 8H 6H AS KS JS 9S 7S JD 8D 7D 6D AC QC JC 9C 8C 7C`,
	`AH KH TH 9H AS KS TS 8S 6S QD 6D AC KC JC TC 9C 8C 6C`,
	`QH JH 8H 7H 6H QS JS 9S 7S AD KD JD TD 9D 8D 7D QC 7C`,
	`AH TH 7H 6H JS TS 8S 7S 6S AD KD QD 9D 8D 7D KC 9C 6C`,
	`KH QH JH 9H 8H AS KS QS 9S JD TD 6D AC QC JC TC 8C 7C`,
	`AH KH 9H 7H 6H AS KS TS 8S JD TD 8D 6D KC QC TC 9C 7C`,
	`QH JH TH 8H QS JS 9S 7S 6S AD KD QD 9D 7D AC JC 8C 6C`,
	`AH KH TH 6H AS KS QS TS 9S 8S AD QD TD 9D TC 8C 7C 6C`,
	`JH TH 9H 8H AS 7S 6S KD QD 9D 8D 6D AC KC QC JC 9C 8C`,
	`QH JH TH 8H TS 9S 7S KD JD TD 9D 7D 6D AC KC QC 7C 6C`,
	`AH KH QH JH 7H 6H QS TS 9S 7S QD JD 9D 7D 6D JC TC 9C`,
	`TH 9H 8H AS KS JS 8S 6S AD KD TD 8D AC KC QC 8C 7C 6C`,
	`KH 9H QS 9S AD KD QD JD TD 8D 7D 6D AC KC QC 8C 7C 6C`,
	`AH 9H TS 9S AD KD QD JD TD 8D 7D 6D AC QC JC 8C 7C 6C`,
	`AH 9H 7H 9S AD KD QD JD TD 7D 6D AC KC JC TC 8C 7C 6C`,
	`9H QS 9S 8S AD KD QD JD 8D 7D 6D KC QC JC TC 8C 7C 6C`,
	`9H KS 9S 6S AD KD QD JD TD 8D 7D 6D AC QC JC TC 7C 6C`,
	`JH 9H 7H 9S AD KD QD JD TD 8D 7D 6D KC JC TC 8C 7C 6C`,
	`9H JS 9S AD KD QD JD TD 8D 7D 6D AC KC JC TC 8C 7C 6C`,
	`9H QS JS 9S AD KD QD JD 8D 7D 6D AC KC QC JC 8C 7C 6C`,
	`TH 9H 9S AD QD JD TD 8D 7D 6D AC KC QC JC TC 8C 7C 6C`,
	`9H QS 9S AD KD QD JD TD 7D 6D AC KC QC JC TC 8C 7C 6C`,
	`9H TS 9S AD KD QD TD 8D 7D 6D AC KC QC JC TC 8C 7C 6C`,
	`9H 9S 8S AD KD QD JD TD 8D 7D 6D AC KC QC JC 8C 7C 6C`,
	`9H 9S AD KD QD JD TD 8D 7D 6D AC KC QC JC TC 8C 7C 6C`,
}

func expected(k int) string {
	var sb strings.Builder
	for i := 0; i < k; i++ {
		sb.WriteString(ans[i])
		sb.WriteString("\n\n")
	}
	return sb.String()
}

func runCase(bin string, k int) error {
	in := fmt.Sprintf("%d\n", k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if out.String() != expected(k) {
		return fmt.Errorf("output mismatch")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD3.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for k := 1; k <= len(ans); k++ {
		if err := runCase(bin, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", k, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
