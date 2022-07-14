package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-clix/cli"
	neat "github.com/itaysk/kubectl-neat/cmd"
)

func main() {
	log.SetFlags(0)

	cmd := cli.Command{
		Use:   "kubectl-neat-diff [file1] [file2]",
		Short: "Remove fields from kubectl diff that carry low / no information",
		Args:  cli.ArgsExact(2),
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		if err := statAndNeatifyPath(args[0]); err != nil {
			return err
		}
		if err := statAndNeatifyPath(args[1]); err != nil {
			return err
		}

		c := exec.Command("diff", "-uN", args[0], args[1])
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func statAndNeatifyPath(path string) error {
	// Stat evaluates symlinks so they are supported (if they lead to dirs or regular files).
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	return neatify(path, fi.Mode())
}

// neatify calls the appropriate function for the file based on FileMode
// neatify does not call Stat() because neatifyDir() aleady has FileInfo for the passed file
func neatify(path string, fm os.FileMode) error {
	if fm.IsDir() {
		return neatifyDir(path)
	}
	if fm.IsRegular() {
		return neatifyFile(path, fm)
	}
	return fmt.Errorf("Found file '%s' which is neither directory nor a regular file. "+
		"Special files like named pipes or device files are not supported.", path)
}

func neatifyDir(dir string) error {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		filename := filepath.Join(dir, fi.Name())
		if err = neatify(filename, fi.Mode()); err != nil {
			return err
		}
	}

	return nil
}

func neatifyFile(filename string, fm os.FileMode) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	n, err := neat.NeatYAMLOrJSON(data)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, []byte(n), fm); err != nil {
		return err
	}
	return nil
}
