package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   s := make([]byte, n)
   t := make([]byte, n)
   // Read strings
   var ss, tt string
   fmt.Fscan(reader, &ss)
   fmt.Fscan(reader, &tt)
   copy(s, ss)
   copy(t, tt)

   const mod = 1000000007
   tot := int64(1)
   noGT := int64(1)
   noLT := int64(1)
   eqOnly := int64(1)

   for i := 0; i < n; i++ {
       var eq, lt, gt int64
       si := s[i]
       ti := t[i]
       if si != '?' && ti != '?' {
           if si == ti {
               eq = 1
           } else if si < ti {
               lt = 1
           } else {
               gt = 1
           }
       } else if si == '?' && ti == '?' {
           eq = 10
           lt = 45
           gt = 45
       } else if si == '?' {
           // ti is digit
           d := int64(ti - '0')
           eq = 1
           lt = d
           gt = 9 - d
       } else {
           // ti == '?' and si is digit
           d := int64(si - '0')
           eq = 1
           lt = 9 - d
           gt = d
       }
       total := (eq + lt + gt) % mod
       tot = tot * total % mod
       noGT = noGT * ((eq + lt) % mod) % mod
       noLT = noLT * ((eq + gt) % mod) % mod
       eqOnly = eqOnly * (eq % mod) % mod
   }
   ans := (tot - noGT - noLT + eqOnly) % mod
   if ans < 0 {
       ans += mod
   }
   fmt.Fprint(writer, ans)
}
