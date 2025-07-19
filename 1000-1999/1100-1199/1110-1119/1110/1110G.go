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
   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       u := make([]int, n-1)
       v := make([]int, n-1)
       for i := 0; i < n-1; i++ {
           var a, b int
           fmt.Fscan(reader, &a, &b)
           u[i] = a - 1
           v[i] = b - 1
       }
       var s string
       fmt.Fscan(reader, &s)
       if n <= 2 {
           fmt.Fprintln(writer, "Draw")
           continue
       }
       cnt := make([]int, n)
       for i := 0; i < n-1; i++ {
           cnt[u[i]]++
           cnt[v[i]]++
       }
       d := 0
       for i := 0; i < n; i++ {
           if cnt[i] > d {
               d = cnt[i]
           }
       }
       if d >= 4 {
           fmt.Fprintln(writer, "White")
           continue
       }
       if d == 3 {
           if n == 4 {
               hasW := false
               for i := 0; i < n; i++ {
                   if s[i] == 'W' {
                       hasW = true
                       break
                   }
               }
               if hasW {
                   fmt.Fprintln(writer, "White")
               } else {
                   fmt.Fprintln(writer, "Draw")
               }
               continue
           }
           cnt2 := make([]int, n)
           for i := 0; i < n-1; i++ {
               if cnt[u[i]] >= 2 {
                   cnt2[v[i]]++
               }
               if cnt[v[i]] >= 2 {
                   cnt2[u[i]]++
               }
           }
           ret := false
           for i := 0; i < n; i++ {
               if cnt[i] == 3 && cnt2[i] >= 2 {
                   ret = true
                   break
               }
           }
           if ret {
               fmt.Fprintln(writer, "White")
               continue
           }
           for i := 0; i < n; i++ {
               if cnt[i] >= 2 && s[i] == 'W' {
                   ret = true
                   break
               }
           }
           if ret {
               fmt.Fprintln(writer, "White")
               continue
           }
           for i := 0; i < n-1; i++ {
               if cnt[u[i]] == 3 && cnt[v[i]] == 1 && s[v[i]] == 'W' {
                   ret = true
                   break
               }
               if cnt[v[i]] == 3 && cnt[u[i]] == 1 && s[u[i]] == 'W' {
                   ret = true
                   break
               }
           }
           if ret {
               fmt.Fprintln(writer, "White")
               continue
           }
           CNT := 0
           for i := 0; i < n; i++ {
               if cnt[i] == 3 {
                   CNT++
               }
           }
           if CNT == 1 {
               for i := 0; i < n; i++ {
                   if s[i] == 'W' {
                       CNT++
                   }
               }
               if CNT == 2 && n%2 == 0 {
                   fmt.Fprintln(writer, "White")
               } else {
                   fmt.Fprintln(writer, "Draw")
               }
               continue
           } else {
               if n%2 == 1 {
                   fmt.Fprintln(writer, "White")
               } else {
                   fmt.Fprintln(writer, "Draw")
               }
               continue
           }
       }
       cnt_w0, cnt_w1 := 0, 0
       for i := 0; i < n; i++ {
           if s[i] == 'W' {
               deg := cnt[i]
               if deg-1 == 0 {
                   cnt_w0++
               } else if deg-1 == 1 {
                   cnt_w1++
               }
           }
       }
       if cnt_w1 > 0 {
           if n == 3 && cnt_w0 == 0 {
               fmt.Fprintln(writer, "Draw")
               continue
           }
           fmt.Fprintln(writer, "White")
           continue
       }
       if cnt_w0 == 2 && n%2 == 1 {
           fmt.Fprintln(writer, "White")
           continue
       } else {
           fmt.Fprintln(writer, "Draw")
           continue
       }
   }
}
