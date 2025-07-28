package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(strs []string) string {
	n := len(strs)
	var freq [6][46]int
	var pref [6][6][46][46]int
	lens := make([]int, n)
	sums := make([]int, n)
	prefixes := make([][]int, n)
	for i := 0; i < n; i++ {
		s := strs[i]
		l := len(s)
		lens[i] = l
		prefix := make([]int, l+1)
		sum := 0
		for j := 0; j < l; j++ {
			d := int(s[j] - '0')
			sum += d
			prefix[j+1] = sum
			if j+1 < l {
				pref[l][j+1][prefix[j+1]][sum]++
			}
		}
		sums[i] = sum
		prefixes[i] = prefix
		freq[l][sum]++
	}
	ans := 0
	for idx := 0; idx < n; idx++ {
		lS := lens[idx]
		sumS := sums[idx]
		preS := prefixes[idx]
		for lenT := 1; lenT <= 5; lenT++ {
			L := lS + lenT
			if L%2 == 1 {
				continue
			}
			half := L / 2
			if half <= lS {
				d := 2*preS[half] - sumS
				if d >= lenT && d <= 9*lenT && d >= 0 && d < 46 {
					ans += freq[lenT][d]
				}
			} else {
				prefLen := half - lS
				if prefLen >= lenT {
					continue
				}
				maxPrefSum := 9 * prefLen
				if maxPrefSum > 45 {
					maxPrefSum = 45
				}
				for ps := 1; ps <= maxPrefSum; ps++ {
					sumT := sumS + 2*ps
					if sumT < lenT || sumT > 9*lenT || sumT >= 46 {
						continue
					}
					ans += pref[lenT][prefLen][ps][sumT]
				}
			}
		}
	}
	return fmt.Sprint(ans)
}

func genCase(rng *rand.Rand) []string {
	n := rng.Intn(5) + 1
	strs := make([]string, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(5) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('0' + rng.Intn(10))
		}
		strs[i] = string(b)
	}
	return strs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		strs := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(strs)))
		for _, s := range strs {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
		expect := solveC(strs)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
