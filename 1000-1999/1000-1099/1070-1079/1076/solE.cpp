#include <bits/stdc++.h>
#include <iostream>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <cmath>
#include <cctype>
#include <algorithm>
#include <vector>
#define lowbit(x) ((x)&(-(x)))
#define rin(i,a,b) for(int i=(a);i<=(b);i++)
#define rec(i,a,b) for(int i=(a);i>=(b);i--)
#define trav(i,a) for(int i=head[(a)];i;i=e[i].nxt)
using std::cin;
using std::cout;
using std::endl;
typedef long long LL;

inline int read(){
	int x=0;char ch=getchar();
	while(ch<'0'||ch>'9') ch=getchar();
	while(ch>='0'&&ch<='9'){x=(x<<3)+(x<<1)+ch-'0';ch=getchar();}
	return x;
}

const int MAXN=300005;
int n,m;
int ecnt,head[MAXN];
int dep[MAXN],siz[MAXN],id[MAXN],tot;
int maxdep;
LL ans[MAXN];
LL b[MAXN];
std::vector<int> vec[MAXN];
struct Edge{
	int to,nxt;
}e[MAXN<<1];

struct Opera{
	int v,d,x;
	friend bool operator < (Opera lf,Opera rt){
		return dep[lf.v]+lf.d>dep[rt.v]+rt.d;
	}
}op[MAXN];

inline void add_edge(int bg,int ed){
	ecnt++;
	e[ecnt].to=ed;
	e[ecnt].nxt=head[bg];
	head[bg]=ecnt;
}

void dfs(int x,int pre,int depth){
	dep[x]=depth;
	maxdep=std::max(maxdep,depth);
	vec[depth].push_back(x);
	id[x]=++tot;
	siz[x]=1;
	trav(i,x){
		int ver=e[i].to;
		if(ver==pre) continue;
		dfs(ver,x,depth+1);
		siz[x]+=siz[ver];
	}
}

inline void Add(int l,int r,int kk){
	for(int i=l;i<=n;i+=lowbit(i)) b[i]+=kk;
	for(int i=r+1;i<=n;i+=lowbit(i)) b[i]-=kk;
}

inline LL Ask(int x){
	LL ret=0;
	for(int i=x;i;i-=lowbit(i)) ret+=b[i];
	return ret;
}

int main(){
	n=read();
	rin(i,2,n){
		int u=read(),v=read();
		add_edge(u,v);
		add_edge(v,u);
	}
	dfs(1,0,1);
	m=read();
	rin(i,1,m){
		op[i].v=read();
		op[i].d=std::min(read(),maxdep-dep[op[i].v]);
		op[i].x=read();
	}
	std::sort(op+1,op+m+1);
	int now=maxdep;
	rin(i,1,m){
		int temp=dep[op[i].v]+op[i].d;
		while(now>temp){
			rin(j,0,(int)vec[now].size()-1){
				ans[vec[now][j]]=Ask(id[vec[now][j]]);
			}
			now--;
		}
		Add(id[op[i].v],id[op[i].v]+siz[op[i].v]-1,op[i].x);
	}
	while(now){
		rin(j,0,(int)vec[now].size()-1){
			ans[vec[now][j]]=Ask(id[vec[now][j]]);
		}
		now--;
	}
	rin(i,1,n) printf("%lld ",ans[i]);
	printf("\n");
	return 0;
}