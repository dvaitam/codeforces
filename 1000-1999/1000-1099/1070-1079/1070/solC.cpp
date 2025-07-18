/*
 ███╗   ███╗ █████╗ ██████╗ ███████╗    ██████╗ ██╗   ██╗       ███╗   ███╗ █████╗ ██████╗  ██████╗██╗  ██╗███████╗██╗     ██╗           ██╗  ██╗██╗██╗
 ████╗ ████║██╔══██╗██╔══██╗██╔════╝    ██╔══██╗╚██╗ ██╔╝██╗    ████╗ ████║██╔══██╗██╔══██╗██╔════╝██║  ██║██╔════╝██║     ██║           ██║  ██║██║██║
 ██╔████╔██║███████║██║  ██║█████╗      ██████╔╝ ╚████╔╝ ╚═╝    ██╔████╔██║███████║██████╔╝██║     ███████║█████╗  ██║     ██║           ███████║██║██║
 ██║╚██╔╝██║██╔══██║██║  ██║██╔══╝      ██╔══██╗  ╚██╔╝  ██╗    ██║╚██╔╝██║██╔══██║██╔══██╗██║     ██╔══██║██╔══╝  ██║     ██║           ██╔══██║██║██║
 ██║ ╚═╝ ██║██║  ██║██████╔╝███████╗    ██████╔╝   ██║   ╚═╝    ██║ ╚═╝ ██║██║  ██║██║  ██║╚██████╗██║  ██║███████╗███████╗███████╗█████╗██║  ██║██║██║
 ╚═╝     ╚═╝╚═╝  ╚═╝╚═════╝ ╚══════╝    ╚═════╝    ╚═╝          ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝╚════╝╚═╝  ╚═╝╚═╝╚═╝
*/

/*
 ඞ
 ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣤⣤⣤⣤⣤⣤⣤⣤⣄⡀
 ⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⡿⠛⠉⠉⠉⠉⠉⠉⠻⢿⣿⣷⡄
 ⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⠋⠀⠀⠀⠀⠀⠀⠀      ⢻⣿⣿⡄
 ⠀⠀⠀⠀⠀⠀⠀⣸⣿⡏⠀⠀⠀   ⣠⣾⣿⣿⣿⠿⠿⠿⢿⣿⣿⣄
 ⠀⠀⠀⠀⠀⠀⠀⣿⣿⠁⠀⠀    ⣿⣿⣯⠁⠀⠀⠀⠀⠀  ⠙⢿⣷⡄
 ⠀⠀⣀⣤⣴⣶⣶⣿⡟⠀⠀⠀   ⣿⣿⣿            ⣿⣷
 ⠀⢰⣿⡟⠋⠉⣹⣿⡇⠀⠀⠀   ⣿⣿⣿⣷⣦⣀⣀⣀⣀⣀⣀⣀⣿⣿
  ⢸⣿⡇   ⣿⣿⡇⠀⠀     ⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣸⣿⡇   ⣿⣿⡇⠀⠀⠀⠀    ⠉⠉⠉⠉⠉⠉⠉⠉⠉⡿⢻⡇
 ⠀⣿⣿    ⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀       ⢸⣿⡇
 ⠀⣿⣿⠀   ⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⢸⣿⡇
 ⠀⣿⣿⠀⠀  ⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⢸⣿⡇
 ⠀⢿⣿   ⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀         ⣿⡇
 ⠀⠸⣿⣦⡀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀       ⣿⣿
 ⠀⠀⠛⢿⣿⣿⣿⣿⡇    ⠀⣠⣿⣿⣿⣿⣄       ⣿⣿
 ⠀⠀⠀⠀⠀⠀⠀⣿⣿⠀⠀⠀⠀⠀⣿⣿⡇⠀⣽⣿⡆     ⢸⣿⡇
 ⠀⠀⠀⠀⠀⠀⠀⣿⣿⠀⠀⠀⠀⠀⣿⣿⡇⠀⢹⣿⡆⠀⠀   ⣸⣿⠇
 ⠀⠀⠀⠀⠀⠀⠀⢿⣿⣦⣄⣀⣠⣴⣿⣿ ⠀⠈⠻⣿⣿⣿⣿⡿⠏
 ⠀⠀⠀⠀⠀⠀⠀⠈⠛⠻⠿⠿⠿⠿⠋⠁
*/

#include<bits/stdc++.h>
#define ll long long
#define pii pair<int, int>
#define pll pair<ll, ll>
#define plll pair<pll, pll>
#define pdd pair<double, double>
#define pu push_back
#define po pop_back
#define fi first
#define se second
#define fifi fi.fi
#define fise fi.se
#define sefi se.fi
#define sese se.se
#define cekcek cout<<'c'<<'e'<<'k'<<endl
using namespace std;

ll N, K, M, l[200005], r[200005], c[200005], p[200005], it, ans, t;
vector<pair<ll, pll>> v;
pll tree[4000005]; //{num, tot price}

bool compare(pair<ll, pll> x, pair<ll, pll> y){
    if(x.fi == y.fi) return x.sefi < y.sefi;
    return x.fi < y.fi;
}

void update(ll tidx, ll tl, ll tr, ll x, ll y){
    if(tl == tr){
        tree[tidx].fi += y;
        tree[tidx].se += y * tl;
        return;
    }
    ll tmid = (tl + tr) / 2;
    if(x <= tmid) update(tidx * 2, tl, tmid, x, y);
    else update(tidx * 2 + 1, tmid + 1, tr, x, y);
    tree[tidx].fi = tree[tidx * 2].fi + tree[tidx * 2 + 1].fi;
    tree[tidx].se = tree[tidx * 2].se + tree[tidx * 2 + 1].se;
    return;
}

ll F(ll tidx, ll tl, ll tr, ll x){
    if(x == 0) return 0;
    if(tl == tr) return min(x, tree[tidx].fi) * tl;
    ll tmid = (tl + tr) / 2;
    if(tree[tidx * 2].fi <= x) return tree[tidx * 2].se + F(tidx * 2 + 1, tmid + 1, tr, x - tree[tidx * 2].fi);
    return F(tidx * 2, tl, tmid, x);
}

int main(){
ios_base::sync_with_stdio(0); cin.tie(NULL); cout.tie(NULL);
cin >> N >> K >> M;
for(ll j = 1; j <= M; j++){
    cin >> l[j] >> r[j] >> c[j] >> p[j];
    v.pu({l[j], {0, j}});
    v.pu({r[j], {1, j}});
    t = max(t, p[j]);
}
sort(v.begin(), v.end(), compare);
it = 0;
for(ll i = 1; i <= N; i++){
    while(it < v.size() && v[it].fi == i && v[it].sefi == 0){
        update(1, 1, t, p[v[it].sese], c[v[it].sese]);
        it++;
    }
    ans += F(1, 1, t, K);
    while(it < v.size() && v[it].fi == i && v[it].sefi == 1){
        update(1, 1, t, p[v[it].sese], -c[v[it].sese]);
        it++;
    }
}
cout << ans << endl;
}