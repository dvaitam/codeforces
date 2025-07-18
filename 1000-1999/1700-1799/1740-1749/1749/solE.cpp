#include <bits/stdc++.h>
#define F first
#define S second
#define pb push_back
#define ppb pop_back
#define fast_io ios::sync_with_stdio(false);cin.tie(NULL);
#define file_io freopen("input.txt","r",stdin);freopen("output.txt","w",stdout);
#define FOR(i,k,n) for(int i = k; i < n; i ++)
#define debf cout<<"(0-0)\n";
#define all(x) x.begin(), x.end()
#define dec(x) //cout << fixed << setprecision(x);
#define pf push_front
#define ppf pop_front
#define dash " ------- "

using namespace std;

typedef long long ll;
typedef pair <int, int> pii;
typedef pair <pii, int> ppi;
typedef pair <int, pii> pip;
typedef pair <ll, ll> pll;
typedef unsigned long long ull;
typedef long double ld;

template <class T> using max_heap = priority_queue <T, vector <T>, less <T> >;
template <class T> using min_heap = priority_queue <T, vector <T>, greater <T> >;

constexpr int MOD = 1e9 + 7, N = 2e5 + 8, INF = 1e9 + 8, LGN = 20, mod = 998244353;
constexpr long double eps = 1e-12, pi = 3.14159265359;

bitset <N * 2> mark;
int t, n, m, dis[N * 2], par[N * 2], dx[4] = {0, 0, 1, -1}, dy[4] = {1, -1, 0, 0}, ututt;
string s[N];

void mkp (int x){
	//////cout << s[1] << " " << s[0];
	while (x % m != 0){
		/*int tmp = x;
		if (s[x / m][x % m] == '#'){
			if (x + m + 1 < m * n && x % m != m - 1 && dis[x] == dis[x + m + 1]){
				tmp = x + m + 1;
			}
			if (x + m - 1 < mn && dis[x] == dis[x + m - 1]){
				tmp = x + m - 1;
			}
			if (x - m + 1 > 0 && x % m != m - 1 && dis[x - m + 1] == dis[x]){
				tmp = x - m + 1;
			}
			if (x - m - 1 > 0 && dis[x] == dis[x - m - 1]){
				tmp = x - m - 1;
			}
		}
		else {
			if (x + m + 1 < m * n && x % m != m - 1 && dis[x] == dis[x + m + 1] + 1){
				tmp = x + m + 1;
			}
			if (x + m - 1 < mn && dis[x] == dis[x + m - 1] + 1){
				tmp = x + m - 1;
			}
			if (x - m + 1 > 0 && x % m != m - 1 && dis[x - m + 1] + 1 == dis[x]){
				tmp = x - m + 1;
			}
			if (x - m - 1 > 0 && dis[x] == dis[x - m - 1] + 1){
				tmp = x - m - 1;
			}
		}*/
		s[x / m][x % m] = '#';
		x = par[x];
	}
	s[x / m][0] = '#';
}

bool isok (int x){
	int a = x / m, b = x % m;
	for (int i = 0; i < 4; i ++){
		//cout << a + dx[i] + 1 << " " << b + dy[i] + 1 << " check shod\n";
		if (a + dx[i] >= 0 && a + dx[i] < n && b + dy[i] >= 0 && b + dy[i] < m && s[a + dx[i]][b + dy[i]] == '#'){
			return false;
		}
	}
	return true;
}

int main(){
	fast_io;
	cin >> t;
	while (t --){
		cin >> n >> m;
		FOR (i, 0, n){
			cin >> s[i];
		}
		//cout << dash << isok (4) << '\n';
		deque <int> q;
		fill (dis, dis + n * m + 1, INF);
		FOR (i, 0, n * m){
			mark[i] = false;
		}
		FOR (i, 0, n){
			dis[i * m] = 0;
			if (!isok (i * m)){
				continue;
			}
			if (s[i][0] != '#'){
				dis[i * m] ++;
			}
			if (dis[i * m]){
				q.pb(i * m);
			}
			else {
				q.pf(i * m);
			}
		}
		while (q.size()){
			int x = q.front();
			q.ppf();
			if (x % m == m - 1 || mark[x]){
				continue;
			}
			mark[x] = true;
			if (x - m + 1 > 0 && dis[x - m + 1] > dis[x] && isok (x - m + 1)){
				par[x - m + 1] = x;
				if (s[(x - m + 1) / m][(x - m + 1) % m] == '#'){
					dis[x - m + 1] = dis[x];
					q.pf(x - m + 1);
				}
				else {
					dis[x - m + 1] = dis[x] + 1;
					q.pb(x - m + 1);
				}
			}
			if (x + m + 1 < m * n && dis[x + m + 1] > dis[x] && isok (x + m + 1)){
				par[x + m + 1] = x;
				if (s[(x + m + 1) / m][(x + m + 1) % m] == '#'){
					dis[x + m + 1] = dis[x];
					q.pf(x + m + 1);
				}
				else {
					dis[x + m + 1] = dis[x] + 1;
					q.pb(x + m + 1);
				}
			}
			if ((x + m - 1) < m * n && (x + m - 1) % m == x % m - 1 && x % m != 0 && dis[x + m - 1] > dis[x] && isok (x + m - 1)){
				par[x + m - 1] = x;
				if (s[(x + m - 1) / m][(x + m - 1) % m] == '#'){
					dis[x + m - 1] = dis[x];
					q.pf(x + m - 1);
				}
				else {
					dis[x + m - 1] = dis[x] + 1;
					q.pb(x + m - 1);
				}
			}
			if (x - m - 1 > 0 && (x - m - 1) % m == x % m - 1 && x % m != 0 && dis[x - m - 1] > dis[x] && isok (x - m - 1)){
				par[x - m - 1] = x;
				if (s[(x - m - 1) / m][(x - m - 1) % m] == '#'){
					dis[x - m - 1] = dis[x];
					q.pf(x - m - 1);
				}
				else {
					dis[x - m - 1] = dis[x] + 1;
					q.pb(x - m - 1);
				}
			}
		}
		int ans = INF, stp = -1;
		for (int i = 0; i < n; i ++){
			if (dis[i * m + m - 1] < ans){
				ans = dis[i * m + m - 1];
				stp = i * m + m - 1;
			}
		}
		if (ans == INF){
			cout << "NO\n";
		}
		else {
			cout << "YES\n";
			mkp (stp);
			for (int i = 0; i < n; i ++){
				cout << s[i] << '\n';
			}
		}
		//for (int i = 0; i < n * m; i ++){
		//	cout << "dis " << i + 1 << " = " << dis[i] << " parent " << i +  +  +  + 1111 << " = " <<  par[i] + 1 << '\n';
		//}
	}
	return 0;
}

// Caught in the moment, not even thinkin' twice
// Everything's frozen, nothing but you and I
// Can't stop my heart from beating
// Why do I love this feeling?