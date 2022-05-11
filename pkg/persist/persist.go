package persist

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pawalt/kvstore/pkg/kv"
)

func WriteOp(f *os.File, w *bufio.Writer, path []string, data []byte) error {
	formatted := "WRITE\t" + strings.Join(path, "/") + "\t" + string(data) + "\n"
	_, err := w.WriteString(formatted)
	if err != nil {
		return fmt.Errorf("error when writing to buffered reader: %v", err)
	}

	// flush to underlying writer
	err = w.Flush()
	if err != nil {
		return fmt.Errorf("error while flushing: %v", err)
	}

	// ensure disk persistence
	err = f.Sync()
	if err != nil {
		return fmt.Errorf("error while syncing: %v", err)
	}

	return nil
}

func Restore(r *bufio.Reader) (kv.KVNode, error) {
	node := kv.NewMapVKNode()

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		// if line is empty continue
		if line == "\n" {
			continue
		}

		line = strings.Trim(line, "\n")

		components := strings.Split(line, "\t")
		if len(components) != 3 {
			return nil, fmt.Errorf("expected 3 parts of line in restore but got: %v", len(components))
		}

		parsedPath, err := kv.ParsePath(components[1])
		if err != nil {
			return nil, fmt.Errorf("error when attempting to parse restore line: %v", err)
		}

		node.Put(parsedPath, []byte(components[2]))
	}

	return node, nil
}
