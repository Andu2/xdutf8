UTF8 stream validator performance could be improved
 - Try using lookup tables rather than manually calculating
 - Try only doing byte-by-byte validation at the edges of the chunk and skip managing state in the middle
