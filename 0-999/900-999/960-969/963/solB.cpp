//Copyright(c)2018 Mstdream
#include<bits/stdc++.h>
using namespace std;

inline void splay(int &v){
	v=0;char c=0;int p=1;
	while(c<'0' || c>'9'){if(c=='-')p=-1;c=getchar();}
	while(c>='0' && c<='9'){v=(v<<3)+(v<<1)+c-'0';c=getchar();}
	v*=p;
}
const int N=400010;
int ans[N],cnt,n;
namespace graph{
	int fir[N],sz,to[N],nxt[N],du[N];
	void add(int x,int y){
		nxt[++sz]=fir[x],fir[x]=sz,to[sz]=y;
		du[y]++;
	}
	void calc(){
		for(int i=1;i<=n;i++){
			if(!du[i])ans[++cnt]=i;
		}
		int l=0;
		while(l!=cnt){
			int v=ans[++l];
			for(int u=fir[v];u;u=nxt[u]){
				du[to[u]]--;
				if(du[to[u]]==0)ans[++cnt]=to[u];
			}
		}
	}
}
int fir[N],sz,to[N],nxt[N],vis[N],f[N];
void add(int x,int y){
	nxt[++sz]=fir[x],fir[x]=sz,to[sz]=y;
}
void dfs(int x,int fa){
	for(int u=fir[x];u;u=nxt[u]){
		if(to[u]!=fa){
			dfs(to[u],x);
		}
	}
	if(fa){
		if(f[x]&1)graph::add(x,fa);
		else graph::add(fa,x),f[fa]++;
	}
}
int main(){
	splay(n);
	for(int i=1,x;i<=n;i++){
		splay(x);
		if(x)add(i,x),add(x,i);
	}
	dfs(1,0);
	if(f[1]&1){
		puts("NO");
		exit(0);
	}
	graph::calc();
	puts("YES");
	for(int i=1;i<=n;i++)printf("%d\n",ans[i]);
}