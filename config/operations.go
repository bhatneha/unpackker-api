package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/nikhilsbhat/unpackker/pkg/backend"
	"gopkg.in/yaml.v2"
)

//Methods implements Pack and Unpack methods to perform packing and unpacking respectively
type Methods interface {
	Pack() string
	Unpack() string
}

//ExecCmd holds the fields required to construct a new cli command
type ExecCmd struct {
	Command string
	Args    []string
	Dir     string
	Writer  io.Writer
}

// GetUinqueID returns an unique ID for the instance called.
func GetUinqueID() string {
	uid := uuid.New()
	return uid.String()
}

//CreateDir creates a directory in the specified path
func CreateDir(path string, uid string) error {
	if err := os.Mkdir(filepath.Join(path, uid), 0777); err != nil {
		return err
	}
	return nil
}

func (e *ExecCmd) getExecCmd() (*exec.Cmd, error) {
	cmd, err := e.getExecutable()
	if err != nil {
		return nil, err
	}
	newArgs := make([]string, 1)
	newArgs = append(newArgs, e.Args...)
	shellCmd := &exec.Cmd{
		Path:   cmd,
		Args:   newArgs,
		Dir:    e.Dir,
		Stdout: e.Writer,
		Stderr: e.Writer,
	}
	return shellCmd, nil
}

func (e *ExecCmd) getExecutable() (string, error) {
	execPath, err := exec.LookPath(e.Command)
	if err != nil {
		return "", err
	}
	return execPath, nil
}

// Pack invokes a CLI to perform packing of an asset
func (p *PackkerInput) Pack() (string, error) {

	dir, _ := filepath.Split(p.PackData.AssetPath)
	p.PackData.Name = fmt.Sprintf("%sunpackker", p.PackData.Name)
	p.PackData.Path = filepath.Join(dir, "/test")
	p.PackData.Environment = "production"
	p.PackData.CleanLocalCache = true
	backend := backend.New()
	backend.Cloud = "fs"
	p.PackData.Backend = backend
	ydata, err := yaml.Marshal(&p.PackData)
	if err != nil {
		log.Println("failed to Marshal", err)
		return "", err
	}
	if err = ioutil.WriteFile(filepath.Join(dir, ".unpackker-config.yaml"), ydata, 0644); err != nil {
		log.Println("failed to write configurations to file", err)
		return "", err
	}
	cmd := NewExecCmd()
	cmd.Command = "unpackker"
	cmd.Dir = dir
	cmd.Args = []string{"generate", "."}
	cmnd, err := cmd.getExecCmd()
	if err != nil {
		log.Println("failed to get the command", err)
		return "", err
	}
	_, err = cmnd.Output()
	if err != nil {
		log.Println("failed to execute the unpakker successfully", err)
		return "", err
	}
	packedAssetPath := filepath.Join(p.PackData.Path, p.PackData.Name+"."+p.PackData.AssetVersion)
	return strings.ReplaceAll(packedAssetPath, ".", "_"), nil
}

//Unpack performs unpacking of an asset
func (u *UnPackkerInput) Unpack() (string, error) {

	return "Unpacked successfully", nil
}

//NewExecCmd initializes a new empty structure ExecCmd
func NewExecCmd() *ExecCmd {
	return &ExecCmd{}
}

// func NewPackMethod() *PackkerInput {
// 	return &PackkerInput{}
// }

// func NewUnpackMethod() *UnPackkerInput {
// 	return &UnPackkerInput{}
// }
