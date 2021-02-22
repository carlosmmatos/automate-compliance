module github.com/carlosmmatos/automate-compliance

go 1.15

require (
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/fatih/set v0.2.1 // indirect
	github.com/opencontrol/compliance-masonry v1.1.6
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/oauth2 v0.0.0-20210210192628-66670185b0cd
	google.golang.org/api v0.40.0
	vbom.ml/util/sortorder v0.0.0-00010101000000-000000000000 // indirect
)

replace vbom.ml/util/sortorder => github.com/fvbommel/sortorder v1.0.1
