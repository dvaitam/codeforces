package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       // skip possible blank line
       fmt.Fscan(in, &n)
       vals := make([]int, n)
       left := make([]int, n)
       right := make([]int, n)
       isLeaf := make([]bool, n)
       for i := 0; i < n; i++ {
           var a int
           fmt.Fscan(in, &a)
           if a >= 0 {
               vals[i] = a
               isLeaf[i] = true
           } else {
               // non-leaf
               var l, r int
               fmt.Fscan(in, &l, &r)
               left[i] = l
               right[i] = r
           }
       }
       // dfs returns multiset and height
       var dfs func(u int) ([]int, int)
       dfs = func(u int) ([]int, int) {
           if isLeaf[u] {
               return []int{vals[u]}, 0
           }
           a, hl := dfs(left[u])
           b, hr := dfs(right[u])
           // merge
           v := append(a, b...)
           sort.Ints(v)
           h := hl
           if hr > h {
               h = hr
           }
           h++
           // remove one element based on height parity
           if h%2 == 1 {
               // maximize: remove smallest
               v = v[1:]
           } else {
               // minimize: remove largest
               v = v[:len(v)-1]
           }
           return v, h
       }
       res, _ := dfs(0)
       if len(res) > 0 {
           fmt.Fprintln(out, res[0])
       } else {
           fmt.Fprintln(out, 0)
       }
   }
}
