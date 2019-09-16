package main

import (
"fmt"
"io"
"io/ioutil"
"net"
"os"
"strings"
"encoding/json"
"strconv"
"bufio"
)

const BUFFERSIZE = 1024
const CONNECTIONACCEPTED = ">> Connection Accepted By The Server"
const CONNECTIONREJECTED = ">> Connection Rejected By The Server"

const DEFAULTREPONAME = "repo_new"
const FILEPATH = "repo_new/polymer.js"
const FILEPATHCONFIG = "repo_new/lgconfig.json"

const CMDCOMMIT = "lg -commit"
const CMDUPDATE = "lg -update"
const CMDCREATE = "lg -create"
const CMDBACKUP = "lg -backup"
const CMDEXIT =   "lg -logout"
const CMDUPTOV =   "lg -uptov"

func main() {
	// server, err := net.Listen("tcp", "10.5.12.114:9999")
	server, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		// fmt.Println("Error listetning: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("\n\n\t***************Welcome**********************\n")
	fmt.Println("\t>> Server started! Waiting for connections...")
	for {
		connection, err := server.Accept()
		if err != nil {
			// fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("\t>> Client Connection Initiated")
		go handleClientConnection(connection)
	}
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func handleClientConnection(connection net.Conn) {
	clientIp := make([]byte, 17)
	connection.Read(clientIp)
	fmt.Print(
		fmt.Sprintf("\t>> Allow Client with ip %s to connect?\n\n\t Yes/No >> ", string(clientIp)))
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if strings.Contains(text, "es") { 
		connection.Write([]byte(CONNECTIONACCEPTED))
		clientCMD := make([]byte, 10)
		connection.Read(clientCMD)
		if strings.Contains(string(clientCMD), CMDUPDATE){
			fmt.Println("\t>> Update Command Recieved")
			fmt.Println(fmt.Sprintf("\t>> Updating user with ip %s to current version", clientIp))
			go sendFileToClient(connection, FILEPATH)
		} else if strings.Contains(string(clientCMD), CMDCREATE){
				fmt.Println("\t>> Create Command Accepted")
				newRepositoryName := make([]byte, 20)
				connection.Read(newRepositoryName)
				fmt.Println("\t>> New Repository Created")
				err := os.Mkdir(DEFAULTREPONAME, 0x777)
				if err != nil {
					// fmt.Println("\t>> Error In Handling the repository creation", err)					
				}

				configFile, err := os.Create(fmt.Sprintf("%s/lgconfig.json", DEFAULTREPONAME))
				newRepositoryName = []byte(strings.Trim(string(newRepositoryName), ":"))
				mapD := map[string]string{"repository_name": string(newRepositoryName), "version": "1"}
  			// mapB, _ := json.Marshal(mapD)
  			json.NewEncoder(configFile).Encode(mapD)
  			configFile.Close()
				if err != nil{
					// fmt.Println("\t>> Error In Handling the repository ", err)					
				}

				clientCMDCOMMIT := make([]byte, 10)
				connection.Read(clientCMDCOMMIT)
				if strings.Contains(string(clientCMDCOMMIT), CMDCOMMIT){
					fmt.Println("\t>> Commit Command Accepted")
					recieveFile(connection)
				}

		} else if strings.Contains(string(clientCMD), CMDCOMMIT){
				fmt.Println("\t>> Commit Command Accepted")
				recieveFile(connection)
		} else if strings.Contains(string(clientCMD), CMDBACKUP){
				fmt.Println("\t>> Backing up and Replicating...")
				sendFileToClient(connection, FILEPATH)
		} else if strings.Contains(string(clientCMD), CMDUPTOV){
			// TODO: make this a generic code
				if strings.Contains(string(clientCMD), "1"){
					fmt.Println("\t>> Updating Client to Version # 1")
					sendFileToClient(connection, fmt.Sprintf("%sv1", FILEPATH))
				} else if strings.Contains(string(clientCMD), "2"){
					fmt.Println("\t>> Updating Client to Version # 2")
					sendFileToClient(connection, fmt.Sprintf("%sv2", FILEPATH))
				} else if strings.Contains(string(clientCMD), "3"){
						fmt.Println("\t>> Updating Client to Version # 3")
						sendFileToClient(connection, fmt.Sprintf("%sv3", FILEPATH))
				}
		} else if strings.Contains(string(clientCMD), CMDEXIT){
				fmt.Println("\t>> Client Exited Consensually.")
		} else {
			fmt.Println("Unknown Command", string(clientCMD))
		}
	} else{
		connection.Write([]byte(CONNECTIONREJECTED))
	}
}

func sendFileToClient(connection net.Conn, filepath string) {
	defer connection.Close()
	file, err := os.Open(filepath)
	if err != nil {
		// fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		// fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)

	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("\t>> Client Updated Successfully")
	return
}

func recieveFile(connection net.Conn){
	
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	
	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	// fmt.Println("Creating commited file ", fmt.Sprintf("repo_new/%s",fileName))
	fmt.Println("\t\t>> Creating commited file ",fileName)
	fileToBeCommited := fmt.Sprintf("repo_new/%s",fileName)

	raw, err2 := ioutil.ReadFile(FILEPATHCONFIG)
	mapReadFromJsonFile := make(map[string]string)

	if err2 != nil{
		fmt.Println("failed loading json file into map")
	}

	json.Unmarshal(raw, &mapReadFromJsonFile)
	// Convert the string version read from json to int
	// flag.Parse()
  // s := flag.Arg(mapReadFromJsonFile["version"])
  // string to int
  i, err := strconv.Atoi(mapReadFromJsonFile["version"])

	fmt.Println("\t\t>> Current Version Before Commit is ", i)

	if _, err := os.Stat(fileToBeCommited); err == nil {
	  // there is an already commited file so change the old file's name
	  fmt.Println("\t\t>> Already Versioned File Exists so Updating Version...")
	  err :=  os.Rename(fileToBeCommited, fmt.Sprintf("%sv%s",fileToBeCommited, strconv.Itoa(i)))

      if err != nil {
        // fmt.Println(err)
      } else {
      	i++
      	mapReadFromJsonFile["version"] = strconv.Itoa(i)
				configFile, err := os.OpenFile(fmt.Sprintf("%s/lgconfig.json", DEFAULTREPONAME), os.O_RDWR, 0600)
				fmt.Println("\t\t>> Current Version After Commit is ", mapReadFromJsonFile["version"])
  			// json.NewEncoder(configFile).Encode(mapReadFromJsonFile)
  			mapEncoded, _ := json.Marshal(mapReadFromJsonFile)
  			configFile.Write(mapEncoded)
  			configFile.Close()

				if err != nil{
					// fmt.Println("\t>> error when updating commit version ", err)					
				}
      }
	}
	newFile, err := os.Create(fileToBeCommited)
	
	if err != nil {
		// fmt.Println("personal error", err)
	}
	defer newFile.Close()
	var receivedBytes int64
	
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("\t>> Server Recieved File Successfully")
}