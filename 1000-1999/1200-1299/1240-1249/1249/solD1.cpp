#include <bits/stdc++.h>
#define rep(i, n) for(int i = 0; i < (int)(n); i++)
#define all(x) (x).begin(),(x).end()
#define rall(x) (x).rbegin(),(x).rend()
#define UNIQUE(v) v.erase( unique(v.begin(), v.end()), v.end() );
using namespace std;
using ll = long long;
template<class T> inline bool chmax(T& a, T b) { if (a < b) { a = b; return 1; } return 0; }
template<class T> inline bool chmin(T& a, T b) { if (a > b) { a = b; return 1; } return 0; }
const ll INF = 1e9;
const ll MOD = 1e9 + 7;

int main(){
    int n, k;
    cin >> n >> k;
    using P = tuple<int, int, int>;
    vector<P> s(n);
    vector<int> x(201);
    rep(i, n){
        int l, r;
        cin >> l >> r;
        s[i] = P(r, l, i+1);
    }
    sort(all(s));

    vector<int> ans;
    rep(i, n){
        int l, r, ind;
        tie(r, l, ind) = s[i];
        bool ok = true;
        //cout << "l = " << l << ", r = " << r << ", ind = " << ind << endl;;
        for(int j=l; j<=r; j++){
            if(x[j]+1 > k) ok = false;
        }
        if(ok){
            for(int j=l; j<=r; j++) x[j]++;
        }
        else ans.push_back(ind);
    }

    cout << ans.size() << endl;
    rep(i, ans.size()) cout << ans[i] << " ";
    cout << endl;
}