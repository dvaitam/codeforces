#include <cstdio>
#include <algorithm>
#include <cstring>
#include <vector>

using namespace std;
int N, T;
char S1[100500], S2[100500];
char ans[100500];
char note(char x, char y) {
  for (char z = 'a'; z <= 'z'; z++)
    if (x != z && y != z)
      return z;
}
int main() {
  scanf("%d%d", &N, &T);
  scanf("%s%s", S1, S2);
  int eq = 0, dif = 0;
  for (int i = 0; i < N; i++)
    if (S1[i] == S2[i]) eq++;
    else dif++;
  
  int sd = 0;
  while (T > dif || 2*T < dif) {
    if (sd < eq) {
      sd++;
      T--;
    } else {
      printf("-1\n");
      return 0;
    }
  }
  if (T > dif || 2*T < dif) {
    printf("-1\n");
    return 0;
  }
  
  int db = 2*T - dif, c = 0;
  //d1 = 2*(dif - T)
  for (int i = 0; i < N; i++) {
    if (S1[i] != S2[i]) {
      if (c < db) ans[i] = note(S1[i], S2[i]);
      else ans[i] = (c%2) ? S1[i] : S2[i];
      c++;
    }
  }
  c = 0;
  for (int i = 0; i < N; i++) {
    if (S1[i] == S2[i]) {
      if (c < sd) ans[i] = 'a' + (S1[i] - 'a' + 1) % 26;
      else ans[i] = S1[i];
      c++;
    }
  }
    
  ans[N] = 0;
  printf("%s\n", ans);
  
  return 0;
}