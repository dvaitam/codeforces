package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var m int
   fmt.Fscan(reader, &m)
   wym := make([]int, m+2)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &wym[i])
   }
   sort.Ints(wym[1 : m+1])
   // match C++ behavior: wym[m+1] defaults to 0
   wym[m+1] = 0

   maxI := 105
   maxJ := m + 2
   dp := make([][]*big.Int, maxI)
   z := make([][]*big.Int, maxI)
   for i := range dp {
       dp[i] = make([]*big.Int, maxJ)
       z[i] = make([]*big.Int, maxJ)
       for j := range dp[i] {
           dp[i][j] = new(big.Int)
           z[i][j] = new(big.Int)
       }
   }
   // initial state dp[1][1][wym[1]] = 1
   dp[1][1].SetBit(dp[1][1], wym[1], 1)

   solve := func(n int) {
       // reconstruct and output solution
       fmt.Println(n)
       sto := make([]int, n+1)
       juz := make([]bool, n+1)
       mac := make([][]int, n+1)
       for i := 1; i <= n; i++ {
           mac[i] = make([]int, n+1)
       }
       kt := m
       ile := 0
       for i := n; i > 1; i-- {
           if z[i][kt].Bit(ile) == 1 {
               sto[i] = wym[kt]
               ile += (i - 1) - wym[kt]
               kt--
           } else {
               sto[i] = wym[kt]
               ile += (i - 1) - wym[kt]
           }
       }
       sto[1] = wym[1]
       // build adjacency
       type pair struct{ val, idx int }
       for h := 1; h <= n; h++ {
           var wek []pair
           for i := 1; i <= n; i++ {
               if !juz[i] {
                   wek = append(wek, pair{sto[i], i})
               }
           }
           sort.Slice(wek, func(i, j int) bool {
               return wek[i].val < wek[j].val
           })
           // reverse all
           for l, r := 0, len(wek)-1; l < r; l, r = l+1, r-1 {
               wek[l], wek[r] = wek[r], wek[l]
           }
           // reverse except first
           if len(wek) > 1 {
               for l, r := 1, len(wek)-1; l < r; l, r = l+1, r-1 {
                   wek[l], wek[r] = wek[r], wek[l]
               }
           }
           // connect edges
           for i := 1; i < len(wek); i++ {
               if i <= wek[0].val {
                   mac[wek[0].idx][wek[i].idx] = 1
               } else {
                   mac[wek[i].idx][wek[0].idx] = 1
                   sto[wek[i].idx]--
               }
           }
           juz[wek[0].idx] = true
       }
       // output matrix
       for i := 1; i <= n; i++ {
           for j := 1; j <= n; j++ {
               fmt.Printf("%d", mac[i][j])
           }
           fmt.Println()
       }
       os.Exit(0)
   }

   // DP transitions
   for i := 1; i <= 100; i++ {
       if dp[i][m].Bit(0) == 1 {
           solve(i)
           return
       }
       for j := 1; j <= m; j++ {
           // stay at j
           bil := wym[j] - i
           var prze *big.Int
           if bil > 0 {
               prze = new(big.Int).Lsh(dp[i][j], uint(bil))
           } else {
               prze = new(big.Int).Rsh(dp[i][j], uint(-bil))
           }
           dp[i+1][j].Or(dp[i+1][j], prze)
           // move to j+1
           bil = wym[j+1] - i
           if bil > 0 {
               prze = new(big.Int).Lsh(dp[i][j], uint(bil))
           } else {
               prze = new(big.Int).Rsh(dp[i][j], uint(-bil))
           }
           dp[i+1][j+1].Or(dp[i+1][j+1], prze)
           z[i+1][j+1].Or(z[i+1][j+1], prze)
       }
   }
   fmt.Println("=(")
}
