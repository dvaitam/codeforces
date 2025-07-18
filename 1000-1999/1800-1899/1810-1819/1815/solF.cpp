/*--------------------------------------------------------------------------------------------------------------------------*/
/* Author = Chauhan Abhishek */
/* Codechef Id = https://www.codechef.com/users/abhishek_9036 */
/* LeetCode Id = https://leetcode.com/abhishekchauhan9036/ */
/* GeeksForGeeks Id = https://auth.geeksforgeeks.org/user/abhishekchauhan9036/profile */
/* GitHub Id = https://github.com/AbhishekChauhan9036 */
/* Language = C++ */
/* Address = Ballia,UP,INDIA (221716)*/
/*--------------------------------------------------------------------------------------------------------------------------*/
#undef _GLIBCXX_DEBUG
#include <bits/stdc++.h>
using namespace std;
#ifdef LOCAL
#include "algo/debug.h"
#else
#define debug(...) 42
#endif
template <typename T> T inverse(T a, T m) {
  T u = 0, v = 1;
  while (a != 0) {
    T t = m / a;
    m -= t * a;
    swap(a, m);
    u -= t * v;
    swap(u, v);
  }
  assert(m == 1);
  return u;
}

template <typename T> class Modular {
public:
  using Type = typename decay<decltype(T::value)>::type;

  constexpr Modular() : value() {}
  template <typename U> Modular(const U &x) { value = normalize(x); }

  template <typename U> static Type normalize(const U &x) {
    Type v;
    if (-mod() <= x && x < mod())
      v = static_cast<Type>(x);
    else
      v = static_cast<Type>(x % mod());
    if (v < 0)
      v += mod();
    return v;
  }

  const Type &operator()() const { return value; }
  template <typename U> explicit operator U() const {
    return static_cast<U>(value);
  }
  constexpr static Type mod() { return T::value; }

  Modular &operator+=(const Modular &other) {
    if ((value += other.value) >= mod())
      value -= mod();
    return *this;
  }
  Modular &operator-=(const Modular &other) {
    if ((value -= other.value) < 0)
      value += mod();
    return *this;
  }
  template <typename U> Modular &operator+=(const U &other) {
    return *this += Modular(other);
  }
  template <typename U> Modular &operator-=(const U &other) {
    return *this -= Modular(other);
  }
  Modular &operator++() { return *this += 1; }
  Modular &operator--() { return *this -= 1; }
  Modular operator++(int) {
    Modular result(*this);
    *this += 1;
    return result;
  }
  Modular operator--(int) {
    Modular result(*this);
    *this -= 1;
    return result;
  }
  Modular operator-() const { return Modular(-value); }

  template <typename U = T>
  typename enable_if<is_same<typename Modular<U>::Type, int>::value,
                     Modular>::type &
  operator*=(const Modular &rhs) {
#ifdef _WIN32
    uint64_t x = static_cast<int64_t>(value) * static_cast<int64_t>(rhs.value);
    uint32_t xh = static_cast<uint32_t>(x >> 32), xl = static_cast<uint32_t>(x),
             d, m;
    asm("divl %4; \n\t" : "=a"(d), "=d"(m) : "d"(xh), "a"(xl), "r"(mod()));
    value = m;
#else
    value = normalize(static_cast<int64_t>(value) *
                      static_cast<int64_t>(rhs.value));
#endif
    return *this;
  }
  template <typename U = T>
  typename enable_if<is_same<typename Modular<U>::Type, long long>::value,
                     Modular>::type &
  operator*=(const Modular &rhs) {
    value = (long long)((int)value * rhs.value % mod());
    return *this;
  }
  template <typename U = T>
  typename enable_if<!is_integral<typename Modular<U>::Type>::value,
                     Modular>::type &
  operator*=(const Modular &rhs) {
    value = normalize(value * rhs.value);
    return *this;
  }

  Modular &operator/=(const Modular &other) {
    return *this *= Modular(inverse(other.value, mod()));
  }

  friend const Type &abs(const Modular &x) { return x.value; }

  template <typename U>
  friend bool operator==(const Modular<U> &lhs, const Modular<U> &rhs);

  template <typename U>
  friend bool operator<(const Modular<U> &lhs, const Modular<U> &rhs);

  template <typename V, typename U>
  friend V &operator>>(V &stream, Modular<U> &number);

private:
  Type value;
};

