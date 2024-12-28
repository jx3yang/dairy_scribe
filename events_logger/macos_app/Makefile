EXE = events_logger_app
SOURCES = $(wildcard *.swift)
SWIFT = swiftc

.PHONY: all clean

all: $(EXE)

$(EXE): $(SOURCES)
	$(SWIFT) $(SOURCES)

clean:
	rm $(EXE)
