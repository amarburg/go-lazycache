package stress_test

import (
  "flag"
  "fmt"
  "net/http"
  "encoding/json"
  "math/rand"
  "sync"
)

var FrameStressCount = flag.Int("frame-stress-count", 50, "Number of frame stress queries to make")
var FrameStressParallelism = flag.Int("frame-stress-parallelism", 5, "Parallelism of frame stress queries")

var wg sync.WaitGroup

func FrameStress( imgUrl string ) error {

  resp,err := http.Get( imgUrl )
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  imgInfo := json.NewDecoder( resp.Body )


  var videoInfo struct {
    NumFrames   int
  }
  err = imgInfo.Decode( &videoInfo )

  if err != nil {
    panic( fmt.Sprintf("Couldn't figure out number of frames: %s", err.Error()))
  } else if videoInfo.NumFrames < 1  {
    panic("Couldn't figure out number of frames")
  }
  fmt.Println("Video has %d frames", videoInfo.NumFrames)

  var urls = make(chan string )

  count := *FrameStressCount
  parallelism := *FrameStressParallelism

  wg.Add( parallelism )

	for i := 0; i < parallelism; i++ {
		go FrameStressWorker(urls)
	}

  for i := 0; i < count; i ++ {
    urls <- fmt.Sprintf("%s/frame/%d", imgUrl, rand.Intn( videoInfo.NumFrames) )
  }

  close(urls)
  fmt.Printf("Waiting for workers to finish...")
  wg.Wait()
  fmt.Printf("done\n")

  return nil
}

func FrameStressWorker(urls chan string) {
	fmt.Println("In random walker")
	for {

    url,ok := <- urls

    if !ok {
      fmt.Println("Channel closed, quitting")
      wg.Done()
      return
    }

		fmt.Println("Random walker Querying URL", url)

		_, err := http.Get(url)
		if err != nil {
			fmt.Printf("%d: ERROR: %s\n", url, err)
			fmt.Printf("Error making request: %s\n", err.Error())
		}


    //
		// defer resp.Body.Close()
		// _, _ := ioutil.ReadAll(resp.Body)


	}

}