template <typename T>
bool operator==(const Modular<T> &lhs, const Modular<T> &rhs) {
  return lhs.value == rhs.value;
}
template <typename T, typename U>
bool operator==(const Modular<T> &lhs, U rhs) {
  return lhs == Modular<T>(rhs);
}
template <typename T, typename U>
bool operator==(U lhs, const Modular<T> &rhs) {
  return Modular<T>(lhs) == rhs;
}

template <typename T>
bool operator!=(const Modular<T> &lhs, const Modular<T> &rhs) {
  return !(lhs == rhs);
}
template <typename T, typename U>
bool operator!=(const Modular<T> &lhs, U rhs) {
  return !(lhs == rhs);
}
template <typename T, typename U>
bool operator!=(U lhs, const Modular<T> &rhs) {
  return !(lhs == rhs);
}

template <typename T>
bool operator<(const Modular<T> &lhs, const Modular<T> &rhs) {
  return lhs.value < rhs.value;
}

template <typename T>
Modular<T> operator+(const Modular<T> &lhs, const Modular<T> &rhs) {
  return Modular<T>(lhs) += rhs;
}
template <typename T, typename U>
Modular<T> operator+(const Modular<T> &lhs, U rhs) {
  return Modular<T>(lhs) += rhs;
}
template <typename T, typename U>
Modular<T> operator+(U lhs, const Modular<T> &rhs) {
  return Modular<T>(lhs) += rhs;
}

template <typename T>
Modular<T> operator-(const Modular<T> &lhs, const Modular<T> &rhs) {
  return Modular<T>(lhs) -= rhs;
}
template <typename T, typename U>
Modular<T> operator-(const Modular<T> &lhs, U rhs) {
  return Modular<T>(lhs) -= rhs;
}
template <typename T, typename U>
Modular<T> operator-(U lhs, const Modular<T> &rhs) {
  return Modular<T>(lhs) -= rhs;
}

template <typename T>
Modular<T> operator*(const Modular<T> &lhs, const Modular<T> &rhs) {
  return Modular<T>(lhs) *= rhs;
}
template <typename T, typename U>
Modular<T> operator*(const Modular<T> &lhs, U rhs) {
  return Modular<T>(lhs) *= rhs;
}
template <typename T, typename U>
Modular<T> operator*(U lhs, const Modular<T> &rhs) {
  return Modular<T>(lhs) *= rhs;
}

template <typename T>
Modular<T> operator/(const Modular<T> &lhs, const Modular<T> &rhs) {
  return Modular<T>(lhs) /= rhs;
}
template <typename T, typename U>
Modular<T> operator/(const Modular<T> &lhs, U rhs) {
  return Modular<T>(lhs) /= rhs;
}
template <typename T, typename U>
Modular<T> operator/(U lhs, const Modular<T> &rhs) {
  return Modular<T>(lhs) /= rhs;
}

template <typename T, typename U>
Modular<T> power(const Modular<T> &a, const U &b) {
  assert(b >= 0);
  Modular<T> x = a, res = 1;
  U p = b;
  while (p > 0) {
    if (p & 1)
      res *= x;
    x *= x;
    p >>= 1;
  }
  return res;
}

template <typename T> bool IsZero(const Modular<T> &number) {
  return number() == 0;
}

template <typename T> string to_string(const Modular<T> &number) {
  return to_string(number());
}

template <typename U, typename T>
U &operator<<(U &stream, const Modular<T> &number) {
  return stream << number();
}

