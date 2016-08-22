#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <unistd.h>

void test_send_signals(int pid)
{
	if (kill(pid, SIGQUIT) == -1)
		perror("failed to send SIGQUIT");
	if (kill(pid, SIGTERM) == -1)
		perror("failed to send SIGTERM");
	if (kill(pid, SIGKILL) == -1)
		perror("failed to send SIGKILL");
	if (kill(pid, SIGALRM) == -1)
		perror("failed to send SIGALRM");
}

void test_file_io()
{
	FILE *fh;
	char buffer[1024];

	printf("reading '/proc/1/cmdline'\n");
	if (!(fh = fopen("/proc/1/cmdline", "r")))
		perror("failed to open '/proc/1/cmdline' for reading");
	else
		fclose(fh);

	printf("writing '/proc/1/cmdline'\n");
	if (!(fh = fopen("/proc/1/cmdline", "w")))
		perror("failed to open '/proc/1/cmdline' for writing");
	else {
		// memset(buffer, 'd', sizeof(buffer));
		// buffer[sizeof(buffer)] = '\0';
		// fprintf(fh, buffer);
		//
		// We do not actually want to write to init's proc tree.
		fclose(fh);
	}

	printf("writing '/some_file'\n");
	if (!(fh = fopen("/some_file", "w")))
		perror("failed to open '/some_file' for writing");
	else {
		memset(buffer, 'c', sizeof(buffer));
		buffer[sizeof(buffer) - 1] = '\0';
		fprintf(fh, buffer);
		fclose(fh);
	}

	printf("reading '/some_file'\n");
	if (!(fh = fopen("/some_file", "r")))
		perror("failed to open '/some_file' for reading");
	else
		fclose(fh);

	printf("deleting '/some_file'\n");
	if (unlink("/some_file") == -1)
		perror("failed to delete '/some_file'");
}

int main(int argc, char** argv)
{
	int pid = atoi(argv[1]);

	printf("parent pid: %d, my pid: %d\n", getppid(), getpid());

	printf("sending out some signals to pid %d\n", pid);
	test_send_signals(pid);

	printf("testing file i/o\n");
	test_file_io();

	return 0;
}
