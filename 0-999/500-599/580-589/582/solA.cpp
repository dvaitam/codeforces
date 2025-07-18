/* Copyright 2015 AcrossTheSky */
#include <iostream>
#include <cstdio>
#include <utility>
#include <cassert>
#include <map>
#include <vector>
#include <deque>
#include <queue>
#include <stack>
#include <set>
#include <cstring>
#include <cstdlib>
#include <cctype>
#include <sstream>
#include <fstream>
#include <string>
#include <cmath>
#include <algorithm>
#define REP(i, a, b) for (int i = (a); i <= (b); ++i)
#define PER(i, a, b) for (int i = (a); i >= (b); --i)
#define RVC(i, c) fot (int i = 0; i < (c).size(); ++i)
#define RED(k, u) for (int k = head[(u)]; k; k = edge[k].next)
#define lowbit(x) ((x) & (-(x)))
#define CL(x, v) memset(x, v, sizeof x)
#define MP std::make_pair
#define PB push_back
#define FR first
#define SC second
#define rank rankk
#define next nextt
#define link linkk
#define index indexx
#define abs(x) ((x) > 0 ? (x) : (-(x)))
using namespace std;
typedef long long LL;
typedef pair<int, int> PII;

template<class T> inline
bool getmin(T *a, const T &b) {
    if (b < *a) {
        *a = b;
        return true;
    }
    return false;
}

template<class T> inline
bool getmax(T *a, const T &b) {
    if (b > *a) {
        *a = b;
        return true;
    }
    return false;
}

template<class T> inline
void read(T *a) {
    char c;
    while (isspace(c = getchar())) {}
    bool flag = 0;
    if (c == '-') flag = 1, *a = 0;
    else
        *a = c - 48;
    while (isdigit(c = getchar())) *a = *a * 10 + c - 48;
    if (flag) *a = -*a;
}
const int mo = 1000000007;
template<class T>
T pow(T a, T b, int c = mo) {
    T res = 1;
    for (T i = 1; i <= b; i <<= 1, a = 1LL * a * a % c) if (b & i) res = 1LL * res * a % c;
    return res;
}
/*======================= TEMPLATE =======================*/
const int N = 100000;
map<int, int> S;
vector<int> ans;
vector<int> tmp;
vector<int> x;
int n;
int main() {
    cin >> n;
    REP(i, 1, n * n) {
        int x;
        read(&x);
        S[x]++;
    }
    while (ans.size() < n) {
        for (int i = 0; i < tmp.size(); ++i) {
            for (int j = 0; j < ans.size(); ++j) 
                S[__gcd(tmp[i], ans[j])]--;
            for (int j = 0; j < x.size(); ++j) 
                S[__gcd(tmp[i], x[j])]--;
        }
        tmp.clear();
        x = ans;
        int xx = 0;
        for (map<int, int>::iterator it = S.begin(); it != S.end(); ++it) 
            if (it -> SC) xx = it -> first;
        tmp.push_back(xx); ans.push_back(xx);
    }
    REP(i, 0, n - 1) printf("%d ", ans[i]);
    cout << endl;
    return 0;
}