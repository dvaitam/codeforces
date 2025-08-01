#include <bits/stdc++.h>
#include<cstdio>
#include<cstring>
#include<cmath>
#include<cstdlib>
#include<algorithm>
#include<string>
#include<vector>
#include<map>
using namespace std;




double pro[2000],dp[2000][200];
long long getn(long long a)
{
	if(a==0) return 0;
	long long b,c=1,ans=0;
	for(int i=0;i<20;i++) 
	{
		if(a/c==0)
		{
			c/=10;
			break;
		}
		c*=10;
	}
	if(a/c!=1) ans+=c;
	else ans+=a%c+1;
	c/=10;
	while(c) 
	{
		ans+=c;
		c/=10;
	}
	return ans;
}
int main()
{
//	freopen("1.txt","r",stdin);
	long long l,r;
	int k,n;
	scanf("%d",&n);
	for(int i=0;i<n;i++)
	{
		scanf("%lld%lld",&l,&r);
		pro[i]=1.0*(getn(r)-getn(l-1))/(r-l+1);
	}
	scanf("%d",&k);
	k=(n*k+99)/100;
	dp[0][0]=1.0;
	for(int i=1;i<=n;i++)
	{
		for(int j=0;j<=i && j<=k;j++)
		{
			if(j==0) dp[i][j]=dp[i-1][j]*(1-pro[i-1]);
			else dp[i][j]=dp[i-1][j-1]*pro[i-1]+dp[i-1][j]*(1-pro[i-1]);
		}
	}
	double ans=0.0;
	for(int i=0;i<k;i++)
		ans=ans+dp[n][i];
	if(ans<0) ans=0;
	if(ans>1.0) ans=1.0;
	printf("%.12lf\n",1.0-ans);
	return 0;
}