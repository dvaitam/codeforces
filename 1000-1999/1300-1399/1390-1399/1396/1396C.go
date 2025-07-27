package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var r1, r2, r3, d int64
	if _, err := fmt.Fscan(reader, &n, &r1, &r2, &r3, &d); err != nil {
		return
	}
	const INF int64 = math.MaxInt64 / 4
	dp0 := -d
	dp1 := INF
	for i := 0; i < n; i++ {
		var a int64
		fmt.Fscan(reader, &a)
		// cost to fully clear stage i: kill normals with pistol, then boss with AWP
		t1 := a*r1 + r3
		// cost to clear normals and injure boss (dirty): either laser or pistol for normals+one shot
		t2 := min(r2, (a+1)*r1)
		// compute new dp
		ndp0 := min(dp0+d+t1, dp1+r3+d+t1)
		ndp1 := min(dp0+d+t2, dp1+d+t2)
		dp0, dp1 = ndp0, ndp1
	}
	// answer is either all cleared, or last dirty then kill last boss
	ans := min(dp0, dp1+r3)
	fmt.Fprint(writer, ans)
}
