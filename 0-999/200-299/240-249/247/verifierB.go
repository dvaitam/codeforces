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

func expandIPv6(s string) string {
	var blocks []string
	if strings.Contains(s, "::") {
		parts := strings.SplitN(s, "::", 2)
		var left, right []string
		if parts[0] != "" {
			left = strings.Split(parts[0], ":")
		}
		if parts[1] != "" {
			right = strings.Split(parts[1], ":")
		}
		missing := 8 - (len(left) + len(right))
		blocks = make([]string, 0, 8)
		blocks = append(blocks, left...)
		for i := 0; i < missing; i++ {
			blocks = append(blocks, "0")
		}
		blocks = append(blocks, right...)
	} else {
		blocks = strings.Split(s, ":")
	}
	for i, b := range blocks {
		if len(b) < 4 {
			blocks[i] = strings.Repeat("0", 4-len(b)) + b
		}
	}
	return strings.Join(blocks, ":")
}

func solveB(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return ""
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(r, &s)
		sb.WriteString(expandIPv6(s))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateCaseB(rng *rand.Rand) string {
	num := rng.Intn(5) + 1 // 1..5 addresses
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", num))
	for i := 0; i < num; i++ {
		blocks := make([]uint16, 8)
		for j := range blocks {
			blocks[j] = uint16(rng.Intn(65536))
		}
		// choose run of zeros to compress with probability
		start := -1
		length := 0
		if rng.Intn(2) == 0 {
			// find longest run of zeros
			l := 0
			bestL := 0
			bestStart := -1
			for j := 0; j < 8; j++ {
				if blocks[j] == 0 {
					l++
					if l > bestL {
						bestL = l
						bestStart = j - l + 1
					}
				} else {
					l = 0
				}
			}
			if bestL > 1 {
				start = bestStart
				length = bestL
			}
		}
		var short []string
		for j := 0; j < 8; {
			if j == start {
				short = append(short, "") // placeholder for ::
				j += length
				continue
			}
			block := fmt.Sprintf("%x", blocks[j])
			short = append(short, strings.TrimLeft(block, "0"))
			if short[len(short)-1] == "" {
				short[len(short)-1] = "0"
			}
			j++
		}
		addr := strings.Join(short, ":")
		if start != -1 {
			if start == 0 {
				addr = "::" + strings.Join(short[1:], ":")
			} else if start+length >= 8 {
				addr = strings.Join(short[:start], ":") + "::"
			} else {
				addr = strings.Join(short[:start], ":") + "::" + strings.Join(short[start+1:], ":")
			}
		}
		// remove empty edges
		addr = strings.ReplaceAll(addr, ":::", "::")
		if strings.HasPrefix(addr, ":") && !strings.HasPrefix(addr, "::") {
			addr = "0" + addr
		}
		if strings.HasSuffix(addr, ":") && !strings.HasSuffix(addr, "::") {
			addr = addr + "0"
		}
		sb.WriteString(addr)
		if i+1 < num {
			sb.WriteByte('\n')
		} else {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseB(rng)
	}
	for i, tc := range cases {
		expect := solveB(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
