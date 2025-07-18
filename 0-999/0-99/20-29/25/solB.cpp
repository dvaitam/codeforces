#include<iostream>
#include<string>
#include<algorithm>
#include<vector>
#include<map>
#include<set>
#include<cstdlib>
#include<cstdio>
#include<cmath>
#include<sstream>
#include<cassert>
#include<queue>
#include<stack>

#define REP(i,b,n) for(int i=b;i<n;i++)
#define rep(i,n)   REP(i,0,n)
#define ALL(C)     (C).begin(),(C).end()
#define pb         push_back
#define mp         make_pair
#define vint       vector<int>
#define FOR(it,o)  for(__typeof((o).begin()) it=(o).begin(); it!=(o).end(); ++it)
#define lli	    long long int
#define ld	    long double

using namespace std;


main(){

  int n;
  cin >> n;
  string c;
  cin >> c;
  rep(i, n){
    if(n - i == 3){
      cout << c[i] << c[i+1] <<c[i+2]<<endl;
      break;
    }
    if(n - i == 2){
      cout << c[i] << c[i+1] << endl;
      break;
    }
    cout << c[i] << c[i+1] << '-';
    i++;
  }
}