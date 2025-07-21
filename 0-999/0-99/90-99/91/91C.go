package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000009

var parent, rankArr, dparity []int
var flag []bool
var bipC int
var pow2 []int64

func find(x int) (root, parity int) {
   if parent[x] == x {
       return x, 0
   }
   r, p := find(parent[x])
   parity = dparity[x] ^ p
   parent[x] = r
   dparity[x] = parity
   return r, parity
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   parent = make([]int, n)
   rankArr = make([]int, n)
   dparity = make([]int, n)
   flag = make([]bool, n)
   for i := 0; i < n; i++ {
       parent[i] = i
       flag[i] = true
   }
   bipC = n
   pow2 = make([]int64, m+1)
   pow2[0] = 1
   for i := 1; i <= m; i++ {
       pow2[i] = (pow2[i-1] << 1) % MOD
   }

   for i := 1; i <= m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       ru, pu := find(u)
       rv, pv := find(v)
       if ru != rv {
           newFlag := flag[ru] && flag[rv] && (pu != pv)
           removed := 0
           if flag[ru] {
               removed++
           }
           if flag[rv] {
               removed++
           }
           // union by rank
           if rankArr[ru] < rankArr[rv] {
               ru, rv = rv, ru
               pu, pv = pv, pu
           }
           parent[rv] = ru
           dparity[rv] = pu ^ pv ^ 1
           if rankArr[ru] == rankArr[rv] {
               rankArr[ru]++
           }
           flag[ru] = newFlag
           if newFlag {
               bipC = bipC - removed + 1
           } else {
               bipC = bipC - removed
           }
       } else {
           // same component, check odd cycle
           if flag[ru] && (pu == pv) {
               flag[ru] = false
               bipC--
           }
       }
       nullity := i - n + bipC
       ans := pow2[nullity] - 1
       if ans < 0 {
           ans += MOD
       }
       fmt.Fprintln(writer, ans)
   }
}
