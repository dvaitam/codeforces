#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef vector<ll> vl;
typedef pair<ll,ll> pl;
typedef vector<pl> vp;
typedef map<ll,ll> mll;
typedef unordered_map<ll,vl> mlv;
#define fr(a,c) for(int a = 0; a < c; a++)
#define f3(a,b,c) for(int a=b;a<c; a++)
#define fl(a,c) for(int a=c; a>=0; a--)
#define pb push_back
#define mp make_pair
#define pob pop_back
#define mod 1000000007
#define F first
#define S second
#define all(x) x.begin(), x.end()
#define sortall(x) sort(all(x))
void addedge(ll x,ll y, mlv & adjlst);
void Hemanth_hp(){
	int n, m;
	cin >> n >> m;
	if (m == 1) cout << 0 << '\n';
	else if (n > m - 1) cout << m << '\n';
	else cout << n + 1 << '\n';
	for (int i = 0; i < min(m - 1, n); i++) {
		for (int j = 0; j < m; j++) {
			cout << (j + i) % m << ' ';
		}
		cout << '\n';
	}
	if (n > m - 1) {
		for (int i = m - 1; i < n; i++) {
			for (int j = 0; j < m; j++) {
				cout << j << ' ';
			}
			cout << '\n';
		}
	}
	return;
}

 
int main(){
ios_base::sync_with_stdio(false);
cin.tie(NULL); cout.tie(NULL);
    int t=1;
    cin>>t;
    while(t--){
        Hemanth_hp();
    }
    return 0;
}