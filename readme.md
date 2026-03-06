# XD UTF-8

Repository for my Go UTF-8 utils.

## Stream Validator (xdutf8/streamvalidator)

Incrementally checks whether the octets passed to it are UTF-8. 
If a character has not finished yet, it will still return that the string is valid, just not complete.

This is useful when receiving a stream of bytes and needing to validate the incoming UTF-8 which may not be complete yet.
The standard library unicode/utf8 package is only capable of validating full strings.

