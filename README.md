# go-lazycache

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/aedd6b5087264d6d8b9e509a99ce6827)](https://www.codacy.com/app/amarburg/go-lazycache?utm_source=github.com&utm_medium=referral&utm_content=amarburg/go-lazycache&utm_campaign=badger)
[![wercker status](https://app.wercker.com/status/34ac5716d8bd050db14e85b8d35b648a/s/master "wercker status")](https://app.wercker.com/project/byKey/34ac5716d8bd050db14e85b8d35b648a)

[![Go Report Card](https://goreportcard.com/badge/github.com/amarburg/go-lazycache)](https://goreportcard.com/report/github.com/amarburg/go-lazycache)



## Todo

[ ] Add integration test for extraction of .png, .jpg and no extension...


## Benchmarking using curl

    repeat 10 { curl -s -o /dev/null -w "%{time_total}\n" -H "Range: bytes=2615776240-2616368015" https://rawdata.oceanobservatories.org/files//RS03ASHS/PN03B/06-CAMHDA301/2017/09/21/CAMHDA301-20170921T211500.mov }
