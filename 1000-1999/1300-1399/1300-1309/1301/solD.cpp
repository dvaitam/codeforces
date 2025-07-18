#include<bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

//#include<boost/algorithm/string.hpp>



//pragmas

#pragma GCC optimize("O3")



//types

#define fastio() ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL)

#define ll long long int

#define ull unsigned long long int

#define vec vector<long long int>

#define pall pair<long long int, long long int>

#define vecpair vector<pair<long long int,long long int>>

#define vecvec(a, i, j) vector<vector<long long int>> a (i, vec (j, 0))

#define vecvecvec(a, i, j, k) vector<vector<vector<long long int>>> dp (i + 1, vector<vector<long long int>>(j + 1, vector<long long int>(k + 1, 0)))



using namespace std;

using namespace __gnu_pbds;



//random stuff

#define all(a) a.begin(),a.end()

#define read(a) for (auto &x : a) cin >> x

#define endl "\n"

#define pb push_back

#define size(x) (long long int)x.size()

#define print(a) for(auto x : a) cout << x << " "; cout << endl

#define sp " " 

ll INF = 9223372036854775807;

typedef tree<long long int, null_type, less<long long int>, rb_tree_tag, tree_order_statistics_node_update> indexed_set;

struct custom_hash {

    static uint64_t splitmix64(uint64_t x) {

        x += 0x9e3779b97f4a7c15;

        x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9;

        x = (x ^ (x >> 27)) * 0x94d049bb133111eb;

        return x ^ (x >> 31);

    }



    size_t operator()(uint64_t x) const {

        static const uint64_t FIXED_RANDOM = chrono::steady_clock::now().time_since_epoch().count();

        return splitmix64(x + FIXED_RANDOM);

    }

};

#define safe_map unordered_map<long long, int, custom_hash>





//debug

void __print(int x) {cerr << x;}

void __print(long x) {cerr << x;}

void __print(long long x) {cerr << x;}

void __print(unsigned x) {cerr << x;}

void __print(unsigned long x) {cerr << x;}

void __print(unsigned long long x) {cerr << x;}

void __print(float x) {cerr << x;}

void __print(double x) {cerr << x;}

void __print(long double x) {cerr << x;}

void __print(char x) {cerr << '\'' << x << '\'';}

void __print(const char *x) {cerr << '\"' << x << '\"';}

void __print(const string &x) {cerr << '\"' << x << '\"';}

void __print(bool x) {cerr << (x ? "true" : "false");}



template<typename T, typename V>

void __print(const pair<T, V> &x) {cerr << '{'; __print(x.first); cerr << ','; __print(x.second); cerr << '}';}

template<typename T>

void __print(const T &x) {int f = 0; cerr << '{'; for (auto &i: x) cerr << (f++ ? "," : ""), __print(i); cerr << "}";}

void _print() {cerr << "]\n";}

template <typename T, typename... V>

void _print(T t, V... v) {__print(t); if (sizeof...(v)) cerr << ", "; _print(v...);}

#ifndef ONLINE_JUDGE

#define debug(x...) cerr << "[" << #x << "] = ["; _print(x)

#define reach cerr<<"reached"<<endl

#else

#define debug(x...)

#define reach 

#endif





/*---------------------------------------------------------------------------------------------------------------------------*/

ll gcd(ll a, ll b) {if (b > a) {return gcd(b, a);} if (b == 0) {return a;} return gcd(b, a % b);}

ll expo(ll a, ll b, ll mod) {ll res = 1; while (b > 0) {if (b & 1)res = (res * a) % mod; a = (a * a) % mod; b = b >> 1;} return res;}

void extendgcd(ll a, ll b, ll*v) {if (b == 0) {v[0] = 1; v[1] = 0; v[2] = a; return ;} extendgcd(b, a % b, v); ll x = v[1]; v[1] = v[0] - v[1] * (a / b); v[0] = x; return;} //pass an arry of size1 3

ll mminv(ll a, ll b) {ll arr[3]; extendgcd(a, b, arr); return arr[0];} //for non prime b

ll mminvprime(ll a, ll b) {return expo(a, b - 2, b);}

bool revsort(ll a, ll b) {return a > b;}

void swap(int &x, int &y) {int temp = x; x = y; y = temp;}

ll combination(ll n, ll r, ll m, ll *fact, ll *ifact) {ll val1 = fact[n]; ll val2 = ifact[n - r]; ll val3 = ifact[r]; return (((val1 * val2) % m) * val3) % m;}

void google(int t) {cout << "Case #" << t << ": ";}

vector<ll> sieve(int n) {int*arr = new int[n + 1](); vector<ll> vect; for (ll i = 2; i <= n; i++)if (arr[i] == 0) {vect.push_back(i); for (ll j = 2 * i; j <= n; j += i)arr[j] = 1;} return vect;}

