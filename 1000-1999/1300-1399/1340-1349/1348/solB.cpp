#include <bits/stdc++.h>

 

using namespace std;

//-D_GLIBCXX_DEBUG -fsanitize=undefined

typedef long long int ll;

typedef pair<int,int> pii;

typedef pair<ll,ll> pll;

typedef long double ld;

typedef vector<int> vi;

typedef vector<ll> vl;

#define mp make_pair

#define pb push_back

#define nl "\n"

#define all(v) v.begin(),v.end()

#define fi first

#define se second

#define debug(x) cout << #x << ": " << x << nl

#define rep(x,start,end) for(int x=(start)-((start)>(end));x!=(end)-((start)>(end));((start)<(end)?x++:x--))

#define sz(x) (int)(x).size()



ll n, m, t, a, b, c, k;

string f, s;

void solve(){

	cin >> n>>k;

	vi v(n);

	set<int> ss;

	rep(i,0,n) {cin>>v[i];ss.insert(v[i]);}

	if (sz(ss)>k){

		cout<<-1<<nl;

		return;

	}

	rep(i,1,n+1){

		if (sz(ss)==k){

			break;

		}

		ss.insert(i);

	}

	cout<<sz(ss)*n<<nl;

	rep(i,0,n){

		for (int j:ss){

			cout<<j<<" ";

		}

	}

	cout<<nl;

}

int TC;

int main(){

	ios::sync_with_stdio(0);

	cin.tie(0);



	cin >> TC;

	while (TC-->0){

		solve();

	}

}