#include<bits/stdc++.h>
using namespace std ;
typedef long long ll ;
const ll N = 200011 ;
int n, k, x=1, c[N], a[N], dis[N];
vector<int> g[N], edges; 
void dfs (int v, int par){
	int k = 0 ; 
	for(int i = 0 ; i < g[v].size() ; i ++){
		int u = g[v][i] ;
		if(u == par)
			continue ;
		dis[u] = dis[v] + 1 ;
		if(g[v].size() <= x and k == c[v])
			k = (k+1)%x ;
		c[u] = k ;
		k = (k+1)%x ;
		dfs(u, v) ;
	}
}
int main ()
{
	scanf("%d%d", &n, &k) ;
	for(int i = 1 ; i < n ; i ++){
		int u, v ;
		scanf("%d%d", &v, &u) ;
		g[v].push_back(u) ;
		g[u].push_back(v) ;
		edges.push_back(u) ;
		edges.push_back(v) ;
	}
	for(int i = 1 ; i <= n ; i ++)
		a[i] = g[i].size() ;
	sort(a+1, a+n+1) ;
	int p = n ;
	for(int i = n-1 ; i >= 1 ; i --){
		while(p >= 1 and a[p] > i)
			p -- ;
		if(n-p <= k)
			x = i ;
	}
	printf("%d\n", x) ;
	for(int i = 0 ; i <= n ; i ++)
		c[i] = -1 ;
	dfs(1, 0) ;
	for(int i = 0 ; i < edges.size() ; i += 2){
		int u = edges[i], v = edges[i+1] ;
		if(dis[u] > dis[v])
			swap(u, v) ;
		printf("%d ", c[v]+1) ;
	}
}