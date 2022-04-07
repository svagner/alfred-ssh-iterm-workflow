package main

import (
	"bufio"
	"github.com/hanjm/errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func GetSSHHostList() (hosts []string, err error) {
	filePath, err := GetSSHNodesFilePath()
	if err != nil {
		return hosts, errors.Errorf(err, "")
	}
	fp, err := os.Open(filePath)
	if err != nil {
		return hosts, errors.Errorf(err, "failed to open file:%s", filePath)
	}
	defer fp.Close()
	sc := bufio.NewScanner(fp)
	sc.Split(bufio.ScanLines)
	const hostPrefix = "host"
	const hostPrefixLen = len(hostPrefix)
	const hostNamePrefix = "hostname"
	const hostNamePrefixLen = len(hostNamePrefix)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "fqdn:") {
                        data := strings.Split(line, ":")
			if len(data) == 2 {
			hosts = append(hosts, data[1])
}
		}
	}
	return hosts, nil
}

func GetSSHNodesFilePath() (filePath string, err error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return "", errors.Errorf(err, "")
	}
	filePath = filepath.Join(homeDir, ".ssh", "nodes")
	return filePath, nil
}

func GetHomeDir() (homeDir string, err error) {
	u, err := user.Current()
	if nil == err {
		return u.HomeDir, nil
	}
	// 环境变量
	if v := os.Getenv("HOME"); v != "" {
		return v, nil
	}
	// shell
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Errorf(nil, "failed to get home from shell, os:%s", runtime.GOOS)
	}
	if v := strings.TrimSpace(string(output)); v != "" {
		return v, nil
	}
	return "", errors.Errorf(nil, "failed to get home, os:%s", runtime.GOOS)
}
