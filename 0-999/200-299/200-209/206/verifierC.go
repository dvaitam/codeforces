package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// reference implementation from 206C1.go
func solveC(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	s1 := []string{""}
	count1 := map[string]int{"": 1}
	t2 := []string{""}
	parent2 := []int{-1}
	count2 := map[string]int{"": 1}
	total := int64(1)
	var out bytes.Buffer
	for i := 0; i < n; i++ {
		var t, v int
		var cs string
		fmt.Fscan(reader, &t, &v, &cs)
		c := cs
		v--
		if t == 1 {
			s := c + s1[v]
			s1 = append(s1, s)
			count1[s]++
			if cnt2, ok := count2[s]; ok {
				total += int64(cnt2)
			}
		} else {
			newID := len(t2)
			parent2 = append(parent2, v)
			tstr := t2[v] + c
			t2 = append(t2, tstr)
			for j := newID; j >= 0; j = parent2[j] {
				prefixLen := len(t2[j])
				s := tstr[prefixLen:]
				count2[s]++
				if cnt1, ok := count1[s]; ok {
					total += int64(cnt1)
				}
				if j == 0 {
					break
				}
			}
		}
		fmt.Fprintln(&out, total)
	}
	return out.String()
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	lines := make([]string, n+1)
	lines[0] = fmt.Sprintf("%d", n)
	size1, size2 := 1, 1
	for i := 0; i < n; i++ {
		t := rng.Intn(2) + 1
		var v int
		if t == 1 {
			v = rng.Intn(size1) + 1
			size1++
		} else {
			v = rng.Intn(size2) + 1
			size2++
		}
		c := string('a' + rune(rng.Intn(26)))
		lines[i+1] = fmt.Sprintf("%d %d %s", t, v, c)
	}
	return strings.Join(lines, "\n") + "\n"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		expect := strings.TrimSpace(solveC(tc))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, strings.TrimSpace(got), tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
