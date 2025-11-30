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

const testcasesB64 = "MTAwCjggNSA3CjggNiAyIDQgNSAzIDEgNwo3IDIgNSA2IDMKNCAyIDEKMiAxIDQgMwo0IDIKNyAyIDYKNiA3IDIgMSA1IDQgMwoyIDUKNyAzIDUKMiA2IDEgMyA1IDcgNAoyIDcgMwo3IDcgNgo2IDIgNyAzIDQgNSAxCjUgNiAyIDQgMSA3IDMKMiAyIDEKMSAyCjEgMgo2IDQgMQoxIDQgNiAyIDMgNQoyIDEgNSA0CjYgMyAzCjMgNCAxIDUgNiAyCjUgMyA2CjcgMyAyCjYgMiAzIDQgNSA3IDEKNSAyIDcKNyA2IDQKMSAyIDYgNyA0IDMgNQoxIDYgNCA1IDIgMwo4IDMgMgo3IDUgNCAyIDggMyA2IDEKMyAyIDcKMiAyIDEKMiAxCjIgMQoyIDIgMQoxIDIKMSAyCjggMiAxCjYgMyA1IDQgOCA3IDIgMQo2IDgKNSAyIDMKMyA0IDUgMiAxCjUgMgozIDIgMgoyIDMgMQoxIDMKMyAzIDIKMSAyIDMKMyAxIDIKOCAyIDcKNCA2IDEgNSA4IDIgNyAzCjEgNAo1IDUgMwoxIDIgNSA0IDMKMiAxIDMgNCA1CjcgMiA1CjEgNyA1IDQgNiAyIDMKMiA1CjYgNiAyCjEgMiA1IDQgMyA2CjYgMSAzIDUgNCAyCjMgMiAxCjEgMiAzCjIgMwoyIDIgMQoyIDEKMiAxCjIgMiAxCjEgMgoyIDEKMiAyIDEKMSAyCjIgMQo4IDggNwozIDQgNiA4IDEgNSAyIDcKMiA2IDMgNCA1IDcgOCAxCjMgMiAxCjEgMyAyCjIgMQoyIDIgMQoyIDEKMiAxCjUgMiAzCjUgNCAxIDMgMgoxIDUKMiAyIDEKMSAyCjIgMQo2IDIgNQoxIDIgMyA1IDQgNgozIDQKMyAyIDIKMSAyIDMKMSAzCjMgMyAyCjEgMyAyCjIgMSAzCjggNiAxCjggNCA1IDcgNiAyIDMgMQo1IDEgOCA3IDQgMgo0IDIgMwo0IDMgMiAxCjIgMwo1IDIgMQoxIDMgMiA0IDUKMiAxCjUgMiA0CjUgMyAyIDEgNAozIDQKNyA2IDQKMiAzIDQgNSAxIDYgNwo2IDMgMiA1IDEgNAo1IDQgNAozIDEgNSAyIDQKNSAxIDIgMwo1IDMgMwoyIDQgMyA1IDEKMSA1IDMKMyAyIDIKMSAzIDIKMSAyCjggMyA1CjggNyA2IDUgNCAyIDEgMwo3IDQgNQozIDMgMgoxIDIgMwozIDIgMQo1IDMgMgo1IDIgNCAzIDEKNCA1IDEKMyAzIDEKMSAzIDIKMiAzIDEKOCAyIDUKMyA0IDUgNyA2IDIgOCAxCjEgNQo4IDMgNwoxIDIgNSA4IDQgNiA3IDMKMSA0IDMKNyA2IDUKNyAzIDUgMiA0IDYgMQoxIDIgNiAzIDcgNQo3IDMgMwo0IDIgMSA2IDMgNyA1CjUgMSA3CjggNSA0CjggMyA0IDUgNiAyIDEgNwo4IDcgNSAyIDMKNSAyIDIKNSAyIDEgNCAzCjQgMwoyIDIgMQoyIDEKMiAxCjYgNiAyCjQgMiA1IDEgMyA2CjIgMSAzIDQgNSA2CjYgNCAxCjQgMiAxIDMgNSA2CjEgNCAyIDUKOCA1IDQKOCAxIDUgMyA0IDIgNyA2CjggNiAzIDUgNwo2IDYgNQo0IDUgMyAxIDYgMgo2IDMgNCAxIDUgMgoyIDIgMQoyIDEKMiAxCjggMyA2CjMgMiA4IDUgNiA0IDcgMQo0IDEgOAoyIDIgMQoyIDEKMSAyCjQgMyAxCjMgMiA0IDEKMyA0IDEKMyAyIDEKMiAxIDMKMyAyCjQgMyAzCjIgNCAzIDEKNCAyIDEKMiAyIDEKMSAyCjEgMgozIDMgMgoxIDMgMgozIDIgMQo2IDIgMQo0IDIgMSAzIDYgNQoyIDUKOCAzIDMKOCAxIDYgNSAzIDcgMiA0CjYgMyA4CjUgMiAxCjEgMyA1IDQgMgozIDQKNCA0IDIKMyAyIDQgMQo0IDEgMyAyCjUgMiA0CjUgMiAxIDQgMwoxIDQKMiAyIDEKMSAyCjEgMgo4IDcgNAozIDUgMSA4IDYgNCAyIDcKNCA2IDUgMSA4IDIgMwo1IDMgMwo1IDIgMyA0IDEKNSAyIDQKMiAyIDEKMSAyCjIgMQo3IDUgMgo3IDIgNiAxIDMgNCA1CjQgNyA1IDYgMgo2IDYgNAo1IDYgMiAzIDEgNAo0IDYgNSAzIDIgMQo0IDQgMgoyIDQgMyAxCjQgMiAxIDMKNCA0IDEKNCAxIDMgMgoyIDMgMSA0CjQgMyAyCjQgMyAxIDIKNCAzIDEKNiA1IDIKMSAyIDQgNiAzIDUKNSA2IDIgMyA0CjcgNiAzCjMgNSAyIDcgMSA2IDQKNiA3IDMgNSAxIDQKNSA1IDEKMiA0IDMgMSA1CjEgMiA0IDUgMwo2IDIgMwo1IDQgMSAyIDMgNgo0IDUKNSA1IDMKMSAzIDUgMiA0CjIgMyA0IDEgNQozIDIgMQoxIDMgMgozIDIKOCA1IDQKNSAxIDMgNiA4IDIgNyA0CjcgNSAxIDQgMgo0IDIgMwozIDQgMSAyCjMgMQo4IDYgNwo2IDIgMSA4IDQgMyA3IDUKNSA4IDEgNyAyIDQKNCAyIDEKMyAyIDQgMQozIDIKNyAyIDIKNSA3IDMgNiAyIDQgMQo2IDcKNCA0IDEKNCAzIDIgMQoxIDMgMiA0CjYgMiAyCjUgNCA2IDIgMSAzCjQgMwo0IDQgMwoxIDQgMiAzCjIgNCAzIDEKMiAyIDEKMiAxCjEgMgo3IDUgMQozIDEgMiA1IDYgNCA3CjMgNyAxIDQgNQo1IDMgNAo0IDMgMiAxIDUKNCAzIDUKNSAzIDIKMSA0IDIgNSAzCjEgMiA1CjYgMyAzCjEgMiAzIDYgNSA0CjMgMiAxCjIgMiAxCjEgMgoxIDIKNCA0IDMKNCAyIDMgMQo0IDEgMiAzCjIgMiAxCjEgMgoxIDIK"