template <typename U, typename T> U &operator>>(U &stream, Modular<T> &number) {
  typename common_type<typename Modular<T>::Type, long long>::type x;
  stream >> x;
  number.value = Modular<T>::normalize(x);
  return stream;
}

using ModType = long long;

struct VarMod {
  static ModType value;
};
ModType VarMod::value;
ModType &md = VarMod::value;
using Mint = Modular<VarMod>;

namespace factorizer {

template <typename T> struct FactorizerVarMod {
  static T value;
};
template <typename T> T FactorizerVarMod<T>::value;

template <typename T> bool IsPrime(T n, const vector<T> &bases) {
  if (n < 2) {
    return false;
  }
  vector<T> small_primes = {2, 3, 5, 7, 11, 13, 17, 19, 23, 29};
  for (const T &x : small_primes) {
    if (n % x == 0) {
      return n == x;
    }
  }
  if (n < 31 * 31) {
    return true;
  }
  int s = 0;
  T d = n - 1;
  while ((d & 1) == 0) {
    d >>= 1;
    s++;
  }
  FactorizerVarMod<T>::value = n;
  for (const T &a : bases) {
    if (a % n == 0) {
      continue;
    }
    Modular<FactorizerVarMod<T>> cur = a;
    cur = power(cur, d);
    if (cur == 1) {
      continue;
    }
    bool witness = true;
    for (int r = 0; r < s; r++) {
      if (cur == n - 1) {
        witness = false;
        break;
      }
      cur *= cur;
    }
    if (witness) {
      return false;
    }
  }
  return true;
}

bool IsPrime(long long n) {
  return IsPrime(n, {2, 325, 9375, 28178, 450775, 9780504, 1795265022});
}

bool IsPrime(int32_t n) { return IsPrime(n, {2, 7, 61}); }

// Jab ulong long version ka need ho tab
// ...........................................
/*
bool IsPrime(ulong long n) {
  if (n < 2) {
    return false;
  }
  vector<uint32_t> small_primes = {2, 3, 5, 7, 11, 13, 17, 19, 23, 29};
  for (uint32_t x : small_primes) {
    if (n == x) {
      return true;
    }
    if (n % x == 0) {
      return false;
    }
  }
  if (n < 31 * 31) {
    return true;
  }
  uint32_t s = __builtin_ctzll(n - 1);
  ulong long d = (n - 1) >> s;
  function<bool(ulong long)> witness = [&n, &s, &d](ulong long a) {
    ulong long cur = 1, p = d;
    while (p > 0) {
      if (p & 1) {
        cur = (__uint128_t) cur * a % n;
      }
      a = (__uint128_t) a * a % n;
      p >>= 1;
    }
    if (cur == 1) {
      return false;
    }
    for (uint32_t r = 0; r < s; r++) {
      if (cur == n - 1) {
        return false;
      }
      cur = (__uint128_t) cur * cur % n;
    }
    return true;
  };
  vector<ulong long> bases_64bit = {2, 325, 9375, 28178, 450775, 9780504,
1795265022}; for (ulong long a : bases_64bit) { if (a % n == 0) { return true;
    }
    if (witness(a)) {
      return false;
    }
  }
  return true;
}
*/

vector<int> least = {0, 1};
vector<int> primes;
int precalculated = 1;

void RunLinearSieve(int n) {
  n = max(n, 1);
  least.assign(n + 1, 0);
  primes.clear();
  for (int i = 2; i <= n; i++) {
    if (least[i] == 0) {
      least[i] = i;
      primes.push_back(i);
    }
    for (int x : primes) {
      if (x > least[i] || i * x > n) {
        break;
      }
      least[i * x] = x;
    }
  }
  precalculated = n;
}

void RunSlowSieve(int n) {
  n = max(n, 1);
  least.assign(n + 1, 0);
  for (int i = 2; i * i <= n; i++) {
    if (least[i] == 0) {
      for (int j = i * i; j <= n; j += i) {
        if (least[j] == 0) {
          least[j] = i;
        }
      }
    }
  }
  primes.clear();
  for (int i = 2; i <= n; i++) {
    if (least[i] == 0) {
      least[i] = i;
      primes.push_back(i);
    }
  }
  precalculated = n;
}

void RunSieve(int n) { RunLinearSieve(n); }

template <typename T>
vector<pair<T, int>> MergeFactors(const vector<pair<T, int>> &a,
                                  const vector<pair<T, int>> &b) {
  vector<pair<T, int>> c;
  int i = 0;
  int j = 0;
  while (i < (int)a.size() || j < (int)b.size()) {
    if (i < (int)a.size() && j < (int)b.size() && a[i].first == b[j].first) {
      c.emplace_back(a[i].first, a[i].second + b[j].second);
      ++i;
      ++j;
      continue;
    }
    if (j == (int)b.size() || (i < (int)a.size() && a[i].first < b[j].first)) {
      c.push_back(a[i++]);
    } else {
      c.push_back(b[j++]);
    }
  }
  return c;
}

template <typename T> vector<pair<T, int>> RhoC(const T &n, const T &c) {
  if (n <= 1) {
    return {};
  }
  if ((n & 1) == 0) {
    return MergeFactors({{2, 1}}, RhoC(n / 2, c));
  }
  if (IsPrime(n)) {
    return {{n, 1}};
  }
  FactorizerVarMod<T>::value = n;
  Modular<FactorizerVarMod<T>> x = 2;
  Modular<FactorizerVarMod<T>> saved = 2;
  T power = 1;
  T lam = 1;
  while (true) {
    x = x * x + c;
    T g = __gcd((x - saved)(), n);
    if (g != 1) {
      return MergeFactors(RhoC(g, c + 1), RhoC(n / g, c + 1));
    }
    if (power == lam) {
      saved = x;
      power <<= 1;
      lam = 0;
    }
    lam++;
  }
  return {};
}

template <typename T> vector<pair<T, int>> Rho(const T &n) {
  return RhoC(n, static_cast<T>(1));
}

template <typename T> vector<pair<T, int>> Factorize(T x) {
  if (x <= 1) {
    return {};
  }
  if (x <= precalculated) {
    vector<pair<T, int>> ret;
    while (x > 1) {
      if (!ret.empty() && ret.back().first == least[x]) {
        ret.back().second++;
      } else {
        ret.emplace_back(least[x], 1);
      }
      x /= least[x];
    }
    return ret;
  }
  if (x <= static_cast<long long>(precalculated) * precalculated) {
    vector<pair<T, int>> ret;
    if (!IsPrime(x)) {
      for (T i : primes) {
        T t = x / i;
        if (i > t) {
          break;
        }
        if (x == t * i) {
          int cnt = 0;
          while (x % i == 0) {
            x /= i;
            cnt++;
          }
          ret.emplace_back(i, cnt);
          if (IsPrime(x)) {
            break;
          }
        }
      }
    }
    if (x > 1) {
      ret.emplace_back(x, 1);
    }
    return ret;
  }
  return Rho(x);
}

template <typename T>
vector<T> BuildDivisorsFromFactors(const vector<pair<T, int>> &factors) {
  vector<T> divisors = {1};
  for (auto &p : factors) {
    int sz = (int)divisors.size();
    for (int i = 0; i < sz; i++) {
      T cur = divisors[i];
      for (int j = 0; j < p.second; j++) {
        cur *= p.first;
        divisors.push_back(cur);
      }
    }
  }
  sort(divisors.begin(), divisors.end());
  return divisors;
}

} // namespace factorizer
/*-----------------------------------------------------------------------------------------------------------------------------------------*/

