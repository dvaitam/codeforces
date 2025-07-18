#include <algorithm>

#include <array>

#include <cmath>

#include <cstring>

#include <chrono>

#include <iostream>

#include <iomanip> 

#include <queue>

#include <map>

#include <numeric>

#include <random>

#include <set>

#include <stack>

#include <string>

#include <utility>

#include <vector>



#define endl '\n'

#define all(x) x.begin(), x.end()

#define rall(x) x.rbegin(), x.rend()

#define fora(a, b, c) for (int i = a; i < b; i += c)

#define forr(a, b, c) for (int i = a; i >= b; i -= c)

#define yes cout << "YES" << endl;

#define no cout << "NO" << endl;

#define cout0 cout << 0 << endl;

#define cout1 cout << -1 << endl;



typedef long long ll;

#define pi std::pair<int, int>

#define pll std::pair<ll, ll>

typedef std::vector<bool> vb;

typedef std::vector<vb> vvb;

typedef std::vector<int> vi;

typedef std::vector<vi> vvi;

typedef std::vector<pi> vpi;

typedef std::vector<ll> vll;

typedef std::vector<pll> vpll;

typedef std::vector<vll> vvll;

typedef std::vector<std::string> vs;

typedef std::vector<vs> vvs;



using namespace std;



int main() {

    std::ios::sync_with_stdio(false);

    std::cin.tie(nullptr);

    cout << setprecision(15);

    int t;

    cin >> t;

    while (t--) {

        int n, k, i;

        cin >> n >> k;

        k %= 4;

        if (!k) { no continue; }

        bool ok = 1;

        yes if (k != 2) for (i = 0; i < n; i += 2) cout << i + 1 << " " << i + 2 << endl;  

        else for (i = 0; i < n; i += 2) { if (ok) cout << i + 2 << " " << i + 1 << endl; else cout << i + 1 << " " << i + 2 << endl; ok = !ok; }

    }

}