package main

import (
    "fmt"
    "math/big"
    "os"
    "os/exec"
    "strings"
)

type solver struct{}

func (solver) solve(input string) (string, error) {
    reader := strings.NewReader(strings.TrimSpace(input))
    var s string
    if _, err := fmt.Fscan(reader, &s); err != nil {
        return "", err
    }

    n := len(s)
    target := make([]int, n)
    for i, c := range s {
        if c == 'A' {
            target[i] = 0
        } else {
            target[i] = 1
        }
    }

    total := big.NewInt(0)
    for x0 := 0; x0 <= 1; x0++ {
        for x1 := 0; x1 <= 1; x1++ {
            prevDP := [2][2]*big.Int{}
            for a := 0; a < 2; a++ {
                for b := 0; b < 2; b++ {
                    prevDP[a][b] = big.NewInt(0)
                }
            }
            prevDP[x0][x1].SetInt64(1)

            for pos := 1; pos <= n-2; pos++ {
                currDP := [2][2]*big.Int{}
                for a := 0; a < 2; a++ {
                    for b := 0; b < 2; b++ {
                        currDP[a][b] = big.NewInt(0)
                    }
                }

                for prev := 0; prev < 2; prev++ {
                    for curr := 0; curr < 2; curr++ {
                        cnt := prevDP[prev][curr]
                        if cnt.Sign() == 0 {
                            continue
                        }
                        for next := 0; next < 2; next++ {
                            tPrev := 0
                            if prev == 0 && curr == 1 {
                                tPrev = 1
                            }
                            tCur := 0
                            if curr == 0 && next == 1 {
                                tCur = 1
                            }
                            if curr^tPrev^tCur != target[pos] {
                                continue
                            }
                            currDP[curr][next].Add(currDP[curr][next], cnt)
                        }
                    }
                }
                prevDP = currDP
            }

            for prev := 0; prev < 2; prev++ {
                for curr := 0; curr < 2; curr++ {
                    cnt := prevDP[prev][curr]
                    if cnt.Sign() == 0 {
                        continue
                    }
                    tPrev := 0
                    if prev == 0 && curr == 1 {
                        tPrev = 1
                    }
                    tLast := 0
                    if curr == 0 && x0 == 1 {
                        tLast = 1
                    }
                    t0 := 0
                    if x0 == 0 && x1 == 1 {
                        t0 = 1
                    }
                    if curr^tPrev^tLast != target[n-1] {
                        continue
                    }
                    if x0^tLast^t0 != target[0] {
                        continue
                    }
                    total.Add(total, cnt)
                }
            }
        }
    }

    return total.String(), nil
}

var testcases = []string{
    "ABBB",
    "BBAABAAAA",
    "AABAAAB",
    "AAAAB",
    "BAA",
    "ABBABBAAA",
    "AAA",
    "BBBBBBAAA",
    "AAABAB",
    "BAABB",
    "BBAABABB",
    "BBABBBBBBB",
    "AABABAA",
    "BBBBBAABBA",
    "AAB",
    "AAA",
    "BAAAABAB",
    "BBABAAAABB",
    "BABAAB",
    "BAAAB",
    "AAB",
    "BAABABB",
    "BBBABAABBB",
    "ABAAAABBA",
    "BBABAAAAA",
    "BBBA",
    "AABAAABA",
    "BAB",
    "BABAAAAAB",
    "BBABAAA",
    "BAA",
    "ABABA",
    "AABABBBB",
    "AABBBBBBAA",
    "BBAAA",
    "ABAAB",
    "BBBBBBA",
    "BBB",
    "AAB",
    "BBAA",
    "ABAB",
    "BBBBB",
    "ABB",
    "ABABB",
    "AABABBBBAB",
    "ABBBAABAB",
    "ABABAABAAB",
    "AABAABABBB",
    "BBBAAA",
    "ABABBB",
    "BBAAABABB",
    "BBBAAAA",
    "BAA",
    "AAB",
    "BAAB",
    "AAAAAAB",
    "BABAAABAAA",
    "BABABAAA",
    "ABAABA",
    "AABBBABA",
    "BBBB",
    "ABBAABA",
    "ABAB",
    "AAABBBB",
    "BBB",
    "AABBB",
    "BABAAB",
    "BAAA",
    "AABBAA",
    "ABAAAA",
    "BBBBAB",
    "AAABAAAA",
    "BBABBB",
    "BAABBAAA",
    "AABBABAAB",
    "ABAABBBB",
    "ABA",
    "AABBAAB",
    "AAAABBA",
    "BBA",
    "AABAAAAB",
    "BBAAAAB",
    "AAAAAA",
    "ABAB",
    "BBAAABABBB",
    "BABAAAAAA",
    "BABBABBABA",
    "BABBABABAB",
    "ABBBAAABAB",
    "ABBBBAB",
    "AAABBAA",
    "BBABBAAAB",
    "BABBBBBAAB",
    "AABBA",
    "ABABABBAB",
    "AAA",
    "BBAABBBA",
    "BAAAABABAB",
    "ABAA",
    "ABAA",
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    s := solver{}

    for idx, tc := range testcases {
        input := strings.TrimSpace(tc) + "\n"
        expected, err := s.solve(input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: oracle error: %v\n", idx+1, err)
            os.Exit(1)
        }

        cmd := exec.Command(bin)
        cmd.Stdin = strings.NewReader(input)
        out, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }

        got := strings.TrimSpace(string(out))
        if got != expected {
            fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
            os.Exit(1)
        }
    }

    fmt.Printf("All %d tests passed\n", len(testcases))
}
