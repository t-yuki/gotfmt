timeout.txt: timeout_test.go
	go test -v -tags=testdata -timeout=1s -run=TestTimeout 2>>$@ 1>>$@ || exit 0
timeout_early.txt: timeout_test.go
	go test -v -tags=testdata -timeout=1us -run=TestTimeout 2>>$@ 1>>$@ || exit 0
sigabrt.txt: sigabrt_test.go
	go test -v -tags="testdata signal" -run=TestSIGABRT 2>>$@ 1>> $@ || exit 0
deadlock.txt: deadlock_test.go
	go test -v -tags=testdata -run=TestDeadlock 2>>$@ 1>> $@ || exit 0
nil.txt: nil_test.go
	go test -v -tags=testdata -run=TestNil 2>>$@ 1>> $@ || exit 0
goroutines100.txt: goroutines100_test.go
	go test -v -tags=testdata -run=TestGoroutines100 2>>$@ 1>> $@ || exit 0
race.txt: race.go
	go run -race race.go 2>>$@ 1>> $@ || exit 0
race_test.txt: race_test.go
	go test -race -v -tags=testdata -run=TestRace 2>>$@ 1>> $@ || exit 0
