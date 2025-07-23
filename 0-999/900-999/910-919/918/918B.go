package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	ipToName := make(map[string]string, n)
	for i := 0; i < n; i++ {
		var name, ip string
		fmt.Fscan(in, &name, &ip)
		ipToName[ip] = name
	}

	for i := 0; i < m; i++ {
		var cmd, ipWithSemicolon string
		fmt.Fscan(in, &cmd, &ipWithSemicolon)
		ip := ipWithSemicolon[:len(ipWithSemicolon)-1]
		name := ipToName[ip]
		fmt.Fprintf(out, "%s %s #%s\n", cmd, ipWithSemicolon, name)
	}
}
