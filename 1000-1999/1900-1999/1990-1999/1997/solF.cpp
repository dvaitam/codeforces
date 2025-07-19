#include <bits/stdc++.h>
#define ll long long
using namespace std;
 
ll fib[55],dp[1005][55005];
const ll mod=998244353;
 
int main() {
	fib[1]=fib[2]=1;
	for (ll i=3;i<=30;i++)fib[i]=fib[i-1]+fib[i-2];
	ll n,x,m;
	cin >> n >> x >> m;
	dp[0][0]=1;
	for (ll i=1;i<=x;i++){
		for (ll j=1;j<=n;j++){
			for (ll l=fib[i];l<=fib[i]*j;l++){
				dp[j][l]+=dp[j-1][l-fib[i]];
				dp[j][l]%=mod;
			}
		}
	} 
	ll ans=0;
	for (ll i=0;i<=fib[x]*n;i++){
		ll t=i,c=0;
		for (ll j=30;j>=1;j--){
			c+=t/fib[j];
			t%=fib[j]; 
		}
		if (c==m)ans+=dp[n][i],ans%=mod;
	}
	cout<< ans;
	return 0;
}