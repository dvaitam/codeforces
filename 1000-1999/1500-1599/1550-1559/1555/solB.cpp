#include <bits/stdc++.h>

#include<ext/pb_ds/assoc_container.hpp>

#include<ext/pb_ds/tree_policy.hpp>



using namespace std;

using namespace __gnu_pbds;

typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> pbds; // find_by_order, order_of_key



#define fastio() ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL)

#define endl "\n"

#define set_bits __builtin_popcountll

#define ff first

#define ss second

#define PI 3.141592653589793238462

#define MOD 1000000007

#define INF 1e18





typedef long long ll;

typedef unsigned long long ull;

typedef long double lld;





/*---------------------------Debuger Begin-------------------------*/

#ifndef ONLINE_JUDGE

#define debug(x) cerr << #x <<" "; _print(x); cerr << endl;

#else

#define debug(x)

#endif

void _print(ll t) {cerr << t;}

void _print(int t) {cerr << t;}

void _print(string t) {cerr << t;}

void _print(char t) {cerr << t;}

void _print(lld t) {cerr << t;}

void _print(double t) {cerr << t;}

void _print(ull t) {cerr << t;}



template <class T, class V> void _print(pair <T, V> p);

template <class T> void _print(vector <T> v);

template <class T> void _print(set <T> v);

template <class T, class V> void _print(map <T, V> v);

template <class T> void _print(multiset <T> v);

template <class T, class V> void _print(pair <T, V> p) {cerr << "{"; _print(p.ff); cerr << ","; _print(p.ss); cerr << "}";}

template <class T> void _print(vector <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T> void _print(set <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T> void _print(multiset <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T, class V> void _print(map <T, V> v) {cerr << "[ "; for (auto i : v) {_print(i); cerr << " ";} cerr << "]";}



/*----------------------------Debuger Ends-----------------------------------*/



ll gcd(ll a, ll b){if (a == 0)return b;return gcd(b % a, a);}

void google(int t) {cout << "Case #" << t << ": ";}

ll countBits(ll n){ll count = 0;while (n){count++;n >>= 1;}return count;}

ll aPowerbModm(ll a, ll b, ll m){if (b == 0)return 1;if (b % 2 == 0){ll t = aPowerbModm(a, (b / 2), m);return (1ll * t * t % m);}else{ll t = aPowerbModm(a, ((b - 1) / 2), m);t = (1ll * t * t % m);return (1ll * a * t % m);}}

bool isPowerof2(ll n){if (!(n & (n - 1))){return true;}return false;}

vector<ll> sieve(int n) {int*arr = new int[n + 1](); vector<ll> vect; for (int i = 2; i <= n; i++)if (arr[i] == 0) {vect.push_back(i); for (int j = 2 * i; j <= n; j += i)arr[j] = 1;} return vect;} 



/* ----------Decimal Precision---------- */

// cout<<fixed<<setprecision(n) -> to fix precision to n decimal places.

// cout<<setprecision(n) -> without fixing

 

 /*------------------PBDS-----------------*/

// order_of_key (val): returns the no. of values less than val

// find_by_order (k): returns the kth largest element.(0-based)

 

/* --------------Builtin Bit Functions-------------- */

// 1) __builtin_clz(x) -> returns the number of zeros at the beginning in the bit representaton of x.

// 2) __builtin_ctz(x) -> returns the number of zeros at the end in the bit representaton of x.

// 3) __builtin_popcount(x) -> returns the number of ones in the bit representaton of x.

// 4) __builtin_parity(x) -> returns the parity of the number of ones in the bit representaton of x.



int dir[4][2]={{1,0},{-1,0},{0,1},{0,-1}};

static bool cmp(pair<ll, ll> &a, pair<ll, ll> &b)

{

        if (a.ff == b.ff)

                return a.ss < b.ss;



        return a.ff < b.ff;

}



void solve(){

    int W,H;

    cin>>W>>H;

    int x1,y1,x2,y2;

    cin>>x1>>y1>>x2>>y2;

    int w,h;

    cin>>w>>h;

    int ans = 1e9;

	if (x2 - x1 + w <= W) {

		ans = min(ans, max(0, w - x1));

		ans = min(ans, max(0, x2 - (W - w)));

	}

	if (y2 - y1 + h <= H) {

		ans = min(ans, max(0, h - y1));

		ans = min(ans, max(0, y2 - (H - h)));

	}

	if (ans == 1e9)

		cout << -1 << endl;

	else

		cout << ans<< endl;

}





int main()

{

        // freopen("input.txt", "r", stdin);

        // freopen("output.txt", "w", stdout);

        fastio();

        ll t=1;

        cin>>t;

        while(t--){

                solve();

        }

      

        return 0; 

}