#include <bits/stdc++.h>
#define pb push_back
#define st first
#define nd second
#define deb(x) cout << x << endl
#define deb2(x) for(auto i: x) cout << i << " "
#define sz(a) int(a.size())
#define all(a) a.begin(), a.end()
 
using namespace std;
 
typedef long long ll;
typedef pair<int, int> pii;
typedef pair<ll, ll> pll;
typedef vector<int> vi;
typedef vector<ll> vll;
 
const int mod = 1e9 + 7;
const int MAXN = 5e3 + 7;
const ll inf = 1e18;
const double pi = double(acos(-1));
const int MAXM=1e5 + 7;
 
struct Cycle{
	int boxmask;
	vll cyclenodes;
};
 
int k;
ll total=0,x;
vector<vll> boxes(15, vll(MAXN));
vll boxsum(15,0),nodes;
map<ll, int> whichbox,processed;
map<ll,ll> nxt;
 
vector<Cycle> validcycles;
 
void solve(){
	cin >> k;
	for(int i=0; i<k; i++){
		int n; cin >> n;
		boxes[i].resize(n);
		for(int j=0; j<n; j++){
			cin >> boxes[i][j];
			whichbox[boxes[i][j]] = i;
			nodes.pb(boxes[i][j]);
			total+=boxes[i][j];
			boxsum[i]+=boxes[i][j];
		}
	}
	if(total%k){
		deb("NO");
		return;
	}
	x=total/k;
	for(int i=0; i<k; i++){
		for(int j=0; j<sz(boxes[i]); j++){
			ll needed = x - boxsum[i] + boxes[i][j];
			if(whichbox.find(needed)!=whichbox.end()) nxt[boxes[i][j]]=needed;
			else nxt[boxes[i][j]]=-1;
		}
	}
	for(auto node: nodes){
		if(processed[node]==1) continue;
		vll cyclenodes;
		int currmask=0;
		bool foundcycle=false;
		map<ll,int>position;
		ll start=node,cur = node;
		while(cur != -1){
			if(position.find(cur)!=position.end()){
				foundcycle = true;
				break;
			}
			position[cur]=sz(cyclenodes);
			cyclenodes.pb(cur);
			cur=nxt[cur];
		}
		if(foundcycle){
			int pos = position[cur];
			vll candcycle;
			int cycleboxmask=0;
			bool valid = true;
			for(int i=pos; i<sz(cyclenodes); i++){
				ll cand=cyclenodes[i];
				int candbox=whichbox[cand];
				if(cycleboxmask & (1<<candbox)){
					valid = false;
					break;
				}
				cycleboxmask |= (1<<candbox);
				candcycle.pb(cand);
			}
			if(valid) validcycles.pb({cycleboxmask, candcycle});
	}
		for(auto v: cyclenodes)
			processed[v]=1;
	}
	vector<vi> cycles_by_mask(1<<k);
	for(int i=0; i<sz(validcycles); i++){
		int m = validcycles[i].boxmask;
		cycles_by_mask[m].pb(i);
	}
	int fullmask = (1<<k)-1;
	vi dp(1<<k,0),parent(1<<k,-1),whichcycle(1<<k,-1),f(1<<k,0),used(1<<k);
	dp[0]=1;
	for(int i=0; i<sz(validcycles); i++) f[validcycles[i].boxmask]=1;
 
	for(int mask = 0; mask < (1 << k); mask++){
		if(!dp[mask])continue;
		int remain = ((1<<k) - 1)^mask;
		for(int s=remain; s>0; s=(s-1)&remain){
			if(f[s]){
				int newmask = mask | s;
				if(!dp[newmask]){
					dp[newmask]=1;
					parent[newmask]=mask;
					used[newmask]=s;
				}
			}
		}		
	}
	if(dp[fullmask]){
		deb("YES");
		int curmask=fullmask;
		vi res, chosensm;
		vector<pii> ans,ansmoves;
		while(curmask){
			int sm = used[curmask];
			chosensm.pb(sm);
			curmask=parent[curmask];
		}
		reverse(all(chosensm));
		for(int sm: chosensm){
			int ci = cycles_by_mask[sm][0];
			Cycle C = validcycles[ci];
			reverse(all(C.cyclenodes));
			int t=sz(C.cyclenodes);
			for(int i=0; i<t; i++){
				ll number=C.cyclenodes[i];
				int src=whichbox[number];
				int dest=whichbox[C.cyclenodes[(i+1)%t]];
				ans.pb({number,dest+1});
			}
		}	
		sort(all(ans), [&](const pii sigma1, const pii sigma2){
			return whichbox[sigma1.st] < whichbox[sigma2.st];
		});
		for(auto gowno: ans){
			cout << gowno.st << " " << gowno.nd << endl;
		}
	} else {
		deb("NO");
	}
}
 
int main(){
	ios_base::sync_with_stdio(0);
	cin.tie();
	cout.tie();
	int t = 1;
	//cin >> t;
	while(t--)
		solve();
}