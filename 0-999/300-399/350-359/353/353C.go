package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read n
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Read array a
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       var ai int64
       fmt.Fscan(reader, &ai)
       a[i] = ai
   }
   // Read binary string s (LSB first)
   var s string
   fmt.Fscan(reader, &s)
   // Precompute prefix sums of a: pre[i] = sum_{k=0..i-1} a[k]
   pre := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       pre[i] = pre[i-1] + a[i-1]
   }
   // Precompute suffix sums of a where s[i]=='1': suf[i] = sum_{k=i..n-1, s[k]=="1"} a[k]
   suf := make([]int64, n+1)
   for i := n - 1; i >= 0; i-- {
       suf[i] = suf[i+1]
       if s[i] == '1' {
           suf[i] += a[i]
       }
   }
   // f(m) is sum of all bits where m has 1: suf[0]
   var ans int64 = suf[0]
   // Try candidates: for each bit i where s[i]==1, set that bit to 0, all lower bits to 1
   for i := n - 1; i >= 0; i-- {
       if s[i] == '1' {
           // sum of higher bits equal to m: suf[i+1]
           // sum of lower bits all 1: pre[i]
           cand := suf[i+1] + pre[i]
           if cand > ans {
               ans = cand
           }
       }
   }
   // Output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   writer.WriteString(strconv.FormatInt(ans, 10))
   writer.WriteByte('\n')
}
