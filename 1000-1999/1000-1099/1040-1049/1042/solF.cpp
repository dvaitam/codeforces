#include <cstdio>
#include <cstring>
#include <iostream>
#include <algorithm>
using namespace std;
#define nc() (p1==p2&&(p2=(p1=buf)+fread(buf, 1, 100000, stdin), p1==p2)?EOF:*p1++)
#define rev(i, p) for(Edge *i = h[p]; i; i = i->next)
#define rep(i, a, b) for(int i = a; i <= b; i ++ )
char buf[100000],*p1,*p2;
inline int rd() {
	int x = 0; char ch = nc();
	while(!isdigit(ch)) ch = nc();
	while(isdigit(ch)) x = (x<<1) + (x<<3) + (ch^48), ch = nc();
	return x;
}
const int N = 1100000;
struct Edge {
	int to, val;
	Edge *next;
}*h[N],e[N<<1];
int vis[N], f[N];
int val, n, ans;
inline void Add_Edge(int u, int v, int w = 1) {
	static int _ = 0;
	Edge *tmp = &e[++_];
	tmp->to = v;
	tmp->val = w;
	tmp->next = h[u];
	h[u] = tmp;
	vis[u] ++;
}
void dfs(int x,int fa) {
	f[x] = vis[x] == 1 ? 0 : -0x3f3f3f3f;
	rev(i, x) {
		if(i->to == fa) continue;
		dfs(i->to,x);
		if(f[x]+f[i->to]+1>val) {
			ans++;
			f[x] = min(f[x], f[i->to]+1);
		}
		else f[x] = max(f[x], f[i->to]+1);
	}
}
int main() {
	n = rd(), val = rd();
	for(int i = 2, x, y; i <= n; i ++ ) {
		x = rd(), y = rd(), Add_Edge(x,y), Add_Edge(y,x);
	}
	for(int i = 1; i <= n; i ++ )
		if(vis[i]!=1){ dfs(i,0); break; }
	printf("%d\n", ++ans);
}