package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {

	scanDirectory, err := prepareTestDirTree("testTree/morphia")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("scan directory is ", scanDirectory)
	defer os.RemoveAll(scanDirectory)

	// "-collection-def-file", "config_test.go",
	// "-out-dir", "."
	// "-collection-def-scan-path", scanDirectory
	oldArgs := os.Args
	os.Args = []string{oldArgs[0], "-out-dir", ".", "-collection-def-scan-path", scanDirectory}
	defer func() { os.Args = oldArgs }()

	cfg, err := NewBuilder().Build(context.Background())
	if err != nil {
		t.Fatal(err.Error())
	}

	cs, err := cfg.FindCollectionToProcess()
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, c := range cs {
		t.Log("found file: ", c)
	}

	fmt.Println(cfg)
}

func prepareTestDirTree(tree string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %v\n", err)
	}

	err = os.MkdirAll(filepath.Join(tmpDir, tree), 0755)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	emptyFile, err := os.Create(filepath.Join(tmpDir, tree, "empty-tpmm.json"))
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	emptyFile.Close()
	return tmpDir, nil
}
