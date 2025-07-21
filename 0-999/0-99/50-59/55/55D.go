package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

var (
   lcmVals    []int
   lcmIndex   map[int]int
   nextLcm    [][]int
   totalMods  = 2520
   divCnt     int
   dp         [][][]int64
   used       [][][]bool
   digits     []int
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func initLcm() {
   // divisors of totalMods
   for i := 1; i <= totalMods; i++ {
       if totalMods%i == 0 {
           lcmVals = append(lcmVals, i)
       }
   }
   // sort divisors (they are in increasing order by construction)
   divCnt = len(lcmVals)
   lcmIndex = make(map[int]int, divCnt)
   for i, v := range lcmVals {
       lcmIndex[v] = i
   }
   // next lcm transitions
   nextLcm = make([][]int, divCnt)
   for i := 0; i < divCnt; i++ {
       nextLcm[i] = make([]int, 10)
       for d := 0; d < 10; d++ {
           if d == 0 {
               nextLcm[i][d] = i
           } else {
               cur := lcmVals[i]
               g := gcd(cur, d)
               nl := cur / g * d
               nextLcm[i][d] = lcmIndex[nl]
           }
       }
   }
   // dp and used arrays: dimensions [pos][lcmIdx][rem]
   // max pos is 19 digits
   dp = make([][][]int64, 20)
   used = make([][][]bool, 20)
   for i := 0; i < 20; i++ {
       dp[i] = make([][]int64, divCnt)
       used[i] = make([][]bool, divCnt)
       for j := 0; j < divCnt; j++ {
           dp[i][j] = make([]int64, totalMods)
           used[i][j] = make([]bool, totalMods)
       }
   }
}

func dfs(pos, lidx, rem int, tight bool) int64 {
   if pos == len(digits) {
       // check divisibility
       if rem%lcmVals[lidx] == 0 {
           return 1
       }
       return 0
   }
   if !tight && used[pos][lidx][rem] {
       return dp[pos][lidx][rem]
   }
   var res int64
   limit := 9
   if tight {
       limit = digits[pos]
   }
   for d := 0; d <= limit; d++ {
       nt := tight && d == limit
       nlidx := nextLcm[lidx][d]
       nrem := (rem*10 + d) % totalMods
       res += dfs(pos+1, nlidx, nrem, nt)
   }
   if !tight {
       used[pos][lidx][rem] = true
       dp[pos][lidx][rem] = res
   }
   return res
}

func solve(n uint64) int64 {
   s := strconv.FormatUint(n, 10)
   digits = make([]int, len(s))
   for i, ch := range s {
       digits[i] = int(ch - '0')
   }
   // clear used
   for i := 0; i < len(digits)+1; i++ {
       for j := 0; j < divCnt; j++ {
           for k := 0; k < totalMods; k++ {
               used[i][j][k] = false
           }
       }
   }
   // compute
   total := dfs(0, lcmIndex[1], 0, true)
   // exclude zero
   return total - 1
}

func main() {
   initLcm()
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   tRaw, _ := reader.ReadString('\n')
   t, _ := strconv.Atoi(tRaw[:len(tRaw)-1])
   for i := 0; i < t; i++ {
       line, _ := reader.ReadString('\n')
       parts := []byte(line)
       var l, r uint64
       fmt.Sscan(string(parts), &l, &r)
       left := uint64(0)
       if l > 0 {
           left = l - 1
       }
       ans := solve(r) - solve(left)
       fmt.Fprintln(writer, ans)
   }
}
