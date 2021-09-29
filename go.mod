module github.com/carlosmmatos/automate-compliance

go 1.15

require (
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/fatih/set v0.2.1 // indirect
	github.com/opencontrol/compliance-masonry v1.1.6
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	google.golang.org/api v0.58.0
	vbom.ml/util/sortorder v0.0.0-00010101000000-000000000000 // indirect
)

replace vbom.ml/util/sortorder => github.com/fvbommel/sortorder v1.0.1
