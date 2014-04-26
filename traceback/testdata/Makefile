timeout.txt: timeout_test.go
	go test -v -tags=testdata -timeout=1s -run=TestTimeout 2>>$@ 1>>$@ || exit 0
timeout_early.txt: timeout_test.go
	go test -v -tags=testdata -timeout=1us -run=TestTimeout 2>>$@ 1>>$@ || exit 0
sigabrt.txt: sigabrt_test.go
	go test -v -tags="testdata signal" -run=TestSIGABRT 2>>$@ 1>> $@ || exit 0
deadlock.txt: deadlock_test.go
	go test -v -tags=testdata -run=TestDeadlock 2>>$@ 1>> $@ || exit 0
