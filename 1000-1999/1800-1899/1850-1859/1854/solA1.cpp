/*
                   الْحَمْدُ لِلَّهِ وَحْدَهُ
    وَالصَّلاَةُ وَالسَّلاَمُ عَلَى مَنْ لاَ نَبِيَّ بَعْدَهُ

    All praise is due to Allah alone,
    and peace and blessings be upon him after
    whom there is no other Prophet.
*/
#include <algorithm>
#include <iostream>
#include <bitset>
#include <cmath>
#include <cstring>
#include <iomanip>
#include <map>
#include <numeric>
#include <queue>
#include <set>
#include <vector>
using namespace std;
// #include<ext/pb_ds/assoc_container.hpp>
// using namespace __gnu_pbds;

typedef long long               ll;
typedef vector<int>             vi;
typedef pair<int, int>          pii;
typedef map<int,int>            mii;
typedef set<int>                si;
//typedef tree<int,null_type,less_equal<int>,rb_tree_tag,tree_order_statistics_node_update>ordered_multiset;

#define cinv(ar)        for(auto &i: ar) cin>>i;
#define coutv(ar)       for(auto i: ar) cout<<i<<' '; cout<<'\n';
#define nl              '\n'
#define in              int i = 0; i < n; ++i
#define in1(x)          int x; cin>>x
#define in2(x,y)        int x, y; cin>>x>>y
#define in3(x,y,z)      int x,y,z; cin>>x>>y>>z
#define print(x)        cout<<x<<'\n';
#define yes             cout<<"YES\n";
#define no              cout<<"NO\n";
#define error           cout<<"-1\n"; return;
#define pu(x)           push_back(x)
#define po              pop_back()
#define B               begin()
#define S               size()
#define mp(x, y)        make_pair(x, y)
#define all(ar)         ar.begin(), ar.end()
#define rall(ar)        ar.rbegin(), ar.rend()
#define Case(x)         cout<<"Case "<<x<<": "
#define debug           cout<<"Ami ekhane!!"<<nl;
#define fixpoint(x)     cout<<fixed<<setprecision(x)
#define UNQ(x)          (x).erase(unique(all(x)), (x).end())
#define MAX(x)          *max_element(all(x))
#define MIN(x)          *min_element(all(x))
#define SUM(x)          accumulate(all(x), 0LL)
#define CNT(x, a)       count(all(x), a)
///..................Graph Moves...................................
//const int dx[] = {+1,-1,+0,+0,-1,+1,-1,+1}; ///King's move
//const int dy[] = {+0,+0,+1,-1,+1,+1,-1,-1};
//const int dx[] = {-2,-2,-1,-1,+1,+1,+2,+2}; ///knight's move
//const int dy[] = {-1,+1,-2,+2,-2,+2,-1,+1};
const int mod           = 1e9+7;
const int inf           = 1e9+9;
inline ll POW(ll n,ll k){ll ans=1;while(k){if(k&1)ans=ans*n;n=n*n;k>>=1;}return ans;}
/// في سبيل الله
const int N = 2e5+5;

void test_case(int T){
    int n; cin>>n;
    vi ar(n);
    ll mx = 0, mn = 0, save;
    for(in){
        cin>>ar[i];
        if(ar[i] > ar[mx]) mx = i;
        if(ar[i] < ar[mn]) mn = i;
    }
    vector<pii> br(33), cr(33);
    if(is_sorted(all(ar))){cout<<"0\n"; return;}
    
    if(ar[mx] > 0){
        br.resize(0);
        int boro = ar[mx];
        while(boro+ar[mn] < 0){
            br.pu(mp(mx+1, mx+1));
            boro *= 2;
        }
        for(int i = 1; i < n; i++){
            if(ar[i] < 0){
                br.pu(mp(i+1, mx+1));
            }
        }
        for(int i = 1; i < n; i++){
            br.pu(mp(i+1, i));
        }
    }
    if(ar[mn] < 0){
        cr.resize(0);
        int choto = ar[mn];
        while(choto + ar[mx] > 0){
            choto *= 2;
            cr.pu(mp(mn+1, mn+1));
        }
        for(int i = n-2; i >= 0; i--){
            if(ar[i] > 0){
                cr.pu(mp(i+1, mn+1));
            }
        }
        for(int i = n-1; i > 0; i--){
            cr.pu(mp(i, i+1));
        }
    }
    vector<pii>ans = br.S < cr.S ? br : cr;
    cout<<ans.size()<<'\n';
    for(auto [a, b]: ans) cout<<a<<' '<<b<<'\n';
}

signed l9_30; signed main(){
    ios::sync_with_stdio(false);
    cin.tie(0);

    l9_30 = true;    cin>>l9_30;
    for(int T = 1; T <= l9_30; T++){
        test_case(T);
    }
    return l9_30 ^ l9_30;
}