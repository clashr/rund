#include <unistd.h>
#include <stdio.h>

int main() {
	printf("--beginning of program--\n");
	int counter = 0;
	pid_t pid = fork();
	if (pid == 0) {
		// child process
		int i = 0;
		sleep(5);
		for (; i < 5; ++i) {
			printf("child process: counter=%d\n", ++counter);
		}
	}
	else if (pid > 0) {
		// parent process
		int j = 0;
		for (; j < 5; ++j) {
			printf("parent process: counter=%d\n", ++counter);
		}
		
		// do not wait for child to finish
	}
	else {
		// fork failed
		printf("fork() failed!\n");
		return 1;
	}
	printf("--end of program--\n");
	return 0;
}
