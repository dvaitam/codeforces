#include<iostream>
#include<algorithm>

using namespace std;

int main(void) {
	int n;
	int arr[100001] = { 0 };
	cin >> n;

	int max_fir = 0;
	int min_fir = 0;
	int max_sec = 0;
	int min_sec = 0;


	for (int i = 0; i < n; i++) {
		int temp = 0;
		scanf("%d", &temp);
		if (min_fir == 0) {
			min_fir = temp;
		}
		else if(min_fir!=0&&min_sec==0){
			if (temp < min_fir) {
				min_sec = min_fir;
				min_fir = temp;
			}
			else {
				min_sec = temp;
			}
		}
		else {
			if (temp < min_fir) {
				min_sec = min_fir;
				min_fir = temp;
			}
			else if (temp == min_fir) {
				min_sec = temp;
			}
			else {
				if (temp < min_sec) {
					min_sec = temp;
				}
			}
		}
		if (max_fir == 0) {
			max_fir = temp;
		}
		else if (max_fir != 0 && max_sec == 0) {
			if (temp > max_fir) {
				max_sec = max_fir;
				max_fir = temp;
			}
			else {
				max_sec = temp;
			}
		}
		else {
			if (temp > max_fir) {
				max_sec = max_fir;
				max_fir = temp;
			}
			else if (temp == max_fir) {
				max_sec = temp;
			}
			else {
				if (temp > max_sec) {
					max_sec = temp;
				}
			}
		}
	}
	int result = 0;
	result = min(max_fir - min_fir, max_fir - min_sec);
	result = min(result, max_sec - min_fir);
	cout << result;
}