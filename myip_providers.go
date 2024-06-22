package main

var (
	myIPProviders = []myIPProvider{
		newCloudflareMyIPProvider(),
		newIPAPIMyIPProvider(),
		newIPifyMyIPProvider(),
		newIPInfoMyIPProvider(),
	}
)
