package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func clean(s string) string {
    var b strings.Builder
    for _, r := range s {
        switch r {
        case '-', ';', '_':
        default:
            if 'A' <= r && r <= 'Z' {
                b.WriteRune(r + ('a' - 'A'))
            } else {
                b.WriteRune(r)
            }
        }
    }
    return b.String()
}

type test struct {
    input    string
    expected string
}

func solve(input string) string {
    lines := strings.Split(strings.TrimSpace(input), "\n")
    if len(lines) < 4 {
        return ""
    }
    orig := []string{strings.TrimSpace(lines[0]), strings.TrimSpace(lines[1]), strings.TrimSpace(lines[2])}
    t := make([]string, 3)
    for i, s := range orig {
        t[i] = clean(s)
    }
    perms := [][3]int{{0,1,2},{0,2,1},{1,0,2},{1,2,0},{2,0,1},{2,1,0}}
    valid := map[string]struct{}{}
    for _, p := range perms {
        valid[t[p[0]]+t[p[1]]+t[p[2]]] = struct{}{}
    }
    n, _ := strconv.Atoi(strings.TrimSpace(lines[3]))
    var out strings.Builder
    for i:=0;i<n;i++ {
        ans := clean(strings.TrimSpace(lines[4+i]))
        if _, ok := valid[ans]; ok {
            out.WriteString("ACC\n")
        } else {
            out.WriteString("WA\n")
        }
    }
    return strings.TrimRight(out.String(), "\n")
}

func generateRandomString(maxLen int) string {
    letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_";
    l := rand.Intn(maxLen) + 1
    var b strings.Builder
    for i := 0; i < l; i++ {
        b.WriteByte(letters[rand.Intn(len(letters))])
    }
    return b.String()
}

func generateTests() []test {
    rand.Seed(43)
    var tests []test
    fixed := []string{
        "a\nb\nc\n1\nabc\n",
        "A-B\nC;D\n_E\n2\nabcd\nXYZ\n",
    }
    for _, f := range fixed {
        tests = append(tests, test{f, solve(f)})
    }
    for len(tests) < 100 {
        s1 := generateRandomString(5)
        s2 := generateRandomString(5)
        s3 := generateRandomString(5)
        cleaned := []string{clean(s1), clean(s2), clean(s3)}
        perms := [][3]int{{0,1,2},{0,2,1},{1,0,2},{1,2,0},{2,0,1},{2,1,0}}
        var answers []string
        for _, p := range perms {
            answers = append(answers, cleaned[p[0]]+cleaned[p[1]]+cleaned[p[2]])
        }
        n := rand.Intn(5) + 1
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%s\n%s\n%s\n%d\n", s1, s2, s3, n))
        for i:=0;i<n;i++ {
            if rand.Intn(2) == 0 {
                // correct answer maybe with random signs/case
                ans := answers[rand.Intn(len(answers))]
                // add random signs and random case
                var sbAns strings.Builder
                for _, ch := range ans {
                    if rand.Intn(3)==0 {
                        sbAns.WriteByte('-')
                    }
                    if rand.Intn(3)==0 {
                        sbAns.WriteByte(';')
                    }
                    if rand.Intn(3)==0 {
                        sbAns.WriteByte('_')
                    }
                    if rand.Intn(2)==0 {
                        sbAns.WriteByte(byte(ch))
                    } else {
                        if 'a' <= ch && ch <= 'z' {
                            sbAns.WriteByte(byte(ch - ('a'-'A')))
                        } else {
                            sbAns.WriteByte(byte(ch))
                        }
                    }
                }
                sb.WriteString(sbAns.String())
            } else {
                // incorrect string
                sb.WriteString(generateRandomString(6))
            }
            sb.WriteByte('\n')
        }
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for i, t := range tests {
        got, err := runBinary(bin, t.input)
        if err != nil {
            fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if got != strings.TrimSpace(t.expected) {
            fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

