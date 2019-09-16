package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"bufio"
	"strconv"
	"strings"
)

const DEFAULTREPONAME = "repo_new"

const BUFFERSIZE = 1024
const CMDCOMMIT = "lg -commit"
const CMDUPDATE = "lg -update"
const CMDCREATE = "lg -create"
const CMDBACKUP = "lg -backup"
const CMDEXIT =   "lg -logout"
const CMDUPTOV =   "lg -uptov"

func main() {
	// connection, err := net.Dial("tcp", "10.5.12.114:9999")
	connection, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		// fmt.Println("personal error", err)
	}
	serverResponse := make([]byte, 36)
	fmt.Println("\n\n\t***************Welcome**********************")
	
	connection.Write([]byte(connection.LocalAddr().String()))
	//Read rejection or acceptance of server response
	connection.Read(serverResponse)
	fmt.Println("\t", string(serverResponse))

	fmt.Print("\t>> ")
	//wait for user command
	reader := bufio.NewReader(os.Stdin)
	cmd, _ := reader.ReadString('\n')
	if strings.Contains(cmd, CMDCOMMIT){
		fmt.Println("\t\t>> commit command executed")
		connection.Write([]byte(CMDCOMMIT))
		sendFileToClient(connection)

	} else if strings.Contains(cmd, CMDCREATE){
		// if create command is given
		// the server will accept 20 bytes of data
		// so that the new project will have a name and folder in the server
		fmt.Println("\t\t>> creating repo command executed")
		connection.Write([]byte(CMDCREATE))
		fmt.Print("\t\t>> Please Enter New Repository Name: ")
		
		newProjectNameReader := bufio.NewReader(os.Stdin)
		newProjectName := make([]byte, 20)
		newProjectNameFromInput,_ := newProjectNameReader.ReadString('\n')
		newProjectNameFromInput = fillString(string(newProjectNameFromInput), 20)
		newProjectName = []byte(newProjectNameFromInput)
		//Write 20 bytes to the connection
		connection.Write([]byte(newProjectName))
		
		fmt.Print("\t>> ")
		newProjectNameReader2 := bufio.NewReader(os.Stdin)
		newProjectNameFromInput2,_ := newProjectNameReader2.ReadString('\n')

		if strings.Contains(newProjectNameFromInput2, CMDCOMMIT){
			connection.Write([]byte(CMDCOMMIT))
			sendFileToClient(connection)
			fmt.Println("\t\t>> File Commited Successfully")
		}
	} else if strings.Contains(cmd, CMDBACKUP){
		fmt.Println("\t\t>> Backing up and Replicating....")
		connection.Write([]byte(CMDBACKUP))
		err := os.Mkdir(DEFAULTREPONAME, 0x777)
		if err != nil {
			// fmt.Println("\t>> error | back up | creating folder ", err)					
		}
	
		recieveFile(connection, true)

	} else if strings.Contains(cmd, CMDUPDATE){
		fmt.Println("\t\t>> updating version command executed")
		connection.Write([]byte(CMDUPDATE))
		recieveFile(connection, false)
	} else if strings.Contains(cmd, CMDUPTOV){
		fmt.Println("\t\t>> rollback command executed")
		connection.Write([]byte(cmd))
		recieveFile(connection, false)
	} else if strings.Contains(cmd, CMDEXIT){
		fmt.Println("\t\t>> exiting")
		connection.Write([]byte(CMDEXIT))
		connection.Close()
	} 
	
}

func recieveFile(connection net.Conn, isFromBackUp bool){
	
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	
	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	
	fileToBeCreatedPath := fileName

	if(isFromBackUp){
		fileToBeCreatedPath = fmt.Sprintf("%s/%s", DEFAULTREPONAME, fileName)
	}
	newFile, err := os.Create(fileToBeCreatedPath)
	
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
	fmt.Println("\t\t>> Update Successfull")
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

func sendFileToClient(connection net.Conn) {
	file, err := os.Open("polymer.js")
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
	fmt.Println("\t\t>> File Sent To Server Successfully")
	return
}