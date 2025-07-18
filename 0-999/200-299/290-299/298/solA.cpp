#include <algorithm>
#include <iostream>
#include <cstdlib>
#include <cstdio>
#include <string>
#include <vector>
#include <queue>
#include <cmath>
#include <ctime>
#include <stack>
#include <set>
#include <map>

#define f first
#define s second
#define sz size()
#define len length()
#define mp make_pair
#define pb push_back
#define ins insert
#define fi(x) freopen(x,"r",stdin)
#define fo(x) freopen(x,"w",stdout)
#define pi pair < int, int >
#define ll long long
#define sit(x) set <x>::iterator
#define mit(x,y) map <x,y>::iterator
#define quit return 0
#define fname ""

using namespace std;

const double PI = acos(-1.0);
const double EPS = (1e-8);
const int MAXN = (int)(2 * 1e5 + 1e3);
const ll INF = (ll)(1e18);

string s;
int n, x, y, c1, c2;

int main() {
//	fi(fname".in");
//	fo(fname".out");
		cin >> n;
      	cin >> s;

      	for (int i=0; i<n; ++i)
      		if (s[i] == 'R') c1++; else
      		if (s[i] == 'L') c2++;

		if (c1 == 0) {
			
			for (int i=n-1; i>=0; --i)
				if (s[i] == 'L') {
					cout << i + 1 << ' ';
					break;
				}

			for (int i=0; i<n; ++i)
				if (s[i] == 'L') {
					cout << i;
					return 0;
				}
			return 0;
		} else if (c2 == 0) {

			for (int i=0; i<n; ++i)
				if (s[i] == 'R') {
					cout << i + 1 << ' ';
					break;
				}			

			for (int i=n-1; i>=0; --i)
				if (s[i] == 'R') {
					cout << i + 2;
					return 0;
				}
			return 0;
		}



		for (int i=0; i<n; ++i)
			if (s[i] == 'R') {
				cout << i + 1 << ' ';
				break;
			}

		for (int i=0; i<n-1; ++i)
			if (s[i] == 'R' && s[i+1] == 'L') {
				cout << i + 1;
				return 0;
			}



	return 0;
}