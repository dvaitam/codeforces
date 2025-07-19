package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
	var cnt [10]int
	for _, b := range data {
		if b >= '0' && b <= '9' {
			cnt[b-'0']++
		}
	}
	// Reserve digits 1, 8, 6, 9 for the middle segment
	cnt[1]--
	cnt[8]--
	cnt[6]--
	cnt[9]--
	// Precomputed permutations of 1,8,6,9 mapped by remainder
	perms := []string{"1869", "6198", "1896", "1689", "1986", "1968", "1698"}
	m := 0
	var sb strings.Builder
	// Build the prefix with remaining digits 9 to 1
	for d := 9; d >= 1; d-- {
		for cnt[d] > 0 {
			m = (m*3 + d) % 7
			sb.WriteByte(byte('0' + d))
			cnt[d]--
		}
	}
	// Append the chosen permutation to make divisible by 7
	sb.WriteString(perms[m])
	// Append all zeros at the end
	for cnt[0] > 0 {
		sb.WriteByte('0')
		cnt[0]--
	}
	// Output result
	fmt.Println(sb.String())
}
