//Code by Rohan Sharma(Naruto_Uchiha-CF)(rs2002-CC)

#include<bits/stdc++.h>

#define jldi_kr ;           ios_base::sync_with_stdio(false);cin.tie(NULL)

#define code_run_kr_bsdk    int main()

#define sirf_ek_baar_kr     solve();

#define hehe ;              using namespace std;

#define test_cases_hai      wt{solve();}

#define f(i,a,n,c)          for(int i=a;i<n;i+=c)

#define fr(i,a,n,c)         for(int i=a;i>n;i-=c)

#define ll                  long long int

#define l(s)                s.size()

#define asort(a)            sort(a,a+n)

#define all(x) 	            (x).begin(), (x).end()

#define dsort(a)            sort(a,a+n,greater<ll>())

#define vasort(v)           sort(v.begin(), v.end());

#define vdsort(v)           sort(v.begin(), v.end(),greater<ll>());

#define vprint(v);		    f(i,0,l(v),1){cout<<v[i]<<' ';}cout<<endl;

#define pb                  push_back

#define ff 				    first

#define ss 				    second

#define vll ;               vector<ll>

#define mpll ;              map<ll,ll>

#define p(x);               cout<<x<<endl;

#define inp(n)              ll n;cin>>n;

#define vin(x);             for(int i=0;i<n;i++){ll u;cin>>u;x.push_back(u);}

#define wt                  ll t;cin>>t;while(t--)

const ll mod = 1e9+7;



hehe;



//--------------------------------------------------------------------------------------------------------------------------------------------

ll gcd(ll a, ll b) {if (b > a) {return gcd(b, a);} if (b == 0) {return a;} return gcd(b, a % b);}

void extendgcd(ll a, ll b, ll*v) {if (b == 0) {v[0] = 1; v[1] = 0; v[2] = a; return ;} extendgcd(b, a % b, v); ll x = v[1]; v[1] = v[0] - v[1] * (a / b); v[0] = x; return;} //pass an arry of size1 3

ll expo(ll a, ll b, ll mod) {ll res = 1; while (b > 0) {if (b & 1)res = (res * a) % mod; a = (a * a) % mod; b = b >> 1;} return res;}

ll mminv(ll a, ll b) {ll arr[3]; extendgcd(a, b, arr); return arr[0];} //for non prime b

ll mminvprime(ll a, ll b) {return expo(a, b - 2, b);}

ll combination(ll n, ll r, ll m, ll *fact, ll *ifact) {ll val1 = fact[n]; ll val2 = ifact[n - r]; ll val3 = ifact[r]; return (((val1 * val2) % m) * val3) % m;}

vector<ll> sieve(int n) {int*arr = new int[n + 1](); vector<ll> vect; for (int i = 2; i <= n; i++)if (arr[i] == 0) {vect.push_back(i); for (int j = 2 * i; j <= n; j += i)arr[j] = 1;} return vect;}

ll mod_add(ll a, ll b, ll m) {a = a % m; b = b % m; return (((a + b) % m) + m) % m;}

ll mod_mul(ll a, ll b, ll m) {a = a % m; b = b % m; return (((a * b) % m) + m) % m;}

ll mod_sub(ll a, ll b, ll m) {a = a % m; b = b % m; return (((a - b) % m) + m) % m;}

ll mod_div(ll a, ll b, ll m) {a = a % m; b = b % m; return (mod_mul(a, mminvprime(b, m), m) + m) % m;}  //only for prime m

ll phin(ll n) {ll number = n; if (n % 2 == 0) {number /= 2; while (n % 2 == 0) n /= 2;} for (ll i = 3; i <= sqrt(n); i += 2) {if (n % i == 0) {while (n % i == 0)n /= i; number = (number / i * (i - 1));}} if (n > 1)number = (number / n * (n - 1)) ; return number;} 

//--------------------------------------------------------------------------------------------------------------------------------------------

#ifndef ONLINE_JUDGE

#define dbg(x) cerr<<#x<<" ";_print_(x);cerr<<endl;

#else

#define dbg(x)

#endif

//--------------------------------------------------------------------------------------------------------------------------------------------

void _print_(ll t) { cerr << t; }

void _print_(int t) { cerr << t; }

void _print_(string t) { cerr << t; }

void _print_(char t) { cerr << t; }

void _print_(long double t) { cerr << t; }

void _print_(double t) { cerr << t; }

//--------------------------------------------------------------------------------------------------------------------------------------------

template <class T, class V> void _print_(pair <T, V> p);

template <class T, class V> void _print_(pair <T, V> p) { cerr << "{"; _print_(p.first); cerr << ","; _print_(p.second); cerr << "}"; }

template <class T> void _print_(set <T> v);

template <class T> void _print_(set <T> v) { cerr << "[ "; for (T i : v) { _print_(i); cerr << " "; } cerr << "]"; }

template <class T, class V> void _print_(map <T, V> v);

template <class T, class V> void _print_(map <T, V> v) { cerr << "[ "; for (auto i : v) { _print_(i); cerr << " "; } cerr << "]"; }

template <class T> void _print_(multiset <T> v);

template <class T> void _print_(multiset <T> v) { cerr << "[ "; for (T i : v) { _print_(i); cerr << " "; } cerr << "]"; }

template <class T> void _print_(vector <T> v);

template <class T> void _print_(vector <T> v) { cerr << "[ "; for (T i : v) { _print_(i); cerr << " "; } cerr << "]"; }

//--------------------------------------------------------------------------------------------------------------------------------------------



void solve()

{

    ll n, t;

    cin >> n >> t;

    ll c = 0;

    f(i, 0, n, 1)

    {

        ll x;

        cin >> x;

        ll r = 0;

        if (t % 2 == 0 && x == t / 2)

        {

            r = (c++) % 2;

        }

        else if (2 * x < t)

        {

            r = 0;

        }

        else

        {

            r = 1;

        }

        cout << r << " ";

    }

    cout << endl;

}



code_run_kr_bsdk

{

    jldi_kr;



    #ifndef ONLINE_JUDGE

    freopen("input.txt", "r", stdin);

    freopen("output.txt", "w", stdout);

    freopen("err.txt", "w", stderr);

    #endif



    // sirf_ek_baar_kr;

    test_cases_hai;



}