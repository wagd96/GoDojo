package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
)

type Args struct {
	A, B int
}

type Response struct {
	Quo, Res int
}

func checkError(err error, on string) {
	if err != nil {
		fmt.Printf("Error on %v: %v\n", on, err)
	}
}

func setValue(reader *bufio.Reader, message *string) (value int) {
	fmt.Printf(*message)
	bytes, _, err := reader.ReadLine()
	checkError(err, "Reading Line")
	stringData := string(bytes)
	temp, err := strconv.ParseInt(stringData, 10, 32)
	checkError(err, "Parsing Int")
	value = int(temp)
	return
}

func setArgs(reader *bufio.Reader, args *Args) {
	var message = "Type a number: "
	args.A = setValue(reader, &message)
	args.B = setValue(reader, &message)
}

func initSlice(reader *bufio.Reader) []int {
	message := "Type slice size: "
	size := setValue(reader, &message)
	slice := make([]int, size)
	for i := range slice {
		var localMessage = fmt.Sprintf("Type the number at position %d : ", i)
		slice[i] = setValue(reader, &localMessage)
	}
	return slice
}

func majorOrMinor(address *string, reader *bufio.Reader, operation *string) {
	var reply int
	client, err := rpc.Dial("tcp", *address)
	checkError(err, "Dialing")
	if err != nil {
		return
	}
	defer client.Close()
	slice := initSlice(reader)
	if err = client.Call(*operation, slice, &reply); err == nil {
		fmt.Printf("%v value in %v is %d\n", *operation, slice, reply)
	}
	checkError(err, fmt.Sprintf("%v Call", *operation))
}

func printMenu() {
	fmt.Printf("\t\tMenu\n")
	fmt.Printf("1. Add\n")
	fmt.Printf("2. Divide\n")
	fmt.Printf("3. Major\n")
	fmt.Printf("4. Minor\n")
	fmt.Printf("5. Exit\n")

	fmt.Printf("Type a option: ")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("A Error has happen and I Survived, the error was: ", err)
		}
	}()
	if len(os.Args) != 2 {
		fmt.Println("Put a server port")
		os.Exit(1)
	}
	var service = os.Args[1]
	var reader = bufio.NewReader(os.Stdin)

	client, err := rpc.Dial("tcp", service)
	checkError(err, "Dialing")

	// Synchronous Call
	if err == nil {
		message := ""
	out:
		for {
			printMenu()
			option := setValue(reader, &message)
			switch option {
			case 1:
				var reply int
				var args = Args{}
				setArgs(reader, &args)
				if err = client.Call("Math.Add", args, &reply); err == nil {
					fmt.Printf("Math %d + %d = %d\n", args.A, args.B, reply)
				}
				checkError(err, "Math.Add Call")
			case 2:
				response := Response{}
				var args = Args{}
				setArgs(reader, &args)
				if err = client.Call("Math.Divide", args, &response); err == nil {
					fmt.Printf("Math %d/%d = %d it's residue is %d \n", args.A, args.B, response.Quo, response.Res)
				}
				checkError(err, "Math.Divide Call")
			case 3:
				operation := "Math.Major"
				majorOrMinor(&service, reader, &operation)
			case 4:
				operation := "Math.Minor"
				majorOrMinor(&service, reader, &operation)
			case 5:
				fmt.Printf("Good Bye!\n")
				break out
			default:
				fmt.Printf("Invalid Option\n")
			}
		}
		defer client.Close()
	}
}
