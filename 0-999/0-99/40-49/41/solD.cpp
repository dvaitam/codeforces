#include <iostream>
#include <math.h>
#include <cmath>
#include <cstring>
#include <algorithm>
#include <cstdio>
#include <string>
using namespace std;

int n,m,k,ans0,ansi;
int mp[101][101],f[101][101][12];
char dir[101][101][12];
void find(int x,int y,int u)
{
	int u0;
	if (x==n)
	{
		printf("%d\n",y);
		return;
	}
	u0=(f[x][y][u]-mp[x][y])%k;
	if (dir[x][y][u]=='R')
		find(x+1,y-1,u0);
	else
		find(x+1,y+1,u0);
	printf("%c",dir[x][y][u]);
}
int main()
{
	char x;
	int u0;
	scanf("%d%d%d",&n,&m,&k);
	k++;
	memset(f,-1,sizeof(f));
	for (int i=1;i<=n;i++)
	{
		getchar();
		for (int j=1;j<=m;j++)
		{
			scanf("%c",&x);
			mp[i][j]=x-'0';
		}
	}
	for (int i=1;i<=m;i++) 
		f[n][i][mp[n][i]%k]=mp[n][i];

	for (int i=n-1;i>0;i--)
		for (int j=1;j<=m;j++)
		{
			for (int u=0;u<k;u++)
			{
				u0=(u+mp[i][j])%k;
				if (j>1&&f[i+1][j-1][u]>=0)
					if (f[i+1][j-1][u]+mp[i][j]>f[i][j][u0])
					{
						dir[i][j][u0]='R';
						f[i][j][u0]=f[i+1][j-1][u]+mp[i][j];
					}
				if (j<m&&f[i+1][j+1][u]>=0)
					if (f[i+1][j+1][u]+mp[i][j]>f[i][j][u0])
					{
						dir[i][j][u0]='L';
						f[i][j][u0]=f[i+1][j+1][u]+mp[i][j];
					}
			}
		}
	ans0=-1;
	for (int i=1;i<=m;i++) 
	{
		if (ans0<f[1][i][0])
		{
			ans0=f[1][i][0];
			ansi=i;
		}
	}
	printf("%d\n",ans0);
	if (ans0>=0)
	find(1,ansi,0);
}