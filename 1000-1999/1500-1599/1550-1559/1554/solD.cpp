#include <iostream>

using namespace std;



int main() {

	// your code goes here

	int t ;

	cin>>t ;

	while(t--)

     {

     	int n ;

     	cin>>n ;

     	string s = "";

     	if(n==1)

        {

        	cout<<"a"<<endl;

        

        }

        else

     	{for(int i = 0 ; i <(n/2) ; i++)

     	s.push_back('a');

     	s.push_back('b');

     	for(int i = 0 ; i < n/2 - 1 ; i++)

     	s.push_back('a');

     	

     if(n%2!=0)

     s.push_back('c');

     	

     	cout<<s<<endl;

     	}

     }

	return 0;

}