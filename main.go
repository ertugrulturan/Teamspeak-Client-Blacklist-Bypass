package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"io/ioutil"
	"strings"
	"path/filepath"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[31m"
	Blue   = "\033[31m"
)

func getHostsFilePath() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "darwin", "linux":
		return "/etc/hosts"
	default:
		return ""
	}
}

func getCacheDir() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(userHome, "AppData", "Roaming", "TS3Client", "cache")
	case "darwin":
		return filepath.Join(userHome, "Library", "Application Support", "TeamSpeak 3", "cache")
	case "linux":
		return filepath.Join(userHome, ".ts3client", "cache")
	default:
		return ""
	}
}

func isRootOrAdmin() bool {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("net", "session")
		err := cmd.Run()
		return err == nil
	case "linux", "darwin":
		return os.Geteuid() == 0
	default:
		return false
	}
}

func appendLinesToHosts(filePath, line1, line2 string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("hosts dosyası okunamadı: %v", err)
	}

	if strings.Contains(string(content), line1) {
		return fmt.Errorf("bypass işlemi zaten uygulanmış. %s mevcut.", line1)
	}

	if strings.Contains(string(content), line2) {
		return nil
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("hosts dosyası açılırken hata oluştu: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + line1 + "\n" + line2 + "\n"); err != nil {
		return fmt.Errorf("hosts dosyasına yazılamadı: %v", err)
	}

	return nil
}

func clearCache() error {
	cacheDir := getCacheDir()
	return os.RemoveAll(cacheDir)
}

func clearTerminal() {
	cmd := exec.Command("clear") 
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	if !isRootOrAdmin() {
		fmt.Println(Red + "Bu programı root/administrator yetkisi ile çalıştırmalısınız!" + Reset)
		time.Sleep(5 * time.Second)
		return
	}

	line1 := "0.0.0.0 blacklist2.teamspeak.com"
	line2 := "0.0.0.0 blacklist.teamspeak.com"

	hostsFilePath := getHostsFilePath()
	if hostsFilePath == "" {
		fmt.Println(Red + "Desteklenmeyen platform." + Reset)
		return
	}

	for {
		clearTerminal()

		fmt.Println(Green + "\n1: Blacklist Bypass (This server is blacklisted. Refusing to connect.)" + Reset)
		fmt.Println(Yellow + "2: Clear Cache (Önbelleği temizle)" + Reset)
		fmt.Println(Blue + "3: Çıkış" + Reset)

		var choice string
		fmt.Print(Green + "Seçiminizi yapın (1/2/3): " + Reset)
		fmt.Scan(&choice)

		switch choice {
		case "1":
			err := appendLinesToHosts(hostsFilePath, line1, line2)
			if err != nil {
				fmt.Printf(Red + "Hata: %v\n" + Reset, err)
			} else {
				fmt.Println(Green + "Satırlar başarıyla hosts dosyasına eklendi." + Reset)
			}
		case "2":
			err := clearCache()
			if err != nil {
				fmt.Printf(Red + "Hata: %v\n" + Reset, err)
			} else {
				fmt.Println(Green + "Cache başarıyla temizlendi." + Reset)
			}
		case "3":
			fmt.Println(Green + "Çıkılıyor..." + Reset)
			return
		default:
			fmt.Println(Red + "Geçersiz seçim." + Reset)
		}

		time.Sleep(2 * time.Second)
	}
}
