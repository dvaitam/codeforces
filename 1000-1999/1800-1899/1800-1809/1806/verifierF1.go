package main

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

// Base64-encoded contents of testcasesF1.txt.
const testcasesF1B64 = "NCA5IDEgNiA3IDEgNQozIDUgMiAzIDUgMQozIDEwIDEgMTAgMTAgMTAKMiA4IDEgMSA1CjIgMiAxIDIgMQo1IDYgMiAxIDYgNSAyIDQKNSA5IDMgNiA2IDIgMSA3CjQgNSAyIDIgNCA1IDEKMyAxMCAyIDcgMyAzCjUgOCA0IDQgNSAyIDMgMQo0IDkgMSAyIDMgNSAzCjUgNyA0IDMgNSA2IDQgNAo1IDUgMyA1IDMgMyAzIDMKNSA5IDEgOCA5IDkgNCA3CjQgOSAzIDcgMSAxIDcKMiA2IDEgMSA0CjUgOSAzIDkgOCA5IDggOQo1IDkgMiAxIDYgNCA1IDcKMiA1IDEgMiA1CjUgMTAgMiA4IDYgOCAzIDgKNCA3IDIgNiA1IDYgNwozIDMgMiAzIDIgMgoyIDYgMSA1IDIKMyA0IDIgMiAxIDEKNSA4IDEgMSAzIDggMyA1CjIgNiAxIDIgNAo0IDQgMiAxIDEgMSAzCjQgOCAyIDMgMSA4IDcKNSAxMCAyIDIgNyA1IDggNwo1IDggMiA3IDQgNyA0IDUKMiA0IDEgMyAzCjQgMTAgMiA2IDQgMyA4CjMgOCAxIDggMyAyCjUgMTAgMyAzIDIgNiAzIDUKNCA5IDMgMiAyIDcgNgo0IDYgMiAxIDIgNiA2CjQgOSAxIDMgOSAzIDgKMyA0IDIgMSAxIDQKMiA4IDEgOCAzCjMgMyAxIDEgMyAxCjUgOSAxIDUgMiAyIDQgNQo1IDYgNCA2IDMgMyA2IDIKNSA3IDIgMiA3IDYgMSAyCjQgOSAzIDggMiA2IDkKNSA3IDQgNyAxIDIgMyAyCjUgNiAyIDEgMiA0IDQgMQo0IDEwIDEgMiA2IDMgNAozIDggMiA2IDggNQo0IDQgMyAyIDMgMiAxCjQgNyAxIDEgMSA0IDUKNSA4IDMgNiA0IDYgOCAzCjIgOSAxIDEgMgoyIDQgMSAxIDMKNCA5IDEgMiAxIDQgMwo0IDcgMyA1IDUgNiAzCjUgNiA0IDYgMSA1IDYgMgoyIDkgMSA3IDMKMyA0IDEgMyAzIDMKNSA3IDQgNCAzIDcgMyA1CjQgNyAzIDMgNCA2IDMKNSA3IDEgMiA3IDYgNCA3CjQgNyAyIDYgNiA3IDcKMyA1IDIgMiAyIDUKMyA0IDIgMiAyIDQKMiA5IDEgNyAzCjIgMiAxIDEgMQozIDUgMSAzIDMgMQoyIDggMSAzIDYKMyA2IDIgNSAzIDEKMyA5IDEgNCAzIDMKNCA2IDMgMyAyIDQgNgo0IDggMSAyIDcgMSAxCjMgOCAyIDQgNiA3CjIgNyAxIDQgNQo1IDcgMiA0IDcgNiAxIDMKNCA1IDEgNSAzIDQgMQozIDUgMiA0IDEgNAo0IDcgMiAxIDQgNyAxCjUgNSAyIDMgMyAzIDUgMgo1IDggMSA2IDMgNyAxIDUKMyA1IDEgNCA1IDQKMyA5IDIgMyA5IDYKMyA0IDEgMiAxIDIKNSAxMCAyIDEgNyA2IDUgNgo0IDcgMyA2IDEgMyA3CjUgOCAyIDggMSA3IDYgNAo0IDYgMSAyIDYgNiAzCjMgNCAyIDQgMiAzCjQgOSAzIDMgNyAxIDIKMyA5IDIgMiAyIDQKMiA1IDEgNCA1CjIgNCAxIDEgNAo0IDUgMyAxIDQgMSAxCjIgNSAxIDQgMgo0IDggMyAyIDMgNiAxCjQgOSAyIDkgNCA1IDEKNSA3IDEgMyAzIDcgNyAzCjIgNiAxIDMgMQoyIDggMSA3IDQKNSA5IDQgNSAxIDMgOCA5Cg=="

type testCase struct {
    n int
    m int
    k int
    arr []int64
}

// Embedded solver logic from 1806F1.go (placeholder algorithm).
func solve(tc testCase) string {
    sum := int64(0)
    for _, v := range tc.arr {
        sum += v
    }
    ans := sum - int64(tc.k)
    if ans < 0 {
        ans = 0
    }
    return strconv.FormatInt(ans, 10)
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesF1B64)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    sc.Split(bufio.ScanWords)
    cases := make([]testCase, 0)
    for {
        if !sc.Scan() {
            break
        }
        n, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("parse n: %v", err)
        }
        if !sc.Scan() {
            return nil, fmt.Errorf("case starting with n=%d missing m", n)
        }
        m, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("parse m: %v", err)
        }
        if !sc.Scan() {
            return nil, fmt.Errorf("case starting with n=%d missing k", n)
        }
        k, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("parse k: %v", err)
        }
        arr := make([]int64, n)
        for i := 0; i < n; i++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("case n=%d truncated at arr[%d]", n, i)
            }
            val, err := strconv.ParseInt(sc.Text(), 10, 64)
            if err != nil {
                return nil, fmt.Errorf("arr[%d]: %v", i, err)
            }
            arr[i] = val
        }
        cases = append(cases, testCase{n: n, m: m, k: k, arr: arr})
    }
    if len(cases) == 0 {
        return nil, fmt.Errorf("no testcases found")
    }
    return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
    if strings.HasSuffix(path, ".go") {
        tmp, err := os.CreateTemp("", "bin*")
        if err != nil {
            return "", nil, err
        }
        tmp.Close()
        if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
            os.Remove(tmp.Name())
            return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
        }
        return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
    }
    return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
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
        fmt.Println("usage: go run verifierF1.go /path/to/binary")
        os.Exit(1)
    }

    cases, err := parseTestcases()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    bin, cleanup, err := buildIfGo(os.Args[1])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer cleanup()

    for idx, tc := range cases {
        var input strings.Builder
        input.WriteString("1\n")
        input.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
        for i, v := range tc.arr {
            if i > 0 {
                input.WriteByte(' ')
            }
            input.WriteString(strconv.FormatInt(v, 10))
        }
        input.WriteByte('\n')

        want := solve(tc)
        got, err := runCandidate(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

