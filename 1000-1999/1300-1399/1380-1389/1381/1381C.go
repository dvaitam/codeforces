package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var rdr = bufio.NewReader(os.Stdin)
var wrt = bufio.NewWriter(os.Stdout)

func nextInt() int {
   var x int
   _, _ = fmt.Fscan(rdr, &x)
   return x
}

type P struct { val, idx int }

func main() {
   defer wrt.Flush()
   T := nextInt()
   for T > 0 {
       T--
       n, x, y := nextInt(), nextInt(), nextInt()
       a := make([]int, n)
       c := make([]bool, n+2)
       for i := 0; i < n; i++ {
           a[i] = nextInt()
           if a[i] >= 0 && a[i] < len(c) {
               c[a[i]] = true
           }
       }
       unused := 0
       for i := 1; i < n+2; i++ {
           if !c[i] {
               unused = i
               break
           }
       }
       y -= x
       b := make([]P, 0, y+2)
       cnt := make([]int, n+2)
       used := make([]bool, n)
       halfLimit := y / 2
       for i := 0; i < n && len(b) < y; i++ {
           if cnt[a[i]]+1 <= halfLimit {
               b = append(b, P{a[i], i})
               cnt[a[i]]++
               used[i] = true
           }
       }
       if (y-len(b) == 1) && (y%2 == 1) {
           xx, yy := -1, -1
           for i := 0; i < n; i++ {
               if !used[i] {
                   xx = i
                   used[i] = true
                   break
               }
           }
           if xx != -1 {
               for i := 0; i < n; i++ {
                   if !used[i] && a[i] != a[xx] {
                       yy = i
                       used[i] = true
                       break
                   }
               }
           }
           if yy == -1 {
               fmt.Fprintln(wrt, "NO")
               continue
           }
           b = append(b, P{a[xx], xx})
           b = append(b, P{unused, yy})
       }
       if len(b) < y || n-len(b) < x {
           fmt.Fprintln(wrt, "NO")
           continue
       }
       ans := make([]int, n)
       for i := range ans {
           ans[i] = 0
       }
       sort.Slice(b, func(i, j int) bool {
           return b[i].val < b[j].val
       })
       bs := len(b)
       half := bs / 2
       for i := 0; i < bs; i++ {
           ni := (i + half) % bs
           ans[b[i].idx] = b[ni].val
       }
       for i := 0; i < n; i++ {
           if ans[i] == 0 {
               if x > 0 {
                   x--
                   ans[i] = a[i]
               } else {
                   ans[i] = unused
               }
           }
       }
       fmt.Fprintln(wrt, "YES")
       for i, v := range ans {
           if i > 0 {
               wrt.WriteByte(' ')
           }
           fmt.Fprint(wrt, v)
       }
       fmt.Fprintln(wrt)
   }
}
