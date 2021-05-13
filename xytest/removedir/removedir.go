package main 
  
import ( 
    "fmt"
//    "os"
    "os/exec"
) 
  
func main() { 
  
    // Remove all the directories and files, does not work
    // Using RemoveAll() function 
    /*
    err := os.RemoveAll("/tmp/data") 
    if err != nil { 
	fmt.Println(err.Error())
    }
    */
   
    rmCmd := exec.Command("sudo", "rm", "-rf", "/tmp/data") 
    _, err := rmCmd.Output()
    if err != nil {
        fmt.Println(err.Error())
    }
    /*
    rmCmd := exec.Command("sudo", "rm", "-rf", "/tmp/data")
    rmCmd.Stderr = os.Stderr
    rmCmd.Stdin = os.Stdin

    out, err := rmCmd.Output()
    if err != nil {
	fmt.Println(err.Error())
    } else {
	fmt.Println(string(out))
    }
    */
} 
