#include<cstdio>
int main() {
  int n,k;
  scanf("%d%d",&n,&k);

  if(n!=k) {
    for(int i=1;i<=n;i++)
      if(i<n-k) printf("%d ",i+1);
      else if(i>n-k) printf("%d ",i);
      else printf("1 ");
  }
  else printf("-1");
}