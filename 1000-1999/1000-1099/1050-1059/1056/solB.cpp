#include<bits/stdc++.h>

using namespace std;

#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>

using namespace __gnu_pbds;

#define endl '\n'
#define MAX

typedef long long ll;
typedef pair<int,int> pii;
//typedef tree<int,null_type,less<int>,rb_tree_tag, tree_order_statistics_node_update> indexed_set;


int main(){
	ios_base::sync_with_stdio(0);
	cin.tie(0);

	//freopen("input.txt", "r", stdin);
	//freopen("output.txt", "w", stdout);
	
	set<ll> pows;
	
	ll p = 2;
	while(p <= 1000){
		pows.insert(p);
		p *= 2;
	}
	
	ll n, m;
	cin >> n >> m;
	ll ans = 0;
	
	
	ll cycle = m;
	
	
	vector<ll> C;
	
//	ll brute = 0;
//	for(ll i = 1; i <= n; i++){
//		for(ll j = 1; j <= n; j++){
//			if((i * i + j * j) % m == 0)
//				brute++;
//		}
//	}
	
	//for(ll i = 1; i <= n; i++)
	//	cout << (i * i) % m << " \n"[i == n];
		
	//cout << "------------" << endl;
 	
	for(ll i = 1; i <= min(n, cycle); i++)
		C.push_back((i * i) % m);
	
	vector<ll> rems(m);
	
	ll sz = (ll)C.size();
	
	ll cant = n / sz;
	ll mod = n % sz;
	
	for(auto &x : C)
		rems[x] += cant;
		
	for(ll i = 0; i < mod; i++)
		rems[C[i]]++;
	
	//for(ll i = 0; i < m; i++)
	//	cout << rems[i] << " \n"[i + 1 == m];
	
	for(ll i = 0; i < m; i++){
		ll opp = -i;
		if(opp < 0)
			opp += m;
		ans += rems[opp] * rems[i];
	}
	
	cout << ans << endl;
//	cout << brute << endl;

	return 0;
}