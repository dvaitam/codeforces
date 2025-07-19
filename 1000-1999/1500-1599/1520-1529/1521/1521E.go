package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

type pair struct {
   first, second int
}

func sqr(a int) int {
   l, r := 0, a
   for l < r {
       m := (l + r) >> 1
       if m*m > a {
           r = m
       } else {
           l = m + 1
       }
   }
   return l
}

func main() {
   defer writer.Flush()
   var T int
   fmt.Fscan(reader, &T)
   for tc := 0; tc < T; tc++ {
       solve()
   }
}

func solve() {
   var m, k int
   fmt.Fscan(reader, &m, &k)
   arr := make([]pair, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &arr[i].first)
       arr[i].second = i + 1
   }
   sort.Slice(arr, func(i, j int) bool {
       return arr[i].first > arr[j].first
   })
   total := m
   maxf := 0
   if k > 0 {
       maxf = arr[0].first
   }
   check := func(a int) bool {
       cnt := a*a - (a/2)*(a/2)
       return a*((a+1)/2) >= maxf && cnt >= total
   }
   l, r := 1, 2*sqr(m)
   for l < r {
       mid := (l + r) >> 1
       if check(mid) {
           r = mid
       } else {
           l = mid + 1
       }
   }
   size := l
   res := make([][]int, size)
   for i := 0; i < size; i++ {
       res[i] = make([]int, size)
   }
   id := 0
   // fill even rows, odd columns
   for i := 0; i < size; i += 2 {
       for j := 1; j < size; j += 2 {
           for id < k && arr[id].first == 0 {
               id++
           }
           if id == k {
               goto DON
           }
           res[i][j] = arr[id].second
           arr[id].first--
       }
   }
   if id == 0 {
       // fill even rows, even columns
       for i := 0; i < size; i += 2 {
           for j := 0; j < size; j += 2 {
               if arr[id].first == 0 {
                   break
               }
               res[i][j] = arr[id].second
               arr[id].first--
           }
       }
       id++
   }
   // fill remaining
   for i := 0; i < size; i++ {
       for j := 0; j < size; j++ {
           for id < k && arr[id].first == 0 {
               id++
           }
           if id == k {
               goto DON
           }
           if res[i][j] != 0 || (i%2 == 1 && j%2 == 1) {
               continue
           }
           res[i][j] = arr[id].second
           arr[id].first--
       }
   }
DON:
   fmt.Fprintln(writer, size)
   for i := 0; i < size; i++ {
       for j := 0; j < size; j++ {
           fmt.Fprint(writer, res[i][j], " ")
       }
       fmt.Fprintln(writer)
   }
}
