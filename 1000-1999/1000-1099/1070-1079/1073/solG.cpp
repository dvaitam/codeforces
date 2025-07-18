/***********Template Starts Here***********/
#include <bits/stdc++.h>

#define pb push_back
#define nl puts ("")
#define sp printf ( " " )
#define phl printf ( "hello\n" )
#define ff first
#define ss second
#define POPCOUNT __builtin_popcountll
#define RIGHTMOST __builtin_ctzll
#define LEFTMOST(x) (63-__builtin_clzll((x)))
#define MP make_pair
#define FOR(i,x,y) for(vlong i = (x) ; i <= (y) ; ++i)
#define ROF(i,x,y) for(vlong i = (y) ; i >= (x) ; --i)
#define CLR(x,y) memset(x,y,sizeof(x))
#define UNIQUE(V) (V).erase(unique((V).begin(),(V).end()),(V).end())
#define MIN(a,b) ((a)<(b)?(a):(b))
#define MAX(a,b) ((a)>(b)?(a):(b))
#define NUMDIGIT(x,y) (((vlong)(log10((x))/log10((y))))+1)
#define SQ(x) ((x)*(x))
#define ABS(x) ((x)<0?-(x):(x))
#define FABS(x) ((x)+eps<0?-(x):(x))
#define ALL(x) (x).begin(),(x).end()
#define LCM(x,y) (((x)/gcd((x),(y)))*(y))
#define SZ(x) ((vlong)(x).size())
#define NORM(x) if(x>=mod)x-=mod;
#define MOD(x,y) (((x)*(y))%mod)
#define ODD(x) (((x)&1)==0?(0):(1))

using namespace std;

typedef long long vlong;
typedef unsigned long long uvlong;
typedef pair < vlong, vlong > pll;
typedef vector<pll> vll;
typedef vector<vlong> vl;
using ll = long long;

const vlong inf = 2147383647;
const double pi = 2 * acos ( 0.0 );
const double eps = 1e-9;
mt19937 rng(std::chrono::duration_cast<std::chrono::nanoseconds>(chrono::high_resolution_clock::now().time_since_epoch()).count());

template<typename S, typename T>
void xmin(S&a, T const&b){if(b<a) a=b;}
template<typename S, typename T>
void xmax(S&a, T const&b){if(b>a) a=b;}

#ifdef DEBUG
     clock_t tStart = clock();
     #define debug(args...) {dbg,args; cerr<<endl;}
     #define timeStamp debug ("Execution Time: ", (double)(clock() - tStart)/CLOCKS_PER_SEC)
     #define bug printf("%d\n",__LINE__);

#else
    #define debug(args...)  // Just strip off all debug tokens
    #define timeStamp
#endif

struct debugger{
    template<typename T> debugger& operator , (const T& v){
        cerr<<v<<" ";
        return *this;
    }
}dbg;

//int knightDir[8][2] = { {-2,1},{-1,2},{1,2},{2,1},{2,-1},{-1,-2},{1,-2},{-2,-1} };
//int dir4[4][2] = {{-1,0},{0,1},{1,0},{0,-1}};

inline vlong gcd ( vlong a, vlong b ) {
    a = ABS ( a ); b = ABS ( b );
    while ( b ) { a = a % b; swap ( a, b ); } return a;
}

vlong ext_gcd ( vlong A, vlong B, vlong *X, vlong *Y ){
    vlong x2, y2, x1, y1, x, y, r2, r1, q, r;
    x2 = 1; y2 = 0;
    x1 = 0; y1 = 1;
    for (r2 = A, r1 = B; r1 != 0; r2 = r1, r1 = r, x2 = x1, y2 = y1, x1 = x, y1 = y ) {
        q = r2 / r1;
        r = r2 % r1;
        x = x2 - (q * x1);
        y = y2 - (q * y1);
    }
    *X = x2; *Y = y2;
    return r2;
}

inline vlong modInv ( vlong a, vlong m ) {
    vlong x, y;
    ext_gcd( a, m, &x, &y );
    x %= m;
    if ( x < 0 ) x += m; //modInv is never negative
    return x;
}

inline vlong power ( vlong a, vlong p ) {
    vlong res = 1, x = a;
    while ( p ) {
        if ( p & 1 ) res = ( res * x );
        x = ( x * x ); p >>= 1;
    }
    return res;
}

inline vlong bigmod ( vlong a, vlong p, vlong m ) {
    vlong res = 1 % m, x = a % m;
    while ( p ) {
        if ( p & 1 ) res = ( res * x ) % m;
        x = ( x * x ) % m; p >>= 1;
    }
    return res;
}

/***********Template Ends Here***********/

