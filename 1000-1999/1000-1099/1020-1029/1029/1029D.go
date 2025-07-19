package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func dem(u int) int {
   cnt := 0
   for u > 0 {
       cnt++
       u /= 10
   }
   return cnt
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n)
   aMod := make([]int, n)
   lens := make([]int, n)
   v := make([][]int, 11)

   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       d := dem(a[i])
       lens[i] = d
       m := a[i] % k
       aMod[i] = m
       v[d] = append(v[d], m)
   }
   for d := 1; d <= 10; d++ {
       if len(v[d]) > 1 {
           sort.Ints(v[d])
       }
   }

   var res int64
   for i := 0; i < n; i++ {
       x := aMod[i]
       for j := 1; j <= 10; j++ {
           x = (x * 10) % k
           arr := v[j]
           if len(arr) == 0 {
               continue
           }
           y := (k - x) % k
           l := sort.Search(len(arr), func(idx int) bool { return arr[idx] >= y })
           r := sort.Search(len(arr), func(idx int) bool { return arr[idx] > y })
           res += int64(r - l)
           if lens[i] == j && aMod[i] == y {
               res--
           }
       }
   }
   fmt.Fprint(writer, res)
}
