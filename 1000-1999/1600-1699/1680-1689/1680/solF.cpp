//I wrote this code 4 u today

#include <bits/stdc++.h>



namespace IO {

    const int DPAIRSIZ = 1 << 18;

    char BB[DPAIRSIZ], *SS = BB, *TT = BB;



    inline char getcha() {

        return SS == TT && (TT = (SS = BB) + fread(BB, 1, DPAIRSIZ, stdin), SS == TT) ? EOF : *SS++;

    }



    template<typename T = int>

    inline T read() {

        T x = 0;

        int fu = 1;

        char c = getcha();

        while (c > 57 || c < 48) {

            if (c == 45) fu = -1;

            c = getcha();

        }

        while (c <= 57 && c >= 48) {

            x = x * 10 + c - 48;

            c = getcha();

        }

        x *= fu;

        return x;

    }



    template<typename T>

    inline void read(T &x) {

        x = 0;

        int fu = 1;

        char c = getcha();

        while (c > 57 || c < 48) {

            if (c == 45) fu = -1;

            c = getcha();

        }

        while (c <= 57 && c >= 48) {

            x = x * 10 + c - 48;

            c = getcha();

        }

        x *= fu;

    }



    template<typename T>

    inline void read(T *bg, T *ed) { while (bg != ed) read(*bg++); }



    inline void read(char &ch) {

        ch = getcha();

        while (ch <= 32) ch = getcha();

    }



    inline void read(char *s) {

        char ch = getcha();

        while (ch <= 32) ch = getcha();

        while (ch > 32) *s++ = ch, ch = getcha();

        *s = '\0';

    }



    inline void sread(char *s) {

        char ch = getcha();

        while (ch < 32) ch = getcha();

        while (ch >= 32) *s++ = ch, ch = getcha();

        *s = '\0';

    }



    inline void pread(char *&s) {

        char ch = getcha();

        while (ch <= 32) ch = getcha();

        while (ch > 32) *s++ = ch, ch = getcha();

        *s = '\0';

    }



    inline void spread(char *&s) {

        char ch = getcha();

        while (ch < 32) ch = getcha();

        while (ch >= 32) *s++ = ch, ch = getcha();

        *s = '\0';

    }



    template<typename T, typename ...Args>

    inline void read(T &x, Args &...args) {

        read(x);

        read(args...);

    }



    char out[DPAIRSIZ], *Out = out;

#define flush() fwrite(out, 1, Out - out, stdout)



    inline void putcha(char x) {

        *Out++ = x;

        if (Out - out >= (DPAIRSIZ)) flush(), Out = out;

    }



    template<typename T>

    inline void fprint(T x) {

        if (x < 0) putcha(45), x = -x;

        if (x > 9) fprint(x / 10);

        putcha(x % 10 + 48);

    }



    inline void print() { putcha(10); }



    template<typename T>

    inline void print(T x) {

        fprint(x);

        putcha(10);

    }



    inline void print(char *ch) {

        while (*ch != '\0') putcha(*(ch++));

        putcha(10);

    }



    inline void put(char *ch) { while (*ch != '\0') putcha(*(ch++)); }



    inline void print(const char *ch) {

        while (*ch != '\0') putcha(*(ch++));

        putcha(10);

    }



    inline void put(const char *ch) { while (*ch != '\0') putcha(*(ch++)); }



    template<typename T, typename ...Args>

    inline void print(T x, Args ...args) {

        fprint(x);

        putcha(32);

        print(args...);

    }



    template<typename ...Args>

    inline void print(const char *ch, Args ...args) {

        while (*ch != '\0') putcha(*(ch++));

        putcha(32);

        print(args...);

    }



    template<typename ...Args>

    inline void print(char *ch, Args ...args) {

        while (*ch != '\0') putcha(*(ch++));

        putcha(32);

        print(args...);

    }



    template<typename T, typename ...Args>

    inline void printl(T x, Args ...args) {

        fprint(x);

        putcha(10);

        printl(args...);

    }



    template<typename ...Args>

    inline void printl(const char *ch, Args ...args) {

        while (*ch != '\0') putcha(*(ch++));

        putcha(10);

        printl(args...);

    }



