#include <bits/stdc++.h>
using namespace std;
using LL = long long;
using point_t = long double; //全局数据类型，可修改为 long long 等

const point_t eps = 1e-8;
const long double PI = acosl(-1);
//const long double PI = numbers::pi_v<long double>;

// 点与向量
template <typename T>
struct point{
    T x, y;

    bool operator==(const point &a) const { return (abs(x - a.x) <= eps && abs(y - a.y) <= eps); }
    bool operator<(const point &a) const{
        if (abs(x - a.x) <= eps)
            return y < a.y - eps;
        return x < a.x - eps;
    }
    bool operator>(const point &a) const { return !(*this < a || *this == a); }
    point operator+(const point &a) const { return {x + a.x, y + a.y}; }
    point operator-(const point &a) const { return {x - a.x, y - a.y}; }
    point operator-() const { return {-x, -y}; }
    point operator*(const T k) const { return {k * x, k * y}; }
    point operator/(const T k) const { return {x / k, y / k}; }
    T operator*(const point &a) const { return x * a.x + y * a.y; } // 点积
    T operator^(const point &a) const { return x * a.y - y * a.x; } // 叉积，注意优先级
    int toleft(const point &a) const
    {
        const auto t = (*this) ^ a;
        return (t > eps) - (t < -eps);
    }                                                             // to-left 测试
    T len2() const { return (*this) * (*this); }                  // 向量长度的平方
    T dis2(const point &a) const { return (a - (*this)).len2(); } // 两点距离的平方

    // 涉及浮点数
    long double len() const { return sqrtl(len2()); }                                                                      // 向量长度
    long double dis(const point &a) const { return sqrtl(dis2(a)); }                                                       // 两点距离
    long double ang(const point &a) const { return acosl(max(-1.0l, min(1.0l, ((*this) * a) / (len() * a.len())))); }      // 向量夹角
    point rot(const long double rad) const { return {x * cosl(rad) - y * sinl(rad), x * sinl(rad) + y * cosl(rad)}; }          // 逆时针旋转（给定角度）
    point rot(const long double cosr, const long double sinr) const { return {x * cosr - y * sinr, x * sinr + y * cosr}; } // 逆时针旋转（给定角度的正弦与余弦）
};

using Point = point<point_t>;

int main(){

#ifdef LOCAL
    freopen("data.in", "r", stdin);
    freopen("data.out", "w", stdout);
#endif

    cin.tie(0);
    cout.tie(0);
    ios::sync_with_stdio(0);

    Point p1, p2;
    cin >> p1.x >> p1.y >> p2.x >> p2.y;
    Point mid = (p1 + p2) / 2;
    Point v = p1 - p2;
    double w = v.len() / 2;
    swap(v.x, v.y);
    v.x = -v.x;
    v = v / v.len();
    int n;
    cin >> n;
    vector<pair<long double, int> > q;
    q.push_back({0, 0});

    auto get = [&](Point p, long double val){

        auto d = [&](long double len){
            Point q = mid + v * len;
            return p.dis(q) - sqrtl(w * w + len * len);
        };

        long double l = -1e12, r = 1e12;
        bool fl = d(l) > val;
        for(int i = 0; i < 100; i++){
            long double mid = (l + r) / 2;
            bool f = d(mid) > val;
            if (f == fl) l = mid;
            else r = mid;
        }
        return r;
    };

    for(int i = 0; i < n; i++){
        int x, y, R;
        cin >> x >> y >> R;
        Point p{x, y};
        long double l = get(p, R);
        long double r = get(p, -R);
        if (l > r) swap(l, r);
        q.push_back({l, 1});
        q.push_back({r, -1});
    }
    sort(q.begin(), q.end());
    int s = 0;
    long double mn = 1e12;
    for(auto [x, y] : q){
        if (s == 0) mn = min(mn, abs(x));
        s += y;
        if (s == 0) mn = min(mn, abs(x));
    }
    cout << fixed << setprecision(20) << sqrtl(mn * mn + w * w) << '\n';

}