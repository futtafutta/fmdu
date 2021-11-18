module insertFmdu

go 1.16

replace bitbucket.org/inadenoita/futalib => ../../../futalib

replace bitbucket.org/inadenoita/fmdu => ../../../fmdu

require (
	bitbucket.org/inadenoita/fmdu v0.0.0
	bitbucket.org/inadenoita/futalib v0.0.0
	github.com/labstack/echo/v4 v4.3.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
)

