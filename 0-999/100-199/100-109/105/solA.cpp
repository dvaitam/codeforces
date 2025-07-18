#include <stdio.h>
#include <iostream>
#include <math.h>
#include <string>
#include <string.h>
#include <stack>
#include <queue>
#include <map>
#include <set>
#include <deque>
#include <vector>
#include <stdlib.h>
#include <memory.h>
#include <algorithm>
#include <sstream>
using namespace std;
typedef long long ll;
map<string,int> q;
int n,m; 
double k;

int main() {
    //freopen("in","r",stdin);
    //freopen("out","w",stdout);
    
    cin >> n >> m >> k;
    
    while (n--) {
	string s; int c;
	cin >> s >> c;
	c = (int) (k * (double) c + 0.00000001);
	if (c >= 100) {
	    q[s] = c;
	}
    }
    while (m--) {
	string s;
	cin >> s;
	if (q.find(s) == q.end()) q[s] = 0;
    }
    
    cout << q.size() << endl;
    for (map<string,int>::iterator i=q.begin();i!=q.end();i++)
	cout << i->first << " " << i->second << endl;
    
    return 0;
    
}