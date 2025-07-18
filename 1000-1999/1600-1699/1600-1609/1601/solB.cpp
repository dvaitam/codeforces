# include <cstdlib>
# include <cstdio>
# include <cassert>
# include <cctype>
# include <cstring>
# include <cmath>
# include <ctime>

# include <algorithm>
# include <functional>
# include <utility>
# include <deque>
# include <queue>
# include <stack>
# include <vector>
# include <set>
# include <unordered_set>
# include <map>
# include <unordered_map>

namespace Main {
  namespace Source {
    typedef   signed int dint; typedef short   signed int hd; typedef long   signed int ld; typedef long long   signed int lld; typedef   signed char cd;
    typedef unsigned int uint; typedef short unsigned int hu; typedef long unsigned int lu; typedef long long unsigned int llu; typedef unsigned char cu;
    typedef long double lf;
    template <typename T> using eq = std::    equal_to <T>;
    template <typename T> using ne = std::not_equal_to <T>;
    template <typename T> using lt = std::less         <T>;
    template <typename T> using le = std::less_equal   <T>;
    template <typename T> using gt = std::greater      <T>;
    template <typename T> using ge = std::greater_equal<T>;
    using std::min; template <typename T> static inline constexpr T& amin(T& a, const T& b) { return a = min(a, b); }
    using std::max; template <typename T> static inline constexpr T& amax(T& a, const T& b) { return a = max(a, b); }
    template <typename T> static inline constexpr const T adif(const T& a, const T& b) { return a < b ? b - a : a - b; }
    template <typename T> static inline constexpr const T mdif(const T& a, const T& b) { return a < b ? T() : a - b; }
    template <typename T> static inline constexpr const T abs(const T& x) { return x < T() ? -x : +x; }
    template <typename T> static inline constexpr const T lowbit(const T& x) { return x bitand (~x + 1); }
    using std::sort; using std::swap;
    namespace IO {
      namespace IO_Buf {
        static char ibuf[1 << 20], obuf[1 << 20]; static const char *il, *ir; static char* op(obuf);
        static inline bool flush_i() { return ir = (il = ibuf) + fread(ibuf, sizeof(char), sizeof ibuf / sizeof(char), stdin), il != ir; }
        static inline void flush_o() { op -= fwrite(obuf, sizeof(char), op - obuf, stdout); } struct IO_Buf_D { compl IO_Buf_D() { flush_o(); } } IO_Buf_D;
        static inline char getchar() { return il == ir && !flush_i() ? EOF : *il++; } static inline char& getchar(char& ch) { return ch = getchar(); }
        static inline void putchar(const char ch) { *op++ = ch; if (op == obuf + sizeof obuf / sizeof(char)) flush_o(); }
        static inline char* gets(char* const); static inline void puts(const char*);
      } using namespace IO_Buf;
      static inline void file_open(const char* const f)
      { static char g[1 << 10]; assert(freopen(strcat(strcpy(g, f), ".in"), "r", stdin)), g[strlen(f)] = '\0', assert(freopen(strcat(g, ".out"), "w", stdout)); }
      static inline char getns() { static char ch; while (isspace(getchar(ch))); return ch; } static inline char& getns(char& ch) { return ch = getchar(); }
      template <typename T> static inline T get() {
        static char ch; hd pn(1); while (not isdigit(getns(ch))) switch (ch) { case '+': pn = +pn; break; case '-': pn = -pn; break; /*default: return T();*/ }
        T r{}; while (r = r * 10 + (ch - '0') * pn, isdigit(getchar(ch))); return r;
      }
      template <typename T> static inline T& get(T& a) { return a = get<T>(); }
      template <typename T> static inline void put(T x, const char sp = 0) {
        if (x < T()) { putchar('-'); } static char t[1 << 10]; hd lg(-1); while (t[++lg] = abs(x % 10) + '0', x = x / 10); while (putchar(t[lg]), lg--);
        if (sp) putchar(sp);
      }
      template <> inline void put(const char ch, const char sp) { putchar(ch); if (sp) putchar(sp); }
      struct iostream {
        template <typename T> inline iostream operator>>(      T& a) { return get(a), iostream(); }
        template <typename T> inline iostream operator<<(const T& x) { return put(x), iostream(); }
      } io;
      namespace IO_Buf {
        static inline char* gets(char* const s)
        { lu len(0); char ch(getns()); while (!isspace(ch) and ch != EOF) s[len++] = ch, getchar(ch); return s[len] = '\0', s; }
        static inline void puts(const char* s) { while (*s) putchar(*s++); putchar('\n'); }
      }
    } using namespace IO;
    namespace Maths {
      template <typename T> static inline constexpr const T gcd(T a, T b) { while (a and b) a < b ? b %= a : a %= b; return a + b; }
    } using namespace Maths;
  } using namespace Source;
  namespace Main { static inline void main(); }
}

signed int main() { Main::Main::main(); return 0; }

namespace Main::Main {
  static constexpr const lu N(300'000);
  static lu n; static lu a[N + 1], b[N + 1]; static lu f[N + 1]; static lu pre[N + 1];
  static inline void output(const lu x) { if (pre[x] <= n) output(pre[x]), put(x, ' '); }
  static inline void main() {
    get(n); for (lu i(1); i <= n; ++i) get(a[i]); for (lu i(1); i <= n; ++i) get(b[i]); memset(f, 0xff, sizeof f), f[n] = 0, pre[n] = n + 1;
    for (lu l(n), r(n); r; r = l) {
      while (l and f[l] == f[r]) --l;
      // fprintf(stderr, "%lu %lu\n", l, r);
      for (lu i(r); i > l; --i) for (lu j(i + b[i] - a[i + b[i]]); !pre[j] and j <= i + b[i]; ++j) f[j] = f[r] + 1, pre[j] = i;
      // for (lu i(0); i <= n; ++i) puts("test"), io << f[i] << ' ' << pre[i] << '\n';
      if (not pre[l]) return puts("-1");
    }
    put(f[0], '\n'), output(0), putchar('\n');
  }
}