signed gen(int T){
    // mt19937 rng(43151);
    auto get_rand = [&](int64_t l, int64_t r){
        return uniform_int_distribution<int64_t>(l, r)(rng);
    };
    auto get_double = [&](double l, double r){
        return uniform_real_distribution<double>(l, r)(rng);
    };
    ofstream o("gen.txt");
    o << T << "\n";
    for(int cas=0;cas<T;++cas){
        /// GEN HERE

        o << "\n";
    }
    o << endl;
    o.close();
    return 0;
}

class Suffix_Array{
    unsigned char mask[8] = { 0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01 };
    #define tget(i) ( (t[(i)/8]&mask[(i)%8]) ? 1 : 0 )
    #define tset(i, b) t[(i)/8]=(b) ? (mask[(i)%8]|t[(i)/8]) : ((~mask[(i)%8])&t[(i)/8])
    #define chr(i) (cs==sizeof(int)?((int*)s)[i]:((unsigned char *)s)[i])
    #define isLMS(i) (i>0 && tget(i) && !tget(i-1))

    // find the start or end of each bucket
    void getBuckets(unsigned char *s, int *bkt, int n, int K, int cs, bool end) {
        int i, sum = 0;
        for (i = 0; i <= K; i++)
            bkt[i] = 0;
        for (i = 0; i < n; i++)
            bkt[chr(i)]++;
        for (i = 0; i <= K; i++) {
            sum += bkt[i];
            bkt[i] = end ? sum : sum - bkt[i];
        }
    }
    void induceSAl(unsigned char *t, int *SA, unsigned char *s, int *bkt, int n, int K, int cs, bool end) {
        int i, j;
        getBuckets(s, bkt, n, K, cs, end);
        for (i = 0; i < n; i++) {
            j = SA[i] - 1;
            if (j >= 0 && !tget(j))
                SA[bkt[chr(j)]++] = j;
        }
    }
    void induceSAs(unsigned char *t, int *SA, unsigned char *s, int *bkt, int n, int K, int cs, bool end) {
        int i, j;
        getBuckets(s, bkt, n, K, cs, end);
        for (i = n - 1; i >= 0; i--) {
            j = SA[i] - 1;
            if (j >= 0 && tget(j))
                SA[--bkt[chr(j)]] = j;
        }
    }
    void SA_IS(unsigned char *s, int *SA, int n, int K, int cs) {
        int i, j;
        unsigned char *t = (unsigned char *) malloc(n / 8 + 1);
        tset(n-2, 0);
        tset(n-1, 1);
        for (i = n - 3; i >= 0; i--)
            tset(i, (chr(i)<chr(i+1) || (chr(i)==chr(i+1) && tget(i+1)==1))?1:0);
        int *bkt = (int *) malloc(sizeof(int) * (K + 1));
        getBuckets(s, bkt, n, K, cs, true);
        for (i = 0; i < n; i++)
            SA[i] = -1;
        for (i = 1; i < n; i++)
            if (isLMS(i))
                SA[--bkt[chr(i)]] = i;
        induceSAl(t, SA, s, bkt, n, K, cs, false);
        induceSAs(t, SA, s, bkt, n, K, cs, true);
        free(bkt);
        int n1 = 0;
        for (i = 0; i < n; i++)
            if (isLMS(SA[i]))
                SA[n1++] = SA[i];
        for (i = n1; i < n; i++)
            SA[i] = -1;
        int name = 0, prev = -1;
        for (i = 0; i < n1; i++) {
            int pos = SA[i];
            bool diff = false;
            for (int d = 0; d < n; d++)
                if (prev == -1 || chr(pos+d) != chr(prev+d) || tget(pos+d) != tget(prev+d)) {
                    diff = true;
                    break;
                } else if (d > 0 && (isLMS(pos+d) || isLMS(prev+d)))
                    break;
            if (diff) {
                name++;
                prev = pos;
            }
            pos = (pos % 2 == 0) ? pos / 2 : (pos - 1) / 2;
            SA[n1 + pos] = name - 1;
        }
        for (i = n - 1, j = n - 1; i >= n1; i--)
            if (SA[i] >= 0)
                SA[j--] = SA[i];
        int *SA1 = SA, *s1 = SA + n - n1;
        if (name < n1)
            SA_IS((unsigned char*) s1, SA1, n1, name - 1, sizeof(int));
        else
            for (i = 0; i < n1; i++)
                SA1[s1[i]] = i;
        bkt = (int *) malloc(sizeof(int) * (K + 1));
        getBuckets(s, bkt, n, K, cs, true);
        for (i = 1, j = 0; i < n; i++)
            if (isLMS(i))
                s1[j++] = i;
        for (i = 0; i < n1; i++)
            SA1[i] = s1[SA1[i]];
        for (i = n1; i < n; i++)
            SA[i] = -1;
        for (i = n1 - 1; i >= 0; i--) {
            j = SA[i];
            SA[i] = -1;
            SA[--bkt[chr(j)]] = j;
        }
        induceSAl(t, SA, s, bkt, n, K, cs, false);
        induceSAs(t, SA, s, bkt, n, K, cs, true);
        free(bkt);
        free(t);
    }
    public:
    int* sa,* inv;
    vector<int> lcp, seg;
    int N;
    private:

