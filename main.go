package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoring = 3
const delay = 5

func main() {
	displaysIntroduction()

	for {
		displaysMenu()
		command := readCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("Exibindo Logs...")
			printLogs()
		case 0:
			fmt.Println("Saindo do Programa.")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando :(")
			os.Exit(-1)
		}
	}
}

func displaysIntroduction() {
	version := 0.1
	fmt.Println("Este programa está na versão", version)
	fmt.Println("")
}

func displaysMenu() {
	fmt.Println("1- Iniciar Monitoramente.")
	fmt.Println("2- Exibir Logs.")
	fmt.Println("0- Sair do Programa.")
}

func readCommand() int {
	var commandRead int
	fmt.Scan(&commandRead)

	return commandRead
}

func startMonitoring() {
	fmt.Println("Monitorando...")

	sites := readFileSystems()

	fmt.Println()

	for i := 0; i < monitoring; i++ {
		fmt.Println("Monitoramento: ", i+1)
		for _, site := range sites {
			testSite(site)
			fmt.Println()
		}

		fmt.Println("")
		time.Sleep(delay * time.Second)
	}
}

func testSite(site string) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("[testSite] Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		registerLog(site, true)
	} else {
		fmt.Println("Site:", site, "está com problemas. Status Code:", resp.StatusCode)
		registerLog(site, false)
	}
}

func readFileSystems() []string {
	var sites []string

	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("[readFileSystems] Ocorreu um erro:", err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		sites = append(sites, line)

		if err == io.EOF {
			break
		}
	}

	file.Close()

	return sites
}

func registerLog(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("[registerLog] Ocorreu um erro:", err)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + "- online:" + strconv.FormatBool(status) + "\n")

	file.Close()
}

func printLogs() {
	file, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("[imprimeLogs] Ocorreu um erro:", err)
	}

	fmt.Println(string(file))

}
