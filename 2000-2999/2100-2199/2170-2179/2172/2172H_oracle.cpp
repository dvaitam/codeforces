#include <cstdint>
#include <vector>
#include <cassert>
 
namespace hsh {
 
    inline namespace hash_mod {
 
        using underlying = std::uint64_t;
 
        constexpr underlying mod = (underlying(1) << 61) - 1; // prime
 
        inline constexpr underlying __plus(underlying a, underlying b) noexcept {
            a += b;
            return a >= mod? a - mod : a;
        }
        inline constexpr underlying __subt(underlying a, underlying b) noexcept {
            return a < b? a + (mod - b) : a - b;
        }
        inline constexpr underlying __mod(__uint128_t x) noexcept {
            return __plus(x >> 61, x & mod);
        }
 
        struct modval {
            underlying val;
 
            friend inline constexpr modval operator + (modval a, modval b) noexcept {
                return {__plus(a.val, b.val)};
            }
            friend inline constexpr modval operator - (modval a, modval b) noexcept {
                return {__subt(a.val, b.val)};
            }
            friend inline constexpr modval operator * (modval a, modval b) noexcept {
                return {__mod(__uint128_t(a.val) * b.val)};
            }
            friend inline constexpr modval fma(modval a, modval b, underlying c) noexcept {
                return {__mod(__uint128_t(a.val) * b.val + c)};
            }
 
            friend inline constexpr modval operator - (modval x) noexcept {
                return x.val? modval{mod - x.val} : x;
            }
 
            inline constexpr modval& operator += (modval b) noexcept {
                return *this = *this + b;
            }
            inline constexpr modval& operator -= (modval b) noexcept {
                return *this = *this - b;
            }
            inline constexpr modval& operator *= (modval b) noexcept {
                return *this = *this * b;
            }
 
            friend inline constexpr bool operator == (modval a, modval b) noexcept {
                return a.val == b.val;
            }
            friend inline constexpr bool operator != (modval a, modval b) noexcept {
                return a.val != b.val;
            }
        };
 
        inline void fill_power(std::size_t n, modval *dst, modval b) noexcept {
            if (!n)return;
            dst[0] = {1};
            for (std::size_t i = 1; i < n; ++i)dst[i] = dst[i - 1] * b;
        }
 
    }
 
    constexpr modval mul = {13331};
 
    struct hash_seq {
 
        std::vector<modval> hs;
        std::vector<modval> p;
 
        hash_seq() = default;
 
        template <typename _Iter>
            void build(_Iter bg, _Iter ed){
                const std::size_t n = std::distance(bg, ed);
                hs.resize(n + 1);
                hs[0] = {};
                for (std::size_t i = 0; i < n; ++i)hs[i + 1] = fma(hs[i], mul, 1 + underlying(bg[i]));
                if (p.size() < n + 1){
                    p.resize(n + 1);
                    fill_power(n + 1, p.data(), mul);
                }
            }
 
        template <typename _Iter>
            hash_seq(_Iter bg, _Iter ed){
                build(bg, ed);
            }
 
        typedef modval range_hash_t;
 
        inline range_hash_t operator () (std::size_t l, std::size_t r) const {
            assert(l <= r && r < hs.size());
            return hs[r] - p[r - l] * hs[l];
        }
 
    };
 
    struct hash_seq_rseq {
 
        std::vector<modval> hs;
        std::vector<modval> rhs;
        std::vector<modval> p;
 
        hash_seq_rseq() = default;
 
        template <typename _Iter>
            void build(_Iter bg, _Iter ed){
                const std::size_t n = std::distance(bg, ed);
                hs.resize(n + 1);
                hs[0] = {};
                for (std::size_t i = 0; i < n; ++i)hs[i + 1] = fma(hs[i], mul, 1 + underlying(bg[i]));
                rhs.resize(n + 1);
                rhs[n] = {};
                for (std::size_t i = n - 1; ~i; --i)rhs[i] = fma(rhs[i + 1], mul, 1 + underlying(bg[i]));
                if (p.size() < n + 1){
                    p.resize(n + 1);
                    fill_power(n + 1, p.data(), mul);
                }
            }
 
