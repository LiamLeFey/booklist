All you need to do is put the files in a directory and run
'go test'.
Couple things:
You need docker installed and it is assumed that dockerd is running.
The user who is running go test needs to be in the docker group. (or sudo go test)
The first time it is run, it might take a while, if the golang docker image is not downloaded.
I've tried to keep a rough history with the commits.
There's still one more test I want to write and I might push it later now that I've figured out the channel/t.Run incompatibility (see comments towards end of booklist_test.go
