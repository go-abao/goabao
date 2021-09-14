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
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func createTar(srcDir, tarFIle string) error {

	if !IsExist(srcDir) {
		return errors.New("srcDir not found")
	}

	/* #nosec */
	c := exec.Command("tar", "-cf", tarFIle, "-C", srcDir, ".")
	if err := c.Run(); err != nil {
		return errors.Wrapf(err, "createTar.exec.Run")
	}
	return nil
}

func tempFileName(prefix, suffix string) (string, error) {
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", errors.Wrapf(err, "tempFileName.randBytes.rand.Read")
	}
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix), nil
}

/* ImageBuild
@dir           					  Dockerfile所在文件夹
@name          					  docker镜像地址以及版本信息
@imageBuildOptionsAuthConfigs    容器镜像服务认证信息配置
*/
func ImageBuild(srcDir, dockerImageURI string, imageBuildOptionsAuthConfigs map[string]types.AuthConfig, timeOutSecond int64) ([]string, error) {

	tarFile, err := tempFileName("docker-", ".image")
	if err != nil {
		return nil, errors.Wrapf(err, "docker buildImage.tarFile")
	}
	defer os.Remove(tarFile)

	if err := createTar(srcDir, tarFile); err != nil {
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

	if err := os.Chdir(srcDir); err != nil {
		return nil, errors.Wrapf(err, "docker buildImage.os.Chdir()")
	}

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

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exist.
func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

// IsDir checks whether the path is a dir,
// it returns false when it's a directory or does not exist.
func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}
