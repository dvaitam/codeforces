package main

import (
   "os"
   "strconv"
)

var (
   n, m    int
   head    []int
   nxt, to []int
   ia, a   []int64
   ans     []int64
   fa, fe, dis []int
   vis, inst   []bool
   tot, eid    int
)

var (
   data []byte
   ptr  int
)

func init() {
   // read all input
   data = make([]byte, 1<<20)
   n, _ := os.Stdin.Read(data)
   data = data[:n]
}

func skip() {
   for ptr < len(data) {
       b := data[ptr]
       if b == ' ' || b == '\n' || b == '\r' || b == '\t' {
           ptr++
           continue
       }
       break
   }
}

func readInt() int {
   skip()
   sign := 1
   if ptr < len(data) && data[ptr] == '-' {
       sign = -1
       ptr++
   }
   var x int
   for ptr < len(data) && data[ptr] >= '0' && data[ptr] <= '9' {
       x = x*10 + int(data[ptr]-'0')
       ptr++
   }
   return x * sign
}

func readInt64() int64 {
   skip()
   sign := int64(1)
   if ptr < len(data) && data[ptr] == '-' {
       sign = -1
       ptr++
   }
   var x int64
   for ptr < len(data) && data[ptr] >= '0' && data[ptr] <= '9' {
       x = x*10 + int64(data[ptr]-'0')
       ptr++
   }
   return x * sign
}

func dfs(x, pre int) {
   vis[x] = true
   inst[x] = true
   for i := head[x]; i != 0; i = nxt[i] {
       e := i >> 1
       if e == pre {
           continue
       }
       u := to[i]
       if !vis[u] {
           fa[u] = x
           fe[u] = e
           a[x]--
           a[u]--
           ans[e]++
           dis[u] = dis[x] + 1
           dfs(u, e)
           if a[u] != 0 {
               ans[e] += a[u]
               a[x] -= a[u]
               a[u] = 0
           }
       } else if inst[u] {
           a[x]--
           a[u]--
           ans[e]++
           if dis[u]%2 == dis[x]%2 {
               eid = e
           }
       }
   }
   inst[x] = false
}

func update(x int, v int64) {
   f := true
   for x != 1 {
       if f {
           ans[fe[x]] -= v
       } else {
           ans[fe[x]] += v
       }
       x = fa[x]
       f = !f
   }
}

// output result and exit
func out() {
   w := os.Stdout
   w.WriteString("YES\n")
   for i := 1; i <= m; i++ {
       w.WriteString(strconv.FormatInt(ans[i], 10))
       w.WriteString("\n")
   }
   os.Exit(0)
}

func main() {
   n = readInt()
   m = readInt()
   head = make([]int, n+1)
   nxt = make([]int, 2*m+2)
   to  = make([]int, 2*m+2)
   ia  = make([]int64, n+1)
   a   = make([]int64, n+1)
   ans = make([]int64, m+2)
   fa  = make([]int, n+1)
   fe  = make([]int, n+1)
   dis = make([]int, n+1)
   vis = make([]bool, n+1)
   inst = make([]bool, n+1)
   tot = 1
   for i := 1; i <= n; i++ {
       a[i] = readInt64()
       ia[i] = a[i]
   }
   for i := 1; i <= m; i++ {
       x := readInt()
       y := readInt()
       tot++
       nxt[tot] = head[x]
       head[x] = tot
       to[tot] = y
       tot++
       nxt[tot] = head[y]
       head[y] = tot
       to[tot] = x
   }
   dis[1] = 1
   dfs(1, -1)
   if a[1] == 0 {
       out()
   }
   if eid != 0 {
       u := to[eid<<1]
       v := to[eid<<1|1]
       var coef int64
       if dis[u]%2 == 1 {
           coef = a[1] / 2
       } else {
           coef = -a[1] / 2
       }
       ans[eid] += coef
       update(u, coef)
       update(v, coef)
       out()
   }
   os.Stdout.WriteString("NO\n")
}
