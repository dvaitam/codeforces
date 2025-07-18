#include <bits/stdc++.h>

#define N 3005

typedef long double db;

using namespace std;

db f[N][N];

int main()
{
	int i,j,k,n,m;
	scanf("%d %d",&n,&m);
	for(i=1;i<=n;++i)for(j=1;j<=m;++j)f[i][j]=1e18;
	for(i=1;i<=m;++i)f[1][i]=0;
	for(i=1;i<n;++i)
	{
		for(j=1;j<=m;++j)
			for(k=1;j+k*(n-i)<=m;++k)
				f[i+1][j+k]=min(f[i+1][j+k],f[i][j]+1.0*j/k);
	}
	db ans=1e18;for(i=1;i<=m;++i)ans=min(ans,f[n][i]);
	printf("%.12Lf\n",n+2*ans);
	return 0;
}