package main

//convert path parameters to fileserver data path
func makepath(hard, app, file string) string {
	return "/srv/fileserver/" + hard + "/" + app + "/" + file
}