ll mod_add(ll a, ll b, ll m = 1000000007) {a = a % m; b = b % m; return (((a + b) % m) + m) % m;}

ll mod_mul(ll a, ll b, ll m = 1000000007) {a = a % m; b = b % m; return (((a * b) % m) + m) % m;}

ll mod_sub(ll a, ll b, ll m = 1000000007) {a = a % m; b = b % m; return (((a - b) % m) + m) % m;}

ll mod_div(ll a, ll b, ll m = 1000000007) {a = a % m; b = b % m; return (mod_mul(a, mminvprime(b, m), m) + m) % m;}  //only for prime m

ll phin(ll n) {ll number = n; if (n % 2 == 0) {number /= 2; while (n % 2 == 0) n /= 2;} for (ll i = 3; i <= sqrt(n); i += 2) {if (n % i == 0) {while (n % i == 0)n /= i; number = (number / i * (i - 1));}} if (n > 1)number = (number / n * (n - 1)) ; return number;} //O(sqrt(N))

void precision(int a) {cout << setprecision(a) << fixed;}

ll ceil_div(ll x, ll y){return (x + y - 1) / y;}

unsigned long long power(unsigned long long x,ll y, ll p){unsigned long long res = 1;x = x % p; while (y > 0){if (y & 1)res = (res * x) % p;y = y >> 1;x = (x * x) % p;}return res;}

unsigned long long modInverse(unsigned long long n,int p){return power(n, p - 2, p);}

ll nCr(ll n,ll r, ll p){if (n < r)return 0;if (r == 0)return 1;unsigned long long fac[n + 1];fac[0] = 1;for (int i = 1; i <= n; i++)fac[i] = (fac[i - 1] * i) % p;return (fac[n] * modInverse(fac[r], p) % p* modInverse(fac[n - r], p) % p)% p;}

ll accumulate(const vec &nums){ll sum = 0; for(auto x : nums) sum += x; return sum;}

ll tmax(ll a, ll b, ll c = 0, ll d = -INF, ll e = -INF, ll f = -INF){return max(a, max(b, max(c, max(d, max(e, f)))));}

int log2_floor(unsigned long long i) {return i ? __builtin_clzll(1) - __builtin_clzll(i) : -1;}

string bin(ll n){return bitset<32>(n).to_string();}

/*--------------------------------------------------------------------------------------------------------------------------*/



ll mod = 1000000007;

vector<pair<ll, string>> mv;



//code starts

int main()

{

    fastio();



    ll n, m, k;

    cin >> m >> n >> k;

    if(k <= 4 * n * m - 2 * (n + m))

    {

        cout << "YES" << endl;

        ll run = 0;

        if(n - 1)

        {

            mv.pb({min(k, n - 1), "R"});

            run += min(k, n - 1);

        }

        for(ll i = n; i > 1 and run < k; i --)

        {

            if((k - run)/3 == 0)

            {   

                if(m - 1)

                {

                    if(k - run == 1)    mv.pb({1, "D"});

                    if(k - run == 2)    mv.pb({1, "DL"});

                    run += k - run;

                }

            }

            else if((k - run)/3 < m - 1)

            {

                ll down = (k - run)/3;  //go down by this much

                mv.pb({(k - run)/3, "DLR"});

                run += ((k - run)/3) * 3;

                if(down and run < k)

                {

                    mv.pb({1, "D"});

                    run ++;

                    if(k - run)

                    {

                        mv.pb({1, "U"});

                        run ++;

                    }

                }

            }   

            else

            {

                if(m - 1)

                {

                    mv.pb({m - 1, "DLR"});

                    run += (m - 1) * 3;            

                    if(run < k) //go back up

                    {

                        mv.pb({min(k - run, m - 1), "U"});

                        run += min(k - run, m - 1);

                    } 

                }

            }

            if(run < k) //go left

            {

                mv.pb({1, "L"});

                run ++;

            }

        }

        if(run < k)

        {

            if(m - 1)

            {

                mv.pb({min(k - run, m - 1), "D"});

                run += min(k - run, m - 1);

                if(run < k)

                {

                    mv.pb({min(k - run, m - 1), "U"});

                    run += min(k - run, m - 1);

                }

            }

        }

        cout << size(mv) << endl;

        for(auto &x : mv)   cout << x.first << sp << x.second << endl;

    }

    else

        cout << "NO" << endl;

}













// There is an idea of a Patrick Bateman. Some kind of abstraction.

// But there is no real me. Only an entity. Something illusory. 

// And though I can hide my cold gaze, and you can shake my hand and

// feel flesh gripping yours, and maybe you can even sense our lifestyles

// are probably comparable, I simply am not there.