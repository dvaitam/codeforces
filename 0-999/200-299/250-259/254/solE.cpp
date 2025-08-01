#include <iostream>
#include <cstdio>
#include <algorithm>
#include <cstring>
#include <vector>
#include <map>
#include <set>

#define Min(a,b) (a<b?a:b)
using namespace std;

typedef pair<int,int> PII;

vector<PII> f[401];

int main()
{
	freopen("input.txt","r",stdin);
	freopen("output.txt","w",stdout);
	int n,v;
	cin>>n>>v;
	int a[401];
	for(int i=0;i<n;i++)
		cin>>a[i];
	int m;
	cin>>m;
	for(int i=0;i<m;i++)
	{
		int l,r,fi;
		cin>>l>>r>>fi;
		l--;r--;
		for(int j=l;j<=r;j++)
			f[j].push_back(make_pair(fi,i+1));
	}
	for(int i=0;i<n;i++)
		sort(f[i].begin(),f[i].end());
	int dp[401][401],tp[401][401];
	memset(dp,0xff,sizeof(dp));
	memset(tp,0,sizeof(tp));
	dp[0][0]=0;

	for(int i=0;i<n;i++)
		for(int j=0;j<=400;j++)
		{
			int x=j+a[i]-v;

			if(dp[i][j]<0||x<0)continue;

			if(dp[i+1][Min(a[i],x)]<dp[i][j])
			{
				dp[i+1][Min(a[i],x)]=dp[i][j];
				tp[i+1][Min(a[i],x)]=j;
			}

			for(int k=0,sum=0;k<f[i].size();k++)
			{
				x-=f[i][k].first;
				if(x<0)break;

				if(dp[i+1][Min(a[i],x)]<dp[i][j]+k+1)
				{
					dp[i+1][Min(a[i],x)]=dp[i][j]+k+1;
					tp[i+1][Min(a[i],x)]=j;
				}
			}
		}

	int ans=0,fj;
	for(int i=0;i<=400;i++)
		if(dp[n][i]>ans)
		{
			ans=dp[n][i];
			fj=i;
		}
	cout<<ans<<endl;
	vector<int> nf;
	for(int i=n;i>0;i--)
	{
		int tf=tp[i][fj];
		nf.push_back(dp[i][fj]-dp[i-1][tf]);
		fj=tf;
	}
	reverse(nf.begin(),nf.end());
	for(int i=0;i<n;i++)
	{
		cout<<nf[i];
		for(int j=0;j<nf[i];j++)
			cout<<" "<<f[i][j].second;
		cout<<endl;
	}
	return 0;
}