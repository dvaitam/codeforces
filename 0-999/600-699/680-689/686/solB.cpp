// In the name of GOD!
// 14 / 11 / 2016
// 1395/ 08 / 24

#include <iostream>
#include <cstring>

using namespace std;

const int N = 100 + 5;

int n, a[N], mark[N];
bool b;

void mk(int cnt)
{
   for(int i = cnt; i < n; i += 2)
      if(a[i] > a[i + 1])
      {
         mark[i] = mark[i + 1] = 1;
         swap(a[i], a[i + 1]);
      }
      else if(a[i] == a[i + 1])
         mark[i] = mark[i + 1] = 2;

   int st = 0;
   for(int i = 0; i <= n + 1; i++)
      if(st && !mark[i])
      {
         cout << st << ' ' << i - 1 << '\n';
         st = 0;
         b = true;
      }
      else if(!st && mark[i] == 1)
         st = i;

   memset(mark, 0, sizeof(mark));
}

int main()
{
   cin >> n;
   for(int i = 1; i <= n; i++)
      cin >> a[i];

   b = true;
   while(b)
   {
      b = false;
      mk(1);
      mk(2);
   }
   return 0;
}