    void make_lcp(const char*s){
        lcp.resize(N);
        int k=0;
        for(int i=0;i<N;++i){
            if(inv[i]!=0){
                for(int j = sa[inv[i]-1];s[i+k]==s[j+k];++k);
                lcp[inv[i]-1]=k;
                if(k)--k;
            }
        }
        lcp[N-1]=0;
        seg.resize(2*N);
        copy(lcp.begin(), lcp.end(), seg.begin()+N);
        for(int i=N-1;i>=0;--i)
            seg[i] = min(seg[2*i], seg[2*i+1]);
    }
    public:
    Suffix_Array(){}
    ~Suffix_Array(){delete[] sa; delete[] inv;}
    void build(string &s, int max_sigma=256){
        N=s.size();
        int *v = new int[N+3];
        SA_IS((unsigned char*)s.c_str(), v, N+1, max_sigma, 1);
        sa = new int[N]; inv = new int[N];
        for(int i=0;i<N;++i){
            sa[i] = v[i+1];
            inv[sa[i]] = i;
        }
        make_lcp(s.c_str());
        delete[] v;
    }
    int get_lcp(int i, int j){
        if(i==j) return N - sa[i];
        int ans = 1e9;
        for(i+=N, j+=N;i<j;i>>=1, j>>=1){
            if(i&1) ans = min(ans, seg[i++]);
            if(j&1) ans = min(ans, seg[--j]);
        }
        return ans;
    }
};

Suffix_Array SA;

vlong solve(vector<int> &A, vector<int> &B) {
    auto cmp = [&](int i, int j) {
        return SA.inv[i] < SA.inv[j];
    };

    sort(ALL(A), cmp);
    sort(ALL(B), cmp);
    map<int, int> M;
    int pos = 0;
    vlong sum = 0, res = 0;
    FOR(i, 0, SZ(A) - 1) {
        int now = SA.inv[A[i]];
        if (i) {
            int x = SA.get_lcp(SA.inv[A[i - 1]], now);
            while (not M.empty()) {
                auto it = --M.end();
                if (it->first > x) {
                    int y = it->second;
                    sum -= (vlong)(it->first - x) * y;
                    M.erase(it);
                    M[x] += y;
                } else {
                    break;
                }
            }
        }
        while (pos < SZ(B) and SA.inv[B[pos]] <= now) {
            int x = SA.get_lcp(SA.inv[B[pos++]], now);
            sum += x;
            M[x]++;
        }
        res += sum;
    }
    M.clear();
    pos = SZ(B) - 1;
    sum = 0;
    ROF(i, 0, SZ(A) - 1) {
        int now = SA.inv[A[i]];
        if (i < SZ(A) - 1) {
            int x = SA.get_lcp(now, SA.inv[A[i + 1]]);
            while (not M.empty()) {
                auto it = --M.end();
                if (it->first > x) {
                    int y = it->second;
                    sum -= (vlong)(it->first - x) * y;
                    M.erase(it);
                    M[x] += y;
                } else {
                    break;
                }
            }
        }
        while (pos >= 0 and SA.inv[B[pos]] > now) {
            int x = SA.get_lcp(now, SA.inv[B[pos--]]);
            sum += x;
            M[x]++;
        }
        res += sum;
    }
    return res;
}

signed main()
{
    #ifdef LOCAL_RUN
    freopen("in.txt", "r", stdin);
    //freopen("out.txt", "w", stdout);
    cin.tie(0); cout.tie(0); ios_base::sync_with_stdio(false);
    int TTT; cin >> TTT; 
	if(TTT < 0) return gen(-TTT);
	while(TTT--){
    #else
    cin.tie(0); cout.tie(0); ios_base::sync_with_stdio(false);
    #endif // LOCAL_RUN

    ///CODE
    int n, q;
    cin >> n >> q;
    string str;
    cin >> str;
    SA.build(str);
    FOR(i, 1, q) {
        int k, l;
        cin >> k >> l;
        vector<int> A(k);
        for (int &x: A) {
            cin >> x;
            x--;
        }
        vector<int> B(l);
        for (int &x: B) {
            cin >> x;
            x--;
        }
        vlong res = solve(A, B);
        cout << res << "\n";
    }

    #ifdef LOCAL_RUN
    }
    #endif // LOCAL_RUN
    return 0;
}