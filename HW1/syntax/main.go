package main

import "fmt"
import "errors"
import "sync"
import "time"

func main() {
  //syntax1()
  //syntax2()
  //syntax3()
  //syntax4()
  //syntax5()
  //syntax6()
  //syntax7()
  //syntax8()
  //syntax9()
  //syntax10()
  //syntax11()
  syntax12()
}


func syntax1() {
  var names [2]string
  primes := [6]int{2, 3, 5, 7, 11, 13}
  fmt.Printf("Hello, Syntax.\n")

  //print these arrays
  for _, value := range names {
    fmt.Printf("%s ", value)
  }
  fmt.Printf("\n")
  for _, value := range primes {
    fmt.Printf("%d ", value)
  }
}

func syntax2() {
  names := make([]string, 0)
  names = append(names, "sparky")
  for _, val := range names {
    fmt.Printf("%s ", val)
  }

  fmt.Printf("\n")

  names2 := make([]string, 8)
  length := len(names2)
  fmt.Printf("%d\n", length)
}

func syntax3() {
  m := make(map[string]int)
  m["ten"] = 10
  fmt.Println(m["ten"])

  ten := m["ten"]
  fmt.Println(ten)

  delete(m, "ten")
  fmt.Println(m["ten"]) // errors

  ten_check, ok := m["ten"]
  fmt.Println(ok) // false
  fmt.Println(ten_check) // 0

}

//The initialization of a struct with three members x, y, and z
type Coordinate struct {
	x int
	y int
	z int
}
//This methods implemented for Coordinate outputs all the contents of a Coordinate
func (c Coordinate) printCurrentPosition() {
  fmt.Printf("(%v, %v, %v)\n", c.x, c.y, c.z)
}

//This two methods are designed to emphasize that a struct can only
//be modified by calling a function/method/interface if its reference/pointer rather than a struct
//is passed as a parameter. That is because when you pass the struct, essentially you
//copy the struct and you pass that copied one to the function/method/interface. Modification on the 
//copied one does not change the value of the originate one.

//doing nothing
func (c Coordinate) addX(d int) {
    c.x += d
}
// yay!
func (c *Coordinate) addX2(d int) {
    c.x += d
}

func syntax4() {
  //create an instance of Coordinate
  origin := Coordinate{x: 0, y: 0, z: 0}
  //print all the content of origin
	fmt.Printf("%+v \n", origin)

  //create another instance of Coordinate but
  //only print the content of member z
  pos2 := Coordinate{x: 12, y: 24, z: 13}
  fmt.Printf("%v \n", pos2.z) // prints 13

  //create a pointer which points to an instance of Coordinate
  //and print its content of member y
  pos3 := &Coordinate{x: 12, y: 24, z: 13}
  fmt.Printf("%v \n", pos3.y) // prints 24

  //create an instance of Coordinate and test the function printCurrentPosition
  pos4 := Coordinate{x: 12, y: 24, z: 13}
  pos4.printCurrentPosition() // prints "(12, 24, 13)"
  

  //verificate that only passing the reference/pointer enables you to 
  //modify a struct in functions/methods/interfaces
  pos4.addX(1)  
  pos4.printCurrentPosition()
  pos4.addX2(1)
  pos4.printCurrentPosition()

}

//An interface can be values of any types which have implemented
//the involved methods. Specifically, Duck is an interface including
//methods walk and talk. Therefore, duck can be any datatypes which 
//have implemented method walk and talk
type Duck interface {
    walk()
    talk() string
}


//This a sample Struct for further illustration
type Number struct {
  value int
}
func syntax5() {
  //anyMap is a map with string as key and empty interface as value. 
  //Recall that an interface can be any datatypes which have implemented 
  //the declaimed methods. As an empty interface does not implement any
  //method, it can be of any datatype.
  anyMap := make(map[string]interface{}) 
  //It can store an integer
  anyMap["one"] = 1
  //It can store a string
  anyMap["two"] = "two"
  //it can also store a struct
  anyMap["three"] = Number{value: 3}


  //You can transfer a variable from an empty interface to a specific datatype 
  //by calling "i.(type)", where i represents the variable to be transferred, 
  //and the type represents the datatype you want to transfer to.
  one := anyMap["one"].(int)
  fmt.Printf("%d", one)

  //But you cannot transfer an interface variable to a datatype which does not
  //match its exact value. For instance, the exact of anyMap["one"] is integer,
  //transferring it to string can cause system to panic.
  //one = anyMap["one"].(string) // This will panic!
}


//Error enables you to handle systemetic errors and facilitates
//understanding of your codes. The following function shows that 
//how to handle errors. This function does not expect the input
//to be zero, so it regards zero as an invalid input. That is, when 
//the input variable is zero, the function returns an error to notify
//caller to process the error fixing.
func mightFail(input int) (int, error) {
    if input == 0 {
        //when the input is zero, return -1 and an error,
        //where the -1 means nothing
        return -1, errors.New("Can't use 0")
    } else {
        //when the input is not zero, executes normally
        //and return nil as error, meaning that there
        //is no error.
        return 10 / input, nil
    }
}

