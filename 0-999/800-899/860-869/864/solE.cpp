#pragma GCC optimize ("O3")
#pragma GCC target ("avx")
#include "bits/stdc++.h" // define macro "/D__MAI"

using namespace std;
typedef long long int ll;

#define xprintf(fmt,...) fprintf(stderr,fmt,__VA_ARGS__)
#define debugv(v) {printf("L%d %s > ",__LINE__,#v);for(auto e:v){cout<<e<<" ";}cout<<endl;}
#define debuga(m,w) {printf("L%d %s > ",__LINE__,#m);for(int x=0;x<(w);x++){cout<<(m)[x]<<" ";}cout<<endl;}
#define debugaa(m,h,w) {printf("L%d %s >\n",__LINE__,#m);for(int y=0;y<(h);y++){for(int x=0;x<(w);x++){cout<<(m)[y][x]<<" ";}cout<<endl;}}
#define ALL(v) (v).begin(),(v).end()
#define repeat(cnt,l) for(auto cnt=0ll;cnt<(l);++cnt)
#define iterate(cnt,b,e) for(auto cnt=(b);cnt!=(e);++cnt)
#define MD 1000000007ll
#define PI 3.1415926535897932384626433832795
#define EPS 1e-12
template<typename T1, typename T2> ostream& operator <<(ostream &o, const pair<T1, T2> p) { o << "(" << p.first << ":" << p.second << ")"; return o; }
template<typename iterator> inline size_t argmin(iterator begin, iterator end) { return distance(begin, min_element(begin, end)); }
template<typename iterator> inline size_t argmax(iterator begin, iterator end) { return distance(begin, max_element(begin, end)); }
template<typename T> T& maxset(T& to, const T& val) { return to = max(to, val); }
template<typename T> T& minset(T& to, const T& val) { return to = min(to, val); }

mt19937_64 randdev(8901016);
inline ll rand_range(ll l, ll h) {
    return uniform_int_distribution<ll>(l, h)(randdev);
}

#define getchar_unlocked getchar
#define putchar_unlocked putchar
namespace {
#define isvisiblechar(c) (0x21<=(c)&&(c)<=0x7E)
    class MaiScanner {
    public:
        template<typename T> void input_integer(T& var) {
            var = 0;
            T sign = 1;
            int cc = getchar_unlocked();
            for (; cc<'0' || '9'<cc; cc = getchar_unlocked())
                if (cc == '-') sign = -1;
            for (; '0' <= cc&&cc <= '9'; cc = getchar_unlocked())
                var = (var << 3) + (var << 1) + cc - '0';
            var = var*sign;
        }
        inline int c() { return getchar_unlocked(); }
        inline MaiScanner& operator>>(int& var) {
            input_integer<int>(var);
            return *this;
        }
        inline MaiScanner& operator>>(long long& var) {
            input_integer<long long>(var);
            return *this;
        }
        inline MaiScanner& operator>>(string& var) {
            int cc = getchar_unlocked();
            for (; !isvisiblechar(cc); cc = getchar_unlocked());
            for (; isvisiblechar(cc); cc = getchar_unlocked())
                var.push_back(cc);
            return *this;
        }
        template<typename IT> void in(IT begin, IT end) {
            for (auto it = begin; it != end; ++it) *this >> *it;
        }
    };
}
MaiScanner scanner;


// ソートしてi番目の値は idx[i]番目の値
// i番目の値は  ソートしてidxr[i]番目の値
template <typename ITER>
void sort_idx_proto(const ITER begin, const ITER end, vector<int> &idx, vector<int> &idxr) {
    size_t n = end - begin;
    idx.resize(n);
    idxr.resize(n);
    for (int i = 0; i < n; ++i) idx[i] = i;
    sort(idx.begin(), idx.end(), [&begin](int l, int r) {return begin[l] < begin[r]; });
    for (int i = 0; i < n; ++i) idxr[idx[i]] = i;
}



template<typename T,typename U>
pair<T,U>& sp_maxset(pair<T, U>& to, const pair<T, U>& val) {
    if (to.first < val.first) to = val;
    return to;
}




ll m, n, kei;

int tt[111]; // time
int dd[111]; // lost
int pp[111]; // value


pair<ll,bitset<111>> dp[2100]; // TIME

int main() {

    scanner >> n;

    repeat(i, n) {
        int t, d, p;
        scanner >> t >> d >> p;
        tt[i] = t;
        dd[i] = d;
        pp[i] = p;
    }

    vector<int> sidx, sidxr;

    sort_idx_proto(dd, dd + n, sidx, sidxr);


    pair<ll, bitset<111>> best;
    best.first = 0;
    best.second.reset();

    repeat(i, n) {
        int x = sidx[i];
        for (int tim = dd[x]-tt[x]-1; 0 <= tim; --tim){
            pair<ll, bitset<111>> s = dp[tim];
            s.first += pp[x];
            s.second.set(x);

            sp_maxset(best,
                sp_maxset(dp[tim + tt[x]], s));
        }
    }

    cout << best.first << endl;

    vector<int> lis;
    repeat(i, n) { if (best.second[i]) lis.push_back(i); }

    cout << lis.size() << endl;

    sort(ALL(lis), [&](int l, int r) { return dd[l] < dd[r]; });
    for (int e : lis) {
        printf("%d ", e + 1);
    }
    cout << endl;



    return 0;
}