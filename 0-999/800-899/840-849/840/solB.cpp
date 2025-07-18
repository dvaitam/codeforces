#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

#define gc getchar()
inline int read()
{
	int x=0,f=1;char c=gc;
	for(;c<48||c>57;c=gc)
		if(c=='-')f=-1;
	for(;c>47&&c<58;c=gc)
		x=(x<<1)+(x<<3)+c-48;
	return x*f;
}

#define maxn 300050
#define maxm 600050
int n,m,d[maxn],x,fl;
int st[maxn],pv[maxm],p1[maxm],pe[maxm],tot;
inline void side(int u,int v,int id)
{
	pv[++tot]=v;
	pe[tot]=id;
	p1[tot]=st[u];
	st[u]=tot;
}

inline void init()
{
	n=read();m=read();
	int s=0;
	for(int i=1;i<=n;++i)
	{
		d[i]=read();
		if(d[i]==-1)x=i;else s+=d[i];
	}
	if(!x&&(s&1)){fl=1;return;}
	for(int i=1;i<=m;++i)
	{
		int u=read(),v=read();
		side(u,v,i);side(v,u,i);
	}
}

int vis[maxn],use[maxn];

void dfs(int x,int e)
{
	vis[x]=1;use[e]=d[x]==1;
	for(int i=st[x];i;i=p1[i])
	{
		int v=pv[i],e2=pe[i];
		if(vis[v])continue;
		dfs(v,e2);
		use[e]^=use[e2];
	}
}

inline void solve()
{
	if(fl){puts("-1");return;}
	if(x)dfs(x,0);else dfs(1,0);
	int k=0;
	for(int i=1;i<=m;++i)
		if(use[i])++k;
	printf("%d\n",k);
	for(int i=1;i<=m;++i)
		if(use[i])printf("%d\n",i);
}

int main()
{
	init();
	solve();
	return 0;
}