type testCase struct {
    n, m, d int
    p       []int
    a       []int
}

func buildIfGo(path string) (string, func(), error) {
    if strings.HasSuffix(path, ".go") {
        tmp, err := os.CreateTemp("", "solbin*")
        if err != nil {
            return "", nil, err
        }
        tmp.Close()
        if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
            os.Remove(tmp.Name())
            return "", nil, fmt.Errorf("build failed: %v\\n%s", err, out)
        }
        return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
    }
    return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesB64)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    sc.Split(bufio.ScanWords)
    if !sc.Scan() {
        return nil, fmt.Errorf("empty data")
    }
    t, err := strconv.Atoi(sc.Text())
    if err != nil {
        return nil, fmt.Errorf("parse t: %v", err)
    }
    cases := make([]testCase, 0, t)
    for i := 0; i < t; i++ {
        var n, m, d int
        for _, dst := range []*int{&n, &m, &d} {
            if !sc.Scan() {
                return nil, fmt.Errorf("truncated at case %d", i+1)
            }
            v, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d: %v", i+1, err)
            }
            *dst = v
        }
        p := make([]int, n)
        for j := 0; j < n; j++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("case %d truncated p", i+1)
            }
            val, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d p[%d]: %v", i+1, j, err)
            }
            p[j] = val
        }
        a := make([]int, m)
        for j := 0; j < m; j++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("case %d truncated a", i+1)
            }
            val, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d a[%d]: %v", i+1, j, err)
            }
            a[j] = val
        }
        cases = append(cases, testCase{n: n, m: m, d: d, p: p, a: a})
    }
    return cases, nil
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// Embedded solver logic from 1778B.go.
func solve(tc testCase) int {
    pos := make([]int, tc.n+1)
    for i, v := range tc.p {
        pos[v] = i + 1
    }
    ans := tc.n + 5
    for i := 0; i < tc.m-1; i++ {
        x := pos[tc.a[i]]
        y := pos[tc.a[i+1]]
        if x > y || y-x > tc.d {
            return 0
        }
        diff := y - x
        cur := diff
        delta := tc.d + 1 - diff
        if y+delta <= tc.n {
            cur = min(cur, delta)
        }
        ans = min(ans, cur)
    }
    if ans > tc.n {
        ans = 0
    }
    return ans
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierB.go /path/to/binary")
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
        fmt.Fprintf(&input, "1\n%d %d %d\n", tc.n, tc.m, tc.d)
        for i, v := range tc.p {
            if i > 0 {
                input.WriteByte(' ')
            }
            input.WriteString(strconv.Itoa(v))
        }
        input.WriteByte('\n')
        for i, v := range tc.a {
            if i > 0 {
                input.WriteByte(' ')
            }
            input.WriteString(strconv.Itoa(v))
        }
        input.WriteByte('\n')

        want := solve(tc)
        gotStr, err := runCandidate(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        got, err := strconv.Atoi(strings.TrimSpace(gotStr))
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d parse output: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}
