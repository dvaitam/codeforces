#include<iostream>
using namespace std;
int main()
{
  int n;
  cin>>n;

  long int arr[n];

  for(int i=0;i<n;i++)
  {
     cin>>arr[i];
     if(arr[i]%2==0) arr[i]=arr[i]-1;
  }

  for(int i=0;i<n;i++)
  {cout<<arr[i]<<" ";}
}