#include <bits/stdc++.h>

using namespace std;



/*-----------------------MACROS-----------------------------------*/



#define fastio()         ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL)

#define mod              1000000007

#define mod1             998244353

#define inf              1e18

#define endl             "\n"

#define sz(x)            (int)x.size()

#define pb               push_back

#define ppb              pop_back

#define mp               make_pair

#define ff               first

#define ss               second

#define pi               3.141592653589793238462

#define all(x)           (x).begin(), (x).end()

#define setbits          __builtin_popcount

#define setbitsll        __builtin_popcountll



/*-----------------------Debugging--------------------------------*/



#ifndef ONLINE_JUDGE

#define debug(x) cerr << #x<<" "; _print(x); cerr << endl;

#else

#define debug(x);

#endif



typedef long long ll;

typedef unsigned long long ull;

typedef long double lld;



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

template <class T> void _print(unordered_set <T> v);

template <class T, class V> void _print(map <T, V> v);

template <class T, class V> void _print(unordered_map <T, V> v);

template <class T> void _print(multiset <T> v);

template <class T> void _print(stack <T> v);

template <class T, class V> void _print(pair <T, V> p) {cerr << "{"; _print(p.ff); cerr << ","; _print(p.ss); cerr << "}";}

template <class T> void _print(vector <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T> void _print(set <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T> void _print(unordered_set <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T> void _print(multiset <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T> void _print(stack <T> v) {cerr << "[ "; while (!v.empty()) {_print(v.top()); cerr << " ";} cerr << "]";}

template <class T, class V> void _print(map <T, V> v) {cerr << "[ "; for (auto i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T, class V> void _print(unordered_map <T, V> v) {cerr << "[ "; for (auto i : v) {_print(i); cerr << " ";} cerr << "]";}





/*--------------------------Algorithms----------------------------*/



ll gcd(ll a, ll b) {if (b > a) {return gcd(b, a);} if (b == 0) {return a;} return gcd(b, a % b);}

void google(int t) {cout << "Case #" << t << ": ";}





/*--------------------------Code begins---------------------------*/



void solve()

{

    int n;

    cin >> n;

    vector<int> p(n);

    for (int& i : p) cin >> i;

    string s;

    cin >> s;

    int z = count(all(s), '0');

    vector<int> x, y;

    for (int i = 0; i < n; i += 1) {

        if (s[i] == '0' and p[i] > z)

            x.push_back(i);

        if (s[i] == '1' and p[i] <= z)

            y.push_back(i);

    }

    debug(x) debug(y)

    assert(x.size() == y.size());

    for (int i = 0; i < (int)x.size(); i += 1) swap(p[x[i]], p[y[i]]);

    for (int i : p) cout << i << " ";

    cout << endl;

}



int main()

{

#ifndef ONLINE_JUDGE

    freopen("error.txt", "w", stderr);

#endif

    fastio();

    int t = 1;

    cin >> t;

    while (t--)

    {

        solve();

    }



    return 0;

}