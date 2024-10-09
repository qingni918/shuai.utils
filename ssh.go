package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

// NewSshClientWithPublicKey ("dev". "192.168.1.49", 22) 新建一个ssh客户端
//
//	host := "47.102.149.203"
//	user := "root"
//
// 环境要求, 需要先进行ssh连接, 存储known_hosts记录, 远端存储id_rsa_pub
func NewSshClientWithPublicKey(user, host string, port int) (*ssh.Client, error) {

	//file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	//if err != nil {
	//	return nil, err
	//}
	//defer file.Close()
	//
	//scanner := bufio.NewScanner(file)
	//var hostKey ssh.PublicKey
	//for scanner.Scan() {
	//	fields := strings.Split(scanner.Text(), " ")
	//	if len(fields) != 3 {
	//		continue
	//	}
	//	if strings.Contains(fields[0], host) {
	//		hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
	//		if err != nil {
	//			return nil, fmt.Errorf("error parsing %q: %v", fields[2], err)
	//		}
	//		break
	//	}
	//}
	//
	//if hostKey == nil {
	//	return nil, fmt.Errorf("no hostkey for %s", host)
	//}

	pkey, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa")
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(pkey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)

	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//sshConfig := ssh.ClientConfig{
	//	User: "dev",
	//	Auth: []ssh.AuthMethod{
	//		ssh.Password("1234"),
	//	},
	//	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	//}

	host = fmt.Sprintf("%s:%d", host, port)
	//log.Println(user, host, signer.PublicKey().Type())
	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, err
	}

	//session, err := client.NewSession()
	//if err != nil {
	//	return nil, err
	//}

	return client, nil
}

// NewSshClientWithPassword 使用账密创建ssh client
func NewSshClientWithPassword(user, password, host string, port int) (*ssh.Client, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	host = fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