/********************************************************************************************************************************************/

template <typename T> struct PrimitiveVarMod {
  static T value;
};
template <typename T> T PrimitiveVarMod<T>::value;

template <typename T, class F>
T GetPrimitiveRoot(const T &modulo, const F &factorize) {
  if (modulo <= 0) {
    return -1;
  }
  if (modulo == 1 || modulo == 2 || modulo == 4) {
    return modulo - 1;
  }
  vector<pair<T, int>> modulo_factors = factorize(modulo);
  if (modulo_factors[0].first == 2 &&
      (modulo_factors[0].second != 1 || modulo_factors.size() != 2)) {
    return -1;
  }
  if (modulo_factors[0].first != 2 && modulo_factors.size() != 1) {
    return -1;
  }
  set<T> phi_factors;
  T phi = modulo;
  for (auto &d : modulo_factors) {
    phi = phi / d.first * (d.first - 1);
    if (d.second > 1) {
      phi_factors.insert(d.first);
    }
    for (auto &e : factorize(d.first - 1)) {
      phi_factors.insert(e.first);
    }
  }
  PrimitiveVarMod<T>::value = modulo;
  Modular<PrimitiveVarMod<T>> gen = 2;
  while (gen != 0) {
    if (power(gen, phi) != 1) {
      continue;
    }
    bool ok = true;
    for (auto &p : phi_factors) {
      if (power(gen, phi / p) == 1) {
        ok = false;
        break;
      }
    }
    if (ok) {
      return gen();
    }
    gen++;
  }
  assert(false);
  return -1;
}

