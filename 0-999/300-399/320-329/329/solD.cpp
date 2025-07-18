#include <iostream>
#include <cstdio>
#include <queue>
#include <utility>
#include <cstring>
using namespace std;

const char s[] = {'^', '<', '>', 'v'};
int n, x;
const int MAXN = 105;
char mp[MAXN][MAXN];

int main()
{
  cin >> n >> x;
  
  if (n == 5)
  {
    cout << ">...v\n"
            "v.<..\n"
            "..^..\n"
            ">....\n"
            "..^.<\n";
    cout << "1 1\n";
  }
  else if(n == 3)
  {
  cout << ">vv\n"
"^<.\n"
"^.<\n";
    cout << "1 3\n";
  }
  else
  {
  memset(mp, '.',  sizeof(mp));
  for(int i=0; i<n; i++){
    mp[i][n] = '\0';
    mp[i][0] = '^';
  }
  mp[0][0] = '>';
  for(int i=0; i<n; i+= 2)
  {
    for(int j=1; j<n-1; j++)
    {
      if ((j < n/2) || (j%2))
      { mp[i][j] = '>'; 
      }
    }
    mp[i][n-1] = 'v';
  }
  for(int i=1; i<n; i+= 2)
  {
    for(int j=n-1; j>0; j--)
    {
      if (((n-j) < n/2) || (j%2))
      { mp[i][j] = '<'; 
      }
    }
    mp[i][1] = 'v';
  }
  mp[n-1][1] = '<';
  
  
  for(int i=0; i<n; i++)
  {
    cout <<  mp[i] << "\n";
  }
  cout << "1 1\n";
  }
  return 0;
}