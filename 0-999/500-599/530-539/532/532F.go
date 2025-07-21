package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007
const BASE = 91138233

func modPow(a, e int) int {
   res := 1
   for e > 0 {
       if e&1 == 1 {
           res = int((int64(res) * int64(a)) % MOD)
       }
       a = int((int64(a) * int64(a)) % MOD)
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   Sbuf := make([]byte, n)
   Tbuf := make([]byte, m)
   fmt.Fscan(reader, &Sbuf)
   fmt.Fscan(reader, &Tbuf)
   S := Sbuf
   T := Tbuf

   // Precompute powers and inverse powers
   pow := make([]int, n+1)
   invPow := make([]int, n+1)
   pow[0] = 1
   invBase := modPow(BASE, MOD-2)
   invPow[0] = 1
   for i := 1; i <= n; i++ {
       pow[i] = int((int64(pow[i-1]) * BASE) % MOD)
       invPow[i] = int((int64(invPow[i-1]) * int64(invBase)) % MOD)
   }

   // Prefix sums for S: HS[c][i] = sum of pow[pos] for S[pos]==c, pos < i
   HS := make([][]int, 26)
   for c := 0; c < 26; c++ {
       HS[c] = make([]int, n+1)
   }
   for i := 0; i < n; i++ {
       ci := int(S[i] - 'a')
       for c := 0; c < 26; c++ {
           HS[c][i+1] = HS[c][i]
       }
       HS[ci][i+1] = HS[ci][i] + pow[i]
       if HS[ci][i+1] >= MOD {
           HS[ci][i+1] -= MOD
       }
   }

   // Hash for T
   tHash := make([]int, 26)
   for j := 0; j < m; j++ {
       tj := int(T[j] - 'a')
       tHash[tj] = int((int64(tHash[tj]) + int64(pow[j])) % MOD)
   }
   // Map tHash value to letters
   tMap := make(map[int][]int)
   for c := 0; c < 26; c++ {
       v := tHash[c]
       tMap[v] = append(tMap[v], c)
   }

   res := make([]int, 0)
   sHash := make([]int, 26)
   visited := make([]bool, 26)
   for i := 0; i + m <= n; i++ {
       inv := invPow[i]
       // compute sHash for window [i, i+m)
       for c := 0; c < 26; c++ {
           raw := HS[c][i+m] - HS[c][i]
           if raw < 0 {
               raw += MOD
           }
           sHash[c] = int((int64(raw) * int64(inv)) % MOD)
           visited[c] = false
       }
       ok := true
       for c := 0; c < 26; c++ {
           if visited[c] {
               continue
           }
           sc := sHash[c]
           tc := tHash[c]
           if sc == tc {
               visited[c] = true
               continue
           }
           // try find d != c with tHash[d] == sc and sHash[d] == tc
           found := false
           for _, d := range tMap[sc] {
               if d != c && !visited[d] && sHash[d] == tc {
                   visited[c] = true
                   visited[d] = true
                   found = true
                   break
               }
           }
           if !found {
               ok = false
               break
           }
       }
       if ok {
           res = append(res, i+1)
       }
   }
   // Output
   fmt.Fprintln(writer, len(res))
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   if len(res) > 0 {
       writer.WriteByte('\n')
   }
}
