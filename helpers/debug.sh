git branch -D debug
wait 
git checkout -B debug
wait
go build -gcflags "all=-N -l"
wait 
dlv exec curgo
