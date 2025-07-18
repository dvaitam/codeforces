/**
 *    author:  AgentPengin ( Độc cô cầu bại )
 *    created: 23.12.2022 10:08:02
 *    too lazy to update time
**/
#include<bits/stdc++.h>

#define EL '\n'
#define fi first
#define se second
#define NAME "TASK"
#define ll long long
#define lcm(a,b) (a/gcd(a,b))*b
#define db(val) "["#val" = " << (val) << "] "
#define bend(v) (v).begin(),(v).end()
#define sz(v) (int)(v).size()
#define ex exit(0)
#define int ll

using namespace std;

const ll mod = 1e9 + 7;
const int inf = 0x1FFFFFFF;
const int MAXN = 2005;

int n, m, Q, sz[MAXN];
set<pair<int,int>> S;
queue<pair<int,int>> q;

bitset<MAXN> foo[MAXN], pre[MAXN];

bool chk(int x, int y) {
	if (x == -1 || y == -1) return 0;
	if (sz[x] > sz[y]) swap(x, y);
	return (foo[x] & foo[y]) != foo[x];
}

void add(int x) {
	set<pair<int,int>>::iterator it1 = S.lower_bound(make_pair(sz[x], x)), it2 = it1;
	it2--;
	S.insert(make_pair(sz[x], x));
	q.push(make_pair(x, (*it2).se));
	q.push(make_pair((*it1).se, x)); 
}

void del(int x) {
	set<pair<int,int>>::iterator it1 = S.lower_bound(make_pair(sz[x], x)), it2 = it1, it3 = it1;
	it3++;
	it2--;
	S.erase(make_pair(sz[x], x));
	q.push(make_pair((*it2).se, (*it3).se));
}

signed main() {
    ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);
    if (ifstream(NAME".inp")) {
        freopen(NAME".inp","r",stdin);
        freopen(NAME".out","w",stdout);
    }
    cin >> n >> m >> Q;
    S.insert(make_pair(1e9, -1));
    S.insert(make_pair(-1, -1));
    for (int i = 1;i <= n;i++) add(i);
    for (int i = 1;i <= m;i++) pre[i] = pre[i - 1], pre[i].flip(i);
    while(Q--) {
    	int x, l, r;
    	cin >> x >> l >> r;
    	del(x);
    	foo[x] ^= pre[l - 1]^pre[r];
    	sz[x] = foo[x].count();
    	add(x);
    	while(!q.empty() && !chk(q.front().fi, q.front().se)) q.pop();
    	if (q.empty()) cout << -1 << '\n';
    	else {
    		int x = q.front().fi, y = q.front().se;
    		if (x > y) swap(x, y);
    		bitset<MAXN> intersect = foo[x]&foo[y];
    		int id1 = (foo[x]^intersect)._Find_first();
    		int id2 = (foo[y]^intersect)._Find_first();
    		if (id1 > id2) swap(id1, id2);
    		cout << x << " " << id1 << " " << y << " " << id2 << '\n';
    	}
    }
    
    
    cerr << "\nTime elapsed: " << 1000 * clock() / CLOCKS_PER_SEC << "ms\n";
    return 0;
}
// agent pengin wants to take apio (with anya-san)