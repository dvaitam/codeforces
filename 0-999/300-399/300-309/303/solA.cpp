#include <stdio.h>      
#include <ctype.h>
#include <math.h>

#include <iomanip>
#include <iostream>
#include <sstream>
#include <utility>
#include <algorithm>
#include <cassert>
#include <string>
#include <vector>
#include <queue>
#include <set>
#include <map>
using namespace std;

typedef vector<int> VI;
typedef long long LL;
typedef pair<int,int> PII;
typedef double LD;

/* CHECKLIST 
 * 1) long longs */

const int DBG = 0, INF = int(1e9);

int main() {
   ios_base::sync_with_stdio(0);
   cout.setf(ios::fixed);

   int n;
   cin >> n;

   if (n % 2 == 0)
      cout << -1 << endl;
   else {
      for (int i = 0; i < n; ++i)
         cout << i << " ";
      cout << endl;
      for (int i = 1; i < n; ++i)
         cout << i << " ";
      cout << 0 << endl;
      for (int i = 1; i < n; i += 2)
         cout << i << " ";
      for (int i = 0; i < n; i += 2)
         cout << i << " ";
      cout << endl;
   }

   return 0;
}