template <typename T> T GetPrimitiveRoot(const T &modulo) {
  return GetPrimitiveRoot(modulo, factorizer::Factorize<T>);
}

/*
void _yes() { cout << "YES" << endl; }
void _no() { cout << "NO" << endl; }
#define endl '\n'
#define ed '\n'

#define int long long
#define rev(v) reverse(v.begin(), v.end())
typedef long long ll;
const int MAX = 1e5 + 10;
const int mod = 1e9 + 7;
#define loop(n) for(long long i=0;i<n;i++)
#define rloop(n) for(long long i=n-1;i>=0;i--)

#define all(x) (x).begin(), (x).end()
ll gcd(ll a, ll b) {
  if (a == 0)
    return b;
  return gcd(b % a, a);
}


void BhagwatGeeta() {
  cout<<"yuy";
}*/


#include <bits/stdc++.h>
using i64 = long long;
std::mt19937 rng;
 
void solve() {
    int n, m;
    std::cin >> n >> m;
    
    std::vector<int> a(n);
    for (int i = 0; i < n; i++) {
        std::cin >> a[i];
    }
    
    std::vector<std::vector<std::pair<int, int>>> adj(n);
    std::vector<std::array<int, 3>> ver(m);
    std::vector f(m, std::array<int, 3>{});
    
    for (int i = 0; i < m; i++) {
        int u, v, w;
        std::cin >> u >> v >> w;
        
        u--, v--, w--;
        ver[i] = {u, v, w};
        adj[u].push_back({i, 0});
        adj[v].push_back({i, 1});
        adj[w].push_back({i, 2});
        
        std::array<int, 3> b{0, 1, 2};
        std::sort(b.begin(), b.end(), [&](int x, int y) {
            return ver[i][x] < ver[i][y];
        });
        
        for (int j = 0; j < 3; j++) {
            a[ver[i][b[j]]] += j + 3;
            f[i][b[j]] += j + 3;
        }
    }
    
    std::priority_queue<std::pair<int, int>> h;
    for (int i = 0; i < n; i++) {
        h.emplace(-a[i], i);
    }
    
    while (!h.empty()) {
        auto [v, x] = h.top();
        h.pop();
        
        if (a[x] != -v) {
            continue;
        }
        
        for (auto [j, t] : adj[x]) {
            int y = ver[j][(t + 1) % 3];
            int z = ver[j][(t + 2) % 3];
            if (a[x] == a[y] && a[x] == a[z]) {
                a[y] += 1;
                f[j][(t + 1) % 3] += 1;
                h.emplace(-a[y], y);
                a[z] += 1;
                f[j][(t + 2) % 3] += 1;
                h.emplace(-a[z], z);
            }
            else if (a[x] == a[y]) {
                a[y] += 2;
                f[j][(t + 1) % 3] += 2;
                h.emplace(-a[y], y);
            }
            else if (a[x] == a[z]) {
                a[z] += 2;
                f[j][(t + 2) % 3] += 2;
                h.emplace(-a[z], z);
            }
        }
    }
    
    for (int i = 0; i < m; i++) {
        auto [u, v, w] = ver[i];
        assert(a[u] != a[v]);
        assert(a[u] != a[w]);
        assert(a[v] != a[w]);
        
        int a = (f[i][0] + f[i][1] - f[i][2]) / 2;
        int b = (f[i][1] + f[i][2] - f[i][0]) / 2;
        int c = (f[i][2] + f[i][0] - f[i][1]) / 2;
        
        assert(1 <= a && a <= 4);
        assert(1 <= b && b <= 4);
        assert(1 <= c && c <= 4);
        std::cout << a << " " << b << " " << c << "\n";
    }
}
 
