//==================================================================
//Author        : Nguyen Trung Hieu - vuondenthanhcong11@yahoo.com
//Problem Name  :
//Discription   :
//Reference from:
//==================================================================

// -------------------- Khai bao thu vien --------------------
#include <set>
#include <bitset>
#include <queue>
#include <deque>
#include <set>
#include <bitset>
#include <queue>
#include <deque>
#include <stack>
#include <sstream>
#include <iostream>
#include <iomanip>

#include <cstdio>
#include <cstdlib>
#include <cmath>
#include <ctime>
#include <cstring>
#include <string>
#include <cassert>

#include <vector>
#include <list>
#include <map>
#include <algorithm>
#include <functional>
#include <numeric>
#include <utility>

using namespace std;

// -------------------- Khai bao cac container --------------------
typedef vector <int> vi;
typedef vector <vi> vvi;
typedef pair <int, int> ii;
typedef vector <ii> vii;
typedef vector <string> vs;

typedef long long int64; //NOTES:int64
typedef unsigned long long uint64; //NOTES:uint64
typedef unsigned uint;

const double pi = acos(-1.0); //NOTES:pi
const double eps = 1e-11; //NOTES:eps
const int MAXI = 0x7fffffff;
const int dx[] = {1, 0, -1, 0};
const int dy[] = {0, 1, 0, -1};
const char dz[] = "SENW";

// -------------------- Dinh nghia lai cac phep toan --------------------
#define forn(i,a,b)     for (int i=(a),_b=(b); i<_b; i++)
#define rforn(i,b,a)    for (int i=(b)-1,_a=(a); i>=_a; i--)
#define reset(a,b)      memset((a),(b),sizeof(a))

#define endline         putchar('\n')
#define space           putchar(' ')
#define EXIT(x)         {cout << x;return 0;}

#define fi              first
#define se              second
#define pb              push_back
#define all(x)          (x).begin(),(x).end()
#define mp(X,Y)         make_pair(X,Y)
#define foreach(i, c)   for(typeof((c).begin()) i = (c).begin(); i != (c).end(); i++)
#define present(c, x)   ((c).find(x) != (c).end())

#define two(x)          (1 << (x))
#define two64(x)        (((int64)(1)) << (x))
#define contain(s, x)   (((s) & two(x)) != 0)
#define contain64(s,x)  (((s) & two64(x)) != 0)

// =========================== Bat dau template ===========================

// -------------------- Cac thao tac nhap - xuat, debug, v.v... --------------------

int in() {
    int x = 0, c;
    for (; (uint)((c = getchar()) - '0') >= 10; ) { if (c == '-') return -in(); if (!~c) throw ~0; }
    do { x = (x << 3) + (x << 1) + (c - '0'); } while ((uint)((c = getchar()) - '0') < 10);
    return x;
}

int64 in64() {
    int64 x = 0, c;
    for (; (uint)((c = getchar()) - '0') >= 10; ) { if (c == '-') return -in(); if (!~c) throw ~0; }
    do { x = (x << 3) + (x << 1) + (c - '0'); } while ((uint)((c = getchar()) - '0') < 10);
    return x;
}

void pin(int n) { char buf[33]; int i = 30; if (n < 0) putchar('-'), n = -n; do { buf[i--] = '0'+ n % 10; n /= 10; } while (n); while (i < 30) putchar(buf[++i]); }

void pin64(int64 n) { char buf[55]; int i = 50; if (n < 0) putchar('-'), n = -n; do { buf[i--] = '0'+ n % 10; n /= 10; } while (n); while (i < 50) putchar(buf[++i]); }

// -------------------- Cac ham co ban --------------------

template<class T> inline T checkMod(T n, T m) {
    return (n % m + m) % m;
}

template<class T> inline T multiplyMod(T a, T b, T m) {
    return (T) ((((int64) (a)*(int64) (b) % (int64) (m))+(int64) (m)) % (int64) (m));
}

template<class T> inline T powerMod(T p, int e, T m)
{
    if (e == 0)return 1 % m;
    else if (e % 2 == 0) {
        T t = powerMod(p, e / 2, m);
        return multiplyMod(t, t, m);
    } else return multiplyMod(powerMod(p, e - 1, m), p, m);
}

template<class T> inline bool isPrimeNumber(T n)
{
    if (n <= 1)return false;
    for (T i = 2; i * i <= n; i++) if (n % i == 0) return false;
    return true;
}

template<class T> inline T sqr(T x) {
    return x*x;
}

template<class T> inline T lowbit(T n) {
    return (n^(n - 1))&n;
}

template<class T> inline int countbit(T n) {
    return (n == 0) ? 0 : (1 + countbit(n & (n - 1)));
}

// -------------------- Cac ham ve toan --------------------

int64 C(int m,int n)
{
    if (n > m) return 0;
    int64 re = 1;
    for(int i = 1; i <= n; i++) re = re * (m - i + 1) / i;
    return re;
}

template<class T> inline T gcd(T a, T b)
{
    if (a < 0)return gcd(-a, b);
    if (b < 0)return gcd(a, -b);
    return (b == 0) ? a : gcd(b, a % b);
}

