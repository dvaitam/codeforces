#pragma optimize("SEX_ON_THE_BEACH")
#pragma GCC optimize("unroll-loops")
#pragma GCC optimize("unroll-all-loops")
#pragma GCC optimize("O3")
 
#pragma GCC optimize("Ofast")
#pragma GCC optimize("fast-math")
//#define _FORTIFY_SOURCE 0 
#pragma GCC optimize("no-stack-protector")
#pragma GCC target("sse,sse2,sse3,ssse3,avx,avx2,popcnt,abm,mmx") 

#include<bits/stdc++.h>
#include <x86intrin.h>

using uint = unsigned int;
using ll = long long int;
using ull = unsigned long long int;
using dd = double;
using ldd = long double;
using pii = std::pair<int, int>;
using pll = std::pair<ll, ll>;
using pdd = std::pair<dd, dd>;
using pld = std::pair<ldd, ldd>;

namespace fast {
    template<typename T>
    T gcd(T a, T b) {
        return gcd(a, b);
    }

    template<>
    unsigned int gcd<unsigned int>(unsigned int u, unsigned int v) {
        int shift;
        if (u == 0)
            return v;
        if (v == 0)
            return u;
        shift = __builtin_ctz(u | v);
        u >>= __builtin_ctz(u);
        do {
            unsigned int m;
            v >>= __builtin_ctz(v);
            v -= u;
            m = (int)v >> 31;
            u += v & m;
            v = (v + m) ^ m;
        } while (v != 0);
        return u << shift;
    }

    template<>
    unsigned long long gcd<unsigned long long>(unsigned long long u, unsigned long long v) {
        int shift;
        if (u == 0)
            return v;
        if (v == 0)
            return u;
        shift = __builtin_ctzll(u | v);
        u >>= __builtin_ctzll(u);
        do {
            unsigned long long m;
            v >>= __builtin_ctzll(v);
            v -= u;
            m = (long long)v >> 63;
            u += v & m;
            v = (v + m) ^ m;
        } while (v != 0);
        return u << shift;
    }

    template<size_t N>
    struct bitset {
        size_t sz = (N + 255) / 256;
        ull arr[(N + 63) / 64];
        ull inv = -1;
        ull last_item_mask = N % 64 ? ((1ULL << (N % 64)) - 1) : inv;

        bitset() {
            static_assert(N > 0, "empty bitset not available");
            reset();
        }

        size_t size() {
            return N;
        }

        size_t count() {
            arr[sz - 1] &= last_item_mask;
            size_t res = 0;
            for (size_t i = 0; i < sz; ++i) {
                res += __builtin_popcountll(arr[i]);
            }
            return res;
        }

        void flip() {
            for (size_t i = 0; i < sz; ++i) {
                arr[i] ^= inv;
            }
        }

        void flip(size_t l) {
            size_t lb = l / 64, lp = l % 64;
            arr[lb] ^= 1ULL << lp;
        }

        void flip(size_t l, size_t r) {
            size_t lb = l / 64, lp = l % 64;
            size_t rb = r / 64, rp = r % 64;
            arr[lb] ^= (1ULL << lp) - 1;
            arr[rb] ^= (rp == 63 ? inv : (1ULL << (rp + 1)) - 1);
            if (lb == rb) return;
            arr[lb] ^= inv;
            ++lb;
            while (lb < rb) {
                arr[lb++] ^= inv;
            }
        }

        bool get(size_t l) {
            size_t lb = l / 64, lp = l % 64;
            return (arr[lb] >> lp) & 1;
        }

        void set() {
            for (size_t i = 0; i < sz; ++i) {
                arr[i] = inv;
            }
        }

        void set(size_t l) {
            size_t lb = l / 64, lp = l % 64;
            arr[lb] |= 1ULL << lp;
        }

        void set(size_t l, size_t r) {
            size_t lb = l / 64, lp = l % 64;
            size_t rb = r / 64, rp = r % 64;
            if (lb == rb) {
                arr[lb] |= ((1ULL << (rp - lp + 1)) - 1) << lp;
                return;
            }
            arr[lb] |= inv << lp;
            arr[rb] |= (rp == 63 ? inv : (1ULL << (rp + 1)) - 1);
            ++lb;
            while (lb < rb) {
                arr[lb++] |= inv;
            }
        }

        void reset() {
            for (int i = 0; i < sz; ++i) {
                arr[i] = 0ULL;
            }
        }

        void reset(size_t l) {
            size_t lb = l / 64, lp = l % 64;
            arr[lb] &= inv ^ (1ULL << lp);
        }

