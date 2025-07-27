package main

import (
   "bufio"
   "fmt"
   "os"
   "math"
)

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   t := make([]int64, n+1)
   x := make([]int64, n+1)
   t[0] = 0
   x[0] = 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &t[i], &x[i])
   }
   // try clone picks at c (0 means no clone)
   for c := 0; c <= n; c++ {
       // build runner picks sequence p: indices != c
       p := make([]int, 0, n+1)
       p = append(p, 0)
       for i := 1; i <= n; i++ {
           if i == c {
               continue
           }
           p = append(p, i)
       }
       ok := true
       m := len(p)
       // check path segments
       for k := 0; k+1 < m; k++ {
           i := p[k]
           j := p[k+1]
           dt := t[j] - t[i]
           // skip segment if covering clone event c
           if c != 0 && i < c && c < j {
               // runner must go via c
               need := abs(x[c]-x[i]) + abs(x[j]-x[c])
               if dt < need {
                   ok = false
                   break
               }
           } else {
               // direct
               if dt < abs(x[j]-x[i]) {
                   ok = false
                   break
               }
           }
       }
       if !ok {
           continue
       }
       // if c>0 and c is last event picked by clone, need runner visit x[c] before t[c]
       if c > 0 {
           // find prev runner pick before c
           prev := 0
           for i := 1; i < c; i++ {
               if i != c {
                   prev = i
               }
           }
           // if no next pick after c, ensure runner reaches c
           // if c < last runner pick j, that segment already checked via i< c < j
           // else c > all p, need check from prev only
           // determine if c is beyond last p
           if prev == p[len(p)-1] {
               // c was last skipped, so no segment covered it
               if t[c]-t[prev] < abs(x[c]-x[prev]) {
                   ok = false
               }
           }
       }
       if ok {
           fmt.Println("YES")
           return
       }
   }
   fmt.Println("NO")
}
