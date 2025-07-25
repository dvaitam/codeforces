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

func solveE(n int, arr []string) int {
	maxLen := 0
	for _, s := range arr {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	maxH := maxLen * 2
	pow := make([]uint64, maxH+1)
	const B uint64 = 1315423911
	pow[0] = 1
	for i := 1; i <= maxH; i++ {
		pow[i] = pow[i-1] * B
	}
	ans := 0
	for l := 0; l < n; l++ {
		wordSet := make(map[string]struct{})
		lengths := make(map[int]struct{})
		localMax := 0
		for r := l; r < n; r++ {
			w := arr[r]
			wordSet[w] = struct{}{}
			L := len(w)
			lengths[L] = struct{}{}
			if L > localMax {
				localMax = L
			}
			Ls := make([]int, 0, len(lengths))
			for x := range lengths {
				Ls = append(Ls, x)
			}
			hashMap := make(map[int]map[uint64]struct{})
			for s := range wordSet {
				var h uint64
				for i := 0; i < len(s); i++ {
					h = h*B + uint64(s[i])
				}
				l0 := len(s)
				m0, ok := hashMap[l0]
				if !ok {
					m0 = make(map[uint64]struct{})
					hashMap[l0] = m0
				}
				m0[h] = struct{}{}
			}
			stable := true
			for x := range wordSet {
				if !stable {
					break
				}
				for y := range wordSet {
					if !stable {
						break
					}
					S := x + y
					tot := len(S)
					hpre := make([]uint64, tot+1)
					for i := 0; i < tot; i++ {
						hpre[i+1] = hpre[i]*B + uint64(S[i])
					}
					dpP := make([]bool, tot+1)
					dpE := make([]bool, tot+1)
					dpP[0] = true
					for i := 0; i < tot; i++ {
						if !dpP[i] {
							continue
						}
						for _, L := range Ls {
							j := i + L
							if j <= tot {
								h := hpre[j] - hpre[i]*pow[L]
								if m0, ok := hashMap[L]; ok {
									if _, ok2 := m0[h]; ok2 {
										dpP[j] = true
									}
								}
							}
						}
					}
					dpE[tot] = true
					for i := tot - 1; i >= 0; i-- {
						for _, L := range Ls {
							j := i + L
							if j <= tot && dpE[j] {
								h := hpre[j] - hpre[i]*pow[L]
								if m0, ok := hashMap[L]; ok {
									if _, ok2 := m0[h]; ok2 {
										dpE[i] = true
										break
									}
								}
							}
						}
					}
					lx := len(x)
					for sft := 1; sft < tot; sft++ {
						if sft == lx {
							continue
						}
						if dpP[sft] && dpE[sft] {
							stable = false
							break
						}
					}
				}
			}
			if stable {
				ans++
			}
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(3) + 1
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(3) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(3))
		}
		arr[i] = string(b)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		sb.WriteString(arr[i])
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	ans := solveE(n, arr)
	return sb.String(), ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		in, expect := genCase(rng)
		outStr, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, in)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(outStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output %q\n", t, outStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", t, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
