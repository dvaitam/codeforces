package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var arr = []string{
	"9C 9D 6S 7S 8S TS JS QS KS AS 6H 7H 8H TH JH QH KH AH",
	"6C 7C 8C TC JC QC KC AC 6D 7D 8D TD JD QD KD AD 9S 9H",
	"9C 9D TD 6S 7S 8S JS QS KS AS 6H 7H 8H TH JH QH KH AH",
	"6C 7C 8C TC JC QC KC AC 6D 7D 8D JD QD KD AD 9S TS 9H",
	"9C 9D JD 6S 7S 8S JS QS KS AS 6H 7H 8H TH JH QH KH AH",
	"6C 7C 8C TC JC QC KC AC 6D 7D 8D TD QD KD AD 9S TS 9H",
	"9C 9D TD JD 6S 7S 8S QS KS AS 6H 7H 8H TH JH QH KH AH",
	"6C 7C 8C TC JC QC KC AC 6D 7D 8D QD KD AD 9S TS JS 9H",
	"9C 9D QD KD AD 8S TS QS KS AS 6H 7H 8H TH JH QH KH AH",
	"6C 7C 8C TC JC QC KC AC 6D 7D 8D TD JD 6S 7S 9S JS 9H",
	"9C 9D JD 6S 7S 8S TS QS KS AS 6H 7H 8H TH JH QH KH AH",
	"6C 7C 8C TC JC QC KC AC 6D 7D 8D TD QD KD AD 9S JS 9H",
	"9C 9D KD AD 7S 8S TS QS KS AS 6H 7H 8H TH JH QH KH AH",
	"6C 7C 8C TC JC QC KC AC 6D 7D 8D TD JD QD 6S 9S JS 9H",
	"9C 9D QD 6S 7S 8S TS JS QS KS 6H 7H 8H TH JH QH KH AH",
	"6C 7C 8C TC JC QC KC AC 6D 7D 8D TD JD KD AD 9S AS 9H",
	"9C 9D TD JD QD 6S 7S 8S KS AS 6H 7H 8H TH JH QH KH AH",
	"6C 7C 8C TC JC QC KC AC 6D 7D 8D KD AD 9S TS JS QS 9H",
	"9D 6S 7S 8S TS JS QS KS AS 6H 7H 8H 9H TH JH QH KH AH",
	"6C 7C 8C 9C TC JC QC KC AC 6D 7D 8D TD JD QD KD AD 9S",
	"8C TC 6D 7D 8D JD QD KD AD 6S 7S 9S JS QS KS 6H 8H TH",
	"6C 7C 9C JC QC KC AC 9D TD 8S TS AS 7H 9H JH QH KH AH",
	"JC QC KC AC 6D 7D 8D 9D TD JS QS KS AS 6H 7H 8H 9H TH",
	"6C 7C 8C 9C TC JD QD KD AD 6S 7S 8S 9S TS JH QH KH AH",
	"8C 9C TC KC AC 6D 7D 8D TD JD QD KD AD 8S JS QS 6H 9H",
	"6C 7C JC QC 9D 6S 7S 9S TS KS AS 7H 8H TH JH QH KH AH",
	"JC QC KC AC 6D 7D 8D 9D TD QS KS AS 6H 7H 8H 9H TH JH",
	"6C 7C 8C 9C TC JD QD KD AD 6S 7S 8S 9S TS JS QH KH AH",
	"7C 9C 6D 7D 8D TD JD QD KD AD 6S 8S JS QS 9H JH KH AH",
	"6C 8C TC JC QC KC AC 9D 7S 9S TS KS AS 6H 7H 8H TH QH",
	"QC KC AC 6D 7D 8D 9D TD JD QS KS AS 6H 7H 8H 9H TH JH",
	"6C 7C 8C 9C TC JC QD KD AD 6S 7S 8S 9S TS JS QH KH AH",
	"6C 8C 9C QC 7D 8D 9D JD QD AD 6S 7S 8S QS KS AS 8H JH",
	"7C TC JC KC AC 6D TD KD 9S TS JS 6H 7H 9H TH QH KH AH",
	"JD QD KD AD 6S 7S 8S 9S TS 6H 7H 8H 9H TH JH QH KH AH",
	"6C 7C 8C 9C TC JC QC KC AC 6D 7D 8D 9D TD JS QS KS AS",
	"6C 7C 8C JC 8D TD QD AD 6S 8S 9S JS KS AS 6H TH KH AH",
	"9C TC QC KC AC 6D 7D 9D JD KD 7S TS QS 7H 8H 9H JH QH",
	"QD KD AD 6S 7S 8S 9S TS JS 6H 7H 8H 9H TH JH QH KH AH",
	"6C 7C 8C 9C TC JC QC KC AC 6D 7D 8D 9D TD JD QS KS AS",
	"6C 7C 8C JC KC 8D 9D TD JD AD 7S JS KS AS 7H 9H KH AH",
	"9C TC QC AC 6D 7D QD KD 6S 8S 9S TS QS 6H 8H TH JH QH",
	"KD AD 6S 7S 8S 9S TS JS QS 6H 7H 8H 9H TH JH QH KH AH",
	"6C 7C 8C 9C TC JC QC KC AC 6D 7D 8D 9D TD JD QD KS AS",
	"6C 9C TC QC KC 8D 9D QD AD 6S 7S 8S 9S JS QS AS JH AH",
	"7C 8C JC AC 6D 7D TD JD KD TS KS 6H 7H 8H 9H TH QH KH",
	"6S 7S 8S 9S TS JS QS KS AS 6H 7H 8H 9H TH JH QH KH AH",
	"6C 7C 8C 9C TC JC QC KC AC 6D 7D 8D 9D TD JD QD KD AD",
	"6C 8C JC QC AC 6D 8D TD JD KD AD 7S 9S QS KS 6H 7H 9H",
}

func expected(k int) string {
	var sb strings.Builder
	for i := 0; i < k; i++ {
		idx := 2 * i
		sb.WriteString(arr[idx])
		sb.WriteByte('\n')
		sb.WriteString(arr[idx+1])
		sb.WriteByte('\n')
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
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for k := 1; k <= 13; k++ {
		if err := runCase(bin, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", k, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
