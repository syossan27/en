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

	// Terminal connected to file descriptor into raw mode
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		foundation.PrintError("Failed to put the terminal into raw mode")
	}
	defer terminal.Restore(fd, oldState)

	// Execute command
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
		terminal.Restore(fd, oldState)
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

	// Get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Create session
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
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

// Load connections
func Load() Connections {
	// Get AES key from .ssh/id_rsa
	key := foundation.GetKey(foundation.KeyPath)

	// Get content from store file
	p, err := ioutil.ReadFile(foundation.StorePath)
	if err != nil {
		foundation.PrintError("Failed to read store file")
	}

	if len(p) == 0 {
		return nil
	}

	// Decrypt content
	dec, err := foundation.Decrypt(key, string(p))
	if err != nil {
		foundation.PrintError("Failed to decrypt connections")
	}

	// Decrypted content convert into yaml
	var cs Connections
	err = yaml.Unmarshal(dec, &cs)
	if err != nil {
		foundation.PrintError("Failed to unmarshal connections yaml")
	}

	return cs
}

func (cs *Connections) Add(name string) {
	// If exist same connection name is error
	if cs.Exist(name) {
		foundation.PrintError("Connection name already exists")
	}

	// Get connection information by prompt
	host, user, password := foundation.AddPrompt()

	// Save connection information to store file
	conn := New(name, host, user, password)
	*cs = append(*cs, *conn)
	save(cs)
}

func (cs *Connections) Update(name string) {
	// Check to exist same connection name in Connections struct
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
	var newConns Connections
	for _, conn := range *cs {
		if conn.Name != name {
			newConns = append(newConns, conn)
		}
	}
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

// Check to exist same connection name
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
	// Get AES key from .ssh/id_rsa
	key := foundation.GetKey(foundation.KeyPath)

	// Open store file
	f, err := os.Create(foundation.StorePath)
	if err != nil {
		foundation.PrintError("Failed to open store file")
	}
	defer f.Close()

	// Convert into yaml
	p, err := yaml.Marshal(cs)
	if err != nil {
		foundation.PrintError("Failed to marshal connections yaml")
	}

	// Encrypt the connections struct converted into yaml
	enc, err := foundation.Encrypt(key, p)
	if err != nil {
		foundation.PrintError("Failed to encrypt connections")
	}

	// Write string to store file
	_, err = f.WriteString(enc)
	if err != nil {
		foundation.PrintError("Failed to write string to store file")
	}
}
