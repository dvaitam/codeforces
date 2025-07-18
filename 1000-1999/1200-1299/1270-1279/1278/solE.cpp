#include <bits/stdc++.h>
#define il inline
#define ll long long
#define Max 500005
using namespace std;
il int read()
{
	char c=getchar();
	int x=0,f=1;
	while(c>'9'||c<'0')
	{
		if(c=='-') f=-1;
		c=getchar();
	}
	while(c>='0'&&c<='9')
	{
		x=x*10+c-'0';
		c=getchar();
	}
	return x*f;
}
struct node
{
	int t,nt;
}e[Max<<1];
int n,l[Max],r[Max],tot,head[Max],sz[Max],sz2[Max];
il void add(int u,int v)
{
	e[++tot].t=v;
	e[tot].nt=head[u];
	head[u]=tot;
}
il bool cmp(int a,int b)
{
	return sz[a]>sz[b];
}
il void dfs(int u,int fa)
{
	sz[u]=1;
	//vector<int> q;
	for(int i=head[u];i;i=e[i].nt)
	{
		int v=e[i].t;
		if(v==fa) continue;
		dfs(v,u);
		sz[u]+=sz[v];
		sz2[u]++;
		//q.push_back(v);
	}
}
il void dfs2(int u,int fa)
{
	int nw2=l[u]+1,nw=l[u]-1;
	for(int i=head[u];i;i=e[i].nt)
	{
		int v=e[i].t;
		if(v==fa) continue;
		r[v]=nw2;
		nw2++;
		l[v]=nw-sz2[v];
		nw=nw-sz[v]*2+1;
	}
	for(int i=head[u];i;i=e[i].nt)
	{
		int v=e[i].t;
		if(v==fa) continue;
		dfs2(v,u);
	}
}
signed main()
{
	n=read();
	for(int i=1;i<n;i++)
	{
		int u=read(),v=read();
		add(u,v),add(v,u);
	}
	dfs(1,0);
	r[1]=2*n;
	l[1]=2*n-sz2[1]-1;
	dfs2(1,0);
	for(int i=1;i<=n;i++) printf("%d %d\n",l[i],r[i]);
}