package main

import (
   "bufio"
   "fmt"
   "os"
)

func prefixFunc(s []byte) []int {
   n := len(s)
   pi := make([]int, n)
   for i := 1; i < n; i++ {
       j := pi[i-1]
       for j > 0 && s[i] != s[j] {
           j = pi[j-1]
       }
       if s[i] == s[j] {
           j++
       }
       pi[i] = j
   }
   return pi
}

func reverseBytes(s []byte) []byte {
   n := len(s)
   r := make([]byte, n)
   for i := 0; i < n; i++ {
       r[i] = s[n-1-i]
   }
   return r
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   a, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   b, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   // strip newline \r?\n?
   a = trimNewline(a)
   b = trimNewline(b)
   n := len(a)
   if len(b) != n {
       fmt.Println("-1 -1")
       return
   }
   A := []byte(a)
   B := []byte(b)
   // AR = reverse(A), BR = reverse(B)
   AR := reverseBytes(A)
   BR := reverseBytes(B)
   // compute max head length h_max such that A[0..h-1] == reverse(B)[0..h-1]
   hMax := 0
   for hMax < n && A[hMax] == BR[hMax] {
       hMax++
   }
   // pi2 for S2_full = AR + '#' + B
   s2 := make([]byte, n+1+n)
   copy(s2, AR)
   s2[n] = '#'
   copy(s2[n+1:], B)
   pi2 := prefixFunc(s2)
   // rolling hash preparation
   const base = 1315423911
   pow := make([]uint64, n+1)
   pow[0] = 1
   for i := 1; i <= n; i++ {
       pow[i] = pow[i-1] * base
   }
   HA := make([]uint64, n+1)
   HB := make([]uint64, n+1)
   for i := 0; i < n; i++ {
       HA[i+1] = HA[i]*base + uint64(A[i])
       HB[i+1] = HB[i]*base + uint64(B[i])
   }
   // try head lengths from hMax down to 1
   for h := hMax; h > 0; h-- {
       if h == n {
           continue
       }
       // t_max at pos in pi2: position = len(AR) + 1 + (n-h-1) = 2*n - h
       pos := 2*n - h
       if pos < len(pi2) {
           tmax := pi2[pos]
           for t := tmax; t > 0; t-- {
               m := n - h - t
               if m < 0 {
                   continue
               }
               // check X: B[0..m-1] vs A[h..h+m-1]
               ok := false
               if m == 0 {
                   ok = true
               } else {
                   hashB := HB[m]
                   hashA := HA[h+m] - HA[h]*pow[m]
                   if hashB == hashA {
                       ok = true
                   }
               }
               if ok {
                   i := h - 1
                   j := n - t
                   fmt.Printf("%d %d", i, j)
                   return
               }
           }
       }
   }
   fmt.Println("-1 -1")
}

func trimNewline(s string) string {
   // remove trailing \r and \n
   for len(s) > 0 {
       last := s[len(s)-1]
       if last == '\n' || last == '\r' {
           s = s[:len(s)-1]
           continue
       }
       break
   }
   return s
}
