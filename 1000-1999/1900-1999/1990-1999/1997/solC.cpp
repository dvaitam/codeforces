#import<iostream>
int n,x;char s[200001];main(){for(std::cin>>n;std::cin>>n>>s;printf("%d\n",x))for(x=n/2;n--;)x+=(s[n]=='(')*2;}