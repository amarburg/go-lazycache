package lazycache

import "testing"

import "github.com/amarburg/go-lazyfs-testfiles/http_server"

//import "fmt"

func TestHttpFS( t *testing.T ) {
  srv := lazyfs_testfiles_http_server.HttpServer()
  defer srv.Stop()

  fs, err := OpenHttpFS( srv.Url )

  if fs == nil || err != nil {
    t.Fatal("Couldn't create HttpFS", err.Error() )
  }

}



func TestHttpFSListing( t *testing.T ) {
  var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/"
  fs, err := OpenHttpFS( OOIRawDataRootURL )

  if fs == nil || err != nil {
    t.Fatal("Couldn't open HttpFS to OOI Raw Data Server", err.Error() )
  }

  testPath := "files/RS03ASHS/PN03B/06-CAMHDA301/2016/"

  listing, err := fs.ReadHttpDir( testPath )

  if err != nil {
    t.Fatal("Couldn't list directory on OOI Raw Data Server", err.Error() )
  }

  if( len( listing.Directories ) != 12) {
    t.Errorf("Didn't get 12 children from %s\n", testPath)
  }

  //fmt.Printf("root contains %d entries\n", len( listing.Children ))

}
