#include <stdio.h>
#include <algorithm>

using namespace std;

pair <int,int> a[1001];

int main(){
  freopen("input.txt","r",stdin);
  freopen("output.txt","w",stdout);
  int n,k;
  scanf("%d%d",&n,&k);
  for(int i=0;i<n;i++){
    scanf("%d",&a[i].first);
    a[i].second = i+1;
  }
  sort(a,a+n);
  int sum = 0;
  printf("%d\n",a[n-k]);
  for(int i=n-k;i<n;i++)
    printf("%d ",a[i].second);
  return 0;
}