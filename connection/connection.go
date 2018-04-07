package connection

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/syossan27/en/foundation"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"gopkg.in/yaml.v2"
)

type Connection struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Connections []Connection

func (c *Connection) Connect() {
	session, err := connect(c.User, c.Password, c.Host, 22)
	if err != nil {
		foundation.PrintError("Failed to connect.\nReason: " + err.Error())
	}
	defer session.Close()

	// ターミナルの標準入力ファイルディスクリプタを
	// 一時的にterminalに準拠させる
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		foundation.PrintError("Failed to put the terminal into raw mode")
	}
	defer terminal.Restore(fd, oldState)

	// excute command
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		foundation.PrintError("Failed to get terminal size")
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		foundation.PrintError("Failed to requests the association of a pty with the session")
	}

	err = session.Shell()
	if err != nil {
		foundation.PrintError("Failed to login shell")
	}

	err = session.Wait()
	if err != nil {
		foundation.PrintError("Failed to command completes successfully")

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

func New(name, host, user, password string) *Connection {
	return &Connection{
		Name:     name,
		Host:     host,
		User:     user,
		Password: password,
	}
}

// 保存ファイルから復号して内容を返す
func Load() Connections {
	// キーファイル（.ssh/id_rsa）からAESキー取得
	key := foundation.GetKey(foundation.KeyPath)

	// 保存ファイルから内容を取得
	p, err := ioutil.ReadFile(foundation.StorePath)
	if err != nil {
		foundation.PrintError("Failed to read store file")
	}

	// 保存ファイル内容が空の場合
	if len(p) == 0 {
		return nil
	}

	// 内容を復号
	dec, err := foundation.Decrypt(key, string(p))
	if err != nil {
		foundation.PrintError("Failed to decrypt connections")
	}

	// 復号した内容をyaml化
	var cs Connections
	err = yaml.Unmarshal(dec, &cs)
	if err != nil {
		foundation.PrintError("Failed to unmarshal connections yaml")
	}

	return cs
}

func (cs *Connections) Add(name string) {
	// 同じコネクション名があった場合、エラー
	if cs.Exist(name) {
		foundation.PrintError("Connection name already exists")
	}

	// プロンプトで取得
	host, user, password := foundation.AddPrompt()

	// コネクション構造体の作成
	conn := New(name, host, user, password)

	// コネクション構造体群にコネクション構造体を追加
	*cs = append(*cs, *conn)

	save(cs)
}

func (cs *Connections) Update(name string) {
	// コネクション構造体群の中に更新対象のコネクションがあるか確認
	conns := *cs
	for key, conn := range conns {
		if conn.Name == name {
			conns[key].Host, conns[key].User, conns[key].Password = foundation.UpdatePrompt(conn.Host, conn.User, conn.Password)
			break
		}
	}

	save(&conns)
}

func (cs *Connections) Delete(name string) {
	// 削除済みのコネクション構造体群を作成
	// newConns := make(Connections, len(*cs)-1)
	var newConns Connections
	for _, conn := range *cs {
		if conn.Name != name {
			newConns = append(newConns, conn)
		}
	}

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	save(&newConns)
}

func (cs *Connections) List() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Host", "User", "Password"})

	for _, conn := range *cs {
		table.Append([]string{
			conn.Name,
			conn.Host,
			conn.User,
			conn.Password[:1] + "*****",
		})
	}

	table.Render()
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

func (cs Connections) Find(name string) Connection {
	var specifiedConn Connection
	for _, conn := range cs {
		if conn.Name == name {
			specifiedConn = conn
		}
	}
	if specifiedConn == (Connection{}) {
		foundation.PrintError("Not found connect name")
	}
	return specifiedConn
}

func save(cs *Connections) {
	// キーファイル（.ssh/id_rsa）からAESキー取得
	key := foundation.GetKey(foundation.KeyPath)

	// 保存ファイルを開く
	f, err := os.Create(foundation.StorePath)
	if err != nil {
		foundation.PrintError("Failed to open store file")
	}
	defer f.Close()

	// yaml化
	p, err := yaml.Marshal(cs)
	if err != nil {
		foundation.PrintError("Failed to marshal connections yaml")
	}

	// yaml化したコネクション構造体群を暗号化
	enc, err := foundation.Encrypt(key, p)
	if err != nil {
		foundation.PrintError("Failed to encrypt connections")
	}

	// 保存ファイルに書き込み
	_, err = f.WriteString(enc)
	if err != nil {
		foundation.PrintError("Failed to write string to store file")
	}
}
