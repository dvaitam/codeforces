#include<bits/stdc++.h>

using namespace std;

#define ll long long

#define vll vector<ll>

#define endl '\n'

int32_t main(){

    ios_base::sync_with_stdio(0);

    cin.tie(0);cout.tie(0);

    ll t;cin>>t;while(t--){

        ll n,k;cin>>n>>k;

        ll l=1,r=n,d=1;

		while(l<=r)cout<<((d^=1)?l++:r--)<<' ';

		cout<<endl;

    }

    return 0;

}