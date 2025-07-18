#pragma GCC target ("avx2")
#pragma GCC optimization ("O3")
#pragma GCC optimization ("unroll-loops")
#include <bits/stdc++.h>
#define owo(i,a, b) for(int i=(a);i<(b); ++i)
#define uwu(i,a, b) for(int i=(a)-1; i>=(b); --i)
#define senpai push_back
#define ttgl pair<int, int>
#define ayaya cout<<"ayaya~"<<endl
 
using namespace std;
using ll = long long;
using ld = long double;
ll MOD = 998244353;
const ll root = 3;
ll binpow(ll a,ll b){ll res=1;while(b){if(b&1)res=(res*a)%MOD;a=(a*a)%MOD;b>>=1;}return res;}
ll modInv(ll a){return binpow(a, MOD-2);}
const int INF = 0x3f3f3f3f;
const int NINF = 0xc0c0c0c0;
const ll INFLL = 0x3f3f3f3f3f3f3f3f;
const ll NINFLL = 0xc0c0c0c0c0c0c0c0;
const int mxN = 100001;
int m, k;
ttgl arr[mxN];
int sqr(int a) {
    int l = 0;
    int r = a;
    while(l < r) {
        int m = (l + r) >> 1;
        if(1LL*m*m > (ll)a)r = m;
        else l = m+1;
    }
    return l;
}
bool check(int a) {
    int cnt = a*a - (a/2)*(a/2);
    return a * ((a+1)/2) >= arr[0].first && cnt >= m;
}
void solve() {
    cin>>m>>k;
    owo(i, 0, k) {
        cin>>arr[i].first;
        arr[i].second =i+1;
    }
    sort(arr, arr+k);
    reverse(arr, arr+k);
    int l = 1;
    int r = 2*sqr(m);
    while(l < r) {
        int m = (l + r) >> 1;
        if(check(m))r = m;
        else l = m+1;
    }
    vector<vector<int>> res(l, vector<int>(l));
    int id = 0;
    for(int i = 0; i < l; i+=2) {
        for(int j = 1; j < l; j += 2) {
            while(id < k && arr[id].first==0)id++;
            if(id==k) goto don;
            //cout<<i<<" "<<j<<" "<<arr[id].second<<" ??\n";
            res[i][j] = arr[id].second;
            arr[id].first--;
        }
    }
    if(id==0) {
        for(int i = 0; i < l; i+= 2) {
            for(int j = 0; j < l; j += 2) {
                if(arr[id].first==0)break;
                res[i][j] = arr[id].second;
                arr[id].first--;
            }
        }
        id++;
    }
    owo(i, 0, l) {
        owo(j, 0, l) {
            while(id < k && arr[id].first==0)id++;
            if(id==k) goto don;
            if(res[i][j] || ((i&1) && (j&1)))continue;
            res[i][j] = arr[id].second;
            arr[id].first--;
        }
    }
    don: ;
    cout<<l<<"\n";
    owo(i, 0, l) {
        owo(j, 0, l) {
            cout<<res[i][j]<<" ";
        }
        cout<<"\n";
    }
}
int main() {
    //freopen("file.in", "r", stdin);
    //freopen("file.out", "w", stdout);
    mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());
    cin.tie(0)->sync_with_stdio(0);
    int T;
    cin>>T;
    owo(tc, 1, T+1) {
        solve();
    }
    return 0;
}