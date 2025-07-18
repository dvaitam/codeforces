#include<cstdio>
#include<cstring>
#include<cmath>
#include<algorithm>
#include<iostream>
#include<stdlib.h>
#include<vector>
#include<queue>
using namespace std;
#define LL long long

const int maxn=500+10;
int n,m,i,j,Aang=0,ans[maxn][maxn],times=0;
int num;
char ch[200];

void print(int k)
{
	num=0;
	while (k>0) ch[++num]=k%10,k/=10;
	while (num)
		putchar(ch[num--]+48);
	putchar(32);
}

int main()
{
#ifdef h10
	freopen("3.in","r",stdin);
	freopen("3.out","w",stdout);
#endif
	scanf("%d%d",&n,&m);
	for (i=1;i<m;i++)
	for (j=1;j<=n;j++)
			ans[j][i]=++times;
	for (i=1;i<=n;i++)
		for (j=m;j<=n;j++)
		{
			ans[i][j]=++times;
			if (j==m) Aang+=times;
		}
	printf("%d\n",Aang);
	for (i=1;i<=n;i++)
	{
		for (j=1;j<=n;j++)
			print(ans[i][j]);
		putchar('\n');
	}
}