template<class T> inline T lcm(T a, T b)
{
    if (a < 0)return lcm(-a, b);
    if (b < 0)return lcm(a, -b);
    return a * (b / gcd(a, b));
}

template<class T> inline T euclide(T a, T b, T &x, T &y)
{
    if (a < 0) {
        T d = euclide(-a, b, x, y);
        x = -x;
        return d;
    }
    if (b < 0) {
        T d = euclide(a, -b, x, y);
        y = -y;
        return d;
    }
    if (b == 0) {
        x = 1;
        y = 0;
        return a;
    } else {
        T d = euclide(b, a % b, x, y);
        T t = x;
        x = y;
        y = t - (a / b) * y;
        return d;
    }
}

template<class T> inline vector<pair<T, int> > factorize(T n)
{
    vector<pair<T, int> > R;
    for (T i = 2; n > 1;) {
        if (n % i == 0) {
            int C = 0;
            for (; n % i == 0; C++, n /= i);
            R.push_back(make_pair(i, C));
        }
        i++;
        if (i > n / i) i = n;
    }
    if (n > 1) R.push_back(make_pair(n, 1));
    return R;
}

template<class T> inline T eularFunction(T n)
{
    vector<pair<T, int> > R = factorize(n);
    T r = n;
    for (int i = 0; i < R.size(); i++)r = r / R[i].first * (R[i].first - 1);
    return r;
}

// -------------------- Cac ham ve ma tran --------------------

const int MaxMatrixSize = 40;

template<class T> inline void showMatrix(int n, T A[MaxMatrixSize][MaxMatrixSize])
{
    for (int i = 0; i < n; i++) {
        for (int j = 0; j < n; j++)cout << A[i][j];
        cout << endl;
    }
}

template<class T> inline void identityMatrix(int n, T A[MaxMatrixSize][MaxMatrixSize])
{
    for (int i = 0; i < n; i++) for (int j = 0; j < n; j++) A[i][j] = (i == j) ? 1 : 0;
}

template<class T> inline void addMatrix(int n, T C[MaxMatrixSize][MaxMatrixSize], T A[MaxMatrixSize][MaxMatrixSize], T B[MaxMatrixSize][MaxMatrixSize])
{
    for (int i = 0; i < n; i++) for (int j = 0; j < n; j++) C[i][j] = A[i][j] + B[i][j];
}

template<class T> inline void subMatrix(int n, T C[MaxMatrixSize][MaxMatrixSize], T A[MaxMatrixSize][MaxMatrixSize], T B[MaxMatrixSize][MaxMatrixSize])
{
    for (int i = 0; i < n; i++) for (int j = 0; j < n; j++) C[i][j] = A[i][j] - B[i][j];
}

template<class T> inline void mulMatrix(int n, T C[MaxMatrixSize][MaxMatrixSize], T _A[MaxMatrixSize][MaxMatrixSize], T _B[MaxMatrixSize][MaxMatrixSize])
{
    T A[MaxMatrixSize][MaxMatrixSize], B[MaxMatrixSize][MaxMatrixSize];
    for (int i = 0; i < n; i++) for (int j = 0; j < n; j++) A[i][j] = _A[i][j], B[i][j] = _B[i][j], C[i][j] = 0;
    for (int i = 0; i < n; i++) for (int j = 0; j < n; j++) for (int k = 0; k < n; k++) C[i][j] += A[i][k] * B[k][j];
}

template<class T> inline void addModMatrix(int n, T m, T C[MaxMatrixSize][MaxMatrixSize], T A[MaxMatrixSize][MaxMatrixSize], T B[MaxMatrixSize][MaxMatrixSize])
{
    for (int i = 0; i < n; i++) for (int j = 0; j < n; j++) C[i][j] = checkMod(A[i][j] + B[i][j], m);
}

template<class T> inline void subModMatrix(int n, T m, T C[MaxMatrixSize][MaxMatrixSize], T A[MaxMatrixSize][MaxMatrixSize], T B[MaxMatrixSize][MaxMatrixSize])
{
    for (int i = 0; i < n; i++) for (int j = 0; j < n; j++) C[i][j] = checkMod(A[i][j] - B[i][j], m);
}

template<class T> inline void mulModMatrix(int n, T m, T C[MaxMatrixSize][MaxMatrixSize], T _A[MaxMatrixSize][MaxMatrixSize], T _B[MaxMatrixSize][MaxMatrixSize])
{
    T A[MaxMatrixSize][MaxMatrixSize], B[MaxMatrixSize][MaxMatrixSize];
    for (int i = 0; i < n; i++) for (int j = 0; j < n; j++) A[i][j] = _A[i][j], B[i][j] = _B[i][j], C[i][j] = 0;
    for (int i = 0; i < n; i++) for (int j = 0; j < n; j++) for (int k = 0; k < n; k++) C[i][j] = (C[i][j] + multiplyMod(A[i][k], B[k][j], m)) % m;
}
// -------------------- Cac ham ve hinh hoc --------------------

double dist(ii a, ii b) {
    return sqrt(sqr((double)(a.fi) - (double)(b.fi)) + sqr((double)(a.se) - (double)(b.se)));
}

