module example.com/postos-anp

go 1.15

require (
	entities v0.0.0
	funcs v0.0.0
	github.com/gocolly/colly/v2 v2.1.0
	mongo v0.0.0-00010101000000-000000000000
	params v0.0.0-00010101000000-000000000000
)

replace params => ./params

replace mongo => ./mongo

replace entities => ./entities

replace funcs => ./funcs
