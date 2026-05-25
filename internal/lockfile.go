package internal

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Package struct {
	Sha  string `yaml:"sha"`
	Path string `yaml:"path"`
	Out  string `yaml:"out"`
	Repo string `yaml:"repo"`
	Ref  string `yaml:"ref"`
}

type Lockfile struct {
	Synced   string             `yaml:"updated"`
	Packages map[string]Package `yaml:"packages"`
}

const (
	LockfileName = "fig-lock.yml"
	TimeFormat   = time.RFC3339
)

func ReadLockfile() (Lockfile, error) {
	var lockfile Lockfile

	data, err := os.ReadFile(LockfileName)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return Lockfile{}, nil
	}

	if err := yaml.Unmarshal(data, &lockfile); err != nil {
		return Lockfile{}, err
	}

	return lockfile, nil
}

func ParseTime(s string) (time.Time, error) {
	t, err := time.Parse(TimeFormat, s)
	if err != nil {
		return t, err
	}

	return t, nil
}

func GetFormattedTime() string {
	return time.Now().Format(TimeFormat)
}

func WriteLockfile(lockfile *Lockfile) error {
	data, err := yaml.Marshal(*lockfile)
	if err != nil {
		return err
	}

	return os.WriteFile(LockfileName, data, 0o644)
}