        template <typename _Iter>
            hash_seq_rseq(_Iter bg, _Iter ed){
                build(bg, ed);
            }
 
        typedef modval range_hash_t;
 
        inline range_hash_t hash(std::size_t l, std::size_t r) const {
            assert(l <= r && r < hs.size());
            return hs[r] - p[r - l] * hs[l];
        }
        inline range_hash_t revhash(std::size_t l, std::size_t r) const {
            assert(l <= r && r < hs.size());
            return rhs[l] - p[r - l] * rhs[r];
        }
 
    };
 
    template <typename _Iter>
        inline modval hash_of(_Iter bg, _Iter ed){
            modval hs = 0;
            for (; bg != ed; ++bg)hs = fma(hs, mul, underlying(*bg));
            return hs;
        }
 
    template <typename _Iter>
        inline modval revhash_of(_Iter bg, _Iter ed){
            modval hs = 0;
            for (; bg != ed; )hs = fma(hs, mul, underlying(*--ed));
            return hs;
        }
 
}
 
//utf-8
 
#include <string>
#include <type_traits>
#include <iterator>
#include <cstdint>
#include <vector>
#include <array>
#include <algorithm>
#include <limits>
#include <cassert>
#include <numeric>
 
#include <vector>
#include <cstdint>
#include <algorithm>
#include <vector>
#include <cstdint>
#include <algorithm>
 
namespace rmq {
 
    using std::vector;
    using size_t = std::uint32_t;
 
    enum rmq_type { rmq_min, rmq_max };
 
    template <typename value_t, rmq_type type_>
        struct rm { };
 
    template <typename value_t>
        struct rm <value_t, rmq_min> {
            inline value_t operator () (const value_t &a, const value_t &b) const {
                return std::min(a, b);
            }
        };
 
    template <typename value_t>
        struct rm <value_t, rmq_max> {
            inline value_t operator () (const value_t &a, const value_t &b) const {
                return std::max(a, b);
            }
        };
 
    template <typename value_t, rmq_type type_>
        struct rmq_solver {
            typedef value_t value_type;
            static constexpr rmq_type type = type_;
            static rm<value_t, type_> mrg;
            vector<vector<value_t> > st;
 
            template <typename cont>
                void build(const cont &c){
                    size_t n = c.size();
                    st.resize(std::__lg(n) + 1);
                    st[0].resize(n);
                    std::copy_n(c.begin(), n, st[0].begin());
                    for (size_t i = 1, _ = std::__lg(n); i <= _; ++i){
                        size_t l = n - (1 << i) + 1;
                        st[i].resize(l);
                        const vector<value_t> &pr = st[i - 1];
                        typedef typename vector<value_t>::const_iterator iter;
                        iter _0 = pr.begin(), _1 = _0 + (1 << (i - 1));
                        for (size_t j = 0; j < l; ++j)
                            st[i][j] = mrg(_0[j], _1[j]);
                    }
                }
 
            inline value_t get(size_t l, size_t r) const {
                size_t t = std::__lg(r - l);
                return mrg(st[t][l], st[t][r - (1 << t)]);
            }
 
        };
 
    template <typename value_t, rmq_type type_>
        rm<value_t, type_> rmq_solver<value_t, type_>::mrg{};
 
}
 
namespace string_algo {
 
    using std::basic_string;
    using std::vector;
    using std::array;
    using size_t = std::uint32_t;
 
    template <typename iter_t, typename iter_tag>
        using ass_iter_tag = std::enable_if<std::is_base_of<iter_tag, typename std::iterator_traits<iter_t>::iterator_category>::value >;
    template <typename iter_t>
        using ass_RAI = ass_iter_tag<iter_t, std::random_access_iterator_tag>;
    #define ass_is_RAI(iter_t) typename = typename ass_RAI<iter_t>::type
 
