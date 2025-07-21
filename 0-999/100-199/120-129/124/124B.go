package main

import (
   "bufio"
   "fmt"
   "os"
)

// nextPerm generates the next lexicographical permutation of a,
// returning false if there is no next permutation.
func nextPerm(a []int) bool {
   n := len(a)
   i := n - 2
   for i >= 0 && a[i] >= a[i+1] {
       i--
   }
   if i < 0 {
       return false
   }
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

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   nums := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &nums[i])
   }
   // initial permutation of digit positions
   perm := make([]int, k)
   for i := 0; i < k; i++ {
       perm[i] = i
   }
   const inf = int64(1<<62 - 1)
   best := inf
   // iterate over all permutations
   for {
       var mn, mx int64
       for idx, s := range nums {
           var v int64
           for _, p := range perm {
               v = v*10 + int64(s[p]-'0')
           }
           if idx == 0 || v < mn {
               mn = v
           }
           if idx == 0 || v > mx {
               mx = v
           }
       }
       diff := mx - mn
       if diff < best {
           best = diff
       }
       if !nextPerm(perm) {
           break
       }
   }
   fmt.Println(best)
}
