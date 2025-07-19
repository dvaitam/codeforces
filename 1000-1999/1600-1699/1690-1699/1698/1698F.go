package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Pair holds two integers
type Pair struct { first, second int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var ta int
   fmt.Fscan(reader, &ta)
   for tc := 0; tc < ta; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &b[i])
       }
       var ans []Pair
       // flip performs a symmetric reversal on a[l..r]
       flip := func(l, r int) {
           ans = append(ans, Pair{l + 1, r + 1})
           for i := l; i <= r; i++ {
               j := r + l - i
               if i < j {
                   a[i], a[j] = a[j], a[i]
               }
           }
       }
       // build adjacent-pair multisets
       ap := make([]Pair, n-1)
       bp := make([]Pair, n-1)
       for i := 0; i+1 < n; i++ {
           x, y := a[i], a[i+1]
           if x > y {
               x, y = y, x
           }
           ap[i] = Pair{x, y}
           x, y = b[i], b[i+1]
           if x > y {
               x, y = y, x
           }
           bp[i] = Pair{x, y}
       }
       sort.Slice(ap, func(i, j int) bool {
           if ap[i].first != ap[j].first {
               return ap[i].first < ap[j].first
           }
           return ap[i].second < ap[j].second
       })
       sort.Slice(bp, func(i, j int) bool {
           if bp[i].first != bp[j].first {
               return bp[i].first < bp[j].first
           }
           return bp[i].second < bp[j].second
       })
       ok := true
       if a[0] != b[0] || a[n-1] != b[n-1] {
           ok = false
       }
       if ok {
           for i := range ap {
               if ap[i] != bp[i] {
                   ok = false
                   break
               }
           }
       }
       if !ok {
           fmt.Fprintln(writer, "NO")
           continue
       }
       can := true
       for i := 1; i < n-1; i++ {
           if a[i] != b[i] {
               // direct adjacent match flip
               found := false
               for j := i; j < n-1; j++ {
                   if a[j] == b[i] && a[j+1] == b[i-1] {
                       flip(i-1, j+1)
                       found = true
                       break
                   }
               }
               if !found {
                   // try two-phase flip
                   win := false
                   for l := i - 1; l < n-1; l++ {
                       hit := false
                       for r := l + 1; r < n; r++ {
                           if a[r] == b[i] && a[r-1] == b[i-1] {
                               hit = true
                           }
                           if a[r] == a[l] && hit {
                               flip(l, r)
                               win = true
                               break
                           }
                       }
                       if win {
                           break
                       }
                   }
                   if win {
                       for j := i; j < n-1; j++ {
                           if a[j] == b[i] && a[j+1] == b[i-1] {
                               flip(i-1, j+1)
                               break
                           }
                       }
                   } else {
                       can = false
                       break
                   }
               }
           }
           if a[i] != b[i] {
               can = false
               break
           }
       }
       if can {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintln(writer, len(ans))
           for _, p := range ans {
               fmt.Fprintln(writer, p.first, p.second)
           }
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
