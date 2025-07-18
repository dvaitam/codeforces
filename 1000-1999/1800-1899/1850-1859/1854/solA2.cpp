#include<bits/stdc++.h>
using namespace std;
using ll=long long;
const int mod=1e9+7;
void solve()
{
	int n;
	cin>>n;
	vector<int>idn,idp;
	int mx=-50,mxid=0,mn=50,mnid=0;
	for(int i=1;i<=n;i++)
	{
		int x;
		cin>>x;
		if(x>0)
		{
			idp.push_back(i);
		}
		else if(x<0)idn.push_back(i);
		if(mx<x)
		{
			mx=x;
			mxid=i;
		}
		if(mn>x)
		{
			mn=x;
			mnid=i;
		}
	}
	if(mx<=0)
	{
		cout<<n-1<<'\n';
		for(int i=n-1;i>=1;i--)
		{
			cout<<i<<' '<<i+1<<'\n';
		}
		return;
	}
	if(mn>=0)
	{
		cout<<n-1<<'\n';
		for(int i=2;i<=n;i++)
		{
			cout<<i<<' '<<i-1<<'\n';
		}
		return;
	}
	vector<pair<int,int>>ans;
	int k=mx;
	while(k+mn<0)
	{
		ans.push_back({mxid,mxid});
		k*=2;
	}
	for(auto x:idn)
	{
		ans.push_back({x,mxid});
	}
	for(int i=2;i<=n;i++)
	{
		ans.push_back({i,i-1});
	}
	if(ans.size()<=31)
	{
		cout<<ans.size()<<'\n';
		for(auto [x,y]:ans)
		{
			cout<<x<<' '<<y<<'\n';
		}
		return;
	}
	ans.clear();
	k=mn;
	while(k+mx>0)
	{
		ans.push_back({mnid,mnid});
		k*=2;
	}
	for(auto x:idp)
	{
		ans.push_back({x,mnid});
	}
	for(int i=n-1;i>=1;i--)
	{
		ans.push_back({i,i+1});
	}
	if(ans.size()<=31)
	{
		cout<<ans.size()<<'\n';
		for(auto [x,y]:ans)
		{
			cout<<x<<' '<<y<<'\n';
		}
		return;
	}
	return;
}
int main()
{
	ios::sync_with_stdio(0);
	cin.tie(0);
	cout.tie(0);
	int t;
	cin>>t;
	while(t--)
		solve();
	return 0;
}