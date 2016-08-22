CFLAGS=-static -O2 -Wall
TARGETS := bad_exit capabilities forking hello testa testa_nonforking zombie

all: $(TARGETS)

install: all
	@cp $(TARGETS) ../rootfs

clean:
	rm -f $(TARGETS)

.PHONY: clean all install
