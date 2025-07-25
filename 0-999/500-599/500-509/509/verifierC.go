package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveSequence(b []int) []string {
	prev := "0"
	res := make([]string, len(b))
	for i, bi := range b {
		L := len(prev)
		found := ""
		ps := make([]int, L+1)
		for j := 0; j < L; j++ {
			ps[j+1] = ps[j] + int(prev[j]-'0')
		}
		for pos := L - 1; pos >= 0 && found == ""; pos-- {
			orig := int(prev[pos] - '0')
			for d := orig + 1; d <= 9; d++ {
				sumPre := ps[pos] + d
				rem := bi - sumPre
				remPos := L - pos - 1
				if rem < 0 || rem > 9*remPos {
					continue
				}
				t := make([]byte, L)
				copy(t, []byte(prev))
				t[pos] = byte('0' + d)
				for k := pos + 1; k < L; k++ {
					after := L - k - 1
					x := rem - 9*after
					if x < 0 {
						x = 0
					}
					if x > 9 {
						x = 9
					}
					t[k] = byte('0' + x)
					rem -= x
				}
				found = string(t)
				break
			}
		}
		if found != "" {
			prev = found
		} else {
			needLen := bi / 9
			if bi%9 != 0 {
				needLen++
			}
			newLen := needLen
			if newLen <= len(prev) {
				newLen = len(prev) + 1
			}
			rem := bi
			t := make([]byte, newLen)
			for k := 0; k < newLen; k++ {
				remPos := newLen - k - 1
				minD := 0
				if k == 0 {
					minD = 1
				}
				x := rem - 9*remPos
				if x < minD {
					x = minD
				}
				if x > 9 {
					x = 9
				}
				t[k] = byte('0' + x)
				rem -= x
			}
			prev = string(t)
		}
		res[i] = prev
	}
	return res
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(5) + 1
	b := make([]int, n)
	for i := range b {
		b[i] = rng.Intn(27) + 1
	}
	return b
}

func runCase(bin string, b []int) error {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(b)))
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Fields(strings.TrimSpace(out.String()))
	expect := solveSequence(b)
	if len(gotLines) != len(expect) {
		return fmt.Errorf("expected %d lines got %d", len(expect), len(gotLines))
	}
	for i, ex := range expect {
		if gotLines[i] != ex {
			return fmt.Errorf("line %d expected %s got %s", i+1, ex, gotLines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		b := generateCase(rng)
		if err := runCase(bin, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
