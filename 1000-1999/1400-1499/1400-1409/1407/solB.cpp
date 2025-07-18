#include<bits/stdc++.h>

#include<ext/pb_ds/assoc_container.hpp>

#include<ext/pb_ds/tree_policy.hpp>

using namespace std;

using namespace __gnu_pbds;

#define ll long long

typedef tree<ll, null_type, less_equal<ll>, rb_tree_tag, tree_order_statistics_node_update> pbds;



void solve(){

ll n; cin >> n;

vector<ll>v(n), b, g1;

for (ll i=0; i<n; i++){

	cin >> v[i];

}

sort(v.begin(), v.end());

ll mx = v.back(); v.pop_back();

b.push_back(mx);

ll a[n+1];

for (ll i=0; i<=n; i++) a[i]=0;

for (int in=1; in<=100; in++){

	ll mxg=1;

	for (ll i=0; i<v.size(); i++){

		if(a[i] != -1){

			ll g = __gcd(v[i], mx);

		if(g == 1){

			g1.push_back(v[i]);

			a[i]=-1;

		}

		else if (g==mx){

			b.push_back(v[i]);

			a[i]=-1;

		}

		else{

			mxg = max(g, mxg);

		}

	}

	}

	if(mxg==1) break;

	mx = mxg;

//	cout << mx <<'\n';

}

for (auto u:g1) b.push_back(u);

// ll m = b[0];

// for (auto u:b) {

// 	m = __gcd(m,u);

// 	cout << m <<' ';

// }

// cout <<'\n';



 for (auto u:b) cout <<u <<' '; cout <<'\n';

}



int main(){

  ios_base::sync_with_stdio(0);

  cin.tie(0);

 

  ll t = 1;  cin >> t;

  while(t--){

    solve();

  }

}