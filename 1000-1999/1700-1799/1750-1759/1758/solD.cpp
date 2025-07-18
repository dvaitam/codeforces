#include <bits/stdc++.h>



using namespace std;

#define ll long long 

#define sz(x) int(x.size())

#define Rashed_To_Get_Accepted ios_base::sync_with_stdio(false);cin.tie(nullptr);cout.tie(nullptr);

const int N = 100005;

const int MOD = 1e9+7 ;//998244353;

// iota(a.begin() , a.end() , x) , x here means start int



// sort(ord.begin(), ord.end(), [&cnt](int x, int y){

	// return cnt[x] > cnt[y];

// });



// stoi(str) to convert string to int  

// dp?, graph?, bs on answer?, stupid observation?

int dx[] = {1,  -1  ,0 , 0};

int dy[] = {0 , 0 , -1 , 1};



void solve()

{

	int n ; cin >> n;

	ll sq = n+100;

	while(sq++){

		ll mx = sq+1;

		vector<ll> ans(n);

		iota(ans.begin() , ans.end() , 1);

		ans.back() = mx;

		ll rem = sq*sq - accumulate(ans.begin() , ans.end() , 0LL);

		ll inc = rem/n  , mod = rem%n;

		if(mod >n-2)continue;

		for(int i =n-1 ;~i ; i--){

			ans[i] += inc;

			if(i && i+1 <n && mod){

				--mod;

				++ans[i];

			}

		}

		for(auto it : ans)cout <<it <<" ";

		cout <<"\n";

		return;

	}

}



int main(){

	

	Rashed_To_Get_Accepted

	int tc;

	tc = 1;

	cin >> tc;

	while(tc--){

		solve();

	}

	return 0;

}