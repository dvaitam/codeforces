#include <bits/stdc++.h>
using namespace std;
 
 
 
int main(int argc, char const *argv[])
{
	// freopen("in", "r", stdin);
	int n;
	cin >> n;
	string s, t;
	cin >> s >> t;
	if(s[0] != s[1] and t[0] != t[1] and s[0] == t[1] and s[1] == t[0]) {
		cout << "YES" << endl;
		string out;
		out += s[0];
		for(char c : {'a', 'b', 'c'}) {
			if(c != s[0] and c != s[1]) {
				out += c;
			}
		}
		out += s[1];
		for(int i = 0; i < n; i++) cout << out[0];
		for(int i = 0; i < n; i++) cout << out[1];
		for(int i = 0; i < n; i++) cout << out[2];
		cout << endl;
		return 0;
	}
	if(s[0] == t[0] and s[1] != t[1] and s[0] != s[1] and t[0] != t[1]) {
		cout << "YES" << endl;
		string out;
		for(char c : {'a', 'b', 'c'}) {
			if(c != s[0]) {
				out += c;
			}
		}
		out += s[0];
		for(int i = 0; i < n; i++) cout << out[0];
		for(int i = 0; i < n; i++) cout << out[1];
		for(int i = 0; i < n; i++) cout << out[2];
		cout << endl;
		return 0;
	}
	if(s[1] == t[1] and s[0] != t[0] and s[0] != s[1] and t[0] != t[1]) {
		cout << "YES" << endl;
		string out;
		out += s[1];
		for(char c : {'a', 'b', 'c'}) {
			if(c != s[1]) {
				out += c;
			}
		}
		for(int i = 0; i < n; i++) cout << out[0];
		for(int i = 0; i < n; i++) cout << out[1];
		for(int i = 0; i < n; i++) cout << out[2];
		cout << endl;
		return 0;
	}
	string p = "abc";
	do {
		int flag = 1;
		string pp = p;
		pp += p[0];
		if(s[0] == pp[0] and s[1] == pp[1]) flag = 0;
		if(s[0] == pp[1] and s[1] == pp[2]) flag = 0;
		if(s[0] == pp[2] and s[1] == pp[3]) flag = 0;
		if(t[0] == pp[0] and t[1] == pp[1]) flag = 0;
		if(t[0] == pp[1] and t[1] == pp[2]) flag = 0;
		if(t[0] == pp[2] and t[1] == pp[3]) flag = 0;
		if(flag) {
			cout << "YES" << endl;
			for(int i = 0; i < n; i++) {
				cout << p;
			} cout << endl;
			// break;
			return 0;
		}
	} while(next_permutation(p.begin(), p.end()));
	cout << "NO" << endl;
	return 0;
	assert(false);
	return 0;
}