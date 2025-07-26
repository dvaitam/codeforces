package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

const MOD = 998244353

var (
   n, k int
   countS2 int
   mu3Map map[string]bool
   mu3List [][]int
   curMu []int
)

// rec generates mu2 sequences: curMu[0:pos] non-increasing >=1, weight sum mu[i]*i <=n
func rec(pos, last, w int) {
   // try v at position pos (1-based)
   for v := 1; v <= last; v++ {
       wi := v * pos
       if w+wi > n {
           break
       }
       curMu[pos-1] = v
       // record this mu2 sequence of length pos
       countS2++
       if k >= 3 {
           // compute mu3 = freq-of-freq of curMu[0:pos]
           // count frequencies of values in curMu[0:pos]
           freq := make(map[int]int)
           for i := 0; i < pos; i++ {
               freq[curMu[i]]++
           }
           // collect counts
           mu3 := make([]int, 0, len(freq))
           for _, c := range freq {
               mu3 = append(mu3, c)
           }
           sort.Sort(sort.Reverse(sort.IntSlice(mu3)))
           // key
           sb := make([]byte, 0, len(mu3)*3)
           for i, x := range mu3 {
               if i > 0 {
                   sb = append(sb, ',')
               }
               sb = strconv.AppendInt(sb, int64(x), 10)
           }
           key := string(sb)
           if !mu3Map[key] {
               mu3Map[key] = true
               // store copy
               muCopy := make([]int, len(mu3))
               copy(muCopy, mu3)
               mu3List = append(mu3List, muCopy)
           }
       }
       // recurse deeper
       rec(pos+1, v, w+wi)
   }
}

// compute partition counts p[i] mod
func partitionSum(n int) int {
   p := make([]int, n+1)
   p[0] = 1
   for i := 1; i <= n; i++ {
       for k := 1; ; k++ {
           // generalized pentagonal: k*(3k-1)/2
           gp := k*(3*k-1)/2
           if gp > i {
               break
           }
           sign := 1
           if k%2 == 0 {
               sign = -1
           }
           p[i] = (p[i] + sign*p[i-gp]) % MOD
           // second: -k*(3k+1)/2
           gp2 := k*(3*k+1)/2
           if gp2 > i {
               continue
           }
           p[i] = (p[i] + sign*p[i-gp2]) % MOD
       }
   }
   // sum p[1..n]
   sum := 0
   for i := 1; i <= n; i++ {
       sum = (sum + p[i]) % MOD
   }
   if sum < 0 {
       sum += MOD
   }
   return sum
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &k)
   if k == 1 {
       fmt.Println(partitionSum(n))
       return
   }
   // prepare for rec
   curMu = make([]int, n)
   if k >= 3 {
       mu3Map = make(map[string]bool)
       mu3List = make([][]int, 0)
   }
   countS2 = 0
   // generate S2 or S3
   rec(1, n, 0)
   if k == 2 {
       ans := countS2 % MOD
       if ans < 0 {
           ans += MOD
       }
       fmt.Println(ans)
       return
   }
   // for k>=3, mu3List is S3 initial
   prev := mu3List
   // iterate for t = 4..k, generating next sets
   for t := 4; t <= k; t++ {
       nextMap := make(map[string]bool)
       nextList := make([][]int, 0, len(prev))
       for _, mu := range prev {
           // freq-of-freq
           freq := make(map[int]int)
           for _, v := range mu {
               freq[v]++
           }
           muNext := make([]int, 0, len(freq))
           for _, c := range freq {
               muNext = append(muNext, c)
           }
           sort.Sort(sort.Reverse(sort.IntSlice(muNext)))
           sb := make([]byte, 0, len(muNext)*3)
           for i, x := range muNext {
               if i > 0 {
                   sb = append(sb, ',')
               }
               sb = strconv.AppendInt(sb, int64(x), 10)
           }
           key := string(sb)
           if !nextMap[key] {
               nextMap[key] = true
               copyMu := make([]int, len(muNext))
               copy(copyMu, muNext)
               nextList = append(nextList, copyMu)
           }
       }
       // if no change, break
       if len(nextList) == len(prev) {
           same := true
           // compare keys existence, since sets
           // we can skip deep compare as lengths equal and mapping structure stable
           // but safe to continue to k
       }
       prev = nextList
       // if stabilized to single state
       if len(prev) == 1 {
           break
       }
   }
   ans := len(prev) % MOD
   fmt.Println(ans)
}
