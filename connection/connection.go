package connection

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Songmu/prompter"
	"github.com/syossan27/en/foundation"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"gopkg.in/yaml.v2"
)

type Connection struct {
	Name        string `yaml:"name"`
	AccessPoint string `yaml:"accessPoint"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
}

type Connections []Connection

func (c *Connection) Connect() {
	session, err := connect(c.User, c.Password, c.AccessPoint, 22)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// ターミナルの標準入力ファイルディスクリプタを
	// 一時的にterminalに準拠させる
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	// excute command
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		log.Fatal(err)
	}

	err = session.Shell()
	if err != nil {
		log.Println(err)
	}

	err = session.Wait()
	if err != nil {
		log.Println(err)
	}
}

func connect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

func New(name, accessPoint, user, password string) *Connection {
	return &Connection{
		Name:        name,
		AccessPoint: accessPoint,
		User:        user,
		Password:    password,
	}
}

// 保存ファイルから復号して内容を返す
func Load() (Connections, error) {
	// キーファイル（.ssh/id_rsa）からAESキー取得
	key := foundation.GetKey(foundation.KeyPath)

	// 保存ファイルから内容を取得
	p, err := ioutil.ReadFile(foundation.StorePath)
	if err != nil {
		return nil, err
	}

	// 保存ファイル内容が空の場合
	if len(p) == 0 {
		return Connections{}, nil
	}

	// 内容を復号
	dec, err := foundation.Decrypt(key, string(p))
	if err != nil {
		return nil, err
	}

	// 復号した内容をyaml化
	var cs Connections
	err = yaml.Unmarshal(dec, &cs)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (cs *Connections) Add(c *Connection) error {
	// 同じコネクション名があった場合、エラー
	if cs.Exist(c.Name) {
		return errors.New("connection name already exists")
	}

	// コネクション構造体群にコネクション構造体を追加
	*cs = append(*cs, *c)

	err := save(cs)
	return err
}

func (cs *Connections) Update(name string) error {
	// コネクション構造体群の中に更新対象のコネクションがあるか確認
	conns := *cs
	for key, conn := range conns {
		if conn.Name == name {
			// 更新内容をプロンプトで取得
			var accessPoint = prompter.Prompt("AccessPoint", conn.AccessPoint)
			var user = prompter.Prompt("User", conn.User)
			var password = prompter.Password("Password")
			if password == "" {
				password = conn.Password
			}
			conns[key].AccessPoint = accessPoint
			conns[key].User = user
			conns[key].Password = password

			break
		}
	}

	err := save(&conns)
	return err
}

func (cs *Connections) Delete(name string) error {
	// コネクション構造体群の中に更新対象のコネクションがあるか確認
	newConns := make(Connections, len(*cs)-1)
	for _, conn := range *cs {
		if conn.Name != name {
			newConns = append(newConns, conn)
		}
	}

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	err := save(&newConns)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func save(cs *Connections) error {
	// キーファイル（.ssh/id_rsa）からAESキー取得
	key := foundation.GetKey(foundation.KeyPath)

	// 保存ファイルを開く
	f, err := os.Create(foundation.StorePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// yaml化
	p, err := yaml.Marshal(cs)
	if err != nil {
		return err
	}

	// yaml化したコネクション構造体群を暗号化
	enc, err := foundation.Encrypt(key, p)
	if err != nil {
		return err
	}

	// 保存ファイルに書き込み
	f.WriteString(enc)
	return nil
}

// 同じコネクション名があるか確認
func (cs Connections) Exist(name string) bool {
	for _, c := range cs {
		if name == c.Name {
			return true
		}
	}
	return false
}
