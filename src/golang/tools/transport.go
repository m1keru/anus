package tools

import (
	"database/sql"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/ssh"
)

// SSHConfig - структура содержащая настройка подключения по ssh
type SSHConfig struct {
	Host, Port, Login, Password, Keypath string
}

// RunOverSSH - выполняет указанную команду на указанном в SSHConfig хосте
func RunOverSSH(config *SSHConfig, command string) ([]byte, error) {
	sshConfig := &ssh.ClientConfig{
		User:    config.Login,
		Auth:    []ssh.AuthMethod{publicKeyFile(config.Keypath)},
		Timeout: 15 * time.Second}
	client, err := ssh.Dial("tcp", config.Host+":"+config.Port, sshConfig)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := client.Close(); err != nil {
			panic("Cannot Close client")
		}
	}()

	out, err := session.CombinedOutput(command)
	if err != nil {
		log.Println(string(out))
		return nil, err
	}
	return out, nil
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("cannot open RSA keyfile:", err)
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Println("cannot parse RSA private key")
		return nil
	}
	return ssh.PublicKeys(key)
}

//GetAnsibleScripts - возвращает массив со всеми скриптами ansible на указанном хосте.
func GetAnsibleScripts(config *SSHConfig) ([]string, error) {
	command := "find project/scripts/ -type f -iname \"*.yml\""
	out, err := RunOverSSH(config, command)
	if err != nil {
		return nil, err
	}
	scripts := strings.Split(string(out), "\n")
	scripts = scripts[:len(scripts)-1]
	if len(scripts) < 1 {
		return nil, errors.New("No scripts on target host: " + config.Host)
	}
	return scripts, nil
}

func StoreScriptsToDB(scripts *[]string, db *sql.DB) (int64, error) {
	var returnValue, tmpResult int64
	sql := "REPLACE INTO ansible_scripts(path) VALUES(?)"
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	for _, script := range *scripts {
		result, err := stmt.Exec(script)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		tmpResult, _ = result.RowsAffected()
		returnValue += tmpResult
	}
	tx.Commit()
	return returnValue, nil

}

var ChanStack = map[string]chan []byte{}

func RunOverSshChan(config *SSHConfig, command string) (string, error) {
	sshConfig := &ssh.ClientConfig{
		User:    config.Login,
		Auth:    []ssh.AuthMethod{publicKeyFile(config.Keypath)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 5 * time.Second}
	client, err := ssh.Dial("tcp", config.Host+":"+config.Port, sshConfig)
	if err != nil {
		return "", err
	}

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}

	//defer func() {
	//	if err := client.Close(); err != nil {
	//		panic("Cannot Close client")
	//	} else {
	//		log.Println("ssh closed")
	//	}
	//}()

	reader, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}

	go session.Run(command)
	//err = session.Run(command)
	//if err != nil {
	//	return "", err
	//}

	chanId := xid.New().String()
	ChanStack[chanId] = make(chan []byte)

	go chanReader(ChanStack[chanId], reader, chanId)

	return chanId, nil
}

func chanReader(ior chan []byte, reader io.Reader, chanId string) {
	b := make([]byte, 128)
	for {
		_, err := reader.Read(b)
		if err == io.EOF {
			ior <- []byte("done")
			delete(ChanStack,chanId)
			close(ior)
			b = nil
			break
		}
		ior <- b
	}
}
