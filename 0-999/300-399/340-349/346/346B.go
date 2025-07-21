package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   s1, s2, virus string
   n1, n2, n3   int
   pi            []int
   nextState     [][]int
   dp            [][][]int
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

// f returns the length of the best sequence starting at i,j with virus prefix length k
func f(i, j, k int) int {
   if i == n1 || j == n2 {
       return 0
   }
   if dp[i][j][k] != -1 {
       return dp[i][j][k]
   }
   res := 0
   // skip s1[i]
   res = max(res, f(i+1, j, k))
   // skip s2[j]
   res = max(res, f(i, j+1, k))
   // take if equal and doesn't complete virus
   if s1[i] == s2[j] {
       c := s1[i]
       nk := nextState[k][c-'A']
       if nk < n3 {
           res = max(res, 1+f(i+1, j+1, nk))
       }
   }
   dp[i][j][k] = res
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &s1)
   fmt.Fscan(reader, &s2)
   fmt.Fscan(reader, &virus)
   n1, n2, n3 = len(s1), len(s2), len(virus)
   // build prefix function for virus
   pi = make([]int, n3)
   for i := 1; i < n3; i++ {
       j := pi[i-1]
       for j > 0 && virus[i] != virus[j] {
           j = pi[j-1]
       }
       if virus[i] == virus[j] {
           j++
       }
       pi[i] = j
   }
   // build automaton
   nextState = make([][]int, n3+1)
   for k := 0; k <= n3; k++ {
       nextState[k] = make([]int, 26)
       if k == n3 {
           // once fully matched, stay in forbidden state
           for c := 0; c < 26; c++ {
               nextState[k][c] = n3
           }
           continue
       }
       for c := 0; c < 26; c++ {
           t := k
           for t > 0 && byte('A'+c) != virus[t] {
               t = pi[t-1]
           }
           if byte('A'+c) == virus[t] {
               t++
           }
           nextState[k][c] = t
       }
   }
   // init dp to -1
   dp = make([][][]int, n1+1)
   for i := 0; i <= n1; i++ {
       dp[i] = make([][]int, n2+1)
       for j := 0; j <= n2; j++ {
           dp[i][j] = make([]int, n3+1)
           for k := 0; k <= n3; k++ {
               dp[i][j][k] = -1
           }
       }
   }
   length := f(0, 0, 0)
   if length <= 0 {
       fmt.Println(0)
       return
   }
   // reconstruct
   i, j, k := 0, 0, 0
   var res []byte
   for length > 0 {
       if i < n1 && f(i, j, k) == f(i+1, j, k) {
           i++
           continue
       }
       if j < n2 && f(i, j, k) == f(i, j+1, k) {
           j++
           continue
       }
       // take
       c := s1[i]
       nk := nextState[k][c-'A']
       if i < n1 && j < n2 && s1[i] == s2[j] && nk < n3 && f(i, j, k) == 1+f(i+1, j+1, nk) {
           res = append(res, c)
           i++
           j++
           k = nk
           length--
       } else {
           // fallback skip i
           i++
       }
   }
   fmt.Println(string(res))
}
