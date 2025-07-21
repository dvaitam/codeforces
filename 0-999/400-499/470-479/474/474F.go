package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   s := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &s[i])
   }
   // build positions map
   pos := make(map[int][]int, n)
   for i, v := range s {
       pos[v] = append(pos[v], i)
   }
   // build logs
   logs := make([]int, n+1)
   for i := 2; i <= n; i++ {
       logs[i] = logs[i/2] + 1
   }
   // build sparse table for gcd
   maxk := logs[n] + 1
   st := make([][]int, maxk)
   st[0] = make([]int, n)
   copy(st[0], s)
   for k := 1; k < maxk; k++ {
       length := 1 << k
       half := length >> 1
       cnt := n - length + 1
       st[k] = make([]int, cnt)
       for i := 0; i < cnt; i++ {
           st[k][i] = gcd(st[k-1][i], st[k-1][i+half])
       }
   }
   // process queries
   var t int
   fmt.Fscan(in, &t)
   for qi := 0; qi < t; qi++ {
       var l, r int
       fmt.Fscan(in, &l, &r)
       l--
       r--
       length := r - l + 1
       k := logs[length]
       g := gcd(st[k][l], st[k][r-(1<<k)+1])
       arr := pos[g]
       // count positions in [l, r]
       cnt := 0
       if len(arr) > 0 {
           lo := sort.Search(len(arr), func(i int) bool { return arr[i] >= l })
           hi := sort.Search(len(arr), func(i int) bool { return arr[i] > r })
           cnt = hi - lo
       }
       eaten := length - cnt
       fmt.Fprintln(out, eaten)
   }
}
