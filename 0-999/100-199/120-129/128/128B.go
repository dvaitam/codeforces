package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil && len(line) == 0 {
       return
   }
   s := strings.TrimSpace(line)
   line, err = reader.ReadString('\n')
   if err != nil && len(line) == 0 {
       return
   }
   k, _ := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
   n := len(s)
   // Build suffix array
   sa := buildSA(s)
   // Build LCP array
   lcp := buildLCP(s, sa)
   var total int64
   for i := 0; i < n; i++ {
       suffixLen := int64(n - sa[i])
       common := int64(lcp[i])
       cnt := suffixLen - common
       if total+cnt >= k {
           t := k - total
           length := int(common + t)
           fmt.Println(s[sa[i] : sa[i]+length])
           return
       }
       total += cnt
   }
   fmt.Println("No such line.")
}

// buildSA returns the suffix array of s
func buildSA(s string) []int {
   n := len(s)
   sa := make([]int, n)
   rank := make([]int, n)
   tmp := make([]int, n)
   for i := 0; i < n; i++ {
       sa[i] = i
       rank[i] = int(s[i])
   }
   for k := 1; k < n; k <<= 1 {
       sort.Slice(sa, func(i, j int) bool {
           a, b := sa[i], sa[j]
           if rank[a] != rank[b] {
               return rank[a] < rank[b]
           }
           ra := -1
           rb := -1
           if a+k < n {
               ra = rank[a+k]
           }
           if b+k < n {
               rb = rank[b+k]
           }
           return ra < rb
       })
       tmp[sa[0]] = 0
       for i := 1; i < n; i++ {
           prev, cur := sa[i-1], sa[i]
           var prev2, cur2 int
           if prev+k < n {
               prev2 = rank[prev+k]
           } else {
               prev2 = -1
           }
           if cur+k < n {
               cur2 = rank[cur+k]
           } else {
               cur2 = -1
           }
           if rank[prev] != rank[cur] || prev2 != cur2 {
               tmp[cur] = tmp[prev] + 1
           } else {
               tmp[cur] = tmp[prev]
           }
       }
       copy(rank, tmp)
       if rank[sa[n-1]] == n-1 {
           break
       }
   }
   return sa
}

// buildLCP returns LCP array for s and given sa
func buildLCP(s string, sa []int) []int {
   n := len(s)
   rank := make([]int, n)
   for i := 0; i < n; i++ {
       rank[sa[i]] = i
   }
   lcp := make([]int, n)
   h := 0
   for i := 0; i < n; i++ {
       if rank[i] > 0 {
           j := sa[rank[i]-1]
           for i+h < n && j+h < n && s[i+h] == s[j+h] {
               h++
           }
           lcp[rank[i]] = h
           if h > 0 {
               h--
           }
       }
   }
   return lcp
}
