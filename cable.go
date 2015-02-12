package cable

import (
	"bufio"
	"errors"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// VERSION constant represents the current Cable version in
// "major.minor.release" notation.
const VERSION = "0.0.1"

// CRLF constant represents the End-of-Line sequence.
const CRLF = "\r\n"

// A FTP contain options to handle the FTP connection.
type FTP struct {
	debug    bool
	addr     string
	conn     net.Conn
	isLogged bool
	passPort int

	reader *bufio.Reader
	writer *bufio.Writer
}

type message struct {
	code    int
	message string
}

// New returns a FTP struct initializated to default values.
func New() *FTP {
	return &FTP{}
}

// Debug set the level of verbosity.
func (ftp *FTP) Debug(debug bool) {
	ftp.debug = debug
}

// Connect perform the dial to the ftp server.
func (ftp *FTP) Connect(addr string) (err error) {
	if ftp.conn, err = net.Dial("tcp", addr); err != nil {
		return err
	}

	ftp.addr = addr
	ftp.writer = bufio.NewWriter(ftp.conn)
	ftp.reader = bufio.NewReader(ftp.conn)

	_, err = ftp.receive()

	if err != nil {
		return err
	}

	return nil
}

// Close closes the current connection, if exists any.
func (ftp *FTP) Close() {
	// We should ensure the connection was established before try to close it.
	if ftp.conn != nil {
		ftp.conn.Close()
		ftp.conn = nil
	}
}

// Login command to FTP service, this methos is sequence of command to perform
// the login operation.
func (ftp *FTP) Login(user string, password string) (err error) {
	if err = ftp.user(user); err != nil {
		return err
	}

	if err = ftp.pass(password); err != nil {
		return err
	}

	ftp.isLogged = true

	return nil
}

// Quit terminates the connection to the FTP server.
func (ftp *FTP) Quit() (err error) {
	defer ftp.Close()

	if err = ftp.send("QUIT"); err != nil {
		return err
	}

	_, err = ftp.receive()

	if err != nil {
		return err
	}

	return nil
}

// User command send the user's name to the ftp server.
// This command belons to the Access Control Command set of commands.
func (ftp *FTP) user(user string) (err error) {
	if len(user) == 0 {
		user = "anonymous"
	}

	if err := ftp.send("USER " + user); err != nil {
		return err
	}

	_, err = ftp.receive()

	if err != nil {
		return err
	}

	return nil
}

// Password command send the user's password to the ftp server.
// This command belons to the Access Control Command set of commands.
func (ftp *FTP) pass(password string) (err error) {
	if len(password) == 0 {
		password = ""
	}

	if err = ftp.send("PASS " + password); err != nil {
		return err
	}

	_, err = ftp.receive()

	if err != nil {
		return err
	}

	return nil
}

// Pasv command requests the server-DTP to "listen" on data port ant wait for a
// connection rather than iniciate one upon receipt of a transfer command.
func (ftp *FTP) Pasv() (err error) {
	if err = ftp.send("PASV"); err != nil {
		return err
	}

	// The response to this command includes the host and port address this
	// server is listening on.
	var msg *message
	msg, err = ftp.receive()

	if err != nil {
		return err
	}

	re, err := regexp.Compile(`\((.*)\)`)
	res := re.FindAllStringSubmatch(msg.message, -1)

	s := strings.Split(res[0][1], ",")

	l1, _ := strconv.Atoi(s[len(s)-2])
	l2, _ := strconv.Atoi(s[len(s)-1])

	ftp.passPort = l1<<8 + l2

	return nil
}

// Pwd comand causes the name of the current working directory to be returned
// in the replay.
func (ftp *FTP) Pwd() (err error) {
	if err = ftp.send("PWD"); err != nil {
		return err
	}

	_, err = ftp.receive()

	if err != nil {
		return err
	}

	return nil
}

// Cwd command allows the user to work with a different directory or dataset
// for file storage or retrieval without altering his login or accounting
// information.
func (ftp *FTP) Cwd(pathname string) (err error) {
	if err = ftp.send("CWD " + pathname + CRLF); err != nil {
		return err
	}

	_, err = ftp.receive()

	if err != nil {
		return err
	}

	return nil
}

func (ftp *FTP) send(cmd string) (err error) {
	// We need to ensure that we have a connection open before send the command.
	if ftp.writer == nil {
		return errors.New("No Connection open.")
	}

	if _, err = ftp.writer.WriteString(cmd + CRLF); err != nil {
		return err
	}

	ftp.writer.Flush()

	return nil
}

func (ftp *FTP) receive() (*message, error) {
	var reply string
	var err error

	if reply, err = ftp.reader.ReadString('\n'); err != nil {
		return &message{}, err
	}

	code, err := strconv.Atoi(reply[:3])
	msg := reply[4:]

	if ftp.debug {
		log.Printf("Code: %d Message: %s", code, msg)
	}

	return &message{code: code, message: msg}, nil
}
