#include<bits/stdc++.h>
#define ll long long
using namespace std;
ll e;
main()
{
	ll w,h,a,b,i;
	cin>>w>>h>>a>>b;
	ll x=__gcd(a,b);
	a=a/x;
	b=b/x;
	ll c=w/a;
	ll d=h/b;
	cout<<min(c,d);
}