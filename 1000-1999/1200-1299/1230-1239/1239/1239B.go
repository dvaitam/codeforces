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
   var s string
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   totl := 0
   for i := 0; i < n; i++ {
       if s[i] == '(' {
           totl++
       }
   }
   totr := n - totl
   if totl != totr {
       fmt.Fprintln(writer, 0)
       fmt.Fprintln(writer, "1 1")
       return
   }
   lastpos := -1
   top := 0
   for i := 0; i < n; i++ {
       if s[i] == '(' {
           top++
       } else if top > 0 {
           top--
       } else {
           lastpos = i
       }
   }
   var t string
   if lastpos < 0 {
       t = s
   } else {
       // rotate to start after last unmatched ')'
       start := lastpos + 1
       if start >= n {
           start = 0
       }
       t = s[start:] + s[:start]
   }
   ans := 0
   top = 0
   for i := 0; i < n; i++ {
       if t[i] == '(' {
           top++
       } else {
           top--
           if top == 0 {
               ans++
           }
       }
   }
   fstans := ans
   anspos1, anspos2 := 1, 1
   rotationStart := 0
   if lastpos >= 0 {
       rotationStart = lastpos + 1
       if rotationStart >= n {
           rotationStart %= n
       }
   }
   // examine each primitive segment
   for i := 0; i < n; {
       // find matching ')' for t[i]
       j := i + 1
       top = 1
       tmp := 0
       for j < n && top != 0 {
           if t[j] == '(' {
               top++
           } else {
               top--
               if top == 1 {
                   tmp++
               }
           }
           j++
       }
       j--
       if ans < tmp+1 {
           ans = tmp + 1
           anspos1 = (rotationStart + i) % n + 1
           anspos2 = (rotationStart + j) % n + 1
       }
       // consider inner splits
       top = 0
       tmp = 0
       st := i + 1
       for u := i + 1; u < j; u++ {
           if t[u] == '(' {
               top++
           } else {
               top--
               if top == 1 {
                   tmp++
               } else if top == 0 {
                   if ans < fstans+tmp+1 {
                       ans = fstans + tmp + 1
                       anspos1 = (rotationStart + st) % n + 1
                       anspos2 = (rotationStart + u) % n + 1
                   }
                   st = u + 1
                   tmp = 0
               }
           }
       }
       i = j + 1
   }
   fmt.Fprintln(writer, ans)
   fmt.Fprintf(writer, "%d %d\n", anspos1, anspos2)
}
