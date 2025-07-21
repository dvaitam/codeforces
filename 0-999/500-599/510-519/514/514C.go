package main

import (
   "bufio"
   "fmt"
   "os"
)

// key for pattern: string length and hash without one character
type key struct {
   length int
   hash   uint64
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   memory := make([]string, n)
   maxLen := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &memory[i])
       if l := len(memory[i]); l > maxLen {
           maxLen = l
       }
   }
   queries := make([]string, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &queries[i])
       if l := len(queries[i]); l > maxLen {
           maxLen = l
       }
   }
   // prepare powers of base
   const base = 1315423911
   pow := make([]uint64, maxLen+1)
   pow[0] = 1
   for i := 1; i <= maxLen; i++ {
       pow[i] = pow[i-1] * base
   }
   // map from pattern key to bitmask of chars at wildcard
   patterns := make(map[key]byte)
   // build from memory
   for _, s := range memory {
       L := len(s)
       // prefix hashes
       h := make([]uint64, L+1)
       for i := 0; i < L; i++ {
           h[i+1] = h[i]*base + uint64(s[i]-'a'+1)
       }
       totalHash := h[L]
       for i := 0; i < L; i++ {
           // hash without s[i]
           len2 := L - i - 1
           // prefix part: h[i]
           // suffix part: from i+1 to end
           suffix := totalHash - h[i+1]*pow[len2]
           ph := h[i]*pow[len2] + suffix
           k := key{L, ph}
           cbit := byte(1 << (s[i] - 'a'))
           patterns[k] |= cbit
       }
   }
   // answer queries
   for _, s := range queries {
       L := len(s)
       h := make([]uint64, L+1)
       for i := 0; i < L; i++ {
           h[i+1] = h[i]*base + uint64(s[i]-'a'+1)
       }
       totalHash := h[L]
       found := false
       for i := 0; i < L && !found; i++ {
           len2 := L - i - 1
           suffix := totalHash - h[i+1]*pow[len2]
           ph := h[i]*pow[len2] + suffix
           k := key{L, ph}
           if mask, ok := patterns[k]; ok {
               // mask has bits of chars at this pos in memory
               want := byte(1 << (s[i] - 'a'))
               if mask&^want != 0 {
                   found = true
                   break
               }
           }
       }
       if found {
           fmt.Fprintln(out, "YES")
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
