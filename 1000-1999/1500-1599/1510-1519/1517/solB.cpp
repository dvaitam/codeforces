#include "bits/stdc++.h"

using namespace std;

#define ll long long

#define nln "\n"

#define vln vector<long long>

#define all(x) x.begin(),x.end()

#define vi vector<int>

#define yes cout<<"YES"<<nln; return;

#define no cout<<"NO"<<nln; return;

#define pb push_back

#define ppb pop_back

#define mo map<ll,ll>

#define vlp vector<pair<ll,ll>> 

#define vs vector<string>

#define fr(x,y) for( ll i = x ; i < y ; ++i)

#define frr(x) for( ll j = 0 ; j < x ; ++j)

#define mn(x) *min_element(all(x))

#define mx(x) *max_element(all(x))







void init_code(){

	#ifndef ONLINE_JUDGE

	freopen("inputf.in","r",stdin);

	freopen("outputf.in","w",stdout);

	#endif

	}

int alp[26];

	ll expo(ll a, ll b , ll m){ll result =1; while(b>0){ if(b&1)result = (result*a) % m; a= a*a%m; b = (b>>1); } return result;}

	ll mminvprime(ll a, ll b) {return expo(a, b - 2, b);}

	void disp(){fr(0,26)cout<<char(i+'a')<<" : "<<alp[i]<<nln;}

	void disp(vln &v){for(auto x:v)cout<<x<<" "; cout<<nln;}

	bool sortbyfir(const pair<int,int> &a,const pair<int,int> &b){return (a.first < b.first);}

	int palind(string s){string t=s;reverse(all(t));for (int i = 0; i < s.size(); ++i){if(t[i]!=s[i]) return -1;}return 1;}

	void sub_min(ll &a,ll &b){ll c= min(a,b);b=b-c;a=c-a;}

	ll apowerb(int a,int b){if(b==0) return 1; if(b&1){return a*apowerb(a,b/2)*apowerb(a,b/2);}else return apowerb(a,b/2)*apowerb(a,b/2);}

	ll mod_fact(ll n, ll mod){if(n >= mod ){ return 0;}else{n = n%mod;ll result =1;if(n==0){return result;}else{for(ll i = 1;i <=n ; i++) {	result = (result*i)%mod;}}return result;}}

	ll mod_mul(ll a,ll b,ll m){a=a%m; b=b%m;return ((((a * b) % m) + m) % m);}

	ll mod_div(ll a, ll b, ll m){a=a%m; b=b%m; return ((mod_mul(a, mminvprime(b,m), m) + m) % m) ;}

	ll Cmod(ll n,ll r,ll mod){if(n>=r){ll a = mod_fact(n,mod);ll b = mod_fact(r,mod);ll c = mod_fact(n-r,mod);ll d = mod_mul(b,c,mod);ll ans = mod_div(a,d,mod);return ans;} else return 0;}

	ll findmex(set<ll> &v){ll x=0,found=1;while(found){ auto it = v.find(x); if(it == v.end()){return x;}else x++;}return x;}

	bool is_palin(vln &v){ll nn = v.size();ll n=nn/2;if(nn&1)n++;fr(0,n){if(v[i] != v[(nn-1)-i]){ return false;}}return true;}

	void prt_stack(stack<int> &s){if(!s.empty()){int x = s.top();s.pop();prt_stack(s);cout<<x<<" ";}}

	void vlninput(vln &v){for (auto& u : v)cin >> u;}

	void vlpinput(vlp &v){for (auto& u : v)cin >> u.first>>u.second;}

	void dispvlp(vlp &v){for (auto& u : v)cout << u.first<<" "<<u.second<<nln;}

	void dispmxq(priority_queue<ll>&mxq){cout<<nln;while(!mxq.empty()){ cout<<mxq.top()<<" "; mxq.pop();}cout<<nln;}

	void dispmxq(priority_queue<ll,vln,greater<ll>>&mnq){cout<<nln;while(!mnq.empty()){ cout<<mnq.top()<<" "; mnq.pop();}cout<<nln;}











////////////////////////	solve /////////////////////////////////////









	// ll binser(vln&v,ll x) {

	// 	// disp(v);

	// 	auto it = lower_bound(all(v),x);

	// 	// cout<<"it"<<*it;

	// 	if(it == v.end()) return -1;

	// 	else{

	// 			if(*it == x) { ll id = it - v.begin(); return id;}

	// 			else return -1;

	// 	}

	// }









	void solve(){

	ll n,m;

	cin>>n>>m;

	vector<vln> v(n,vln(m,0));

	fr(0,n) frr(m) cin>>v[i][j];



	priority_queue<ll> mnp;

	// map<pair<int,int>,pair<int,int>> mp;

	fr(0,n) frr(m) { mnp.push(v[i][j]); if(mnp.size() > m) mnp.pop(); }  

	ll k = 0;

	for(auto &vec:v) {sort(all(vec),greater<ll>());}

	ll M=m;

	while(m>0) {

		ll x = mnp.top();

		// cout<<x<<nln;

		for(int i=0;i<n;i++) {

			// ll bins = binser(v[i],x);

				auto it = find(v[i].begin()+k,v[i].end(),x);

			if(it != v[i].end()) {

				ll bins = it-v[i].begin();

				// cout<<"x = "<<x<<" at " << bins <<" to "<<k<<nln;

				// cout<<bins<<nln;

				// swap

				ll temp = v[i][k];

				v[i][k] = v[i][bins];

				v[i][bins] = temp;

				k++;

				mnp.pop();

				break;

			}

		}

		m--;

	}

// cout<<nln;

	fr(0,n){ frr(M) cout<<v[i][j]<<" "; cout<<nln;}

    	}



































//////////////////////////////     main   ////////////////////////////////////





int main(){

// init_code();

ios_base::sync_with_stdio(false);

cin.tie(NULL);

	int t;

	cin>>t;

	while(t--){

	solve();

	// cout<<nln;

	}

	return 0;

}