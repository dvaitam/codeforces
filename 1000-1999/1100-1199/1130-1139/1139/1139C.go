package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

var parent []int

// readInt reads an integer from bufio.Reader
func readInt(r *bufio.Reader) int {
   num := 0
   sign := 1
   b, err := r.ReadByte()
   for err == nil && (b < '0' || b > '9') && b != '-' {
       b, err = r.ReadByte()
   }
   if err != nil {
       return 0
   }
   if b == '-' {
       sign = -1
       b, err = r.ReadByte()
   }
   for err == nil && b >= '0' && b <= '9' {
       num = num*10 + int(b-'0')
       b, err = r.ReadByte()
   }
   return num * sign
}

// find returns the representative of x with path compression
func find(x int) int {
   if parent[x] != x {
       parent[x] = find(parent[x])
   }
   return parent[x]
}

// modPow computes a^b mod mod
func modPow(a, b int64) int64 {
   var res int64 = 1
   a %= mod
   for b > 0 {
       if b&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       b >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   n := readInt(reader)
   k := readInt(reader)
   parent = make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
   }
   for i := 1; i < n; i++ {
       u := readInt(reader)
       v := readInt(reader)
       w := readInt(reader)
       if w == 0 {
           parent[find(u)] = find(v)
       }
   }
   size := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       size[find(i)]++
   }
   var ans int64 = modPow(int64(n), int64(k))
   for i := 1; i <= n; i++ {
       if size[i] > 0 {
           ans = (ans - modPow(size[i], int64(k)) + mod) % mod
       }
   }
   fmt.Fprint(writer, ans)
}
