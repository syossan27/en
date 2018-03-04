package connection

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

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
func Load(key []byte, path string) (Connections, error) {
	// 保存ファイルから内容を取得
	p, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 保存ファイル内容が空の場合
	if len(p) == 0 {
		return Connections{}, nil
	}

	// 内容を復号
	dec, err := Decrypt(key, string(p))
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

// AES-256で復号
func Decrypt(key []byte, encrypted string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := data[:aes.BlockSize]
	src := data[aes.BlockSize:]
	dst := make([]byte, len(src))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(dst, src)

	return dst, nil
}

// TODO: SaveというよりAdd? その場合Save処理を切り分けたほうがいいかも
func (cs *Connections) Save(c *Connection, key []byte, path string) error {
	// 同じコネクション名があった場合、エラー
	if cs.Exist(c.Name) {
		return errors.New("connection name already exists")
	}

	// 保存ファイルを開く
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// コネクション構造体群にコネクション構造体を追加
	*cs = append(*cs, *c)

	// yaml化
	p, err := yaml.Marshal(cs)
	if err != nil {
		return err
	}

	// yaml化したコネクション構造体群を暗号化
	enc, err := Encrypt(key, p)
	if err != nil {
		return err
	}

	// 保存ファイルに書き込み
	f.WriteString(enc)

	return nil
}

func (cs *Connections) Update(key []byte, path string) error {
	// 保存ファイルを開く
	f, err := os.Create(path)
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
	enc, err := Encrypt(key, p)
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

// テキストをAES-256で暗号化して、base64でエンコード
func Encrypt(key []byte, data []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)
	encoded := base64.StdEncoding.EncodeToString(cipherText)
	return encoded, nil
}
