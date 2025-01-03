#include "uncompress.h"

#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>

void uncompress(const char *inputFile, const char *outputFile) {
  FILE *in = fopen(inputFile, "r");
  FILE *out = fopen(outputFile, "w");

  if (!in || !out) {
    perror("Error opening file");
    exit(EXIT_FAILURE);
  }

  char ch;
  int count = 0;

  while ((ch = fgetc(in)) != EOF) {
    if (isdigit(ch)) {
      count = count * 10 + (ch - '0'); // Handle multi-digit counts
    } else {
      if (count > 0) {
        for (int i = 0; i < count; i++) {
          fputc(ch, out); // Write 'count' copies of the character
        }
      }
      count = 0;
    }
  }

  fclose(in);
  fclose(out);
  printf("File uncompressed successfully.\n");
}