    // 前缀函数
    template <typename value_t, typename output_iter_t, ass_is_RAI(output_iter_t)>
        void prefix_function(const basic_string<value_t> &s, output_iter_t pi){//[0,pi_i)=[i-pi_i,i)
            if (s.empty())return;
            size_t n = s.size();
            *pi = 0;
            for (size_t i = 1, j = 0; i != n; ++i){
                while (j && s[i] != s[j])j = pi[j - 1];
                if (s[i] == s[j])++j;
                pi[i] = j;
            }
        }
    // KMP 算法
    template <typename value_t, typename output_iter_t, ass_is_RAI(output_iter_t)>
        vector<size_t> kmp(const basic_string<value_t> &text, const basic_string<value_t> &pat, output_iter_t pi){//返回匹配位置(l point)
            if (text.size() < pat.size())return {};
            if (text.size() == pat.size()){if (text == pat)return {0};return {};}
            if (pat.empty())return {};
            prefix_function(pat, pi);
            size_t n = text.size(), m = pat.size();
            vector<size_t> ret;
            for (size_t i = 0, j = 0; i != n; ++i){
                while (j && text[i] != pat[j])j = pi[j - 1];
                if (text[i] == pat[j])++j;//[i-m+1,i]=[0,j)
                if (j == m)ret.emplace_back(i - m + 1), j = pi[j - 1];
            }
            return ret;
        }
 
    // z 函数
    template <typename value_t, typename output_iter_t, ass_is_RAI(output_iter_t)>
        void z_function(const basic_string<value_t> &s, output_iter_t z){//z_i = lcp(s,s[i:])
            if (s.empty())return;
            size_t n = s.size();
            *z = n;
            for (size_t i = 1, l = 0, r = 1; i != n; ++i){//[l,r)
                if (i < r && z[i - l] < r - i)z[i] = z[i - l];
                else {
                    if (i < r)z[i] = r - i; else z[i] = 0;
                    while (i + z[i] < n && s[z[i]] == s[i + z[i]])++z[i];
                }
                if (i + z[i] > r)l = i, r = i + z[i];
            }
        }
 
    // 回文自动机
    template <size_t _val_max>
        struct pam {
            static constexpr size_t value_max = _val_max;// 字符集为 [0,value_max) 内的整数
            basic_string <size_t> str;
            struct node {
                size_t father;
                array<size_t, value_max> children;
                size_t length;
                size_t fail;
                size_t count;// 以当前节点对应位置结束的回文子串数量
                size_t total;// 当前节点对应子串出现次数
                size_t endpos;// [endpos - length, endpos)
                node(size_t _len = 0) : father{}, children{}, length{_len}, fail{}, count{}, total{}, endpos{} {}
            };
            vector<node> tree;
            size_t last;
            pam() : str{}, tree{}, last{1} {
                tree.reserve(2);
                tree.emplace_back();// 偶根
                tree.emplace_back(-1);// 奇根
                tree[0].fail = 1;
            }
            void clear(){
                str.clear();
                tree.clear();
                tree.reserve(2);
                tree.emplace_back();// 偶根
                tree.emplace_back(-1);// 奇根
                tree[0].fail = 1;
                last = 1;
            }
            size_t _get_fail(size_t now){
                const size_t length = str.size();
                while (length < tree[now].length + 2 || str[length - tree[now].length - 2] != str.back())
                    now = tree[now].fail;
                return now;
            }
            size_t push_back(size_t ch){// 返回 以当前节点对应位置结束的回文子串数量
                str.push_back(ch);
                size_t now = _get_fail(last);
                if (!tree[now].children[ch]){
                    tree.emplace_back(tree[now].length + 2);
                    tree.back().father = now;
                    tree.back().fail = tree[_get_fail(tree[now].fail)].children[ch];
                    tree.back().count = tree[tree.back().fail].count + 1;
                    tree.back().endpos = str.size();
                    tree[now].children[ch] = tree.size() - 1;
                }
                last = tree[now].children[ch];
                ++tree[last].total;
                return tree[last].count;
            }
            void build_total(){
                for (size_t i = tree.size() - 1; i > 1; --i)
                    tree[tree[i].fail].total += tree[i].total;
            }
            size_t init_node(size_t str_length){
                return str_length & 1;
            }
            size_t next(size_t now, size_t ch){
                return tree[now].children[ch];
            }
            bool is_psubstr_of(const basic_string <size_t> &t){// 判断是否为回文子串
                size_t nw = init_node(t.size()), l = (t.size() - 1) >> 1, r = t.size() - 1 - l;
                while (~l){
                    if (t[l] != t[r])return false;
                    nw = next(nw, t[l]), --l, ++r;
                    if (!nw)return false;
                }
                return true;
            }
            template <typename func_t>// 遍历所有非空回文子串
                void for_each_psubstr(func_t func){
                    for (size_t i = 2, _ = tree.size(); i != _; ++i)
                        func(tree[i].endpos - tree[i].length, tree[i].endpos, tree[i].count, tree[i].total);// [left,right), count, total
                }
        };
 
