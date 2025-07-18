#include <bits/stdc++.h>

#ifdef DEBUG
#define debug(...) fprintf(stderr, __VA_ARGS__)
#else
#define debug(...)
#endif

#ifdef __WIN32
#define LLFORMAT "I64"
#define Rand() ((rand() << 15) + rand())
#else
#define LLFORMAT "ll"
#define Rand() (rand())
#endif

using namespace std;

const int maxn = 1e6 + 10;

void scf(int &x) {
	char c = getchar();
	while(c < '0' || c > '9') c = getchar();
	x = 0;
	while(c >= '0' && c <= '9') x = x * 10 + c - '0', c = getchar();
	return;
}

int n, ch[maxn][2], typ[maxn], f[maxn][2], g[maxn], par[maxn];

void dfs(int u) {
	if(ch[u][0]) dfs(ch[u][0]);
	if(ch[u][1]) dfs(ch[u][1]);
	if(ch[u][0] || ch[u][1]) {
		if(typ[u] == 3) g[u] = !g[ch[u][0]];
		else if(typ[u] == 0) g[u] = g[ch[u][0]] & g[ch[u][1]];
		else if(typ[u] == 1) g[u] = g[ch[u][0]] | g[ch[u][1]];
		else g[u] = g[ch[u][0]] ^ g[ch[u][1]];
	}
	return;
}

int F(int u, int x) {
	int &t = f[u][x];
	if(~t) return t;
	if(u == 1) return t = x;
	int p = par[u];
	if(typ[p] == 3) return t = F(p, x ^ 1);
	if(typ[p] == 2) return t = F(p, x ^ g[ch[p][0] + ch[p][1] - u]);
	if(typ[p] == 1) return t = F(p, x | g[ch[p][0] + ch[p][1] - u]);
	if(typ[p] == 0) return t = F(p, x & g[ch[p][0] + ch[p][1] - u]);
}

int main() {
	scf(n);
	memset(f, -1, sizeof f);
	for (int i = 1; i <= n; ++i) {
		char c = getchar();
		while(c < 'A' || c > 'Z') c = getchar();
		if(c == 'I') scf(g[i]);
		else if(c == 'N') scf(ch[i][0]), par[ch[i][0]] = i, typ[i] = 3;
		else scf(ch[i][0]), scf(ch[i][1]), par[ch[i][0]] = par[ch[i][1]] = i, typ[i] = (c == 'A' ? 0 : (c == 'O' ? 1 : 2));
	}
	dfs(1);
	for (int u = 1; u <= n; ++u) f[u][g[u]] = g[1];
	for (int i = 1; i <= n; ++i) if(!ch[i][0] && !ch[i][1]) putchar('0' + F(i, g[i] ^ 1));
	putchar('\n');
	return 0;
}