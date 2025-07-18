///Dieser Code ist meiner geliebten Sofa gewidmet







#include <bits/stdc++.h>



using namespace std;



///#include "add_segment_tree.hpp"

///#include "max_segment_tree.hpp"

///#include "int_mod.hpp" ///somehow doesn't work linking



namespace INT_MOD {

    constexpr int mod = 998244353;

    struct intmod {

        int value;

        intmod(int);

        intmod(long long);

        intmod();

        friend intmod operator+ (const intmod&, const intmod&);

        friend intmod operator- (const intmod&, const intmod&);

        friend intmod operator* (const intmod&, const intmod&);

        friend intmod operator+ (const intmod&, const int&);

        friend intmod operator- (const intmod&, const int&);

        friend intmod operator* (const intmod&, const int&);

        friend intmod& operator += (intmod&, const intmod&);

        friend intmod& operator -= (intmod&, const intmod&);

        friend intmod& operator *= (intmod&, const intmod&);

        friend intmod& operator += (intmod&, const int&);

        friend intmod& operator -= (intmod&, const int&);

        friend intmod& operator *= (intmod&, const int&);

        friend intmod operator- (const intmod&);



        intmod to_power(int power) const;

        intmod reversed() const;

        friend intmod operator/ (const intmod&, const intmod&);

        friend intmod& operator/= (intmod&, const intmod&);

        friend std::ostream& operator<< (std::ostream&, const intmod&);

    };



    void fill_f(int n, std::vector<intmod>& v);

    void fill_p2(int n, std::vector<intmod>& v);

    void fill_p3(int n, std::vector<intmod>& v);

    void fill_pk(int n, int k, std::vector<intmod>& v);

    void fill_r(int n, std::vector<intmod>& v);

}







namespace INT_MOD{



    inline int shorten(int value_) {

        return (value_ >= mod)? (value_ - mod): (value_ < 0? value_ + mod: value_);

    }





    intmod::intmod(int value_) {

        ///value_ in (-mod, 2 * mod)

        value = shorten(value_);

    }



    intmod::intmod(long long value_) {

        value = shorten(value_ % mod);

    }



    intmod::intmod() {

        value = 0;

    }



    std::ostream& operator<< (std::ostream& out, const intmod& c) {

        out << c.value;

        return out;

    }



    intmod operator+ (const intmod& a, const intmod& b) {

        return intmod(a.value + b.value);

    }



    intmod operator- (const intmod& a, const intmod& b) {

        return intmod(a.value - b.value);

    }



    intmod operator* (const intmod& a, const intmod& b) {

        return intmod(a.value * 1LL * b.value);

    }



    intmod operator+ (const intmod& a, const int& b) {

        return intmod(a.value + b);

    }



    intmod operator- (const intmod& a, const int& b) {

        return intmod(a.value - b);

    }



    intmod operator* (const intmod& a, const int& b) {

        return intmod(a.value * 1LL * b);

    }



    intmod& operator+= (intmod& a, const intmod& b) {

        a.value = shorten(a.value + b.value);

        return a;

    }



    intmod& operator-= (intmod& a, const intmod& b) {

        a.value = shorten(a.value - b.value);

        return a;

    }



    intmod& operator*= (intmod& a, const intmod& b) {

        a.value = a.value * 1LL * b.value % mod;

        return a;

    }



    intmod& operator+= (intmod& a, const int& b) {

        a.value = shorten(a.value + b);

        return a;

    }



    intmod& operator-= (intmod& a, const int& b) {

        a.value = shorten(a.value - b);

        return a;

    }



    intmod& operator*= (intmod& a, const int& b) {

        a.value = shorten(a.value * 1LL * b % mod);

        return a;

    }



    intmod intmod::to_power(int power) const {

        power %= mod - 1;

        intmod a = 1;

        intmod b = intmod(value);

        for (int bit = 0; bit < 31; bit++) {

            if (power & (1 << bit))

                a *= b;

            b *= b;

        }



        return a;

    }



    intmod intmod::reversed() const{

        return intmod(value).to_power(mod - 2);

    }



    intmod operator/ (const intmod& a, const intmod& b) {

        return a * b.reversed();

    }



    intmod& operator/= (intmod& a, const intmod& b) {

        intmod c = a * b.reversed();

        a.value = c.value;

        return a;

    }



    intmod operator- (const intmod& a) {

        return intmod() - a;

    }





    void fill_f(int n, std::vector<intmod>& v) {

        v = std::vector<intmod>(n + 1);

        v[0] = intmod(1);

        for (int i = 1; i <= n; i++) {

            v[i] = v[i - 1] * i;

        }

    }