    // 后缀排序
    template <typename value_t, typename output_iter_t, ass_is_RAI(output_iter_t)>
        void suffix_sort(const basic_string<value_t> &s, output_iter_t sa,
                std::enable_if_t<std::is_unsigned<value_t>::value>* = nullptr){// 返回的 sa_i 在 [0,n) 中
            if (s.empty())return;
            size_t n = s.size();
            size_t *nrk, *c, *id = new size_t [n];
            size_t m;
            {
                size_t maxv = *max_element(s.begin(), s.end());
                if (maxv >= n << 2){
                    c = new size_t [m = n << 1] {};
                    nrk = new size_t [m] {};
                    auto __pred = [&s](size_t u, size_t v){ return s[u] < s[v]; };
                    std::iota(sa, sa + n, 0);
                    std::sort(sa, sa + n, __pred);
                    for (size_t i = 0; i < n; ++i)nrk[i] = std::lower_bound(sa, sa + n, i, __pred) - sa;
                } else {
                    c = new size_t [m = std::max(maxv + 1, n << 1)] {};
                    nrk = new size_t [m] {};
                    std::copy_n(s.begin(), n, nrk);
                    for (size_t i = 0; i < n; ++i)++c[nrk[i]];
                    for (size_t i = 1; i <= maxv; ++i)c[i] += c[i - 1];
                    for (size_t i = n - 1; ~i; --i)sa[--c[nrk[i]]] = i;
                }
            }
            for (size_t i = 0, l = 1; l <= n && m != n; ++i, l <<= 1){
                {
                    size_t t = 0;
                    for (size_t j = n - l; j < n; ++j)id[t++] = j;
                    for (size_t j = 0; j < n; ++j)if (sa[j] >= l)id[t++] = sa[j] - l;
                }
                std::fill_n(c, m, 0);
                for (size_t j = 0; j < n; ++j)++c[nrk[j]];
                for (size_t j = 1; j < m; ++j)c[j] += c[j - 1];
                for (size_t j = n - 1; ~j; --j)sa[--c[nrk[id[j]]]] = id[j];
                std::swap(c, nrk);
                std::fill_n(nrk + n, n, size_t(-1));
                nrk[*sa] = 0;
                for (size_t j = m = 1; j < n; ++j){
                    if (c[sa[j]] != c[sa[j - 1]] || c[sa[j] + l] != c[sa[j - 1] + l])++m;
                    nrk[sa[j]] = m - 1;
                }
            }
            delete []nrk; delete []c; delete []id;
        }
    template <typename value_t, typename output_iter_t, ass_is_RAI(output_iter_t)>
        void suffix_sort(const basic_string<value_t> &s, output_iter_t sa,
                std::enable_if_t<std::is_signed<value_t>::value &&
                    std::is_integral<value_t>::value && std::is_unsigned<
                        typename std::make_unsigned<value_t>::type>::value>* = nullptr){// 返回的 sa_i 在 [0,n) 中
            using U_t = typename std::make_unsigned<value_t>::type;
            basic_string<U_t> ss; ss.resize(s.size());
            std::transform(s.begin(), s.end(), ss.begin(),
                [](value_t x){ return U_t(x) - std::numeric_limits<value_t>::min(); });
            suffix_sort<U_t, output_iter_t>(ss, sa);
        }
 
