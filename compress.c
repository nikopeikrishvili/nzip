#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#define MAX_UNIQUE 256 // Maximum unique hex values (0x00 to 0xFF for 1 byte)
bool isUnique(int value, int uniqueHex[], int size) {
  for (int i = 0; i < size; i++) {
    if (uniqueHex[i] == value) {
      return false; // Value already exists
    }
  }
  return true;
}

void addUnique(int value, int uniqueHex[], int *size) {
  if (isUnique(value, uniqueHex, *size)) {
    uniqueHex[*size] = value;
    (*size)++;
  }
}
void charToHex(char c, char hex[3]) {
  const char *hexDigits = "0123456789ABCDEF";

  hex[0] = hexDigits[(c >> 4) & 0xF]; // High nibble
  hex[1] = hexDigits[c & 0xF];        // Low nibble
  hex[2] = '\0';                      // Null-terminate the string
}
void compress(const char *inputFile, const char *outputFile) {
  uint32_t frequency[MAX_UNIQUE] = {0};
  FILE *in = fopen(inputFile, "rb");   // Open in binary mode
  FILE *out = fopen(outputFile, "wb"); // Write in binary mode
  int uniqueHex[MAX_UNIQUE];
  int size = 0; // Current size of the uniqueHex array

  if (!in || !out) {
    perror("Error opening file");
    exit(EXIT_FAILURE);
  }

  unsigned char currentChar, prevChar;
  int count = 0;

  // Read the first byte
  if (fread(&prevChar, sizeof(unsigned char), 1, in) == 0) {
    fclose(in);
    fclose(out);
    return;
  }

  count = 1;

  // Read the rest of the file byte by byte
  while (fread(&currentChar, sizeof(unsigned char), 1, in) > 0) {
    // printf("Current Char %02x\n,", currentChar);
    // addUnique(currentChar, &uniqueHex, &size);
    frequency[currentChar]++;
    count++;
  }

  printf("Hex Value  | Frequency\n");
  printf("-----------------------\n");
  for (int i = 0; i < MAX_UNIQUE; i++) {
    if (frequency[i] > 0) {
      printf("0x%02X      | %u\n", i, frequency[i]);
    }
  }

  // Write the last character and its count

  printf("Count of all characters %d\n", count);
  printf("Count of all unique characters %d\n", size);
  fclose(in);
  fclose(out);
  printf("File compressed successfully.\n");
}
