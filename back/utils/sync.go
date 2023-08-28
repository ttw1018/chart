package utils

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func sync(localDir, remoteDir, ip, user, port string) error {
	keyPath := "/Users/tiwent/.ssh/id_rsa"

	addr := ip + ":" + port

	signer := getPrivateKey(keyPath)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("Fail to dia: ", err)
	}
	defer client.Close()

	sftpclient, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal("Failed to creat SFTP client: ", err)
	}

	defer sftpclient.Close()

	err = syncRemoteToLocalDirectory(sftpclient, localDir, remoteDir)
	if err != nil {
		log.Fatal("failed to sync directoried: ", err)
	}
	fmt.Println("Syncing directories successfully")
	return err
}

func syncLocalToRemoteDirectory(sftpClient *sftp.Client, localDir string, remoteDir string) error {
	err := filepath.Walk(localDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(localDir, path)
		remotePath, _ := filepath.Rel(remoteDir, relPath)

		if info.IsDir() {
			err := sftpClient.MkdirAll(remotePath)
			if err != nil {
				return fmt.Errorf("failed to create remote directory: %v", err)
			}
		} else {
			err := syncLocalToRemoteFile(sftpClient, path, remotePath)
			if err != nil {
				return err
			}
		}
		return nil

	})
	return err
}

func syncRemoteToLocalDirectory(sftpClient *sftp.Client, localDir string, remoteDir string) error {
	if _, err := os.Stat(localDir); os.IsNotExist(err) {
		err := os.MkdirAll(localDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create local directory: %v", err)
		}
		fmt.Printf("create directory %s\n", localDir)
	}

	entries, err := sftpClient.ReadDir(remoteDir)
	if err != nil {
		return fmt.Errorf("failed to list remote directory: %v", err)
	}

	for _, entry := range entries {
		remotePath := filepath.Join(remoteDir, entry.Name())
		localPath := filepath.Join(localDir, entry.Name())
		if entry.IsDir() {
			err = syncRemoteToLocalDirectory(sftpClient, localPath, remotePath)
			if err != nil {
				return err
			}
		} else {
			if local_info, err1 := os.Stat(localPath); err1 == nil {
				remote_info, _ := sftpClient.Stat(remotePath)
				if remote_info.Size() == local_info.Size() && local_info.ModTime().After(remote_info.ModTime()) {
					continue
				}
			}
			err = syncRemoteToLocalFile(sftpClient, localPath, remotePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func syncLocalToRemoteFile(sftpClient *sftp.Client, localPath, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %v", err)
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to sync local file: %v", err)
	}
	return nil
}

func syncRemoteToLocalFile(sftpClient *sftp.Client, localPath, remotePath string) error {
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %v", err)
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("failed to sync remote file: %v", err)
	}

	fmt.Printf("copy %s -> %s\n", remotePath, localPath)
	return nil
}

func deleteRemoteFile(sftpClient *sftp.Client, remotePath string) error {
	err := sftpClient.Remove(remotePath)
	if err != nil {
		return fmt.Errorf("failed to delete remote file: %v", err)
	}
	return nil
}

func deleteLocalFile(localPath string) error {
	err := os.Remove(localPath)
	if err != nil {
		return fmt.Errorf("failed to delete local file: %v", err)
	}
	return nil
}

func getPrivateKey(keyPath string) ssh.Signer {

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal("Failed to read private kye file: ", err)
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	return signer
}
