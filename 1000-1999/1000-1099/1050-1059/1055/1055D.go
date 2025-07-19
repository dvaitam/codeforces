package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func commonSuffix(a, b string) string {
   la, lb := len(a), len(b)
   i := 0
   for i < la && i < lb && a[la-1-i] == b[lb-1-i] {
       i++
   }
   if i == 0 {
       return ""
   }
   return a[la-i:]
}

func commonPrefix(a, b string) string {
   la, lb := len(a), len(b)
   i := 0
   for i < la && i < lb && a[i] == b[i] {
       i++
   }
   return a[:i]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   v1 := make([]string, n)
   v2 := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v1[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v2[i])
   }

   var a1, a2, sp1, sp2 string
   firstDiff := false
   for i := 0; i < n; i++ {
       if v1[i] == v2[i] {
           continue
       }
       s1, s2 := v1[i], v2[i]
       // find first and last differing positions
       x1 := 0
       for x1 < len(s1) && s1[x1] == s2[x1] {
           x1++
       }
       x2 := len(s1) - 1
       for x2 >= 0 && s1[x2] == s2[x2] {
           x2--
       }
       ca1 := s1[x1 : x2+1]
       ca2 := s2[x1 : x2+1]
       newSp1 := s1[:x1]
       newSp2 := s1[x2+1:]
       if !firstDiff {
           firstDiff = true
           a1, a2 = ca1, ca2
           sp1, sp2 = newSp1, newSp2
       } else {
           if a1 != ca1 || a2 != ca2 {
               fmt.Fprintln(writer, "NO")
               return
           }
           sp1 = commonSuffix(sp1, newSp1)
           sp2 = commonPrefix(sp2, newSp2)
       }
   }
   if !firstDiff {
       // no differences
       fmt.Fprintln(writer, "YES")
       fmt.Fprintln(writer, "")
       fmt.Fprintln(writer, "")
       return
   }
   ans1 := sp1 + a1 + sp2
   ans2 := sp1 + a2 + sp2
   // validation
   for i := 0; i < n; i++ {
       s := v1[i]
       pos := strings.Index(s, ans1)
       if pos == -1 {
           if v1[i] != v2[i] {
               fmt.Fprintln(writer, "NO")
               return
           }
       } else {
           replaced := s[:pos] + ans2 + s[pos+len(ans1):]
           if replaced != v2[i] {
               fmt.Fprintln(writer, "NO")
               return
           }
       }
   }
   fmt.Fprintln(writer, "YES")
   fmt.Fprintln(writer, ans1)
   fmt.Fprintln(writer, ans2)
}
