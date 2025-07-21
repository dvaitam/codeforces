package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/big"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type test struct {
    input    string
    expected string
}

func toRoman(n int) string {
    vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
    syms := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
    var sb strings.Builder
    for i := 0; i < len(vals); i++ {
        for n >= vals[i] {
            sb.WriteString(syms[i])
            n -= vals[i]
        }
    }
    return sb.String()
}

func solve(input string) string {
    in := bufio.NewReader(strings.NewReader(input))
    var a int
    var bstr string
    if _, err := fmt.Fscan(in, &a, &bstr); err != nil {
        return ""
    }
    var cstr string
    fmt.Fscan(in, &cstr)
    num := new(big.Int)
    num.SetString(strings.TrimSpace(cstr), a)
    if bstr == "R" {
        val := int(num.Int64())
        return toRoman(val)
    }
    b, _ := strconv.Atoi(bstr)
    out := strings.ToUpper(num.Text(b))
    out = strings.TrimLeft(out, "0")
    if out == "" { out = "0" }
    return out
}

func randNumber(max int64) int64 { return rand.Int63n(max) }

func generateTests() []test {
    rand.Seed(44)
    var tests []test
    fixed := []string{
        "2 10\n1010\n",
        "10 R\n1990\n",
        "16 2\nFF\n",
    }
    for _, f := range fixed {
        tests = append(tests, test{f, solve(f)})
    }
    for len(tests) < 100 {
        a := rand.Intn(24)+2
        useRoman := rand.Intn(3)==0
        var bstr string
        if useRoman {
            bstr = "R"
        } else {
            b := rand.Intn(24)+2
            bstr = strconv.Itoa(b)
        }
        var value int64
        if useRoman {
            value = randNumber(300000)+1
        } else {
            value = randNumber(1_000_000_000)
        }
        cstr := strings.ToUpper(strconv.FormatInt(value, a))
        if rand.Intn(4)==0 {
            cstr = "0"+cstr
        }
        inp := fmt.Sprintf("%d %s\n%s\n", a, bstr, cstr)
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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

