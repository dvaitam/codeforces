#include <string>
#include <vector>
#include <sstream>
#include <iostream>
#include <algorithm>

#define pb push_back
#define rep(i,s,t) for(int i = (s), _t = (t); i < _t; ++i)
#define debug(x) cerr << #x << " = " << (x) << endl

using namespace std;

typedef vector<int> ivec;
// EOT

// ABCDEFGHIJKLMNOPQRSTUVWXYZ
string S = "AEFHIKLMNTVWXYZ";

int main()
{
	rep (i, 21, 51)
	{
		bool a = (min(i, 25) + i) % (2 + i % 3) > 0;
		cout << a << endl;
	}
	return 0;
}