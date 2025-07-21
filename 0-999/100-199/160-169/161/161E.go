package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   // Precompute primes up to 100000 and organize by length and prefix
   const maxV = 100000
   isPrime := make([]bool, maxV)
   for i := 2; i < maxV; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i < maxV; i++ {
       if isPrime[i] {
           for j := i * i; j < maxV; j += i {
               isPrime[j] = false
           }
       }
   }
   // primesByLen[n] = list of prime strings of length n (with leading zeros)
   primesByLen := make([][]string, 6)
   for v := 2; v < maxV; v++ {
       if !isPrime[v] {
           continue
       }
       s := fmt.Sprintf("%05d", v)
       // for each possible length 1..5, take suffix of rightmost len
       for n := 1; n <= 5; n++ {
           str := s[5-n:]
           // skip primes with leading zero for row1 only handled in logic
           primesByLen[n] = append(primesByLen[n], str)
       }
   }
   // prefixMap[n][L][prefix] -> list of prime strings of length n with given prefix length L
   prefixMap := make([][]map[string][]string, 6)
   for n := 1; n <= 5; n++ {
       prefixMap[n] = make([]map[string][]string, n+1)
       for L := 0; L <= n; L++ {
           prefixMap[n][L] = make(map[string][]string)
       }
       for _, p := range primesByLen[n] {
           for L := 1; L < n; L++ {
               pre := p[:L]
               prefixMap[n][L][pre] = append(prefixMap[n][L][pre], p)
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // Process each safe
   for ti := 0; ti < t; ti++ {
       var pi string
       fmt.Fscan(reader, &pi)
       n := len(pi)
       // matrix for partial fill
       var mat [5][5]byte
       // fill row1 and col1
       for j := 0; j < n; j++ {
           mat[0][j] = pi[j]
           mat[j][0] = pi[j]
       }
       // DFS to count
       var dfs func(i int) int64
       dfs = func(i int) int64 {
           if i > n {
               return 1
           }
           idx := i - 1
           // build prefix of length i-1 for row i
           pre := make([]byte, idx)
           for j := 0; j < idx; j++ {
               pre[j] = mat[idx][j]
           }
           key := string(pre)
           candidates := prefixMap[n][idx][key]
           var cnt int64
           for _, pstr := range candidates {
               // ensure first row matches given pi only for i=1, but row1 is fixed
               // assign suffix for row i and symmetric
               for j := idx; j < n; j++ {
                   mat[idx][j] = pstr[j]
                   if j > idx {
                       mat[j][idx] = pstr[j]
                   }
               }
               cnt += dfs(i + 1)
           }
           return cnt
       }
       // Start from row 2
       var result int64
       if n == 1 {
           result = 1
       } else {
           // prefix for row2 length1 is mat[1][0] which was set to pi[1]
           // dfs(2) will handle row2
           result = dfs(2)
       }
       fmt.Fprintln(writer, result)
   }
}
