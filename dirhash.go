/*
   Package dirhash provides a function to compute the sha256 hash of a directory. The algorithm
   which produces this hash is deterministic, and thus it will always yield the same hash value
   for an identical directory structure.

   The algorithm works as follows: all files and subdirectories in the directory to be hashed
   are listed. The files are hashed using SHA256 and the subdirectories are hashed recursively
   using the algorithm described here. These hash values are assembled into a "pseudo-file"
   which looks like this:

       8C1E0D4467DC345BCBE4122CB5F3A872A596FF7B5BB360B1A545FEF5991296AC "bar"
       0CE63AFC1E92EE82744300A778E523B9F42A53FE99201BD39FB8E2DE82965297 "empty"
       A2E5BE5B8170F0B419A304B948D5711E9C555EB74A280EDA1AF6D53BB46478C5 "foo"
       =
       F5F12CF4210548CB4794FA08DD099186F5C4B3424BDC6535F1E63C2EBCD882BE "asd.txt"
       EAD9E82A649437D8A03BE6756862DC2B058976B565440FDAE81FBD9960128B4E "baz.txt"

   To be more precise, the file consists of two sections, directories followed by files, with
   a '=' symbol on its own line separating them. Each section consists of zero or more lines,
   where each line represents a single file or directory. The lines are ordered based on the
   literal bytes of the escaped filename. Each line begins with the hash of the item being
   listed (presented in capitalized hexadecimal), followed by a space, followed by the escaped
   filename enclosed in quotes (and then the line is terminated with a single '\n' character,
   as is the line containing just the '=' sign).

   Filenames are escaped following the familiar rules:
     * Replace all '\' with '\\'
     * Replace all '"' with '\"'
   no other characters are escaped, as these changes are sufficient to unambiguously store
   any filename.
*/
package dirhash

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

// HashDir performs the directory hashing algorithm described previously.
func HashDir(path string) ([]byte, error) {
	// Open whatever's at the given path
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the info corresponding to whatever we opened
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Error out if it isn't a directory
	if !info.IsDir() {
		return nil, errors.New("not a directory")
	}

	// Get the full list of directory contents
	contents, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}

	// Iterate over the contents of the directory accumulating hashes recursively
	var dirs = make(map[string]string)
	var files = make(map[string]string)
	for _, x := range contents {
		if x.IsDir() {
			hash, err := HashDir(path + "/" + x.Name())
			if err != nil {
				return nil, err
			}
			dirs[x.Name()] = fmt.Sprintf("%X", hash)
		} else {
			hash, err := HashFile(path + "/" + x.Name())
			if err != nil {
				return nil, err
			}
			files[x.Name()] = fmt.Sprintf("%X", hash)
		}
	}

	// Create lists of all subdirectories and files in alphabetical order
	var dirPaths []string
	for k, _ := range dirs {
		dirPaths = append(dirPaths, k)
	}
	sort.Strings(dirPaths)

	var filePaths []string
	for k, _ := range files {
		filePaths = append(filePaths, k)
	}
	sort.Strings(filePaths)

	// Create the special "file" representing the directory's contents
	var pseudoFile string
	for _, dirPath := range dirPaths {
		pseudoFile += dirs[dirPath] + " \"" + escape(dirPath) + "\"\n"
	}
	pseudoFile += "=\n"
	for _, filePath := range filePaths {
		pseudoFile += files[filePath] + " \"" + escape(filePath) + "\"\n"
	}
	log.Printf("Hashing directory:\n\"\"\"\n%s\"\"\"\n", pseudoFile)

	// Hash this special file
	hasher := sha256.New()
	_, err = hasher.Write([]byte(pseudoFile))
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func escape(x string) string {
	return strings.NewReplacer("\\", "\\\\", "\"", "\\\"").Replace(x)
}

// HashFile ought to yield the same hash values as the unix 'sha256sum' utility.
func HashFile(path string) ([]byte, error) {
	// Read whatever's at the given path
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Feed the file contents into an SHA256 hash
	hasher := sha256.New()
	_, err = hasher.Write(contents)
	if err != nil {
		return nil, err
	}

	// And return the hash output
	return hasher.Sum(nil), nil
}
