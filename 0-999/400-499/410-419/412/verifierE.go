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

func expectedE(s string) string {
	n := len(s)
	L := make([]int, n)
	PL := make([]int, n)
	isLetter := func(c byte) bool { return c >= 'a' && c <= 'z' }
	isDigit := func(c byte) bool { return c >= '0' && c <= '9' }
	isLocal := func(c byte) bool { return isLetter(c) || isDigit(c) || c == '_' }
	isDomain := func(c byte) bool { return isLetter(c) || isDigit(c) }
	for i := 0; i < n; i++ {
		if isLocal(s[i]) {
			if i > 0 {
				L[i] = L[i-1] + 1
			} else {
				L[i] = 1
			}
		}
		if isLetter(s[i]) {
			if i > 0 {
				PL[i] = PL[i-1] + 1
			} else {
				PL[i] = 1
			}
		} else if i > 0 {
			PL[i] = PL[i-1]
		}
	}
	R1 := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		if isDomain(s[i]) {
			if i+1 < n {
				R1[i] = R1[i+1] + 1
			} else {
				R1[i] = 1
			}
		}
	}
	R2 := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		if isLetter(s[i]) {
			if i+1 < n {
				R2[i] = R2[i+1] + 1
			} else {
				R2[i] = 1
			}
		}
	}
	DotValue := make([]int64, n)
	for i := 0; i < n-1; i++ {
		if s[i] == '.' {
			DotValue[i] = int64(R2[i+1])
		}
	}
	PDOV := make([]int64, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			PDOV[i] = PDOV[i-1] + DotValue[i]
		} else {
			PDOV[i] = DotValue[i]
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		if s[i] != '@' {
			continue
		}
		if i == 0 || i+1 >= n {
			continue
		}
		leftLen := L[i-1]
		if leftLen < 1 {
			continue
		}
		start := i - leftLen
		end := i - 1
		var letters int
		if start > 0 {
			letters = PL[end] - PL[start-1]
		} else {
			letters = PL[end]
		}
		if letters == 0 {
			continue
		}
		domMax := 0
		if i+1 < n {
			domMax = R1[i+1]
		}
		if domMax < 1 {
			continue
		}
		dLow := i + 2
		if dLow >= n {
			continue
		}
		dHigh := i + 1 + domMax
		if dHigh >= n {
			dHigh = n - 1
		}
		var rightSum int64
		if dLow > 0 {
			rightSum = PDOV[dHigh] - PDOV[dLow-1]
		} else {
			rightSum = PDOV[dHigh]
		}
		if rightSum <= 0 {
			continue
		}
		ans += int64(letters) * rightSum
	}
	return fmt.Sprint(ans)
}

func runCase(bin string, s string) error {
	input := s + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedE(s)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func genCase(rng *rand.Rand) string {
	l := rng.Intn(20) + 1
	alphabet := "abcdefghijklmnopqrstuvwxyz0123456789_@."
	var sb strings.Builder
	for i := 0; i < l; i++ {
		sb.WriteByte(alphabet[rng.Intn(len(alphabet))])
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := genCase(rng)
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
