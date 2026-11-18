module github.com/ivmerk/goshop/iam

go 1.27.3

require (
	github.com/ivmerk/goshop/platform v0.0.0
	github.com/ivmerk/goshop/shared v0.0.0
)

replace github.com/ivmerk/goshop/shared => ../shared

replace github.com/ivmerk/goshop/platform => ../platform
