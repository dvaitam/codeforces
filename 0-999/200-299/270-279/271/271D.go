package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var s string
    if _, err := fmt.Fscanln(reader, &s); err != nil {
        return
    }
    var goodStr string
    fmt.Fscanln(reader, &goodStr)
    var k int
    fmt.Fscanln(reader, &k)
    n := len(s)
    good := make([]bool, 26)
    for i := 0; i < 26 && i < len(goodStr); i++ {
        if goodStr[i] == '1' {
            good[i] = true
        }
    }
    // Precompute prefix hashes and powers
    const base uint64 = 1315423911
    H := make([]uint64, n+1)
    P := make([]uint64, n+1)
    P[0] = 1
    for i := 0; i < n; i++ {
        H[i+1] = H[i]*base + uint64(s[i]-'a'+1)
        P[i+1] = P[i] * base
    }
    seen := make(map[uint64]struct{})
    var result int
    for i := 0; i < n; i++ {
        bad := 0
        for j := i; j < n; j++ {
            if !good[s[j]-'a'] {
                bad++
                if bad > k {
                    break
                }
            }
            // Compute hash of s[i:j]
            hj := H[j+1] - H[i]*P[j-i+1]
            if _, exists := seen[hj]; !exists {
                seen[hj] = struct{}{}
                result++
            }
        }
    }
    fmt.Println(result)
}
