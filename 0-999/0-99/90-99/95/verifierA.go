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
	"unicode"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randString(rng *rand.Rand, n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	var b bytes.Buffer
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		fmt.Fprintln(&b, randString(rng, rng.Intn(5)+1))
	}
	fmt.Fprintln(&b, randString(rng, rng.Intn(10)+1))
	letter := 'a' + rng.Intn(26)
	fmt.Fprintf(&b, "%c\n", letter)
	return b.String()
}

// solveReference implements the logic correctly:
// 1. Case-insensitive match for forbidden strings.
// 2. If a character is covered, it MUST be replaced.
// 3. Priority: Maximize lucky letters, then minimize lexicographically.
//    If original != lucky, replace with lucky.
//    If original == lucky, replace with 'a' (or 'b' if 'a' == lucky), respecting case.
func solveReference(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	scan := func() string {
		if scanner.Scan() {
			return scanner.Text()
		}
		return ""
	}

	nStr := scan()
	if nStr == "" {
		return ""
	}
	var n int
	fmt.Sscan(nStr, &n)

	forbidden := make([]string, n)
	for i := 0; i < n; i++ {
		forbidden[i] = scan()
	}

	w := scan()
	luckyStr := scan()
	if len(luckyStr) == 0 {
		return ""
	}
	luckyRune := unicode.ToLower([]rune(luckyStr)[0])

	wLen := len(w)
	covered := make([]bool, wLen)

	for _, s := range forbidden {
		sLen := len(s)
		if sLen > wLen {
			continue
		}
		for i := 0; i <= wLen-sLen; i++ {
			if strings.EqualFold(w[i:i+sLen], s) {
				for j := 0; j < sLen; j++ {
					covered[i+j] = true
				}
			}
		}
	}

	var result strings.Builder
	result.Grow(wLen)

	for i := 0; i < wLen; i++ {
		originalChar := rune(w[i])
		if !covered[i] {
			result.WriteRune(originalChar)
			continue
		}

		isUpper := unicode.IsUpper(originalChar)
		var targetLucky rune
		if isUpper {
			targetLucky = unicode.ToUpper(luckyRune)
		} else {
			targetLucky = luckyRune
		}

		if originalChar != targetLucky {
			result.WriteRune(targetLucky)
		} else {
			var replacement rune
			if isUpper {
				replacement = 'A'
			} else {
				replacement = 'a'
			}
			if replacement == originalChar {
				replacement++
			}
			result.WriteRune(replacement)
		}
	}
	return result.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierA.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		candOut, cErr := runBinary(candidate, input)
		
		refOut := solveReference(input)
		
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\ninput:\n%s", t+1, cErr, input)
			os.Exit(1)
		}
		if candOut != refOut {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:%s\nactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