    void fill_p2(int n, std::vector<intmod>& v) {

        v = std::vector<intmod>(n + 1);

        v[0] = intmod(1);

        for (int i = 1; i <= n; i++) {

            v[i] = v[i - 1] * 2;

        }

    }



    void fill_p3(int n, std::vector<intmod>& v) {

        v = std::vector<intmod>(n + 1);

        v[0] = intmod(1);

        for (int i = 1; i <= n; i++) {

            v[i] = v[i - 1] * 3;

        }

    }



    void fill_pk(int n, int k, std::vector<intmod>& v) {

        v = std::vector<intmod>(n + 1);

        v[0] = intmod(1);

        for (int i = 1; i <= n; i++) {

            v[i] = v[i - 1] * k;

        }

    }



    void fill_r(int n, std::vector<intmod>& v) {

        v = std::vector<intmod>(n + 1);

        v[1] = intmod(1);

        for (int i = 2; i <= n; i++) {

            v[i] = -v[mod % i] * (mod / i);

        }

    }



    void fill_rf(int n, std::vector<intmod>& v) {

        fill_r(n, v);

        v[0] = intmod(1);

        for (int i = 2; i <= n; i++) {

            v[i] *= v[i - 1];

        }

    }

}





using namespace INT_MOD;

typedef intmod imt;







struct Space {

    vector<int> space = vector<int>(30, 0);



    void add(int v) {

        if (v == 0) return;

        for (int bit = 29; bit >= 0; bit--) {

            if (((1 << bit) & v) == 0)

                continue;

            if (space[bit] == 0) {

                space[bit] = v;

                return;

            }

            else {

                v ^= space[bit];

            }

        }

    }



    void merge (const Space& pspace) {

        for (auto v: (pspace.space)) {

            add(v);

        }

    }



    int get_max() {

        int ans = 0;

        for (int bit = 29; bit >= 0; bit--) {

            if (space[bit] == 0 || (ans & (1 << bit))) {

                continue;

            }

            else {

                ans ^= space[bit];

            }

        }

        return ans;

    }

};





void print_vv(const vector<vector<bool>>& vv) {

    for (int i = 0; i < vv.size(); i++) {

        for (int j = 0; j < vv[i].size(); j++) {

            cout << vv[i][j] << ' ';

        }

        cout << '\n';

    }

}



void print_v(const vector<int>& v) {

    for (auto e: v) {

        cout << e << ' ';

    }

    cout << '\n';

}



template<class T>

vector<T> read_v (size_t s) {

    vector<T> v(s);

    for (T& t: v) cin >> t;

    return v;

}





void solve() {

    int n, k, x;

    cin >> n >> k >> x;

    vector<vector<int>> groups;

    vector<bool> used(n + 1);

    int xor_sum = 0;

    for (int i = 1; i <= n; i++) {

        xor_sum ^= i;

    }

    if (k % 2 && xor_sum != x) {

        cout << "NO\n";

        return;

    }



    if (k % 2 == 0 && xor_sum != 0) {

        cout << "NO\n";

        return;

    }



    for (int i = 1; i <= n; i++) {

        if (i == x) {

            used[x] = true;

            groups.push_back({x});

            continue;

        }

        if (used[i]) continue;



        if ((x ^ i) <= n) {

            groups.push_back({i, x ^ i});

            used[i] = used[x ^ i] = true;

        }

    }



    if (groups.size() >= k) {

        cout << "YES\n";

        for (int h = k; h < groups.size(); h++) {

            for (auto e: groups[h]) {

                groups[0].push_back(e);

            }

        }



        for (int i = 1; i <= n; i++) {

            if (!used[i]) {

                groups[0].push_back(i);

            }

        }



        for (int y = 0; y < k; y++) {

            cout << groups[y].size() << ' ';

            for (auto e: groups[y]) {

                cout << e << ' ';

            }

            cout << '\n';

        }

        return;

    }



    if (groups.size() == k - 1) {

        cout << "YES\n";

        groups.push_back({});

        for (int i = 1; i <= n; i++) {

            if (!used[i]) {

                groups[k - 1].push_back(i);

            }

        }



        for (int y = 0; y < k; y++) {

            cout << groups[y].size() << ' ';

            for (auto e: groups[y]) {

                cout << e << ' ';

            }

            cout << '\n';

        }

        return;

    }



    cout << "NO\n";

    return;





}



signed main() {

    ios_base::sync_with_stdio(false); cin.tie(0); cout.tie(0);

    int testcases = 1;

    cin >> testcases;

    while(testcases--) solve();

    return 0;

}