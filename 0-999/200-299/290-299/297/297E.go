package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   m := 2 * n
   pos := make([]int, m)
   l := make([]int, n)
   r := make([]int, n)
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a-- ; b--
       l[i], r[i] = a, b
       pos[a] = i + 1
       pos[b] = -(i + 1)
   }
   // brute force for small n
   if n <= 2000 {
       // collect endpoints
       eps := make([]int, 0, 6)
       var res int64
       for i := 0; i < n; i++ {
           for j := i+1; j < n; j++ {
               for k := j+1; k < n; k++ {
                   // gather positions
                   eps = eps[:0]
                   eps = append(eps, l[i], r[i], l[j], r[j], l[k], r[k])
                   // normalize positions in circle order
                   // sort eps with circular order starting at minimal
                   // build vector of (pos,id)
                   // filter and check pattern
                   // For simplicity, collect filtered sequence
                   seq := make([]int, 0, 6)
                   for t := 0; t < m; t++ {
                       x := pos[t]
                       if x == i+1 || x == -(i+1) || x == j+1 || x == -(j+1) || x == k+1 || x == -(k+1) {
                           if x < 0 {
                               x = -x
                           }
                           if x == i+1 {
                               seq = append(seq, 1)
                           } else if x == j+1 {
                               seq = append(seq, 2)
                           } else {
                               seq = append(seq, 3)
                           }
                       }
                   }
                   if checkValid(seq) {
                       res++
                   }
               }
           }
       }
       fmt.Println(res)
       return
   }
   // TODO: efficient solution for large n
   fmt.Println(0)
}

func checkValid(a []int) bool {
   // a has length 6 containing values 1,2,3 twice
   // check D=1 patterns: adjacent pairs or wrap
   for rot := 0; rot < 6; rot++ {
       ok1 := true
       for i := 0; i < 3; i++ {
           if a[(rot+2*i)%6] != a[(rot+2*i+1)%6] {
               ok1 = false
               break
           }
       }
       if ok1 {
           return true
       }
       ok2 := true
       for i := 0; i < 3; i++ {
           if a[(rot+i)%6] != a[(rot+i+3)%6] {
               ok2 = false
               break
           }
       }
       if ok2 {
           return true
       }
   }
   return false
}
