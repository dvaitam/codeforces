#include <bits/stdc++.h>
using namespace std;
#define all(x)          (x).begin(),(x).end()
#define MOD             1000000007
#define case(tc)        cout << "Case " << tc << ": ";
#define close           << "\n";
typedef int64_t ll;
typedef uint64_t ull;
typedef long double lld;
typedef pair<int,int> pii;
typedef pair<ll,ll> pll;

void debug_(int val) {cerr << val;}
void debug_(ll val) {cerr << val;}
void debug_(ull val) {cerr << val;}
void debug_(lld val) {cerr << val;}
void debug_(double val) {cerr << val;}
void debug_(char val) {cerr << val;}
void debug_(string val) {cerr << val;}
template <class T, class V> void debug_(pair <T, V> p) {cerr << "{"; debug_(p.first); cerr << ","; debug_(p.second); cerr << "}";}
template <class T> void debug_(vector <T> v) {cerr << "[ "; for (T i : v) {debug_(i); cerr << " ";} cerr << "]";}
template <class T, class V> void debug_(map <T, V> v) {cerr << "[ "; for (auto i : v) {debug_(i); cerr << " ";} cerr << "]";}
template <class T> void debug_(set <T> v) {cerr << "[ "; for (T i : v) {debug_(i); cerr << " ";} cerr << "]";}
template <class T> void debug_(multiset <T> v) {cerr << "[ "; for (T i : v) {debug_(i); cerr << " ";} cerr << "]";}
#ifdef necromancer
#define debug(x) cerr << #x <<" "; debug_(x); cerr << endl;
#else
#define debug(x);
#endif

void test(int tc){
    int n,m; cin >> n >> m;
    vector<tuple<int,int,int>> a(m);
    for(int i = 0; i < m; ++i){
        int x,w; cin >> x >> w;
        a[i] = make_tuple(w,x,i+1);
    }

    sort(all(a));
    sort(a.begin(), a.begin()+2*n, [&](tuple<int,int,int> xx, tuple<int,int,int> yy){
        return get<1>(xx) < get<1>(yy);
    });

    int sum = 0;
    for(int i = 0; i < 2*n; ++i){
        sum += get<0>(a[i]);
    }

    cout << sum close 

    int left = 0;
    int right = 2*n-1;
    while(left < right){
        cout << get<2>(a[left]) << " " << get<2>(a[right]) close 
        left++;
        right--;
    }
    cout close 

}

int main() {
    ios_base::sync_with_stdio(false), cin.tie(0), cout.tie(0);
#ifdef necromancer
    freopen("./input.txt",  "r", stdin);
    // freopen("./output.txt", "w", stdout);
    // freopen("./error.txt",  "w", stderr);
#endif
    int testcase = 1;
    cin >> testcase;
    for (int tc = 0; tc < testcase; tc++) {
        test(tc);
    }
}