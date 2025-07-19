package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wtr := bufio.NewWriter(os.Stdout)
   defer wtr.Flush()
   var n int
   if _, err := fmt.Fscan(rdr, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(rdr, &a[i])
   }
   var wid int64 = 0
   ok := true
   for i := 1; i < n; i++ {
       if a[i] == a[i-1] {
           ok = false
           break
       }
       if a[i] == a[i-1]+1 || a[i-1] == a[i]+1 {
           continue
       }
       if a[i] > a[i-1] {
           d := a[i] - a[i-1]
           if wid != 0 {
               if d != wid {
                   ok = false
                   break
               }
           } else {
               wid = d
           }
       } else {
           d := a[i-1] - a[i]
           if wid != 0 {
               if d != wid {
                   ok = false
                   break
               }
           } else {
               wid = d
           }
       }
   }
   if wid == 0 {
       wid = 1000000000
   }
   if ok {
       for i := 1; i < n; i++ {
           w := false
           if a[i] == a[i-1]+1 && (a[i]-1)%wid != 0 {
               w = true
           }
           if a[i-1] == a[i]+1 && (a[i-1]-1)%wid != 0 {
               w = true
           }
           if a[i]-a[i-1] == wid || a[i-1]-a[i] == wid {
               w = true
           }
           if !w {
               ok = false
               break
           }
       }
   }
   if !ok {
       fmt.Fprintln(wtr, "NO")
   } else {
       fmt.Fprintln(wtr, "YES")
       // fixed number of rows; columns = wid
       fmt.Fprintf(wtr, "%d %d\n", 1000000000, wid)
   }
}
