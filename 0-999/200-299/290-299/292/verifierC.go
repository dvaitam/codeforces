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

func enumerate(digits []int) []string {
	must := 0
	for _, d := range digits {
		must |= 1 << d
	}
	var s [20]int
	var cur [4]int
	var ans []string
	var length int

	var part func(pos, cnt int)
	part = func(pos, cnt int) {
		if pos == length {
			if cnt == 4 {
				ans = append(ans, fmt.Sprintf("%d.%d.%d.%d", cur[0], cur[1], cur[2], cur[3]))
			}
			return
		}
		if cnt >= 4 {
			return
		}
		curval := 0
		for i := 0; i < 3; i++ {
			if pos+i >= length {
				break
			}
			curval = curval*10 + s[pos+i]
			if curval > 255 {
				break
			}
			if i > 0 && s[pos] == 0 {
				break
			}
			cur[cnt] = curval
			part(pos+i+1, cnt+1)
		}
	}

	var goRec func(pos, mask int)
	goRec = func(pos, mask int) {
		if pos >= 2 && pos <= 6 && mask == must {
			for iter := 0; iter < 2; iter++ {
				length = 2*pos - iter
				for i := 0; i < pos; i++ {
					s[length-1-i] = s[i]
				}
				part(0, 0)
			}
		}
		if pos == 6 {
			return
		}
		for d := 0; d < 10; d++ {
			if must&(1<<d) != 0 {
				s[pos] = d
				goRec(pos+1, mask|1<<d)
			}
		}
	}

	goRec(0, 0)
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	perm := rng.Perm(10)
	digits := make([]int, n)
	for i := 0; i < n; i++ {
		digits[i] = perm[i]
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, d := range digits {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", d))
	}
	sb.WriteByte('\n')

	ans := enumerate(digits)
	var out strings.Builder
	fmt.Fprintf(&out, "%d\n", len(ans))
	for _, ip := range ans {
		out.WriteString(ip)
		out.WriteByte('\n')
	}
	return sb.String(), strings.TrimSpace(out.String())
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
