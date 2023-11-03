.POSIX:
PREFIX ?= /usr/local
MANPREFIX ?= $(PREFIX)/share/man
GO ?= go
GOFLAGS ?=
VERSION = v0.2.0

all: calendar

calendar:
	$(GO) build -ldflags "-X main.Version=$(VERSION)" $(GOFLAGS) -o build/calendar
	scdoc < doc/calendar.1.scd | sed "s/VERSION/$(VERSION)/g" > build/calendar.1
	scdoc < doc/calendar-config.5.scd | sed "s/VERSION/$(VERSION)/g" > build/calendar-config.5

clean:
	rm -rf build

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	mkdir -p $(DESTDIR)$(MANPREFIX)/man1
	mkdir -p $(DESTDIR)$(MANPREFIX)/man5
	install -m755 build/calendar \
		$(DESTDIR)$(PREFIX)/bin/calendar
	install -m644 build/calendar.1 \
		$(DESTDIR)$(MANPREFIX)/man1/calendar.1
	install -m644 build/calendar-config.5 \
		$(DESTDIR)$(MANPREFIX)/man5/calendar-config.5

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/calendar
	rm -f $(DESTDIR)$(MANPREFIX)/man1/calendar.1
	rm -f $(DESTDIR)$(MANPREFIX)/man5/calendar-config.5

.DEFAULT_GOAL := all

.PHONY: all calendar clean install uninstall
