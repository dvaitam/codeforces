#include <bits/stdc++.h>

#include <ext/random>

using namespace std;



//#define labeltc

#define endl "\n"

#define MAXN 200001

#define ll long long

#define MOD 1000000007

//#define MOD 998244353 //853?

#define INFLL 1000000000001000000LL

#define INFI 1001000000

#define pii pair<int,int>

#define pll pair<ll,ll>

#define fi first

#define sc second

#define m_p make_pair

#define p_b push_back

#define l_b lower_bound

#define u_b upper_bound

#define vi vector<int>

#define vll vector<ll>

#define sp(x) x << " "

#define rand_num(x,y) uniform_int_distribution<ll>((ll)x,(ll)y)(rng)

#define lsb(x) (x&(-x))

#define dgt(x) (x-'0')

#define all(x) x.begin(),x.end()

#define pans(x) ((x)? "YES " : "NO ")

template<typename T> bool ckmin(T& a, const T& b) {return (a>b)? a = b, 1 : 0;}

template<typename T> bool ckmax(T& a, const T& b) {return (a<b)? a = b, 1 : 0;}



__gnu_cxx::sfmt19937 rng(std::chrono::steady_clock::now().time_since_epoch().count());



ll power(ll x, ll e, ll m = LONG_LONG_MAX){

    //e %= phi[m];

    ll r = 1;

    while(e>0){

        if(e%2) r = (r*x)%m;

        x = (x*x)%m;

        e >>= 1;

    }

    return r;

}



template <typename T, typename U>

ostream& operator<<(ostream& os, const pair<T,U>& v){

    os << sp(v.fi) << v.sc;

    return os;

}

template <typename T>

ostream& operator<<(ostream& os, const vector<T>& v){

    for (int i = 0; i < v.size(); ++i) {os << sp(v[i]);}

    return os;

}

template <typename T>

ostream& operator<<(ostream& os, const set<T>& v){

    for (auto it : v) {os << sp(it);}

    return os;

}

template <typename T, typename U>

ostream& operator<<(ostream& os, const map<T,U>& v){

    for (auto it : v) {os << it << "\n";}

    return os;

}



void setIO(string s) {

    freopen((s+".in").c_str(),"r",stdin);

    freopen((s+".out").c_str(),"w",stdout);

}



int uwu[501][501];



void precomp(){

    uwu[0][0] = 7;

    uwu[0][1] = 6;

    uwu[0][2] = 9;

    uwu[1][0] = 8;

    uwu[1][1] = 2;

    uwu[1][2] = 5;

    uwu[2][0] = 1;

    uwu[2][1] = 4;

    uwu[2][2] = 3;

    //769

    //825

    //143

    //R: 134259678

    //Q: 12345678

}



void solve(int tc){

    #ifdef labeltc

    cout << "Test case " << tc << ": ";

    #endif

    

    int n; cin >> n;

    if(n<3){

        cout << -1;

        return;

    }

    for (int i=0; i<3; i++) for (int j=0; j<3; j++) uwu[i][j] += n*n-9;

    int i = 3;

    int j = 0;

    int cur = n*n-9;

    int mode = 0;

    while(cur){

        uwu[i][j] = cur;

        if(mode == 0) j++;

        else if(mode == 1) j--;

        else if (mode == 2) i++;

        else i--;

        if(i==j) mode = 3-mode;

        else if(i<0){

            i=0;

            j++;

            mode = 2;

        }

        else if(j<0){

            i++;

            j=0;

            mode = 0;

        }

        cur--;

    }

    for (int i=0; i<n; i++){

        for (int j=0; j<n; j++) cout << sp(uwu[i][j]);

        cout << endl;

    }

    return;

}



signed main(){

    //setIO();

    ios_base::sync_with_stdio(0);

    cin.tie(NULL);

    precomp();

    int t=1;

    //cin >> t;

    for (int i=1; i<=t; i++) solve(i);

    return 0;

}