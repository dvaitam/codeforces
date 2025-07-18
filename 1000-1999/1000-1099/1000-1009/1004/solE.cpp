//Sonya and Ice Cream
#include <bits/stdc++.h>
using namespace std;
inline char gc()
{
    static char now[1<<16],*s,*t;
    if(s==t) {t=(s=now)+fread(now,1,1<<16,stdin); if(s==t) return EOF;}
    return *s++;
}
inline int read()
{
    int x=0; char ch=gc();
    while(ch<'0'||'9'<ch) ch=gc();
    while('0'<=ch&&ch<='9') x=x*10+ch-'0',ch=gc();
    return x;
}
const int N=1e5+10;
int n,s,fa[N],a[N],m,h[N],edCnt=0,mx,rt,q[N];
int dis[N],mxd[N],ans=2e9; bool mark[N];
struct edge{int v,nxt,val;} ed[N<<1];
inline void dfs(int x)
{
	if(dis[x]>dis[mx]) mx=x; 
	for(int i=h[x];i;i=ed[i].nxt)
	{
		int y=ed[i].v; if(y==fa[x]) continue;
		fa[y]=x;dis[y]=dis[x]+ed[i].val;dfs(y);
	}
}
inline void dfs2(int x)
{
	if(mxd[x]>mxd[mx]) mx=x;
	for(int i=h[x];i;i=ed[i].nxt)
	{
		int y=ed[i].v;
		if(y==fa[x]||mark[y]) continue;
		mxd[y]=mxd[x]+ed[i].val,dfs2(y);
	}
}
int main()
{
	n=read(); s=read();
	for(int i=1;i<n;++i)
	{
		int x=read(),y=read(),val=read();
		ed[++edCnt].v=y;ed[edCnt].nxt=h[x];h[x]=edCnt;ed[edCnt].val=val;
		ed[++edCnt].v=x;ed[edCnt].nxt=h[y];h[y]=edCnt;ed[edCnt].val=val;
	}
	mx=1;dfs(1);rt=mx;fa[rt]=0;dis[rt]=0;dfs(rt);
	while(mx) a[++m]=mx,mark[mx]=1,mx=fa[mx];
	for(int i=1;i<=m;++i) mx=a[i],dfs2(a[i]),mxd[a[i]]=mxd[mx];
	int l=1,r=1,h=1,t=0;
	for(;l<=m;++l)
	{
		while(h<=t&&q[h]<l) ++h;
		while(r<=m&&r-l<s)
		{
			while(h<=t&&mxd[a[r]]>=mxd[a[q[t]]]) --t; 
			q[++t]=r;++r;
		}
		ans=min(ans,max(max(dis[a[1]]-dis[a[l]],dis[a[r-1]]),mxd[a[q[h]]]));
	}
	printf("%d\n",ans);
	return 0;
}