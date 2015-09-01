package main

import(
	"fmt"
	"bufio"
//	"errors"
	"net"
	"log"
	"strings"
	"strconv"
	"io"
)

type FTP struct{

	conn net.Conn
	addr string
	debug  bool
	reader *bufio.Reader
	writer *bufio.Writer
}

func (ftp *FTP)Connect(host string,port int)(error){
	var err error
	ftp.addr = host
	address := fmt.Sprintf("%s:%d",host,port)
	if ftp.conn,err = net.Dial("tcp",address);err != nil{
		return err
	}
	ftp.writer = bufio.NewWriter(ftp.conn)
	ftp.reader = bufio.NewReader(ftp.conn)
	
	return nil;
}

func (ftp *FTP) Response()(Message string,err error){
	Message,err = ftp.reader.ReadString('\n')
	if ftp.debug{
		log.Printf("* from respon*[+]  %s",Message )
	}
	return 
	
}

func (ftp *FTP)Request(cmd string)(  error){
	cmd += "\r\n"
	
	if ftp.debug{
		log.Printf("*from rqst*[+]  %s",cmd)	
	}
	if _,err := ftp.writer.WriteString(cmd);err!=nil{
		return err	
	}
	if err := ftp.writer.Flush();err != nil{
		return err
	}
	
	return nil
}
func (ftp *FTP)Action(cmd string)( Message string ){
	
	 ftp.Request( cmd )
	Message,_ =  ftp.Response()
	return 

}
func (ftp *FTP)Login(username string,password string)bool{
	ftp.Action("USER " + username)
	ftp.Action("PASS " + password)
	return true
}
func (ftp *FTP)Pwd(){
	ftp.Action("PWD")
}
func (ftp *FTP)Exit(){
	ftp.Action("QUIT")
}
// Make FTP Server switching to passive mod
func (ftp *FTP)Passive()int{
	mess := ftp.Action("PASV")
	start := strings.Index(mess,"(")
	end :=  strings.Index(mess,")")
	pli := strings.Split(mess[start+1:end],",")
	l1,_ :=strconv.Atoi(pli[len(pli)-1])
	l2,_ := strconv.Atoi(pli[len(pli)-2])
	return l2*256 + l1
	
}
//Switch FTP server's mode
func (ftp *FTP)Type(mode string){
	ftp.Action("Type "+mode)
}

// return array about file,each element is file 's line
func (ftp *FTP)RetFl(filename string)(file [] string){
	ftp.Type("I")
	passp := ftp.Passive()
	ftp.Request("RETR "+filename)

	nwcon := ftp.NewCon(passp)
	ftp.Response()
	reader := bufio.NewReader(nwcon)
	for{
		line,err := reader.ReadString('\n')
		if err == io.EOF{
			break
		}else if err != nil{
			return
		}
		file = append(file,string(line))
	}
	nwcon.Close()
	ftp.Response()
	return 
}
// List all files and folders
func (ftp *FTP)List()(file []string){
	passp := ftp.Passive()
	ftp.Request("LIST")
	nwcon := ftp.NewCon(passp)
	ftp.Response()
	reader := bufio.NewReader(nwcon)
	ftp.Response()
	for{
	line,err := reader.ReadString('\n')
	if err == io.EOF{
		break
	}else if err != nil{
		return
	}
	file = append(file,string(line))
	nwcon.Close()
	}
	return 
}

// Switch into a new folder
func (ftp *FTP)Cd(path string){
	ftp.Action("CWD "+path)
}
func (ftp *FTP)NewCon(port int )(conn net.Conn){
	address := fmt.Sprintf("%s:%d",ftp.addr,port)
	if ftp.debug{
		log.Printf("new connect to %s ",address)
	}
	conn,_ = net.Dial("tcp",address)
	return
}
// Get infomation of FTP server
func (ftp *FTP)SerIf(){
	ftp.Action("SYST")
}


func (ftp *FTP)Rename(filename string,newname string){
	ftp.Action("RNFR "+filename)
	ftp.Action("RNTO "+newname)

}
func (ftp *FTP)Debug(value bool){
	ftp.debug = value
}
func main(){
	username := "username"
	password := "pass"
	ftp := new(FTP)
	ftp.Debug(true)
	ftp.Connect("192.168.0.1",21)
	message,_ := ftp.Response()
	fmt.Println(message)
	ftp.Login(username,password)
	ftp.Pwd()
	ff:= ftp.RetFl("test2.txt")
	list := ftp.List()
//	ftp.Cd("..")
//	ftp.Pwd()
//	ftp.Rename("testfolder/test0.txt","test0.txt")
	ftp.SerIf()
	fmt.Println(ff)
	fmt.Println(list)
	ftp.Exit()

}