    template <typename value_t>
        struct string_suffix {
            typedef value_t value_type;
            vector<size_t> sa, rk, height;//height_0 = 0   height_i = lcp(sa_i,sa_{i-1})
            rmq::rmq_solver<size_t, rmq::rmq_min> rmqs;
            void build(const basic_string<value_t> &s){//l=1e6, value=char, ~1s
                size_t n = s.size();
                sa.resize(n), rk.resize(n), height.resize(n);
                suffix_sort(s, sa.begin());
                for (size_t i = 0; i < n; ++i)rk[sa[i]] = i;
                height[0] = 0;
                for (size_t i = 0, k = 0; i < n; ++i){
                    if (rk[i] == 0)continue;
                    if (k)--k;
                    size_t mx = std::min(n - i, n - sa[rk[i] - 1]);
                    while (k < mx && s[i + k] == s[sa[rk[i] - 1] + k])++k;
                    height[rk[i]] = k;
                }
                rmqs.build(height);
            }
            string_suffix() = default;
            string_suffix(const basic_string<value_t> &s){ build(s); }
            inline size_t lcp(size_t u, size_t v) const {
                if (u == v)return sa.size() - u;
                u = rk[u], v = rk[v];
                return u > v? rmqs.get(v + 1, u + 1) : rmqs.get(u + 1, v + 1);
            }
            struct substr_t {
                const string_suffix* ss;
                size_t l, r;
                enum cmp_t {
                    less_no_prefix,//a < b, a not prefix of b
                    prefix_of,//a is prefix of b
                    same,//a=b
                    include,//b is prefix of a
                    greater_no_include,//a>b, b not prefix of a
                };
                inline cmp_t cmp_with(const substr_t &o) const {
                    assert(ss == o.ss);
                    size_t lc = ss->lcp(l, o.l);
                    size_t L = r - l, oL = o.r - o.l;
                    if (L == oL && lc >= r - l)return same;
                    if (L < oL && lc >= L)return prefix_of;
                    if (L > oL && lc >= oL)return include;
                    return ss->rk[l + lc] < ss->rk[o.l + lc]? less_no_prefix : greater_no_include;
                }
                inline bool operator < (const substr_t &o) const {
                    return cmp_with(o) < same;
                }
                inline bool operator > (const substr_t &o) const {
                    return cmp_with(o) > same;
                }
                inline bool operator <= (const substr_t &o) const {
                    return cmp_with(o) <= same;
                }
                inline bool operator >= (const substr_t &o) const {
                    return cmp_with(o) >= same;
                }
                inline bool operator == (const substr_t &o) const {
                    return cmp_with(o) == same;
                }
                inline bool operator != (const substr_t &o) const {
                    return cmp_with(o) != same;
                }
            };
            substr_t substr(size_t l, size_t r) const {//[l,r)
                return substr_t{this, l, r};
            }
            struct suffix_t {
                const string_suffix* ss;
                size_t p;
                bool operator < (const suffix_t &o) const {
                    assert(ss == o.ss);
                    return ss->rk[p] < ss->rk[o.p];
                }
                bool operator > (const suffix_t &o) const {
                    assert(ss == o.ss);
                    return ss->rk[p] > ss->rk[o.p];
                }
                explicit operator substr_t () const {
                    return substr_t{ss, p, ss->sa.size()};
                }
            };
            suffix_t operator [] (size_t p) const {
                return suffix_t{this, p};
            }
        };
 