func syntax6() {
    result, err := mightFail(0)
    //the caller check whether any errors are returned
    if err != nil {
        //if an error is returned, print the debug message
        fmt.Printf("Error!")
    } else {
        //if an error is not returned, executes normally
      fmt.Println(result)
    }
}

//A simple function that prints Hello, Golang!
func f() {
  fmt.Printf("Hello, Golang!\n")
}


func syntax7() {
  //create a Goroutine to run function f
  go f()
}

//The following code trys to send a message to a channel
//and receive messages from it in one scope simultaneously. This can not work 
//because you have to execute the sending and receiving in different threads
//in most cases. The default channels are unbuffered, meaning that
//they will only accept sends if there is a corresponding receive, vice versa. In other words,
//channel will block until it find a pair of matching sender and receiver from different threads. 
//The following code will cause a deadlock.
func syntax8() {
  //Initialize a unbuffered channel
  ch := make(chan int)
  v1 := 1
  //Send a message to the channle,
  //this channel will block until it finds a receiver from another thread
  ch <- v1
  //The following code will not be executed
  v2 := <- ch
  fmt.Printf("%d\n", v2)
}

//An appropiate way to use channel is performing sending and receiving
//messages simultaneously in different threads.

//The following function has two parameters, an integer array s and a channel c,
//aiming to calulate the sum of all items in s, and
//send the result to channel c.
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}
//The following function delegates two Goroutines to 
//compute the sum of an integer array, with each Goroutine
//handles a part of the job and subsequently combine the results.
//Different from the previous code, the following code has different 
//threads to send and receive message simultaneously. Specifically,
//two Goroutines send messages to the channel, and the main function 
//receives the messages from the channel. This scheme enables sending and 
//receiving to perform simultaneously, which will not block the process.
func syntax9() {
  s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c

	fmt.Println(x, y, x+y)
}


//If you still want to handle sending and receiving in one scope,
//you can try buffered channel. Different from an unbuffered channel, 
//a buffered channel does not require finding a pair of matching 
//sender and receivers. Instead, it store a messsage in a queue when
//someone sends messages to it, and return the first element in the front if 
//someone receives messages from it. This scheme does not block the thread.
func syntax10() {
  //Initialize a buffered channel
  ch := make(chan int, 100)
  v1 := 1
  //Send a message to the channle,
  ch <- v1
  //ok is ture if the channel is not empty and a message is returned
  //ok is false if the channel is empty and no message is returned
  v2, ok := <- ch
  if ok {
    fmt.Printf("%d\n", v2)
  } else {
    fmt.Printf("Empty channel")
  }
}

//Mutex ensures the atomic modification of a variable cross different Goroutines.
//That is, only one routines can access a variable at a time. The insight of Mutex
//is to lock the variable first, and release it after reading or modifying operations.

//The following function aims to increase the input parameter v by 1 and then send it 
//to a buffered channel for main function to output the result. 
func increaseWithMutexLock(v *int, lock *sync.Mutex, ch chan int) {
	lock.Lock()//the lock is taken by current scope
  //the following code can only be executed by one routine.
  res := *v + 1
  //send the result to the channel
  ch <- res
  time.Sleep(1000 * time.Millisecond)
  *v = res
	lock.Unlock()//the lock is released
}

//The test demonstrates that with Mutex, we can delegate 10 routines to increase 
//the variable a from 1 to 12 correctly. Only one routine can modify a at a time,
//which ensures the atomic modification of a. The results will output integers 
//from 2 to 12 gradually.
func syntax11() {
  a := 1
  //initialize a Mutex lock using library sync
  var lock sync.Mutex
  //initialize a channel to output the results
  ch := make(chan int, 100)
  //create 10 routines to increate a one by one
  for i := 0; i <= 10; i++ {
    go increaseWithMutexLock(&a, &lock, ch)
  }
  //output the results received from channel
  for i := 0; i <= 10; i++ {
    res := <- ch
    fmt.Printf("%d\n", res)
  }
}


//Without Mutex, serveral routines may increase the variable simultaneously,
//leading to an incorrect result.
func increaseWithoutMutexLock(v *int, ch chan int) {
  res := *v + 1
  ch <- res
  time.Sleep(1000 * time.Millisecond)
  *v = res
}

func syntax12() {
  a := 1
  ch := make(chan int, 100)
  for i := 0; i <= 10; i++ {
    go increaseWithoutMutexLock(&a, ch)
  }
  for i := 0; i <= 10; i++ {
    res := <- ch
    fmt.Printf("%d\n", res)
  }
}
