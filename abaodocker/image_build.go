/*
 * @Author: lorock
 * @Github: https://github.com/lorock
 * @Date: 2021-09-14 15:13:14
 * @LastEditors: lorock
 * @LastEditTime: 2021-09-15 16:38:58
 * @FilePath: /goabao/abaodocker/image_build.go
 * @Description:
 */
package abaodocker

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-abao/goabao/abaois"
	"github.com/pkg/errors"
)

/**
 * @description:
 * @param {string} srcDockerfileDir
 * @param {DockerImageBaseConfig} dockerImageBaseConfig
 * @param {AuthConfigMap} imageBuildOptionsAuthConfigs
 * @param {int64} timeOutSecond
 * @return {*}
 */
func ImageBuild(srcDockerfileDir string, dockerImageBaseConfig DockerImageBaseConfig, imageBuildOptionsAuthConfigs AuthConfigMap, timeOutSecond int64) ([]string, error) {

	tarFile, err := tempFileName("docker-", ".image")
	if err != nil {
		return nil, errors.Wrapf(err, "docker buildImage.tarFile")
	}
	defer os.Remove(tarFile)

	if err := createTar(srcDockerfileDir, tarFile); err != nil {
		return nil, errors.Wrapf(err, "docker buildImage.createTar")
	}

	/* #nosec */
	dockerFileTarReader, err := os.Open(tarFile)
	if err != nil {
		return nil, errors.Wrapf(err, "docker buildImage.dockerFileTarReader")
	}
	defer dockerFileTarReader.Close()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, errors.Wrapf(err, "docker buildImage.cli")
	}

	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOutSecond)*time.Second)
	defer cancel()

	buildArgs := make(map[string]*string)

	PWD, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrapf(err, "docker buildImage.PWD")
	}

	defer os.Chdir(PWD)

	if err := os.Chdir(srcDockerfileDir); err != nil {
		return nil, errors.Wrapf(err, "docker buildImage.os.Chdir()")
	}
	dockerImageUrl := []string{
		dockerImageBaseConfig.RegistryServerURI,
		dockerImageBaseConfig.ImageNamespace,
		dockerImageBaseConfig.ImageName,
	}

	dockerImageURI := strings.Join(dockerImageUrl, "/") + ":" + dockerImageBaseConfig.ImageVersion

	resp, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Dockerfile:  "./Dockerfile",
			Tags:        []string{dockerImageURI},
			NoCache:     true,
			Remove:      true,
			PullParent:  true,
			BuildArgs:   buildArgs,
			AuthConfigs: imageBuildOptionsAuthConfigs,
		})

	if err != nil {
		return nil, errors.Wrapf(err, "docker buildImage.resp")
	}

	defer resp.Body.Close()

	var messages []string

	rd := bufio.NewReader(resp.Body)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return messages, errors.Wrapf(err, "docker buildImage.bufio.NewReader()")
		}
		messages = append(messages, string(n))
	}

	return messages, nil
}

/**
 * @description:
 * @param {*} srcDir
 * @param {string} tarFIle
 * @return {*}
 */
func createTar(srcDir, tarFIle string) error {

	if !abaois.IsExist(srcDir) {
		return errors.New("srcDir not found")
	}

	/* #nosec */
	c := exec.Command("tar", "-cf", tarFIle, "-C", srcDir, ".")
	if err := c.Run(); err != nil {
		return errors.Wrapf(err, "createTar.exec.Run")
	}
	return nil
}

/**
 * @description:
 * @param {*} prefix
 * @param {string} suffix
 * @return {*}
 */
func tempFileName(prefix, suffix string) (string, error) {
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", errors.Wrapf(err, "tempFileName.randBytes.rand.Read")
	}
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix), nil
}
