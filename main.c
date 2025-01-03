#include "compress.h"
#include "uncompress.h"
#include <stdio.h>  // For standard I/O functions like printf, fopen, fclose
#include <stdlib.h> // For standard library functions like exit
#include <string.h> // For string handling functions like strcmp

// Function prototypes
void compress(const char *inputFile, const char *outputFile);

int main(int argc, char *argv[]) {

  if (argc != 4) {
    printf("Usage: %s <compress|uncompress> <input_file> <output_file>\n",
           argv[0]);
    return EXIT_FAILURE;
  }

  if (strcmp(argv[1], "compress") == 0) {
    compress(argv[2], argv[3]);
  } else if (strcmp(argv[1], "uncompress") == 0) {
    uncompress(argv[2], argv[3]);
  } else {
    printf("Invalid option. Use 'compress' or 'uncompress'.\n");
    return EXIT_FAILURE;
  }

  return EXIT_SUCCESS;
}