int main() {
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);
    
    int t;
    std::cin >> t;
    
    while (t--) {
        solve();
    }
    
    return 0;
}
/*
signed main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
  int kk = 1;
  cin >> kk;
  while (kk--) {
    BhagwatGeeta();
  }
  return 0;
}
*/
/*
aaanoka
qdlftcaasd
aahvjnaefv
kakjsdcedg
#include <bits/stdc++.h>
using namespace std;
#define int long long
int mod=1000000007;
int calc(int i,int j){
    int first;
  if(i%2==0){
     first=i/2;
     first= ((first%mod)*((i-1)%mod))%mod;
  }
  else{
      first=(i-1)/2;
     first= ((first%mod)*((i)%mod))%mod;
  }
  int second;
   if(j%2==0){
    second=j/2;
     second= ((second%mod)*((j+1)%mod))%mod;
  }
  else{
     second=(j+1)/2;
     second= ((second%mod)*((j)%mod))%mod;
  }
  return (second-first+mod)%mod;
}
signed main() {
int t;
cin>>t;
while(t--){
    int n,m,k;
    cin>>n>>m>>k;
    map<int,int> row;
    map<int,int> col;
    while(k--)
    {
        int q;
        cin>>q;
        int x,v;
        cin>>x>>v;
        if(q==0){
            if(row.find(x)!=row.end())row[x]=(row[x]*v)%mod;
            else row[x]=v;
        }else{
            if(col.find(x)!=col.end())col[x]=(col[x]*v)%mod;
            else col[x]=v;
        }
    }
    int var=0;
    int cons=0;
        int sum= calc(1,n*m);

    for(auto itr=col.begin();itr!=col.end();itr++){
        var=var+itr->second-1;
        cons+=((itr->second-1)*(itr->first))%mod;
        var=var%mod;
        cons=cons%mod;
        int val=((calc(0,n-1)*m)%mod + (n*(itr->first))%mod)%mod;
        sum =sum + (val*(itr->second-1))%mod;
        sum=(sum+mod)%mod;
    }

    var=(var*m)%mod;


    for(auto itr=row.begin();itr!=row.end();itr++){
        int rowno=itr->first;
        int rowc=itr->second;
        int sumr=calc((rowno-1)*m+1,rowno*m)+((var*(rowno-1))%mod)+cons;
        sumr=sumr%mod;
        int add=((sumr%mod)*(rowc-1)%mod)%mod;
        sum=sum+add;
        sum=(sum+mod)%mod;

    }
    cout<<sum<<"\n";
}


 int n=4;
  int arr[4];
  for(int i=0;i<n;i++)
  {
      cin>>arr[i];
  }
  if((arr[0]+1)%arr[1]==(arr[2]+1)%arr[3])
  {
    cout<<1<<endl;

  }else{
    int e=arr[1]*arr[3];
     e=e/__gcd(arr[1],arr[3]);
    e-=arr[0]%arr[1];
    cout<<e<<endl;
  }
}*/