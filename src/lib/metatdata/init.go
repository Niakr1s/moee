package metadata

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"
)

//go:embed metadata.py
var metadata_py string
var metadata_py_filepath string

func init() {
	err := exec.Command("pip", "install", "mutagen").Run()
	if err != nil {
		log.Fatalf("pip not found, please, install python version 3.10+: %v", err)
	}

	metadata_py_filepath = filepath.Join(os.TempDir(), "metadata.py")
	err = os.WriteFile(metadata_py_filepath, []byte(metadata_py), 0666)
	if err != nil {
		log.Fatalf("couldn't write %s", metadata_py_filepath)
	}
}

func runMetadataPy(args ...string) error {
	args = append([]string{metadata_py_filepath}, args...)
	cmd := exec.Command("python", args...)

	buf := log.Writer()
	cmd.Stdout = buf
	cmd.Stderr = buf

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error in runMetadataPy: %v", err)
	}
	return nil
}
