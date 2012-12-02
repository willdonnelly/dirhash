dirhash
=======

Package dirhash provides a function to compute the sha256 hash of a
directory. The algorithm which produces this hash is deterministic, and
thus it will always yield the same hash value for an identical directory
structure.

The algorithm works as follows: all files and subdirectories in the
directory to be hashed are listed. The files are hashed using SHA256 and
the subdirectories are hashed recursively using the algorithm described
here. These hash values are assembled into a "pseudo-file" which looks
like this:

    8C1E0D4467DC345BCBE4122CB5F3A872A596FF7B5BB360B1A545FEF5991296AC "bar"
    0CE63AFC1E92EE82744300A778E523B9F42A53FE99201BD39FB8E2DE82965297 "empty"
    A2E5BE5B8170F0B419A304B948D5711E9C555EB74A280EDA1AF6D53BB46478C5 "foo"
    =
    F5F12CF4210548CB4794FA08DD099186F5C4B3424BDC6535F1E63C2EBCD882BE "asd.txt"
    EAD9E82A649437D8A03BE6756862DC2B058976B565440FDAE81FBD9960128B4E "baz.txt"

To be more precise, the file consists of two sections, directories
followed by files, with a '=' symbol on its own line separating them.
Each section consists of zero or more lines, where each line represents
a single file or directory. The lines are ordered based on the literal
bytes of the escaped filename. Each line begins with the hash of the
item being listed (presented in capitalized hexadecimal), followed by a
space, followed by the escaped filename enclosed in quotes (and then the
line is terminated with a single '\n' character, as is the line
containing just the '=' sign).

Filenames are escaped following the familiar rules:

    * Replace all '\' with '\\'
    * Replace all '"' with '\"'

no other characters are escaped, as these changes are sufficient to
unambiguously store any filename.

[package documentation](http://go.pkgdoc.org/github.com/willdonnelly/dirhash)
