#include <bits/stdc++.h>
using namespace std;

#define prime 100003

int main(){
  int N,M;
  scanf("%d %d",&N,&M);
  printf("100003 100003\n");
  for (int i=1;i<N-1;i++){
    printf("%d %d 1\n",i,i+1);
  }
  printf("%d %d %d\n",N-1,N,prime-N+2);
  bool done = false;
  int cnt = M-N+1;
  if (!cnt) return 0;
  for (int i=1;i<N;i++){
    for (int j=i+2;j<=N;j++){
      printf("%d %d 100004\n",i,j);
      cnt--;
      if (!cnt) return 0;
    }
  }
}