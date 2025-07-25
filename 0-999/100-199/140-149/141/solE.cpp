#include <vector>
#include <algorithm>
#include <iostream>
#include <cstdio>
#include <cstring>
#include <queue>
#include<set>

using namespace std;

int fa[1010];

int fin(int n)
{
	if(n==fa[n])
		return n;
	return fa[n]=fin(fa[n]);
}
int usd[100010];
int u[100010],v[100010],k[100010];
int main()
{
	int n,m,i,j;

	scanf("%d%d",&n,&m);
	for(i=1;i<=n;i++)
		fa[i]=i;

	char s[10];
	if(n%2==0)
	{
		puts("-1");
		return 0;
	}
	int num=0;
	for(i=1;i<=m;i++)
	{
		scanf("%d%d%s",&u[i],&v[i],s);

		if(fin(u[i])-fin(v[i]))
		{
			fa[fin(u[i])]=fin(v[i]);
			usd[i]=1;
			if(s[0]=='S')
				num++;
		}
		k[i]=s[0];
	}
	if(num*2+1==n)
		;
	else
	{
		for(i=1;i<=n;i++)
			fa[i]=i;
		if(num*2+1<n)
		{
			for(i=1;i<=m;i++)if(k[i]=='S'&&usd[i])
				fa[fin(u[i])]=fin(v[i]);
			else
				usd[i]=0;
			for(i=1;i<=m;i++)if(k[i]=='S'&&fin(u[i])-fin(v[i]))
			{
				fa[fin(u[i])]=fin(v[i]);
				usd[i]=1;
				num++;
				if(num*2+1==n)
					break;
			}
			for(i=1;i<=m;i++)if(k[i]-'S'&&fin(u[i])-fin(v[i]))
			{
				fa[fin(u[i])]=fin(v[i]);
				usd[i]=1;
			}
		}
		else
		{
			for(i=1;i<=m;i++)if(k[i]-'S'&&usd[i])
				fa[fin(u[i])]=fin(v[i]);
			else
				usd[i]=0;
			for(i=1;i<=m;i++)if(k[i]-'S'&&fin(u[i])-fin(v[i]))
			{
				fa[fin(u[i])]=fin(v[i]);
				usd[i]=1;
				num--;
				if(num*2+1==n)
					break;
			}
			for(i=1;i<=m;i++)if(k[i]=='S'&&fin(u[i])-fin(v[i]))
			{
				fa[fin(u[i])]=fin(v[i]);
				usd[i]=1;
			}
		}
	}
	if(num*2+1==n)
	{
		printf("%d\n",n-1);
		for(i=1;i<=m;i++)if(usd[i])
			printf("%d ",i);
		puts("");
	}
	else
		puts("-1");
}