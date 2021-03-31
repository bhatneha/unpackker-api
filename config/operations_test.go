package config

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getExecCmd(t *testing.T) {
	t.Run("should retrun valid exec.Cmd", func(t *testing.T) {
		expected := &exec.Cmd{
			Path:   "/home/bhatneha/go/bin/unpackker",
			Args:   []string{"generate", "."},
			Dir:    "/tmp/75260522-42e0-4992-9005-b2dc762f54d5",
			Stdout: os.Stdout,
			Stderr: os.Stdout,
		}

		newCmd := NewExecCmd()
		newCmd.Command = "unpackker"
		newCmd.Args = []string{"generate", "."}
		newCmd.Dir = "/tmp/75260522-42e0-4992-9005-b2dc762f54d5"
		newCmd.Writer = os.Stdout
		actual, err := newCmd.getExecCmd()

		assert.Nil(t, err)
		assert.Equal(t, *expected, *actual)
	})

	t.Run("should not return valid command but return error", func(t *testing.T) {
		expected := &exec.Cmd{
			Path:   "/home/bhatneha/go/bin/unpackker",
			Args:   []string{"generate", "."},
			Dir:    "/tmp/75260522-42e0-4992-9005-b2dc762f54d5",
			Stdout: os.Stdout,
			Stderr: os.Stdout,
		}

		cmd := NewExecCmd()
		cmd.Command = "unpacer"
		cmd.Args = []string{"generate", "."}
		cmd.Dir = "/tmp/75260522-42e0-4992-9005-b2dc762f54d5"
		cmd.Writer = os.Stdout
		actual, err := cmd.getExecCmd()

		assert.NotNil(t, err)
		assert.NotEqual(t, expected, actual)
	})
}

func Test_getExecutable(t *testing.T) {
	t.Run("should return the path of the library is present", func(t *testing.T) {
		expected := "/home/bhatneha/go/bin/unpackker"

		e := NewExecCmd()
		e.Command = "unpackker"
		actual, err := e.getExecutable()

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
