#include<iostream>
#include<cstdio>
#include<cmath>
#include<cstring>
#include<algorithm>
using namespace std;

typedef long long LL;

inline LL read()
{
	char c=getchar();LL num=0,f=1;
	for(;!isdigit(c);c=getchar())
		f=c=='-'?-1:f;
	for(;isdigit(c);c=getchar())
		num=num*10+c-'0';
	return num*f;
}

const int N=2e5+5;

struct Edge
{
	LL u,v,w;
	bool operator < (const Edge &A) const
	{
		return w<A.w;
	}
}edge[N*2];

LL a[N],fa[N];

int find(int x)
{
	return x==fa[x]?x:fa[x]=find(fa[x]);
}

int main()
{
	int n=read(),m=read();
	LL mn=1234678951222311ll,pos;
	for(int i=1;i<=n;++i)
	{
		a[i]=read();
		if(a[i]<mn)
			mn=a[i],pos=i;
	}
	for(int i=1;i<=m;++i)
		edge[i].u=read(),edge[i].v=read(),edge[i].w=read();
	int c=m;
	for(int i=1;i<=n;++i)
	{
		if(i==pos) continue;
		edge[++c].u=i,edge[c].v=pos,edge[c].w=mn+a[i];
	}
	sort(edge+1,edge+c+1);
	for(int i=1;i<=n;++i)
		fa[i]=i;
	LL ans=0,u,v,w,x,y;
	for(int i=1;i<=c;++i)
	{
		u=edge[i].u,v=edge[i].v,w=edge[i].w;
		x=find(u),y=find(v);
		if(x==y) continue;
		fa[x]=y;ans+=w;
	}
	cout<<ans;
	return 0;
}