        void reset(size_t l, size_t r) {
            size_t lb = l / 64, lp = l % 64;
            size_t rb = r / 64, rp = r % 64;
            if (lb == rb) {
                arr[lb] &= inv ^ ((1ULL << (rp - lp + 1) - 1) << lp);
                return;
            }
            arr[lb] &= (1ULL << lp) - 1;
            arr[rb] &= (rp == 63 ? 0 : (inv << (rp + 1)));
        }

        bitset& operator|=(bitset<N>& other) {
            for (size_t i = 0; i < sz; ++i) {
                arr[i] |= other.arr[i];
            }
            return *this;
        }

        bitset& operator&=(bitset<N>& other) {
            for (size_t i = 0; i < sz; ++i) {
                arr[i] &= other.arr[i];
            }
            return *this;
        }

        bitset& operator^=(bitset<N>& other) {
            for (size_t i = 0; i < sz; ++i) {
                arr[i] ^= other.arr[i];
            }
            return *this;
        }

        bitset& operator<<=(size_t shift) {
            size_t sb = shift / 64, sp = shift % 64;
            //big shift per blocks
            if (sb) {
                size_t i;
                for (i = sz - 1; i >= sb; --i) {
                    arr[i] = arr[i - sb];
                }
                for (; i > 0; --i) {
                    arr[i] = 0;
                }
                arr[0] = 0;
            }
            //small shift 
            if (sp) {
                size_t sp_rev = 64 - sp;
                for (size_t i = sz - 1; i > 0; --i) {
                    arr[i] <<= sp;
                    arr[i] |= arr[i - 1] >> sp_rev;
                }
                arr[0] <<= sp;
            }

            return *this;
        }

        bitset& operator>>=(size_t shift) {
            size_t sb = shift / 64, sp = shift % 64;
            //big shift per blocks
            if (sb) {
                size_t i;
                for (i = 0; i + sb < sz; ++i) {
                    arr[i] = arr[i + sb];
                }
                for (; i < sz; ++i) {
                    arr[i] = 0;
                }
            }
            //small shift 
            if (sp) {
                size_t sp_rev = 64 - sp;
                for (size_t i = 0; i < sz - 1; ++i) {
                    arr[i] >>= sp;
                    arr[i] |= (arr[i + 1] & ((1ULL << sp) - 1)) << sp_rev;
                }
                arr[sz - 1] >>= sp;
            }

            return *this;
        }

        size_t next_zero(size_t l) {
            arr[sz - 1] |= ~last_item_mask;
            size_t lb = l / 64, lp = l % 64;
            if ((~arr[lb]) >> lp) {
                return lb * 64 + __builtin_ctzll((~arr[lb]) >> lp) + lp;
            }
            ++lb;
            while (lb < sz) {
                if (~arr[lb]) {
                    return lb * 64 + __builtin_ctzll(~arr[lb]);
                }
                ++lb;
            }
            return size();
        }

        size_t first_zero() {
            arr[sz - 1] |= ~last_item_mask;
            for (size_t i = 0; i < sz; ++i) {
                if (~arr[i]) {
                    return i * 64 + __builtin_ctzll(~arr[i]);
                }
            }
            return size();
        }

        size_t next_one(size_t l) {
            arr[sz] &= last_item_mask;
            size_t lb = l / 64, lp = l % 64;
            if (arr[lb] >> lp) {
                return lb * 64 + __builtin_ctzll(arr[lb] >> lp) + lp;
            }
            ++lb;
            while (lb < sz) {
                if (arr[lb]) {
                    return lb * 64 + __builtin_ctzll(arr[lb]);
                }
                ++lb;
            }
            return size();
        }

        size_t first_one() {
            arr[sz] &= last_item_mask;
            for (size_t i = 0; i < sz; ++i) {
                if (arr[i]) {
                    return i * 64 + __builtin_ctzll(arr[i]);
                }
            }
            return size();
        }

        size_t last_one() {
            arr[sz - 1] &= last_item_mask;
            for (size_t i = sz - 1; i > 0; --i) {
                if (arr[i]) {
                    return i * 64 + 63 - __builtin_clzll(arr[i]);
                }
            }
            if (arr[0])
                return 63 - __builtin_clzll(arr[0]);
            return -1;
        }

        size_t last_zero() {
            arr[sz - 1] |= ~last_item_mask;
            for (size_t i = sz - 1; i > 0; --i) {
                if (~arr[i]) {
                    return i * 64 + 63 - __builtin_clzll(~arr[i]);
                }
            }
            if (~arr[0])
                return 63 - __builtin_clzll(~arr[0]);
            return -1;
        }
    };

    template<size_t N>
    std::ostream& operator<<(std::ostream& os, bitset<N> other) {
        for (int i = 0; i < other.size(); ++i) {
            os << other.get(i);
        }
        return os;
    }
}
 
