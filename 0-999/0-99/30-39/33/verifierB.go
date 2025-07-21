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

const infB = 1000000000

func solveB(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	sLine, _ := reader.ReadString('\n')
	tLine, _ := reader.ReadString('\n')
	s := strings.TrimSpace(sLine)
	t := strings.TrimSpace(tLine)
	if len(s) != len(t) {
		return "-1\n"
	}
	var n int
	fmt.Fscan(reader, &n)
	d := make([][]int, 26)
	for i := 0; i < 26; i++ {
		d[i] = make([]int, 26)
		for j := 0; j < 26; j++ {
			if i == j {
				d[i][j] = 0
			} else {
				d[i][j] = infB
			}
		}
	}
	for i := 0; i < n; i++ {
		var uStr, vStr string
		var l int
		fmt.Fscan(reader, &uStr, &vStr, &l)
		u := int(uStr[0] - 'a')
		v := int(vStr[0] - 'a')
		if l < d[u][v] {
			d[u][v] = l
		}
	}
	for k := 0; k < 26; k++ {
		for i := 0; i < 26; i++ {
			for j := 0; j < 26; j++ {
				if d[i][k] < infB && d[k][j] < infB {
					if d[i][j] > d[i][k]+d[k][j] {
						d[i][j] = d[i][k] + d[k][j]
					}
				}
			}
		}
	}
	m := len(s)
	result := make([]byte, m)
	ans := 0
	for i := 0; i < m; i++ {
		si := int(s[i] - 'a')
		ti := int(t[i] - 'a')
		if d[si][ti] > d[ti][si] {
			si, ti = ti, si
		}
		best := d[si][ti]
		for j := 0; j < 26; j++ {
			if d[si][j] < infB && d[ti][j] < infB {
				if d[si][j]+d[ti][j] < best {
					best = d[si][j] + d[ti][j]
				}
			}
		}
		if best >= infB {
			return "-1\n"
		}
		c := ti
		if best != d[si][ti] {
			for j := 0; j < 26; j++ {
				if d[si][j] < infB && d[ti][j] < infB && d[si][j]+d[ti][j] == best {
					c = j
					break
				}
			}
		}
		result[i] = byte('a' + c)
		ans += best
	}
	return fmt.Sprintf("%d\n%s\n", ans, string(result))
}

func genTestB() (string, string) {
	l := rand.Intn(5) + 1
	letters := []rune("abcdef")
	var s, t strings.Builder
	for i := 0; i < l; i++ {
		s.WriteRune(letters[rand.Intn(len(letters))])
		t.WriteRune(letters[rand.Intn(len(letters))])
	}
	n := rand.Intn(5)
	var sb strings.Builder
	sb.WriteString(s.String() + "\n")
	sb.WriteString(t.String() + "\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		a := letters[rand.Intn(len(letters))]
		b := letters[rand.Intn(len(letters))]
		w := rand.Intn(5) + 1
		fmt.Fprintf(&sb, "%c %c %d\n", a, b, w)
	}
	input := sb.String()
	out := solveB(input)
	return input, out
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTestB()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