    template<typename ...Args>

    inline void printl(char *ch, Args ...args) {

        while (*ch != '\0') putcha(*(ch++));

        putcha(10);

        printl(args...);

    }



    template<typename T>

    inline void sprint(T x) {

        fprint(x);

        putcha(32);

    }



    template<typename T, typename ...Args>

    inline void sprint(T x, Args ...args) {

        fprint(x);

        putcha(32);

        sprint(args...);

    }



    template<typename T>

    inline void sprint(T *bg, T *ed) { while (bg != ed) sprint(*bg++); }



    template<typename T>

    inline void print(T *bg, T *ed) {

        while (bg != ed) sprint(*bg++);

        putcha(10);

    }



    template<typename T>

    inline void printl(T *bg, T *ed) { while (bg != ed) print(*bg++); }



    class AutoFlush {

    public:

        ~AutoFlush() { flush(); }

    } __AutoFlush;

} // namespace IO

using namespace IO;

#define vc vector



#define nd node*

#define pnd pair<nd, nd>



using namespace std;

typedef long long ll;

typedef vector<ll> vll;

typedef pair<ll, ll> pll;

typedef vc<pll> vpll;

typedef vc<vll> vvll;

typedef vc<vpll> vvpll;



template<const ll MOD>

struct mod_mul : std::multiplies<const ll> {

    ll operator()(const ll a, const ll b) {

        return (a * b) % MOD;

    }

};





template<typename T>

inline void sort(T &a) {

    sort(a.begin(), a.end());

}



template<typename T>

inline void unique(T &a) {

    a.resize(unique(a.begin(), a.end()) - a.begin());

}



template<typename T>

inline void reverse(T &a) {

    reverse(a.begin(), a.end());

}



const ll INF = 9023372036854775808ll;

const ll MOD = 1000000007ll;



vc<pair<int, int>> g[1000005];



bool used[1000005];

bool color[1000005];

int same[1000005];

int diff[1000005];

int h[1000005];



vc<int> d[1000005];



int cc = 0;

bool kk;



void cfs(int v, int p = -1) {

    used[v] = true;

    for (auto [to, i]: g[v]) {

        if (i == p) continue;

        if (!used[to]) color[to] = !color[v], cfs(to, i), h[to] = h[v] + 1, d[v].push_back(to);

        else {

            if (h[to] > h[v]) {

                if (color[v] == color[to]) ++same[to], --same[v], ++cc, kk = color[v];

                else ++diff[to], --diff[v];

            }

        }

    }

}



bool ans = true;



int vv;



pair<int, int> dfs(int v) {

    pair<int, int> x = {same[v], diff[v]};

    for (auto to: d[v]) {

        auto tmp = dfs(to);

        x.first += tmp.first, x.second += tmp.second;

    }

    if (x.first == cc && !x.second) ans = true, vv = v, kk = !color[v];

    return x;

}



void hu(int v) {

    color[v] ^= 1;

    for (auto to: d[v]) {

        hu(to);

    }

}



int main() {

    int t;

    read(t);

    while (t--) {

        int n, m;

        read(n, m);

        cc = 0;

        ans = false;

        kk = false;

        for (int v = 1; v <= n; ++v)

            d[v].clear(), h[v] = 0, diff[v] = 0, same[v] = 0, color[v] = false, g[v].clear(), used[v] = false;

        for (int i = 0; i < m; ++i) {

            int v, u;

            read(v, u);

            g[v].emplace_back(u, i);

            g[u].emplace_back(v, i);

        }

        cfs(1);

        for (int v = 1; v <= n; ++v) used[v] = false;

        if (cc <= 1) {

            print("YES");

            for (int v = 1; v <= n; ++v) {

                putcha((color[v] ^ (!kk)) + '0');

            }

            putcha('\n');

            continue;

        }

        dfs(1);

        if (ans) {

            hu(vv);

            print("YES\n");

            for (int v = 1; v <= n; ++v) {

                putcha((color[v] ^ (!kk)) + '0');

            }

            putcha('\n');

        } else {

            print("NO\n");

        }

    }

}