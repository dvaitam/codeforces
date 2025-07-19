package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 12345

var (
   m         int
   crat      []int
   typeChar  []int
   need      [26]bool
   state     [][]int
   g         [][]int
   nx        [][]int
   tot       int
)

// generate all state vectors
func gen(v []int, t int) {
   if t == m {
       tmp := make([]int, m)
       copy(tmp, v)
       state = append(state, tmp)
       return
   }
   for i := 0; i < crat[t]; i++ {
       gen(append(v, i), t+1)
   }
}

// fast exponentiation of transition
func rec(n int64) []int {
   if n == 0 {
       ret := make([]int, tot)
       ret[0] = 1
       return ret
   }
   half := rec(n >> 1)
   ret := make([]int, tot)
   // square
   for i := 0; i < tot; i++ {
       vi := half[i]
       if vi == 0 {
           continue
       }
       for j := 0; j < tot; j++ {
           ret[g[i][j]] = (ret[g[i][j]] + vi*half[j]) % mod
       }
   }
   // multiply by one character if odd
   if n&1 == 1 {
       tmp := make([]int, tot)
       for j := 0; j < 26; j++ {
           if !need[j] {
               continue
           }
           tmp[nx[0][j]]++
       }
       ret1 := make([]int, tot)
       for i := 0; i < tot; i++ {
           vi := ret[i]
           if vi == 0 {
               continue
           }
           for j := 0; j < tot; j++ {
               if tmp[j] == 0 {
                   continue
               }
               ret1[g[i][j]] = (ret1[g[i][j]] + vi*tmp[j]) % mod
           }
       }
       return ret1
   }
   return ret
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   fmt.Fscan(reader, &n, &m)
   crat = make([]int, m)
   typeChar = make([]int, m)
   for i := 0; i < m; i++ {
       var cs string
       var ci int
       fmt.Fscan(reader, &cs, &ci)
       typeChar[i] = int(cs[0] - 'A')
       crat[i] = ci
       need[typeChar[i]] = true
   }
   // remove duplicate types with same crat
   if m > 0 {
       leaf := make([]bool, m)
       for i := 0; i < m; i++ {
           if !leaf[i] {
               for j := i + 1; j < m; j++ {
                   if typeChar[j] == typeChar[i] && crat[j] == crat[i] {
                       leaf[j] = true
                   }
               }
           }
       }
       newType := make([]int, 0, m)
       newCrat := make([]int, 0, m)
       // always include first
       newType = append(newType, typeChar[0])
       newCrat = append(newCrat, crat[0])
       for i := 1; i < m; i++ {
           if !leaf[i] {
               newType = append(newType, typeChar[i])
               newCrat = append(newCrat, crat[i])
           }
       }
       typeChar = newType
       crat = newCrat
       m = len(typeChar)
   }
   // build states
   gen([]int{}, 0)
   tot = len(state)
   // precompute transitions
   mn := make([]int, m+1)
   mn[m] = 1
   if m > 0 {
       mn[m-1] = crat[m-1]
       for i := m - 2; i >= 0; i-- {
           mn[i] = mn[i+1] * crat[i]
       }
   }
   // combine states
   g = make([][]int, tot)
   for i := 0; i < tot; i++ {
       g[i] = make([]int, tot)
       for j := 0; j < tot; j++ {
           num := 0
           for t := 0; t < m; t++ {
               num += ((state[i][t] + state[j][t]) % crat[t]) * mn[t+1]
           }
           g[i][j] = num
       }
   }
   nx = make([][]int, tot)
   for i := 0; i < tot; i++ {
       nx[i] = make([]int, 26)
       for j := 0; j < 26; j++ {
           if !need[j] {
               continue
           }
           num := 0
           for t := 0; t < m; t++ {
               add := 0
               if typeChar[t] == j {
                   add = 1
               }
               num += ((state[i][t] + add) % crat[t]) * mn[t+1]
           }
           nx[i][j] = num
       }
   }
   // compute result
   res := rec(n)
   ans := 0
   for i := 0; i < tot; i++ {
       ok := make([]bool, 26)
       for t := 0; t < m; t++ {
           if state[i][t] == 0 {
               ok[typeChar[t]] = true
           }
       }
       add := true
       for j := 0; j < 26; j++ {
           if need[j] && !ok[j] {
               add = false
               break
           }
       }
       if add {
           ans = (ans + res[i]) % mod
       }
   }
   fmt.Println(ans)
}
