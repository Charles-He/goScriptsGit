package main

import (
        "github.com/tarm/serial"
		"fmt"
		"os"
		"log"
		"regexp"
		"bufio"
		// "runtime"
		"sync"
		// "bytes"
		// "reflect"
)

// runtime.GOMAXPROCS(2)
var wg sync.WaitGroup

// for serial data
// var serialData = []byte{}

func main() {
	
    wg.Add(2)


	// capture "port" and initial parameter from SCID.
	port := os.Args[1]
	parameter := os.Args[2]

	
	
	// open Serial Port
	// c := &serial.Config{Name: port, Baud: 115200, ReadTimeout: time.Second * 5}
	c := &serial.Config{Name: port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	
	// send the init pararmeter to board
	_, err = s.Write([]byte(parameter))
	if err != nil {
		log.Fatal(err)
	}
	
	// process serial port
	serialChan := make(chan []byte) 
	
	go func () {
		defer wg.Done()
		
		messagebyte := make([]byte, 128)
		
		go processSerialData(serialChan)
		
		for {
	
			n, err := s.Read(messagebyte)
			if err != nil {
				log.Fatal(err)
			}
			
			serialChan <- messagebyte[:n]
			
		}
	}()
	
	// communication(s, serialPort, quit)
	
	// stdin read and send
	
	ch := make(chan string)
	go func(ch chan string) {
		defer wg.Done()
		
		// Uncomment this block to actually read from stdin
		reader := bufio.NewReader(os.Stdin)
		for {
			x, err := reader.ReadString('\n')
			if err != nil { // Maybe log non io.EOF errors, if you want
				close(ch)
				return
			}
			
			_, err = s.Write([]byte(x))
			if err != nil {
				log.Fatal(err)
			}
			
			// if receive "stop", then quit
			match, _ := regexp.MatchString("stop",x)
			if match == true {
				os.Exit(3)
			}

			// fmt.Println("okokok")
		}
		
		// Simulating stdin
		// ch <- "A line of text"
		close(ch)
	}(ch)
	
	// fmt.Println("Done, stdin must be closed")
	
	// runtime tail
	// fmt.Println("Waiting To Finish")
    wg.Wait()

    // fmt.Println("\nTerminating Program")
    
}

func processSerialData (serialChan chan []byte) {
	var serialData = []byte{}
	var currentMessage = []byte{}
	
	// debug
	// for _, n := range(messagebyte) {
	//	fmt.Printf("%q", n) 
	// }
			
	for {
		serialData = append (serialData, <- serialChan...)
		fmt.Printf("%q\n", serialData)
		
		for {
			// find the header, if no header then discard
			if len(serialData) > 3 && serialData[0] == 0xfe {
				// fmt.Println(len(serialData))
				// fmt.Printf("If: %q\n", serialData)
				length := int(serialData[2])
				if length+3 <= len(serialData) {
					for i:=0;i<=length-1;i++ {
						if serialData[3+i] != 0xfe {
							currentMessage = append (currentMessage, serialData[3+i])
						}
					}
					message := string(currentMessage)
					fmt.Printf("%s\n", message)
		
					cutLength := len(currentMessage)+3
					serialData = serialData[cutLength:]
					currentMessage = []byte{}
				} else {
					break
				}
			
		
			} else if len(serialData) > 0 && serialData[0] == 0xfe {
				// fmt.Printf("Else if: %q\n", serialData)
				break
			} else {
				// fmt.Printf("Else: %q\n", serialData)
				badMark := 0
				for i:=0;i<=len(serialData)-1;i++ {
						if serialData[i] == 0xfe {
							badMark = i
							break
						}
				}
				if badMark != 0 {
					serialData = append (serialData[badMark:])
				} else {
					serialData = []byte{}
					break
				}
			}
			
		}
	}
/*	
	if message[0] == 0xfe && len(message) > 3 {
		length := int(message[2])
				
		if length == n-3 {
			message = string(message[3:])
		} else if length < n-3 {
			message = string(message[3:length+2])
			message1 = string(messagebyte[3:n])
		}

	} else {
				message2 := string(message[:])
				message = message1 + message2
				message1 = ""
			}
			
			// fmt.Println("the message is:", message)
			if prem != message && message != "" {
				fmt.Printf("%s\n", message)
				
			}
			prem = message
			
	fmt.Printf("%s\n", message)
*/
}
	






