package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases (string and expected per line, separated by space).
const embeddedTestcases = `eebc No
deaeadcebb No
eeddbbbe No
aabeaca No
deddd No
dbcaabdbcd No
dedec No
edebcaceb No
eeeabe No
caadd No
ca No
bacddaa No
eadececeba No
aaaee No
b Yes
cecbacc No
bddded No
eaeecdbcdc No
cecadecad No
ebacdccecd No
e Yes
a Yes
cdceec No
cbc Yes
eccdaa No
bcebcbcbda No
ec No
bdbacb No
dcbaaebceb No
caece No
dce No
dcdcd Yes
dadbbadeed No
badececba No
cabaaebdea No
d Yes
be No
baeed No
e Yes
cb No
edacb No
aeab No
cbad No
dacbceeeda No
caabaaaa No
aaeedcbc No
cd No
ecccbcd No
be No
d Yes
eb No
c Yes
eedaedac No
cdddabbe No
eadbd No
acc No
cadaedace No
aeadbbdae No
ea No
bacaaad No
ecaae No
ebaeaeaec No
babbbdedcc No
dcedadebdb No
eeedbdb No
add No
debbcbbee No
becdee No
cbcacddbbe No
bcdbdd No
bdeeadadad No
babc No
bcbb No
acbacbdaaa No
cc Yes
c Yes
ecaaccdd No
abeddbec No
ca No
adecaec No
dccace No
eadecaceb No
bcd No
ae No
ced No
cbddcbaab No
eedcaccda No
ade No
beccd No
eacdabce No
ae No
bb Yes
ddbeebdbee No
ebb No
ccadd No
ceecdec No
aebaabdb No
dbbebaeda No
cbbdaeaace No
eade No`

// Embedded reference solution (from 784D.go).
func solve(s string) string {
	ok := true
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			ok = false
			break
		}
	}
	if ok {
		return "Yes"
	}
	return "No"
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	total := 0
	passed := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Println("invalid test case:", line)
			continue
		}
		s := parts[0]
		input := s + "\n"
		expected := solve(s)
		total++
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Case %d: %v\n", total, err)
			fmt.Printf("Output: %s\n", got)
			continue
		}
		if got == expected {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", total, expected, got)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%d/%d cases passed\n", passed, total)
}
