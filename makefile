# Compiler and flags
CC = gcc
CFLAGS = -Wall -Wextra -O2

# Target executable name
TARGET = compress_tool

# Source files
SRC = main.c compress.c uncompress.c

# Build the program
all: $(TARGET)

$(TARGET):
	$(CC) $(CFLAGS) -o $(TARGET) $(SRC)

# Clean up the executable
clean:
	rm -f $(TARGET)

# Phony targets
.PHONY: all clean
