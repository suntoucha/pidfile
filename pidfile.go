package pidfile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type PidFile struct {
	name string
}

func Init(fname string) *PidFile {
	if strings.Trim(fname, " ") == "" {
		return nil
	}
	return &PidFile{name: fname}
}

// Removing PID file
func (this *PidFile) Remove() error {
	return os.Remove(this.name)
}

// Checking pid file, if there is no errors, creates it and returns nil
func (this *PidFile) TryLock() error {
	// В случае возникновения ошшибок - считаем, что идет синхронизация
	// Если PID файл не обнаружен либо если не найдет процесс с таким PID
	//  - разрешается синхронизация. В любых других случаях синхронизация
	// отменяется.
	if p, err := os.Open(this.name); err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return fmt.Errorf("Error occurred while reading PID file: %s\n", err.Error())
		}
		p, err = os.Create(this.name)
		if err != nil {
			return fmt.Errorf("Error occurred while creating PID file: %s\n", err.Error())
		}
		defer p.Close()

		p.WriteString(fmt.Sprintf("%d", syscall.Getpid()))
	} else {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(p)
		if err != nil {
			return fmt.Errorf("Error occurred while reading pid file: %s\n", err.Error())
		}
		_, err = exec.Command("kill", "-0", buf.String()).Output()
		if err == nil {
			return fmt.Errorf("Found running instance\n")
		} else {
			fmt.Printf("Process with PID [%s] was not found\n", buf.String())
		}
	}
	return nil
}
