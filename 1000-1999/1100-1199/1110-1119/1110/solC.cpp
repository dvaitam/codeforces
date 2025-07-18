#include<iostream>
#include<bits/stdc++.h>
#include<algorithm>
typedef long long int ll;
typedef unsigned long long int ull;
typedef long double ld;
typedef std::pair<int, int> pii;
typedef std::pair<ll, ll> pll;
typedef std::pair<double, double> pdd;
using namespace std;


#define fastio() ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0)
#define loop(i,a,b) for(ll i=a;i<b;i++)
#define test() ull t;cin>>t;while(t--)
#define pb push_back
#define mkp make_pair
#define nl cout<<"\n"
#define MOD 1000000007
#define trace(s,a) cout<<s<<a<<" ",nl
#define f first
#define s second
#define oa(a,n) loop(i,0,n)cout<<a[i]<<" ";nl;
#define ov(a) loop(i,0,a.size()) cout<<a[i]<<" ";nl;

ll gcd(ll a,ll b)
{
	if(b==0)
		return a;
	else
		return gcd(b,a%b);
}

ll power[27]={};
int main() 
{
	fastio();
	#ifndef ONLINE_JUDGE
	freopen("input.txt", "r", stdin);
	freopen("output.txt", "w", stdout);
	#endif
	power[0]=1;
	ll arr[]={0,1,1,5,1,21,1,85,73,341,89,1365,1,5461,4681,21845,1,87381,1,349525,299593,1398101,178481,5592405,1082401};
	loop(i,1,27)
	{
		power[i]=power[i-1]*2;
	}
	test()
	{
		ll n;
		cin>>n;
		ll len=0;
		while(n>=power[len])
		{
			len++;
		}
		ll ans;
		if(n==power[len]-1)
		{
			ans = arr[len-1];
		}
		else
		{
			ans = power[len]-1;
		}
		cout<<ans,nl;
	}

	return 0;
}