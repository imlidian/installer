
package main  
  
import (  
	"bytes"
	"os"	
	// "encoding/json"  
	//"mime/multipart"	
    "log"  
    "net/http"  
	"fmt"  
	//"io"
	"flag"
	
)  

var tag ="latest"
var registry="hub.opshub.sh"
var namespace ="containerops"
var api_url ="https://build.opshub.sh/assembling/build?"
func main() {
	//url := "https://build.opshub.sh/assembling/build?image=test-java-gradle-testng&tag=latest&registry=hub.opshub.sh&namespace=containerops"


    //port := flag.String("port", ":8080", "http listen port")
    var image string
    flag.StringVar(&image, "image", "test_image", "image name")
 
    flag.Parse()
 
  //  fmt.Println("port:", *port)
	fmt.Println("image:", image)
	 imagename:= image
	 url := api_url
     buf := bytes.NewBufferString(url)
	 buf.Write([]byte("image="))
	 buf.Write([]byte(imagename))
	 buf.Write([]byte("&tag="))
	 buf.Write([]byte(tag))
	 buf.Write([]byte("&tag="))
	 buf.Write([]byte(tag))
	 buf.Write([]byte("&registry="))
	 buf.Write([]byte(registry))
	 buf.Write([]byte("&namespace="))
	 buf.Write([]byte(namespace))
	 fmt.Println(buf.String()) //hello roc
		UploadBinaryFile("./coala.tar",buf.String())
	}
// Upload binary file to the Dockyard service.
func UploadBinaryFile(filePath, url string) error {
		
	if f, err := os.Open(filePath); err != nil {
		return err
	} else {
		defer f.Close()
		if req, err := http.NewRequest(http.MethodPost,
			url, f); err != nil {
			return err
		} else {
			req.Header.Set("Content-Type", "text/plain")

			client := &http.Client{}
			if resp, err := client.Do(req); err != nil {
				return err
			} else {
				defer resp.Body.Close()

				switch resp.StatusCode {
				case http.StatusOK:
					{
						body := &bytes.Buffer{}
						_, err := body.ReadFrom(resp.Body)
					 	 if err != nil {
							log.Fatal(err)
						}
					 	resp.Body.Close()
						fmt.Println(resp.StatusCode)
						fmt.Println(resp.Header)
						fmt.Println(body)
						return nil						
					}
				case http.StatusBadRequest:
					return fmt.Errorf("Binary upload failed.")
				case http.StatusUnauthorized:
					return fmt.Errorf("Action unauthorized.")
				default:
					return fmt.Errorf("Unknown error.")
				}
			
			}
		}
	}

	return nil
}

<<<<<<< HEAD
	//init()
	func init() {
	
	}
	
	//main()
	func main() {
	url := "https://build.opshub.sh/assembling/build?image=test-java-gradle-testng&tag=latest&registry=hub.opshub.sh&namespace=containerops"
	UploadBinaryFile("./checkstyle.tar",url)
		/////////////////////////////////
		// if err := cmd.RootCmd.Execute(); err != nil {
		// 	fmt.Fprintf(os.Stderr, err.Error())
		// 	os.Exit(1)
		// }
//https://build.opshub.sh/assembling/build?image=test-java-gradle-testng&tag=latest&registry=hub.opshub.sh&namespace=containerops
	// extraParams := map[string]string{
	// 	"image":       "testtestng",
	// 	"tag":      "latest",
	// 	"registry": "hub.opshub.sh",
	// 	"namespace":"containerops",
	// }
	//request, err := newfileUploadRequest("https://google.com/upload", extraParams, "file", "/tmp/doc.pdf")	
	//url := "https://build.opshub.sh/assembling/build"
	//url := "https://build.opshub.sh/assembling/build?image=test-java-gradle-testng&tag=latest&registry=hub.opshub.sh&namespace=containerops"
	
	// resp ,err :=newfileUploadRequest(url,extraParams,"file", "checkstyle.tar")

	// defer resp.Body.Close()
	// 		b, err := ioutil.ReadAll(resp.Body)
	// 		if err != nil {
	// 			log.Println("http.Do failed,[err=%s][url=%s]", err, url,b)
	// 		}

///////////////////////////////////////////////////////
	// extraParams := map[string]string{
	// 	"image":       "testtestng",
	// 	"tag":      "latest",
	// 	"registry": "hub.opshub.sh",
	// 	"namespace":"containerops",
	// }
	// 	UploadBinaryFile(extraParams);

		/////////////////
			// request, err := newfileUploadRequest(url, extraParams, "file","./cpd.tar")
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// client := &http.Client{}
			// resp, err := client.Do(request)
			// if err != nil {
			// 	log.Fatal(err)
			// } else {
			// 	body := &bytes.Buffer{}
			// 	_, err := body.ReadFrom(resp.Body)
			//   if err != nil {
			// 		log.Fatal(err)
			// 	}
			//   resp.Body.Close()
			// 	fmt.Println(resp.StatusCode)
			// 	fmt.Println(resp.Header)
			// 	fmt.Println(body)
			// }
	}


	func newfileUploadhttp(params map[string]string){
		
		target_url := "https://build.opshub.sh/assembling/build"
		
	//	target_url := "http://localhost:9200/prod/aws"
		body_buf := bytes.NewBufferString("")
		body_writer := multipart.NewWriter(body_buf)
		jsonfile := "cpd.tar"
		file_writer, err := body_writer.CreateFormFile("upfile", jsonfile)

		for key, val := range params {
			_ = body_writer.WriteField(key, val)
		}

		if err != nil {
		fmt.Println("error writing to buffer")
		return
		}
		fh, err := os.Open(jsonfile)
		if err != nil {
		fmt.Println("error opening file")
		return
		}
		io.Copy(file_writer, fh)
		body_writer.Close()
		//http.Post(target_url, "application/json", body_buf)
		
		resp, err := http.Post(target_url, "application/x-tar", body_buf)
		if err != nil {
			log.Fatal(err)
		} else {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(resp.Body)
		  if err != nil {
				log.Fatal(err)
			}
		  resp.Body.Close()
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Header)
			fmt.Println(body)
		}

		return 
	}
	
=======
>>>>>>> 734a33ff6065febc90a200cbd4acb08f6c3df2a0

	
