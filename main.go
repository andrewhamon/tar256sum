package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"sort"
)

func main() {
	gz, err := gzip.NewReader(os.Stdin)

	if err != nil {
		log.Fatal(err)
	}

	tr := tar.NewReader(gz)
	h := sha256.New()
	entryHashes := make(map[string]string)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		hdrHash := hashHeader(hdr, h)
		if _, exists := entryHashes[hdrHash]; exists {
			log.Fatal("Encountered duplicate tar headers")
		}
		entryHashes[hdrHash] = hashContents(tr, h)
	}

	for _, entryHash := range sortedMapKeys(entryHashes) {
		contentHash := entryHashes[entryHash]

		writeStringOrFatal(h, contentHash)
	}

	fmt.Printf("%x  -\n", h.Sum(nil))
}

func hashHeader(hdr *tar.Header, h hash.Hash) string {
	writeOrFatal(h, []byte{hdr.Typeflag})
	writeStringOrFatal(h, hdr.Name)
	writeStringOrFatal(h, hdr.Linkname)

	b8 := make([]byte, 9)
	writeInt64OrFatal(h, hdr.Size, b8)
	writeInt64OrFatal(h, hdr.Mode, b8)
	writeInt64OrFatal(h, int64(hdr.Uid), b8)
	writeInt64OrFatal(h, int64(hdr.Gid), b8)

	writeStringOrFatal(h, hdr.Gname)
	writeStringOrFatal(h, hdr.Uname)

	writeInt64OrFatal(h, hdr.ModTime.UnixMicro(), b8)
	writeInt64OrFatal(h, hdr.AccessTime.UnixMicro(), b8)
	writeInt64OrFatal(h, hdr.ChangeTime.UnixMicro(), b8)

	for _, key := range sortedMapKeys(hdr.PAXRecords) {
		writeStringOrFatal(h, hdr.PAXRecords[key])
	}

	result := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	return result
}

func hashContents(tr *tar.Reader, h hash.Hash) string {
	_, err := io.Copy(h, tr)
	if err != nil {
		log.Fatal(err)
	}
	result := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	return result
}

func writeStringOrFatal(w io.Writer, s string) {
	writeOrFatal(w, []byte(s))
}

func writeInt64OrFatal(w io.Writer, i int64, b []byte) {
	binary.PutVarint(b, i)
	writeOrFatal(w, b)
}

func writeOrFatal(w io.Writer, b []byte) {
	written, err := w.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	if written != len(b) {
		log.Fatal("Did not write all bytes")
	}
}

func sortedMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
