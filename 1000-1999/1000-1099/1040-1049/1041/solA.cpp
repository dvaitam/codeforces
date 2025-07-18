#include <iostream>
#include <bits/stdc++.h>
#include <set>
using namespace std;

#define ll long long
#define endl "\n"
#define pll pair<ll, ll>
#define ppll pair<ll, pair<ll, ll>>
#define mp make_pair
#define pb push_back

#define sz(x) x.size()
#define fr(i,a,b) for(ll i=a;i<=b;i++)
#define repr(i,n) for(ll i=0ll;i<n;i++)

#define fast_io ios_base::sync_with_stdio(false)
#define accuracy long long int precision = numeric_limits<double>::digits10

#define N 1005

ll n;
int main(){
	fast_io;
	accuracy;
	cin>>n;
	vector<ll> a(n);
	repr(i,n) cin>>a[i];
	sort(a.begin(),a.end());
	ll c=0;
	repr(i,sz(a)-1){
		c+=a[i+1ll]-a[i]-1;
	}
	cout<<c<<endl;
}