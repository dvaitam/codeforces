#include <bits/stdc++.h>

 

using i64 = long long;

 

constexpr int N = 4E5;

 

struct SegmentTree {

    const int n;

    std::vector<int> sum, pre;

    SegmentTree(int n) : n(n), sum(4 << std::__lg(n)), pre(4 << std::__lg(n)) {}

    SegmentTree(std::vector<int> init) : SegmentTree(init.size()) {

        std::function<void(int, int, int)> build = [&](int p, int l, int r) {

            if (r - l == 1) {

                sum[p] = init[l];

                pre[p] = init[l];

                return;

            }

            int m = (l + r) / 2;

            build(2 * p, l, m);

            build(2 * p + 1, m, r);

            pull(p);

        };

        build(1, 0, n);

    }

    void pull(int p) {

        sum[p] = sum[2 * p] + sum[2 * p + 1];

        pre[p] = std::min(pre[2 * p], sum[2 * p] + pre[2 * p + 1]);

    }

    void modify(int p, int l, int r, int x, const int &v) {

        if (r - l == 1) {

            sum[p] = pre[p] = v;

            return;

        }

        int m = (l + r) / 2;

        if (x < m) {

            modify(2 * p, l, m, x, v);

        } else {

            modify(2 * p + 1, m, r, x, v);

        }

        pull(p);

    }

    void modify(int p, const int &v) {

        modify(1, 0, n, p, v);

    }

    int rangeQuery(int p, int l, int r, int x, int y) {

        if (l >= y || r <= x) {

            return 0;

        }

        if (l >= x && r <= y) {

            return sum[p];

        }

        int m = (l + r) / 2;

        return rangeQuery(2 * p, l, m, x, y) + rangeQuery(2 * p + 1, m, r, x, y);

    }

    int rangeQuery(int l, int r) {

        return rangeQuery(1, 0, n, l, r);

    }

    int find(int p, int l, int r, int x, int y, int &v) {

        if (l >= y || r <= x) {

            return -1;

        }

        if (l >= x && r <= y && pre[p] > v) {

            v -= sum[p];

            return -1;

        }

        if (r - l == 1) {

            assert(sum[p] <= v);

            return r;

        }

        int m = (l + r) / 2;

        int res = find(2 * p, l, m, x, y, v);

        if (res == -1) {

            res = find(2 * p + 1, m, r, x, y, v);

        }

        return res;

    }

    int find(int l, int r, int v) {

        return find(1, 0, n, l, r, v);

    }

};

 

int main() {

    std::ios::sync_with_stdio(false);

    std::cin.tie(nullptr);

    

    int n;

    std::cin >> n;

    

    std::vector<int> t(N);

    

    std::map<int, int> ranges;

    

    SegmentTree seg(std::vector(N, -1));

    

    for (int i = 0; i < n; i++) {

        int a;

        std::cin >> a;

        a--;

        

        std::vector<std::pair<int, int>> X, Y;

        

        if (t[a]) {

            t[a] = 0;

            seg.modify(a, -1);

            auto it = std::prev(ranges.upper_bound(a));

            auto [l, r] = *it;

            

            X.emplace_back(l, r);

            

            ranges.erase(it);

            

            int j = seg.find(l, r, -2);

            

            if (l < j - 2) {

                ranges[l] = j - 2;

                Y.emplace_back(l, j - 2);

            }

            if (j < r) {

                ranges[j] = r;

                Y.emplace_back(j, r);

            }

        } else {

            t[a] = 1;

            seg.modify(a, 1);

            

            auto it = ranges.upper_bound(a);

            

            if (it != ranges.begin() && a < std::prev(it)->second) {

                it--;

                

                auto [l, r] = *it;

                

                X.emplace_back(l, r);

                it = ranges.erase(it);

                

                r += 2;

                

                if (it != ranges.end() && it->first == r) {

                    X.push_back(*it);

                    r = it->second;

                    ranges.erase(it);

                }

                

                Y.emplace_back(l, r);

                ranges[l] = r;

            } else {

                if (a % 2 == 1) {

                    a--;

                }

                int l = a, r = a + 2;

                

                if (it != ranges.begin() && std::prev(it)->second == l) {

                    it--;

                    l = it->first;

                    X.push_back(*it);

                    it = ranges.erase(it);

                }

                

                if (it != ranges.end() && it->first == r) {

                    X.push_back(*it);

                    r = it->second;

                    ranges.erase(it);

                }

                

                Y.emplace_back(l, r);

                ranges[l] = r;

            }

        }

        

        std::cout << X.size() << "\n";

        for (auto [l, r] : X) {

            std::cout << l + 1 << " " << r << "\n";

        }

        std::cout << Y.size() << "\n";

        for (auto [l, r] : Y) {

            std::cout << l + 1 << " " << r << "\n";

        }

    }

    

    return 0;

}