package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var K int64
   if _, err := fmt.Fscan(reader, &n, &K, &m); err != nil {
       return
   }
   // Build Fibonacci words until length >= 200
   P := make([]string, 30)
   P[0], P[1] = "0", "1"
   BS := 0
   for i := 0; ; i++ {
       P[i+2] = P[i] + P[i+1]
       if len(P[i+1]) >= 200 {
           BS = i + 1
           break
       }
   }
   AA := P[BS]
   BB := P[BS+1]
   // extend a few more
   for i := 0; i < 5; i++ {
       P[BS+i+2] = P[BS+i] + P[BS+i+1]
   }
   // small case
   if n <= BS+3 {
       s := P[n]
       L := len(s)
       U := make([]string, L)
       for i := 0; i < L; i++ {
           U[i] = s[i:]
       }
       sort.Strings(U)
       r := U[int(K-1)]
       if len(r) > m {
           fmt.Println(r[:m])
       } else {
           fmt.Println(r)
       }
       return
   }
   // precompute substrings of P[BS+4]
   P3 := P[BS+3]
   P4 := P[BS+4]
   sz3, sz4 := len(P3), len(P4)
   uniq := make(map[string]struct{})
   for i := 0; i < sz4; i++ {
       maxl := min(sz4-i, 200)
       for l := 1; l <= maxl; l++ {
           uniq[P4[i:i+l]] = struct{}{}
       }
   }
   // map substrings to indices in lex order
   keys := make([]string, 0, len(uniq))
   for s := range uniq {
       keys = append(keys, s)
   }
   sort.Strings(keys)
   cnt := len(keys)
   R := make([]string, cnt+1)
   MNum := make(map[string]int, cnt)
   for i, s := range keys {
       idx := i + 1
       MNum[s] = idx
       R[idx] = s
   }
   // counts
   C1 := make([]int64, cnt+1)
   C2 := make([]int64, cnt+1)
   TC := make([]int64, cnt+1)
   // count in P4
   for i := 0; i < sz4; i++ {
       maxl := min(sz4-i, 200)
       for l := 1; l <= maxl; l++ {
           C2[MNum[P4[i:i+l]]]++
       }
   }
   // count in P3
   for i := 0; i < sz3; i++ {
       maxl := min(sz3-i, 200)
       for l := 1; l <= maxl; l++ {
           C1[MNum[P3[i:i+l]]]++
       }
   }
   // last substrings of P4
   Last := make([]int, 201)
   for i := 1; i <= 200; i++ {
       Last[i] = MNum[P4[sz4-i:]]
   }
   // cross boundary mappings
   sA, sB := len(AA), len(BB)
   Num0 := make([]int, 201)
   Num1 := make([]int, 201)
   for i := 1; i < 200; i++ {
       s0 := BB[sB-i:]
       p0 := AA[:200-i]
       Num0[i] = MNum[s0+p0]
       s1 := BB[sB-i:]
       p1 := BB[:200-i]
       Num1[i] = MNum[s1+p1]
   }
   // dynamic update
   INF := int64(2e18)
   parity := (BS + 5) & 1
   for i := BS + 5; i <= n; i++ {
       for j := 1; j <= cnt; j++ {
           sum := C1[j] + C2[j]
           if sum < INF {
               TC[j] = sum
           } else {
               TC[j] = INF
           }
       }
       if i&1 == parity {
           for j := 1; j < 200; j++ {
               TC[Last[j]]--
               TC[Num0[j]]++
           }
       } else {
           for j := 1; j < 200; j++ {
               TC[Last[j]]--
               TC[Num1[j]]++
           }
       }
       for j := 1; j <= cnt; j++ {
           C1[j] = C2[j]
           C2[j] = TC[j]
       }
   }
   // find K-th
   var acc int64
   for i := 1; i <= cnt; i++ {
       acc += C2[i]
       if acc >= K {
           ans := R[i]
           if len(ans) > m {
               fmt.Println(ans[:m])
           } else {
               fmt.Println(ans)
           }
           return
       }
   }
}