namespace someUsefull {
    template<typename T1, typename T2>
    inline void checkMin(T1& a, T2 b) {
        if (a > b)
            a = b;
    }
 
    template<typename T1, typename T2>
    inline void checkMax(T1& a, T2 b) {
        if (a < b)
            a = b;
    }

    template<typename T1, typename T2>
    inline bool checkMinRes(T1& a, T2 b) {
        if (a > b) {
            a = b;
            return true;
        }
        return false;
    }

    template<typename T1, typename T2>
    inline bool checkMaxRes(T1& a, T2 b) {
        if (a < b) {
            a = b;
            return true;
        }
        return false;
    }
}
 
namespace operators {
    template<typename T1, typename T2>
    std::istream& operator>>(std::istream& in, std::pair<T1, T2>& x) {
        in >> x.first >> x.second;
        return in;
    }
 
    template<typename T1, typename T2>
    std::ostream& operator<<(std::ostream& out, std::pair<T1, T2> x) {
        out << x.first << " " << x.second;
        return out;
    }
 
    template<typename T1>
    std::istream& operator>>(std::istream& in, std::vector<T1>& x) {
        for (auto& i : x) in >> i;
        return in;
    }
 
    template<typename T1>
    std::ostream& operator<<(std::ostream& out, std::vector<T1>& x) {
        for (auto& i : x) out << i << " ";
        return out;
    }
}
 
//name spaces
using namespace std;
using namespace operators;
using namespace someUsefull;
//end of name spaces
 
//defines
#define ff first
#define ss second
#define all(x) (x).begin(), (x).end()
#define rall(x) (x).rbegin(), (x).rend()
#define NO {cout << "NO"; return;}
#define YES {cout << "YES"; return;}
//end of defines

//debug defines
#ifdef HOME
    #define debug(x) cerr << #x << " " << (x) << endl;
    #define debug_v(x) {cerr << #x << " "; for (auto ioi : x) cerr << ioi << " "; cerr << endl;}
    #define debug_vp(x) {cerr << #x << " "; for (auto ioi : x) cerr << '[' << ioi.ff << " " << ioi.ss << ']'; cerr << endl;}
    #define debug_v_v(x) {cerr << #x << "/*\n"; for (auto ioi : x) { for (auto ioi2 : ioi) cerr << ioi2 << " "; cerr << '\n';} cerr << "*/" << #x << endl;}
    int jjj;
    #define wait() cin >> jjj;
    #define PO cerr << "POMELO" << endl;
    #define OL cerr << "OLIVA" << endl;
    #define gen_clock(x) cerr << "Clock " << #x << " created" << endl; ll x = clock(); 
    #define check_clock(x) cerr << "Time spent in " << #x << ": " << clock() - x << endl; x = clock();
    #define reset_clock(x) x = clock();
    #define say(x) cerr << x << endl;
#else
    #define debug(x) 0;
    #define debug_v(x) 0;
    #define debug_vp(x) 0;
    #define debug_v_v(x) 0;
    #define wait() 0;
    #define PO 0;
    #define OL 0;
    #define gen_clock(x) 0;
    #define check_clock(x) 0;
    #define reset_clock(x) 0;
    #define say(x) 0;
#endif // HOME

void solve(int test) {
    int l;
    vector<int> cnt(3);
    cin >> cnt >> l;
    vector<pair<int, pair<int, int>>> ans;
    vector<int> lt(3, 0);
    for (int t = 1; t <= l; ++t) {
        int mi = -1;
        int then_task = -1;
        for (int i = 0; i < 3; ++i) {
            int ht = t - lt[i];
            int can = -1;
            for (int j = min(2, ht - 2); j >= 0; --j) {
                if (cnt[j]) {
                    can = j;
                    break;
                }
            }
            if (can == -1) continue;
            if (mi == -1) {
                mi = i;
                then_task = can;
                continue;
            }
            if (then_task < can) {
                mi = i;
                then_task = can;
                continue;
            }
            if (then_task > can) {
                continue;
            }
            if (lt[mi] < lt[i]) {
                mi = i;
            }
        }
        if (mi == -1) continue;
        cnt[then_task]--;
        lt[mi] = t;
        ans.push_back({mi + 1, {t - then_task - 2, t}});
    }
    cout << ans.size() << '\n';
    for (auto i : ans) cout << i << '\n';
}

signed main() {
    ios_base::sync_with_stdio(false);
    cout.tie(0);
    cin.tie(0);
    //freopen("file.in", "r", stdin);
    //freopen("file.out", "w", stdout);

    int t = 1;
    //cin >> t;
    for (int i = 0; i < t; ++i) {
        solve(i+1);
        cout << '\n';
        //PO;
    }
    return 0;
}
/*

*/