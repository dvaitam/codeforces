package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read input string
   line, _, err := reader.ReadLine()
   if err != nil {
       return
   }
   s := string(line)
   n := len(s)
   if n == 0 {
       fmt.Println("0 0")
       return
   }
   // Build composite array C of length M = 2*n-1
   M := 2*n - 1
   C := make([]byte, M)
   for i := 0; i < M; i++ {
       if i&1 == 0 {
           C[i] = s[i/2]
       } else {
           // boundary information: 1 if s[j]!=s[j+1]
           j := i / 2
           if s[j] != s[j+1] {
               C[i] = '1'
           } else {
               C[i] = '0'
           }
       }
   }
   // Manacher's algorithm for odd-length palindromes in C
   rad := make([]int, M)
   center, right := 0, -1
   for i := 0; i < M; i++ {
       var k int
       if i > right {
           k = 1
       } else {
           mirror := center*2 - i
           // rad[mirror] may exceed right-i+1
           if rad[mirror] < right-i+1 {
               k = rad[mirror] + 1
           } else {
               k = right - i + 1
           }
       }
       // expand around i with radius k
       for i-k >= 0 && i+k < M && C[i-k] == C[i+k] {
           k++
       }
       rad[i] = k - 1
       if i+rad[i] > right {
           center = i
           right = i + rad[i]
       }
   }
   var evenAns, oddAns int64
   // For each center k, count valid substrings
   for k := 0; k < M; k++ {
       rk := rad[k]
       // minimal d matching parity
       d0 := k & 1
       if rk < d0 {
           continue
       }
       // largest D <= rk with D%2 == k%2
       D := rk
       if D&1 != k&1 {
           D--
       }
       // count of d values: (D - d0)/2 + 1
       cnt := int64((D - d0)/2 + 1)
       if k&1 == 0 {
           // center on s positions => odd-length substrings in s
           oddAns += cnt
       } else {
           // center on boundary positions => even-length substrings in s
           evenAns += cnt
       }
   }
   // Output: even-length then odd-length counts
   fmt.Printf("%d %d", evenAns, oddAns)
}