    // 马拉车算法
    template <typename value_t, typename output_iter_t, ass_is_RAI(output_iter_t)>
        void manacher(const basic_string<value_t> &s, output_iter_t even, output_iter_t odd){
            // even_i 为 s_{i-even_i+1} ~ s_{i+even_i} 为回文串的最大值 （0<=i<n-1）， odd_i 为 s_{i-odd_i} ~ s_{i+odd_i} 为回文串的最大值 （0<=i<n）
            if (s.empty())return;
            size_t n = s.size();
            if (n == 1){*odd = 0;return;}
            if (s[0] == s[1])*even = 1;else *even = 0;
            for (size_t i = 1, l = 1 - *even, r = *even; i < n - 1; ++i){// [l,r] 为满足 ((l+r-1)/2) < i 且 r 最大的回文子串
                even[i] = l > r || i >= r? (size_t)0 : std::min<size_t>(r - i, even[l + r - i - 1]);
                while (i >= even[i] && i + even[i] + 1 < n && s[i - even[i]] == s[i + even[i] + 1])++even[i];
                if (i + even[i] > r)l = i - even[i] + 1, r = i + even[i];
            }
            *odd = 0;
            for (size_t i = 1, l = 0, r = 0; i < n; ++i){
                odd[i] = i >= r? (size_t)0 : std::min<size_t>(r - i, odd[l + r - i]);
                while (i >= odd[i] + 1 && i + odd[i] + 1 < n && s[i - odd[i] - 1] == s[i + odd[i] + 1])++odd[i];
                if (i + odd[i] > r)l = i - odd[i], r = i + odd[i];
            }
        }
 
    // 最小表示法
    template <typename iter_t, ass_is_RAI(iter_t)>
        iter_t minimal_string(iter_t bg, iter_t ed){
            if (bg == ed)return bg;
            iter_t i = bg, j = bg + 1;
            while (i != ed && j != ed){
                if (i > j)swap(i, j);
                size_t k = 0;
                while (j + k != ed && *(i + k) == *(j + k))++k;
                if (j + k == ed)break;
                if (*(i + k) < *(j + k))j += k + 1;
                else i += k + 1;
                if (i == j)++j;
            }
            return std::min(i, j);
        }
 
}
 
#include <bits/stdc++.h>
using namespace std;
using mod = hsh::modval;
int k;
long long t;
char a[1 << 19];
int L, C;
mod bv[1 << 19]; // section hash value begin at i
hsh::hash_seq hs[1 << 18];
string_algo::string_suffix<hsh::underlying> ss;
bool lss(int u, int v){
    if (u == v)return 0;
    int i = ss.lcp(u, v);
    if (i >= C)return 0;
    assert(bv[u + i] != bv[v + i]);
    auto &h1 = hs[(u + i) % C], &h2 = hs[(v + i) % C];
    int b1 = (u + i) / C, b2 = (v + i) / C;
    // for (int j = 0; j < L; ++j){
    //     char U = h1(b1 + j, b1 + j + 1).val, V = h2(b2 + j, b2 + j + 1).val;
    //     if (U < V)return 1;
    //     if (U > V)return 0;
    // }
    // assert(0);
    int ll = 0, rr = L - 1, j = L; // first ne char
    while (ll <= rr){
        int mm = ll + rr >> 1;
        if (h1(b1, b1 + mm + 1) == h2(b2, b2 + mm + 1))ll = mm + 1;
        else j = mm, rr = mm - 1;
    }
    assert(j < L);
    return a[u + i + j * C] < a[v + i + j * C];
}
int main(){
    // freopen(".in", "r", stdin);
    ios::sync_with_stdio(0), cin.tie(0), cout.tie(0);
    cin >> k >> t; t %= k;
    for (int i = 0; i < 1 << k; ++i)cin >> a[i];
    copy_n(a, 1 << k, a + (1 << k));
    L = 1 << t, C = 1 << (k - t);
    string bf; bf.resize(L + L);
    for (int i = 0; i < C; ++i){
        for (int j = 0; j < L + L; ++j)bf[j] = a[i + j * C];
        hs[i].build(bf.begin(), bf.end());
        for (int j = 0; j < L; ++j)bv[i + j * C] = bv[i + j * C + (1 << k)] = hs[i](j, j + L);
    }
    ss.build(basic_string<hsh::underlying>(reinterpret_cast<hsh::underlying*>(bv), 1 << k + 1));
    int p = 0;
    for (int x = 1; x < 1 << k; ++x)
        if (lss(x, p))p = x;
    for (int i = 0; i < C; ++i)
        for (int j = 0; j < L; ++j)cout << a[p + i + j * C];
    return 0;
}