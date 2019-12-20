package main

import (
	"bufio"
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strings"
)

var (
	configPath  string
	decryptPath string
	sessions    []Session
)

type Session struct {
	Name     string
	HostName string
	UserName string
	Password string
}

func init() {
	flag.Parse()

	flag.StringVar(&configPath, "config-path", "WinSCP.ini", "path to config file")
	flag.StringVar(&decryptPath, "decrypt-path", "WinSCP_decrypt.ini", "log file name")
}

func main() {
	if !fileExists(configPath) {
		log.Fatalf("Файл %s не найден", configPath)
	}

	if fileExists(decryptPath) {
		if err := os.Remove(decryptPath); err != nil {
			log.Fatalln("Не удалось удалить старый результат выполнения", err)
		}
	}

	decryptIni()

	if len(sessions) == 0 {
		log.Fatalln("Ничего не спарсилось...")
	}

	if err := saveResult(); err != nil {
		log.Fatalln("Не удалось сохранить результат", err)
	}

	log.Printf("Файл %s успешно обработан", configPath)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func decryptIni() {
	cfg, err := ini.InsensitiveLoad(configPath)
	if err != nil {
		panic(err)
	}

	for _, c := range cfg.Sections() {
		if c.HasKey("Password") {
			name := strings.TrimPrefix(c.Name(), "sessions\\")
			name = strings.ReplaceAll(name, "%20", " ")

			sessions = append(sessions, Session{
				Name:     name,
				HostName: c.Key("HostName").Value(),
				UserName: c.Key("UserName").Value(),
				Password: decrypt(c.Key("HostName").Value(), c.Key("UserName").Value(), c.Key("Password").Value()),
			})
		}
	}
}

func saveResult() error {
	file, err := os.Create(decryptPath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	for _, session := range sessions {
		_, _ = fmt.Fprintln(
			w,
			fmt.Sprintf(
				"[%s]\nHostName=%s\nUserName=%s\nPassword=%s\n",
				session.Name,
				session.HostName,
				session.UserName,
				session.Password,
			),
		)
	}

	return w.Flush()
}
