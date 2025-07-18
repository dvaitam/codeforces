#include <bits/stdc++.h> 

using namespace std;

char a[100];
int n, m;

int main() {
	for (char c = 'a'; c <= 'z'; ++c) {
		a[++m] = c;	
	} 
	for (char c = 'A'; c <= 'Z'; ++c) {
		a[++m] = c;	
	} 
	m = 0;
  int n; 
  cin >> n;
  if (n == 1) {
  	cout << "a a" << endl;
  	return 0;
  }
  string s = "";
  while (n > 1) {
  	--n;
  	// cout << n << endl;
  	if (n % 2 == 0) {
  		++m;
  		s += a[m];
  		s += a[m];
  		n /= 2;
  	} else {
  		++m;
  		s += a[m];
  	}
  }
  string t = "";
  for (int i = 1; i <= m; ++i) {
  	s += a[i];
  	t += a[i];
  }
  cout << s << " " << t << endl;
  return 0;
}