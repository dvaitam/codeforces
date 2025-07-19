package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func check(d int64, A1, A2 []int64) (bool, []int64) {
   n := len(A1)
   b1 := make([]int64, n)
   b2 := make([]int64, n)
   copy(b1, A1)
   copy(b2, A2)
   ans := make([]int64, 0, n)
   // phase 1: greedily match large elements
   for len(b1) > 0 {
       x := b1[len(b1)-1]
       y := b2[len(b2)-1]
       if x <= d && y <= d {
           break
       }
       if x > y {
           // need x - d in b2
           target := x - d
           idx := sort.Search(len(b2), func(i int) bool { return b2[i] >= target })
           if idx == len(b2) || b2[idx] != target {
               return false, nil
           }
           // remove x from b1, remove target from b2
           b1 = b1[:len(b1)-1]
           b2 = append(b2[:idx], b2[idx+1:]...)
           ans = append(ans, x)
       } else {
           // y > x
           target := y - d
           idx := sort.Search(len(b1), func(i int) bool { return b1[i] >= target })
           if idx == len(b1) || b1[idx] != target {
               return false, nil
           }
           b2 = b2[:len(b2)-1]
           b1 = append(b1[:idx], b1[idx+1:]...)
           ans = append(ans, -target)
       }
   }
   // remaining pairs must sum to d
   // b1 and b2 now sorted, lengths equal
   m := len(b1)
   for i := 0; i < m; i++ {
       if b1[i]+b2[m-1-i] != d {
           return false, nil
       }
       ans = append(ans, b1[i])
   }
   if len(ans) != n {
       return false, nil
   }
   return true, ans
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var tt int
   fmt.Fscan(in, &tt)
   const K = int64(1000000000)
   for tt > 0 {
       tt--
       var n int
       fmt.Fscan(in, &n)
       A1 := make([]int64, n)
       A2 := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &A1[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &A2[i])
       }
       sort.Slice(A1, func(i, j int) bool { return A1[i] < A1[j] })
       sort.Slice(A2, func(i, j int) bool { return A2[i] < A2[j] })
       // trivial check
       ok := true
       sum0 := A1[0] + A2[n-1]
       for i := 0; i < n; i++ {
           if A1[i]+A2[n-1-i] != sum0 {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(out, "YES")
           for i, v := range A1 {
               if i > 0 {
                   out.WriteByte(' ')
               }
               fmt.Fprint(out, v)
           }
           fmt.Fprintln(out)
           fmt.Fprintln(out, 0, sum0)
           continue
       }
       // equal arrays
       eq := true
       for i := 0; i < n; i++ {
           if A1[i] != A2[i] {
               eq = false
               break
           }
       }
       if eq {
           fmt.Fprintln(out, "YES")
           for i, v := range A1 {
               if i > 0 {
                   out.WriteByte(' ')
               }
               fmt.Fprint(out, v)
           }
           fmt.Fprintln(out)
           fmt.Fprintln(out, 0, 0)
           continue
       }
       found := false
       // general cases
       for i := 0; i < n && !found; i++ {
           var d int64
           if A2[n-1] > A1[n-1] {
               d = A2[n-1] - A1[i]
           } else {
               d = A1[n-1] - A2[i]
           }
           if d < 0 {
               continue
           }
           ok, ans := check(d, A1, A2)
           if ok {
               fmt.Fprintln(out, "YES")
               // print house positions
               for j, v := range ans {
                   if j > 0 {
                       out.WriteByte(' ')
                   }
                   fmt.Fprint(out, v+K)
               }
               fmt.Fprintln(out)
               if A2[n-1] > A1[n-1] {
                   fmt.Fprintln(out, K, d+K)
               } else {
                   fmt.Fprintln(out, K, d+K)
               }
               found = true
           }
       }
       if !found {
           fmt.Fprintln(out, "NO")
       }
   }
}
