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

// Base64-encoded contents of testcasesA.txt.
const testcasesA64 = "NyA3ICstLS0tLS0gKystKystKyAtKystLSstIC0tKy0tLSsgKystKy0tKyAtKysrKystICsrLS0rLS0KMiA5IC0rLS0rLS0rLSArKysrLS0rKysKMyAxICsgLSAtCjkgNCArLS0tIC0tKy0gKy0tKyArKy0rICstKy0gLSsrKyArKysrICsrKy0gKy0rKwoxMCAxICsgKyArIC0gKyArICsgLSArIC0KMiA0ICstLS0gKystKwoxMCAyIC0rIC0tIC0rICsrICsrIC0tICstICsrIC0tIC0tCjcgNSArKy0rLSArLSsrLSAtLS0rLSAtLSsrKyAtKysrLSAtLSstLSArKy0rLQozIDggLSsrLS0tLSsgLSsrKystLS0gKysrKysrKysKNSA1ICsrLS0rICstLSstICstKystICsrKystIC0tKy0tCjcgNiArKy0tKysgKy0tLS0rIC0rKy0rKyAtLS0rLSsgLSsrLSsrIC0tLS0tKyArLS0tLSsKOCAyIC0tICstIC0rICstICsrICsrICsrIC0rCjIgNyAtLS0rLSstICstKy0rLSsKMiAxIC0gKwoyIDggLS0rKysrKysgLS0tKystKy0KNyA5IC0tLS0rKystLSAtKystKy0tLSsgKysrKysrLSstIC0rKy0tLSsrKyAtKy0tLS0rKysgKy0tLSstLS0tIC0tKysrKy0tKwoxIDggLSstKystLSsKMTAgOCAtLSsrLSsrKyArLSstKy0tKyArLS0rLSstLSArKystLS0rLSAtKysrLS0rKyArKy0tKystLSArKy0rKy0rKyAtLS0rLS0tKyAtKystKy0tKyArKy0rKysrKwoyIDEgLSAtCjEgMTAgKysrLSstLS0rKwo0IDYgLS0tLS0rICsrLSstKyArKy0rKysgLS0tKysrCjcgMiArKyAtLSAtKyAtLSAtKyArLSAtKwo3IDcgKy0tLSstLSAtKysrLSsrIC0tKy0tKysgLS0rLS0rLSArLSstLSsrIC0rLSstLSsgLSstLSsrLQoxIDEgLQo3IDkgLS0rLS0tKysrICstLS0rKysrLSAtLSsrKy0tLS0gLS0rKysrLS0tICstLS0tKy0tLSAtKy0tKy0rKysgLS0tKysrLSstCjEgMiArLQo1IDQgKy0rKyAtLS0tICstLS0gKy0rLSArKystCjYgNyAtLS0rKy0rICstKysrKy0gKysrLSstLSAtLSstLS0tIC0rLS0rKysgKy0rLSsrKwo3IDQgLSstKyAtKystIC0tLS0gKysrLSArKy0rICstLS0gLS0tLQo0IDEgLSArICsgLQo4IDggKy0tLS0rKysgKysrLSstKy0gLSsrKy0rKy0gLS0tKy0rLS0gLS0tKysrKysgLSsrLSstKysgLS0rLSstKysgKy0tKy0tLSsKMyAxIC0gLSArCjMgNyArKysrKy0rICsrKy0tKy0gKystLS0tKwo0IDkgKy0tLSstKy0rIC0tLSsrKy0rLSAtKy0tKy0rLSsgKy0rLSsrKysrCjQgNiArKy0tKysgLS0rLS0rIC0tLSstLSAtKystLSsKMSA5ICsrKy0rKysrKwo5IDMgKy0rIC0rKyAtLS0gKystICstLSArLSsgLSstIC0tLSAtKy0KOSA1IC0rLSstIC0rLS0tICsrLSstIC0rLS0tICsrLSstICsrLSsrICstLS0tIC0tLS0rICsrKy0rCjYgMSAtIC0gKyArIC0gLQozIDkgKysrKy0tLS0tICstLS0tLSstKyArKy0tLSstKysKNiA5ICsrLS0rKystKyAtKy0rLSsrKy0gKy0tLS0tLS0rIC0rLSstKy0tLSAtKy0tKy0tKy0gLSsrKy0tKysrCjIgOSAtKy0rLS0rLS0gKystLSsrKystCjEgOSArLSstKysrLSsKNSA2IC0tLSstKyArKystKysgLSstKy0rIC0rKy0rKyArLS0tLS0KOSA2ICsrLSsrKyArLSsrLSsgLS0rKy0rIC0rLS0tLSAtLSsrLS0gKy0rLSsrICstKystKyAtLS0tKysgLS0rKystCjYgMiAtKyArLSArKyAtLSArKyArKwozIDQgLSstLSArKystIC0rKy0KNCAzIC0tLSArLSsgKystICsrKwoxIDYgKysrKysrCjcgMSAtIC0gLSArICsgLSArCjkgNyArKy0tKy0tICstLSstKy0gLSstLS0rLSArLSstLSsrICstKy0tLSsgLS0rKy0tLSArLSstLSstIC0rKy0rLS0gLSstLSstLQo3IDQgLSsrKyAtKystIC0tKysgLSsrKyArKysrICstLS0gLSstLQozIDYgLSstLSstIC0tLS0tLSArKysrKysKMTAgMTAgKy0rLSsrKystLSAtKy0tLS0rKysrICsrLSstLS0rKy0gKy0rLSstLSstKyAtKy0tLS0tKystICstKy0rKy0tLSsgLS0rKy0rLSsrKyArKysrLS0rLS0tICsrKystLS0rKy0gLS0rKy0tKy0tKwo0IDUgKy0tLS0gKystLS0gKy0tLS0gLS0tKy0KNSAyIC0rIC0rICstIC0rICstCjcgOCAtKy0tLSstKyArLS0tKy0rKyAtKy0rLSsrKyAtKysrKy0rKyArKystKystKyArKystLS0rLSArLSstKy0tLQo4IDMgKy0rICstKyArKy0gKy0rIC0rKyAtKy0gLSstIC0rLQo1IDcgKy0rLSstLSAtKy0rLS0rIC0rLS0tKysgLS0tLSsrKyAtKysrLSstCjEwIDMgLS0rIC0tKyArLSsgLSstICstKyAtKy0gKystIC0rKyArKysgKy0rCjIgOSAtLS0tKystLS0gLS0tLSsrKysrCjcgOCArLSsrLSsrLSAtKysrLS0tKyArLSsrLS0tKyArKystKy0tKyAtLSstLS0rLSArKystLSstLSAtLSsrLSsrKwoyIDcgKystLS0rLSAtLSsrLSsrCjggMSAtIC0gKyAtICsgLSArIC0KMTAgMiArKyAtKyArLSArKyArKyAtKyAtLSArLSAtKyArLQo2IDMgLSsrIC0tLSArKysgLSsrIC0rLSAtKy0KOSA5IC0rKy0tKy0rLSAtKysrLSstLS0gKysrLSstLS0rIC0tLS0rLS0rKyArLS0rKystKysgKy0rLS0rKy0tIC0tLS0rKy0rKyArKy0rKy0rKysgLSsrKysrLSstCjQgMSAtICsgLSAtCjQgNyAtKysrKysrICstLS0tLS0gLSsrKysrKyArLSstKy0rCjEwIDMgLS0tIC0tKyArLSsgLS0tIC0tKyAtLS0gLS0tIC0tKyArLSsgLSsrCjEwIDkgLS0tLS0tKystIC0rKysrLSsrKyArLS0rKysrLSsgLS0tLSstKysrICstLSstKy0rKyAtKy0rKystKysgKy0tLSsrLSstIC0tLS0tKy0rKyAtLS0tKy0rKysgKy0tKy0tLSstCjkgOSArKysrKystLS0gKy0rKy0rKystICsrLSstLS0rLSArKysrKy0tLS0gKy0rKysrLSsrICsrLSstKy0rKyAtKy0tLS0tKysgLS0tLS0tKysrIC0tKysrKystKwoyIDkgLS0tLS0tLS0rICstLS0rKystKwo2IDcgKysrKy0rLSArLS0tLS0rICstLSsrLSsgLS0tKystKyAtLS0tKy0tICstLS0rKy0KNCA5IC0tKy0rLSsrLSArKy0tKystKysgLSstLS0rLS0tICsrKysrLSstLQo4IDcgLS0rKystKyAtLS0tLS0rICsrKystKysgKysrLSstKyArLS0rKysrICsrLS0tKy0gKysrLSsrKyAtLSstKy0rCjEgMyArKy0KMyAzIC0rLSArKy0gLS0rCjUgNCArLS0tIC0tKy0gLSsrLSArKystICstLSsKMyAyIC0tIC0tICstCjcgNyArKystLS0rIC0rLSsrLSsgLSsrKysrKyAtLSsrKystICsrKysrLS0gLS0rLS0rKyArLSsrKysrCjggOCArLS0tLS0rKyAtLSsrKy0rKyAtLSstKystLSArLSsrLS0tLSArKy0rLSsrLSAtKystLSsrKyArLS0rKy0tKyAtKystKysrKwo4IDUgLS0rKysgKy0rKysgKysrKysgKy0rLS0gLS0rKysgLSstLS0gLS0rLSsgKysrKysKOCA2ICstLSstLSArLSstLSsgLS0rLSstICstLSstLSAtKysrLSsgLS0rLS0tIC0rLS0rKyArLS0rKy0KNSA2ICstLSsrKyAtLS0rLS0gLS0tKystICsrKy0tKyAtKy0rLSsKNSA1ICstKysrICstLS0rICstLS0rIC0tKysrICstKy0rCjkgNSArKy0rKyArKy0rKyAtLSstKyAtKysrLSArLSstLSArLS0rLSArLSstLSAtLS0tLSAtKysrLQo2IDEgLSAtIC0gKyArICsKOCA5ICsrLS0rKy0tLSAtLSsrKystKysgLS0rKy0rLS0rICstLS0tKysrLSAtLS0rLS0tKy0gLS0rLS0tLS0tIC0tKy0rLS0tLSArKy0rLSstKysKOCAyICsrICsrICsrIC0tICstIC0tIC0rICstCjMgMiArLSAtLSAtLQoxIDEwIC0tLSsrLSstLSsKNyAyIC0rICsrIC0rIC0tICsrICsrICsrCjggMiArKyArKyArKyArKyArLSAtKyArLSAtKwo5IDIgKy0gLSsgKysgKysgLS0gLS0gLSsgKysgLS0KNSA2IC0rKysrLSArKysrLS0gKy0tKysrIC0tKystKyAtLSsrKy0KNyAyIC0tIC0rIC0rIC0rIC0tICstICstCjIgNiAtKy0rKy0gLSsrLSsrCjcgNSArLSsrKyAtKy0rKyArKy0tLSAtLSsrLSAtKy0rLSArKy0tLSArKy0tLQo0IDUgKysrLSsgKystKy0gLS0tLS0gLSstLS0K"

type testCase struct {
    n, k int
    members []string
}

// Embedded solution logic from 1818A.go: count occurrences of the first string.
func solve(tc testCase) string {
    base := tc.members[0]
    cnt := 0
    for _, s := range tc.members {
        if s == base {
            cnt++
        }
    }
    return strconv.Itoa(cnt)
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesA64)
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
            return nil, fmt.Errorf("case with n=%d missing k", n)
        }
        k, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("parse k: %v", err)
        }
        members := make([]string, n)
        for i := 0; i < n; i++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("case n=%d truncated at member %d", n, i)
            }
            members[i] = sc.Text()
        }
        cases = append(cases, testCase{n: n, k: k, members: members})
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
        fmt.Println("usage: go run verifierA.go /path/to/binary")
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
        fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
        for i, s := range tc.members {
            if i > 0 {
                input.WriteByte(' ')
            }
            input.WriteString(s)
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

