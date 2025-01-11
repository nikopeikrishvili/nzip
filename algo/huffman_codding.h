
#ifndef NZIP_HUFFMAN_CODDING_H
#define NZIP_HUFFMAN_CODDING_H
#include <stdlib.h>
struct MinHeapNode {

  // single input character
  char data;

  // Frequency of the character
  size_t frequency;

  // Pointers of left and right child nodes
  struct MinHeapNode *left, *right;
};

struct MinHeap {

  // Current size of min heap
  size_t size;

  // Capacity of min heap
  size_t capacity;

  // Array of minheap node pointers
  struct MinHeapNode **array;
};

#endif
