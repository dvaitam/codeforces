#include <stdio.h>
#include <sstream>
#include <iostream>
#include <fstream>
#include <algorithm>
#include <vector>
#include <set>
#include <map>
#include <stack>
#include <memory.h>
#include <queue>
#include <string>
#include <string.h>
#include <cmath>
#include <utility>
#include <bitset>
#include <time.h>
#include <climits>

#define PI 3.1415926535897932384626433832795
#define sqr(x) ((x)*(x))
#define OUT_RT cerr << (float(clock()) / CLOCKS_PER_SEC) << endl

typedef long long LL;
typedef long double LD;
typedef unsigned long long ULL;

using namespace std;
int a,b,c;
int s;
LD x, y, z;

int main() {
//	freopen("input.txt","r",stdin);
//	freopen("output.txt","w",stdout);

	cin >> s;
	cin >> a >> b >> c;
	if (a == 0 && b == 0 && c == 0) {
		cout << "0 0 0\n";
		return 0;
	}
	if (c == 0) {
		cout << LD(s) * a / (a + b) << " " << LD(s) * b / (a + b) << " " << 0 << endl;
		return 0;
	}

	y = (LD(s) * c * b) / (c * b + c * a + c * c);
	x = (s * a - y * a) / (c + a);
	z = s - x - y;
	cout.precision(20);
	cout << fixed;
	cout << x << " " << y << " " << z << endl;

	return 0;
}