package cable

import (
	"net"
	"testing"
)

func TestConnect(t *testing.T) {
	ftp := New()
	ftp.Debug(true)

	defer ftp.Quit()

	if err := ftp.Connect("ftp.debian.org:21"); err != nil {
		t.Errorf("%s", err)
	}
}

func TestConnectBadHostName(t *testing.T) {
	ftp := New()
	ftp.Debug(true)

	defer ftp.Quit()

	if err := ftp.Connect("nohost:21"); err == nil {
		t.Error("This test should fail.")
	}
}

func TestConnectAnonymousLogin(t *testing.T) {
	ftp := New()
	ftp.Debug(true)

	defer ftp.Quit()

	if err := ftp.Connect("ftp.debian.org:21"); err != nil {
		if terr, ok := err.(net.Error); ok && !terr.Timeout() {
			t.Error("This should fail with connection timeout error.")
		}
	}

	if err := ftp.Login("", ""); err != nil {
		t.Errorf("%s", err)
	}
}

func TestPassiveCommand(t *testing.T) {
	ftp := New()
	ftp.Debug(true)

	defer ftp.Quit()

	if err := ftp.Connect("ftp.debian.org:21"); err != nil {
		t.Errorf("%s", err)
	}

	if err := ftp.Login("", ""); err != nil {
		t.Errorf("%s", err)
	}

	if err := ftp.Pasv(); err != nil {
		t.Errorf("%s", err)
	}

	t.Logf("Passive open port: %d", ftp.passPort)
}

func TestPwdCommand(t *testing.T) {
	ftp := New()
	ftp.Debug(true)

	defer ftp.Quit()

	if err := ftp.Connect("ftp.debian.org:21"); err != nil {
		t.Errorf("%s", err)
	}

	if err := ftp.Login("", ""); err != nil {
		t.Errorf("%s", err)
	}

	if err := ftp.Pwd(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestCwdCommand(t *testing.T) {
	ftp := New()
	ftp.Debug(true)

	defer ftp.Quit()

	if err := ftp.Connect("ftp.debian.org:21"); err != nil {
		t.Errorf("%s", err)
	}

	if err := ftp.Login("", ""); err != nil {
		t.Errorf("%s", err)
	}

	if err := ftp.Cwd("/debian"); err != nil {
		t.Errorf("%s", err)
	}
}
