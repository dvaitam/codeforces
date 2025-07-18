#include <bits/stdc++.h>




namespace zawa {

using i16 = std::int16_t;
using i32 = std::int32_t;
using i64 = std::int64_t;
using i128 = __int128_t;

using u8 = std::uint8_t;
using u16 = std::uint16_t;
using u32 = std::uint32_t;
using u64 = std::uint64_t;

using usize = std::size_t;

} // namespace zawa


namespace zawa {

void SetFastIO() {
    std::cin.tie(nullptr)->sync_with_stdio(false);
}

void SetPrecision(u32 dig) {
    std::cout << std::fixed << std::setprecision(dig);
}

} // namespace zawa
using namespace zawa;

int main() {
    SetFastIO();

    int t; std::cin >> t;
    while (t--) {
        int n; std::cin >> n;
        std::vector<long long> a(n);
        for (auto& x : a) std::cin >> x;
        std::vector<long long> sum(n + 1);
        for (int i{} ; i < n ; i++) sum[i + 1] = sum[i] + a[i];
        // for (auto x : sum) std::cout << x << ' ';
        // std::cout << std::endl;
        std::vector<long long> mins{sum};
        // for (auto x : mins) std::cout << x << ' ';
        // std::cout << std::endl;
        for (int i{n - 1} ; i >= 0 ; i--) {
            mins[i] = std::min(mins[i + 1], sum[i]);
        }
        long long max{sum[n]}, ans{(long long)1e18};
        for (int i{} ; i < n ; i++) {
            long long v{sum[n] + std::max(0LL, sum[i] - mins[i])};
            if (max < v) {
                max = v;
                ans = sum[i];
            }
        }
        std::cout << ans << '\n';
    }
}