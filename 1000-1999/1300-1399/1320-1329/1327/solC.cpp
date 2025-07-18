#include <bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>

// #include <atcoder/all>

// using namespace atcoder;

using namespace __gnu_pbds;

using namespace std;

using uns = unsigned;

using ll = long long;

using ld = long double;

using vb = std::vector<bool>;

using vvb = std::vector<vb>;

using vc = std::vector<char>;

using vvc = std::vector<vc>;

using vi = std::vector<int>;

using vvi = std::vector<vi>;

using vvvi = std::vector<vvi>;

using vll = std::vector<ll>;

using vvll = std::vector<vll>;

using vvvll = std::vector<vvll>;

using vld = std::vector<ld>;

using vvld = std::vector<vld>;

using vvvld = std::vector<vvld>;

using pll = std::pair<ll, ll>;

using vpll = std::vector<pll>;

using pii = std::pair<int, int>;

using vpii = std::vector<pii>;

using vu = std::vector<uns>;

using vs = std::vector<std::string>;

using ordered_set = tree<pll, null_type, std::less<pll>, rb_tree_tag, tree_order_statistics_node_update>;



#define siz(x) (ll) x.size()

#define all(v) (v).begin(), (v).end()



#ifdef DEBUG

#include "dbg.hpp"

#else

#define err(...)

#define deb(...)

#endif



const ll mod = 998244353;

//using mint = modint998244353;

const ll maxN = 1e5+5;

void run_brute(){

    //cout<<"\n................\n";

}

template<class T> inline T POW(T a,ll n){

    T r=1;

    for (; n>0; n>>=1,a*=a){

         if (n&1)r*=a;

    } 

    return r;

}

inline ll POW(int a,ll n){ return POW((ll)a,n); }



template<class T> vector<T> powers(T m,ll n){

	vector<T> ret(n+1,1);

	for (ll i=1;i<=n;++i) ret[i]=ret[i-1]*m;

	return ret;

}

#include <bits/stdc++.h>

using namespace std;

using ll = long long;

void solve()

{   

    ll n,m,k;

    cin>>n>>m>>k;  

    vector<array<ll,2>> src(k),dest(k);

    for(int i=0;i<k;i++){

        cin>>src[i][0]>>src[i][1];

    }

    for(int i=0;i<k;i++){

        cin>>dest[i][0]>>dest[i][1];

    }



    // (RRRR)(D)(LLLL)

    // 2*n*m + (m-1)

    // 2*n*m+m-1

    





    /*

    If the chip is located by the wall of the board, and the action chosen by Petya moves it towards the wall, then the chip remains in its current position.

    */

    string right = string(m-1,'R');

    string left  = string(m-1,'L');

    string ans = "";

    

    for(int i=1;i<=n-1;i++){

        ans+='U';

    }

    

    for(int i=1;i<=m-1;i++){

        ans+='L';

    }



    for(int i=1;i<=n;i++){

        if(i&1){

            ans+=right;

            ans+='D';

        }else{

           ans+=left;

           ans+='D';

        }

    }



    cout<<ans.size()<<"\n";

    cout<<ans<<"\n";

}

int main()

{

#ifndef ONLINE_JUDGE

    freopen("input.txt", "r", stdin);

    freopen("output.txt", "w", stdout);

#endif

    std::ios_base::sync_with_stdio(false), std::cin.tie(nullptr);

    int t = 1;



    //cout << fixed << setprecision(20);

    //cin >> t;

    // precompute(105);

    for (int i = 1; i <= t; i++)

    {

        // cout << "Case #" << i << ": ";

        solve();

       // run_brute();

    }

    return 0;

}