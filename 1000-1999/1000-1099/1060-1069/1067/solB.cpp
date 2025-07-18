#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

#define gc getchar()
inline int read(){
	int x=0,f=0;char c=gc;
	for(;c<48||c>57;c=gc)
		if(c=='-')f=1;
	for(;c>47&&c<58;c=gc)
		x=x*10+c-48;
	return f?-x:x;
}
#define io read()

#define maxn 100005
int n,k,st[maxn],pv[maxn<<1],p1[maxn<<1],tot;
int deg[maxn];
inline void side(int u,int v){
	pv[++tot]=v;
	p1[tot]=st[u];
	st[u]=tot;
}

inline void init(){
	n=io;k=io;
	for(int i=1;i<n;++i){
		int u=io,v=io;
		side(u,v);side(v,u);
		++deg[u];++deg[v];
	}
}

int q[maxn],d[maxn],vis[maxn];
inline void solve(){
	int l=1,r=0,fl=1;
	for(int i=1;i<=n;++i)
		if(deg[i]==1)q[++r]=i,vis[i]=1;
	
	while(l<=r){
		int x=q[l++];
		for(int i=st[x];i;i=p1[i]){
			int v=pv[i];
			if(vis[v]){
				if(d[v]==d[x]-1||d[v]==d[x]+1)continue;
				fl=0;break;
			}
			d[v]=d[x]+1;
			q[++r]=v;
			vis[v]=1;
		}
	}
	if(d[q[r]]!=k)fl=0;
	if(deg[q[r]]<3)fl=0;
	for(int i=1;i<=n;++i)
		if(i!=q[n]&&deg[i]!=1&&deg[i]<4)fl=0;
	puts(fl?"YES":"NO");
}

int main(){
	init();
	solve();
	return 0;
}