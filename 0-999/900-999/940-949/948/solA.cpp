#include <bits/stdc++.h>

using namespace std;

const int MAXN = 500 + 7;
char str[MAXN][MAXN];

const int dx[] = {0, 0, -1, 1};
const int dy[] = {-1, 1, 0, 0};

int main () {
  int n, m;
  scanf("%d %d", &n, &m);
  for(int i = 0; i < n; ++i) {
    scanf("%s", str + i);
  }
  bool ans = true;
  for(int i = 0; i < n; ++i) {
    for(int j = 0; j < m; ++j) {
      if(str[i][j] == 'W') {
        for(int k = 0; k < 4; ++k) {
          int xx = i + dx[k];
          int yy = j + dy[k];
          if(xx >=0 && xx < n && yy >= 0 && yy < m) {
            if(str[xx][yy] == 'S') {
              ans = false;
            } else if(str[xx][yy] == '.') {
              str[xx][yy] = 'D';
            }
          }
        }
      }
    }
  }
  if(ans == false) {  
    puts("NO");
  } else {
    puts("YES");
    for(int i = 0; i < n; ++i) {
      printf("%s\n", str[i]);
    }
  }
  return 0;
}