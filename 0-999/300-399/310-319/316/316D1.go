package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func nextPerm(a []int) bool {
   n := len(a)
   // find largest i such that a[i] < a[i+1]
   i := n - 2
   for i >= 0 && a[i] >= a[i+1] {
       i--
   }
   if i < 0 {
       return false
   }
   // find j>i s.t. a[j] > a[i]
   j := n - 1
   for a[j] <= a[i] {
       j--
   }
   a[i], a[j] = a[j], a[i]
   // reverse a[i+1:]
   for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
       a[l], a[r] = a[r], a[l]
   }
   return true
}

func validPerm(p, cap []int) bool {
   n := len(p)
   vis := make([]bool, n)
   for i := 0; i < n; i++ {
       if vis[i] {
           continue
       }
       // start cycle at i
       cur := i
       cycle := []int{}
       for !vis[cur] {
           vis[cur] = true
           cycle = append(cycle, cur)
           cur = p[cur]
       }
       k := len(cycle)
       if k <= 1 {
           continue
       }
       // count capacities
       b0, b1, b2 := 0, 0, 0
       for _, v := range cycle {
           if cap[v] <= 0 {
               b0++
           } else if cap[v] == 1 {
               b1++
           } else {
               b2++
           }
       }
       if b0 > 0 {
           return false
       }
       if b1 > 2 {
           return false
       }
       // need at least k-2 nodes with cap>=2 for k>=3
       if k >= 3 && b2 < k-2 {
           return false
       }
       // for k==2, b0==0 ensures both cap>=1
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   caps := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &caps[i])
   }
   // initial permutation 0..n-1
   p := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
   }
   var cnt int64 = 0
   for {
       if validPerm(p, caps) {
           cnt++
       }
       if !nextPerm(p) {
           break
       }
   }
   fmt.Println(cnt % mod)
}
