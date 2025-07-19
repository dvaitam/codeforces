package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       solve(in, out)
   }
}

func solve(in *bufio.Reader, out *bufio.Writer) {
   var m int
   fmt.Fscan(in, &m)
   b := make([]int, m)
   cnt := [5]int{}
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i])
       c := b[i]
       if c > 4 {
           c = 4
       }
       cnt[c]++
   }
   if cnt[4] > 0 || cnt[3] >= 2 || (cnt[3] >= 1 && cnt[2] >= 1) || cnt[2] >= 3 {
       fmt.Fprintln(out, -1)
       return
   }
   v := make([]int, m)
   a := make([][]int, 0)
   // initial zero row
   row0 := make([]int, m)
   copy(row0, v)
   a = append(a, row0)
   check := make([]bool, m)
   if cnt[2] > 0 {
       p1, p2 := -1, -1
       for i := 0; i < m; i++ {
           if b[i] == 2 {
               if p1 != -1 {
                   p2 = i
               } else {
                   p1 = i
               }
           }
       }
       if p2 != -1 {
           insert2 := func(p1, p2, c1, c2 int) {
               v[p1] = c1
               v[p2] = c2
               row := make([]int, m)
               copy(row, v)
               a = append(a, row)
           }
           insert2(p1, p2, 2, 2)
           insert2(p1, p2, 0, 1)
           insert2(p1, p2, 2, 1)
           insert2(p1, p2, 0, 2)
           insert2(p1, p2, 1, 1)
           insert2(p1, p2, 2, 0)
           insert2(p1, p2, 1, 2)
           insert2(p1, p2, 1, 0)
           check[p1] = true
           check[p2] = true
       } else {
           insert1 := func(p, c int) {
               v[p] = c
               row := make([]int, m)
               copy(row, v)
               a = append(a, row)
           }
           insert1(p1, 2)
           insert1(p1, 1)
           check[p1] = true
       }
       for i := 0; i < m; i++ {
           if check[i] {
               continue
           }
           if len(a)%2 == 1 {
               f2(&a, i)
           } else {
               f3(&a, i)
           }
       }
   } else {
       insert1 := func(p, c int) {
           v[p] = c
           row := make([]int, m)
           copy(row, v)
           a = append(a, row)
       }
       if cnt[3] > 0 {
           p := 0
           for i := 0; i < m; i++ {
               if b[i] == 3 {
                   p = i
               }
           }
           insert1(p, 3)
           insert1(p, 1)
           insert1(p, 2)
           check[p] = true
       } else {
           insert1(0, 1)
           check[0] = true
       }
       for i := 0; i < m; i++ {
           if !check[i] {
               f1(&a, i)
           }
       }
   }
   for _, row := range a {
       mx := 0
       for _, x := range row {
           if x > mx {
               mx = x
           }
       }
       if mx > 0 {
           for j, x := range row {
               if j > 0 {
                   out.WriteByte(' ')
               }
               fmt.Fprint(out, x)
           }
           out.WriteByte('\n')
       }
   }
}

func f1(a *[][]int, x int) {
   old := *a
   n := len(old)
   for i := 0; i < n; i++ {
       row := make([]int, len(old[i]))
       copy(row, old[i])
       *a = append(*a, row)
   }
   for i := 0; i < n; i += 2 {
       (*a)[i+1][x] = 1
       (*a)[i+n][x] = 1
   }
   (*a)[n/2], (*a)[n] = (*a)[n], (*a)[n/2]
}

func f2(a *[][]int, x int) {
   old := *a
   n := len(old)
   for i := 0; i < n; i++ {
       row := make([]int, len(old[i]))
       copy(row, old[i])
       *a = append(*a, row)
   }
   for i := 1; i < n; i += 2 {
       (*a)[i+1][x] = 1
       (*a)[i+n][x] = 1
   }
   (*a)[n][x] = 1
}

func f3(a *[][]int, x int) {
   old := *a
   n := len(old)
   for i := 0; i < n; i++ {
       row := make([]int, len(old[i]))
       copy(row, old[i])
       *a = append(*a, row)
   }
   for i := 0; i < n; i += 2 {
       (*a)[i+1][x] = 1
       (*a)[i+n][x] = 1
   }
   // swap a[n] and a[2*n-1]
   (*a)[n], (*a)[2*n-1] = (*a)[2*n-1], (*a)[n]
}
