#include <algorithm>
#include <array>
#include <bitset>
#include <cassert>
#include <chrono>
#include <cmath>
#include <cstring>
#include <functional>
#include <iomanip>
#include <iostream>
#include <map>
#include <numeric>
#include <queue>
#include <random>
#include <set>
#include <vector>
using namespace std;

// http://www.open-std.org/jtc1/sc22/wg21/docs/papers/2016/p0200r0.html
template<class Fun> class y_combinator_result {
    Fun fun_;
public:
    template<class T> explicit y_combinator_result(T &&fun): fun_(std::forward<T>(fun)) {}
    template<class ...Args> decltype(auto) operator()(Args &&...args) { return fun_(std::ref(*this), std::forward<Args>(args)...); }
};
template<class Fun> decltype(auto) y_combinator(Fun &&fun) { return y_combinator_result<std::decay_t<Fun>>(std::forward<Fun>(fun)); }


template<typename A, typename B> ostream& operator<<(ostream &os, const pair<A, B> &p) { return os << '(' << p.first << ", " << p.second << ')'; }
template<typename T_container, typename T = typename enable_if<!is_same<T_container, string>::value, typename T_container::value_type>::type> ostream& operator<<(ostream &os, const T_container &v) { os << '{'; string sep; for (const T &x : v) os << sep << x, sep = ", "; return os << '}'; }

void dbg_out() { cerr << endl; }
template<typename Head, typename... Tail> void dbg_out(Head H, Tail... T) { cerr << ' ' << H; dbg_out(T...); }
#ifdef NEAL_DEBUG
#define dbg(...) cerr << "(" << #__VA_ARGS__ << "):", dbg_out(__VA_ARGS__)
#else
#define dbg(...)
#endif

template<typename T_vector>
void output_vector(const T_vector &v, bool add_one = false, int start = -1, int end = -1) {
    if (start < 0) start = 0;
    if (end < 0) end = int(v.size());

    for (int i = start; i < end; i++)
        cout << v[i] + (add_one ? 1 : 0) << (i < end - 1 ? ' ' : '\n');
}


bool perfect_square(int64_t x) {
    int64_t s = int64_t(sqrt(x + 0.5L));
    return s * s == x;
}

int main() {
    ios::sync_with_stdio(false);
#ifndef NEAL_DEBUG
    cin.tie(nullptr);
#endif

    int64_t N;
    cin >> N;
    vector<bool> included(N + 1, true);

    auto standard = [&]() -> void {
        for (int64_t n = N; n % 4 != 0; n--)
            included[n] = false;

        included[(N - N % 4) / 2] = false;
    };

    if (N % 4 == 3) {
        if (perfect_square(N + 1))
            included[N / 2 + 1] = included[N] = false;
        else if (perfect_square(2 * (N / 2) * (N / 2 - 1)))
            included[N / 2 - 2] = included[N] = false;
        else if (perfect_square((N / 2 - 1) * N))
            included[N / 2 - 2] = included[N - 2] = false;
        else
            included[2] = included[N / 2] = included[N] = false;
    } else if (N % 4 == 2) {
        if (perfect_square(2 * (N / 2 + 1)))
            included[N / 2 + 1] = false;
        else if (perfect_square(2 * (N / 2) * (N / 2 - 1)))
            included[N / 2 - 2] = false;
        else
            included[2] = included[N / 2] = false;
    } else {
        standard();
    }

    included[1] = true;
    vector<int64_t> remaining, removed;

    for (int64_t i = 1; i <= N; i++)
        if (included[i])
            remaining.push_back(i);
        else
            removed.push_back(i);

    cout << remaining.size() << '\n';
    output_vector(remaining);
    // cout << removed << endl;
}