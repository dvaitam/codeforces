package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   sort.Ints(arr)
   s, e := 0, 10000
   bst := -1
   for s <= e {
       m := (s + e) / 2
       ok, _, _ := try(m, arr, n)
       if ok {
           bst = m
           e = m - 1
       } else {
           s = m + 1
       }
   }
   // compute st and init for best m
   _, st, initVal := try(bst, arr, n)
   // output
   fmt.Println(bst)
   fmt.Printf("%d %d", initVal, st)
}

// try checks if given m is feasible; returns (ok, k, initVal)
func try(m int, arr []int, n int) (bool, int, int) {
   for k := 0; k <= 20000; k++ {
       mn := arr[0] - m
       mx := arr[0] + m
       ok := true
       for i := 1; i < n; i++ {
           mn += k
           mx += k
           mxa := arr[i] + m
           mna := arr[i] - m
           if mx > mxa {
               mx = mxa
           }
           if mn < mna {
               mn = mna
           }
           if mn > mx {
               ok = false
               break
           }
       }
       if ok {
           initVal := mn - (n-1)*k
           return true, k, initVal
       }
   }
   return false, 0, 0
}
