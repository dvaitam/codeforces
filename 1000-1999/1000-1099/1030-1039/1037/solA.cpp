#include "bits/stdc++.h"

using namespace std;
using ll = long long;
using ull = unsigned long long;

typedef pair<ll, ll> P;
typedef pair<ll, P> PP;

const ll MOD = 1e9 + 7;
const ll INF = 9e18;
const double eps = -1e10;
const int di[4] = { 1,0,-1,0 }, dj[4] = { 0,1,0,-1 };

#define ALL(x) (x).begin(),(x).end()
#define pb push_back
#define eb emplace_back
#define fr first
#define sc second

ll n;

int main() {
	cin >> n;
	for (int i = 0;i < 32;i++) {
		if (n < (1 << i)) {
			cout << i << endl;
			return 0;
		}
	}
	return 0;
}