double distR(ii a, ii b) {
    return sqr((double)(a.fi) - (double)(b.fi)) + sqr((double)(a.se) - (double)(b.se));
}

int cross(ii c, ii a, ii b) {
    return (a.fi - c.fi)*(b.se - c.se)-(b.fi - c.fi)*(a.se - c.se);
}

int crossOper(ii c, ii a, ii b)
{
    double t = (a.fi - c.fi)*(b.se - c.se)-(b.fi - c.fi)*(a.se - c.se);
    if (fabs(t) <= eps) return 0;
    return (t < 0) ? -1 : 1;
}

bool isIntersect(ii a, ii b, ii c, ii d)
{
    return crossOper(a, b, c) * crossOper(a, b, d) < 0 && crossOper(c, d, a) * crossOper(c, d, b) < 0;
}

bool isMiddle(double s, double m, double t) {
    return fabs(s - m) <= eps | fabs(t - m) <= eps | (s < m) != (t < m);
}

// -------------------- Cac ham ve chuyen doi giua cac kieu --------------------

bool isUpperCase(char c) {
    return c >= 'A' && c <= 'Z';
}

bool isLowerCase(char c) {
    return c >= 'a' && c <= 'z';
}

bool isLetter(char c) {
    return c >= 'A' && c <= 'Z' | c >= 'a' && c <= 'z';
}

bool isDigit(char c) {
    return c >= '0' && c <= '9';
}

char toLowerCase(char c) {
    return (isUpperCase(c)) ? (c + 32) : c;
}

char toUpperCase(char c) {
    return (isLowerCase(c)) ? (c - 32) : c;
}

template<class T> string toString(T n) {
    ostringstream ost;
    ost << n;
    ost.flush();
    return ost.str();
}

int toInt(string s) {
    int r = 0;
    istringstream sin(s);
    sin >> r;
    return r;
}

int64 toInt64(string s) {
    int64 r = 0;
    istringstream sin(s);
    sin >> r;
    return r;
}

double toDouble(string s) {
    double r = 0;
    istringstream sin(s);
    sin >> r;
    return r;
}

template<class T> void stoa(string s, int &n, T A[]) {
    n = 0;
    istringstream sin(s);
    for (T v; sin >> v; A[n++] = v);
}

template<class T> void atos(int n, T A[], string &s) {
    ostringstream sout;
    for (int i = 0; i < n; i++) {
        if (i > 0)sout << ' ';
        sout << A[i];
    }
    s = sout.str();
}

template<class T> void atov(int n, T A[], vector<T> &vi) {
    vi.clear();
    for (int i = 0; i < n; i++) vi.push_back(A[i]);
}

template<class T> void vtoa(vector<T> vi, int &n, T A[]) {
    n = vi.size();
    for (int i = 0; i < n; i++)A[i] = vi[i];
}

template<class T> void stov(string s, vector<T> &vi) {
    vi.clear();
    istringstream sin(s);
    for (T v; sin >> v; vi.push_back(v));
}

template<class T> void vtos(vector<T> vi, string &s) {
    ostringstream sout;
    for (int i = 0; i < vi.size(); i++) {
        if (i > 0)sout << ' ';
        sout << vi[i];
    }
    s = sout.str();
}

// -------------------- Cac ham ve phan so --------------------

template<class T> struct Fraction {
    T a, b;
    Fraction(T a = 0, T b = 1);
    string toString();
};

template<class T> Fraction<T>::Fraction(T a, T b) {
    T d = gcd(a, b);
    a /= d;
    b /= d;
    if (b < 0) a = -a, b = -b;
    this->a = a;
    this->b = b;
}

template<class T> string Fraction<T>::toString() {
    ostringstream sout;
    sout << a << "/" << b;
    return sout.str();
}

template<class T> Fraction<T> operator+(Fraction<T> p, Fraction<T> q) {
    return Fraction<T > (p.a * q.b + q.a * p.b, p.b * q.b);
}

template<class T> Fraction<T> operator-(Fraction<T> p, Fraction<T> q) {
    return Fraction<T > (p.a * q.b - q.a * p.b, p.b * q.b);
}

template<class T> Fraction<T> operator*(Fraction<T> p, Fraction<T> q) {
    return Fraction<T > (p.a * q.a, p.b * q.b);
}

template<class T> Fraction<T> operator/(Fraction<T> p, Fraction<T> q) {
    return Fraction<T > (p.a * q.b, p.b * q.a);
}
// =========================== Ket thuc template ===========================

#define N 700005

int n;
ii A[N];

int main(){
    freopen ("input.txt", "r", stdin);
    freopen ("output.txt", "w", stdout);
    int cnt[5001];
    reset(cnt, 0);
    n = 2 * in();
    forn (i, 0, n) {
        int val = in();
        A[i] = mp(val, i + 1);
        cnt[val]++;
    }

    forn (i, 0, 5001)
        if (cnt[i] % 2) {
            pin(-1);
            return 0;
        }

    sort(A, A + n);
    for (int i = 0; i < n; i += 2) {
        pin(A[i].se), space, pin(A[i + 1].se), endline;
    }
}