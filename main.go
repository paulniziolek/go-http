package main

func main() {
	server := NewServer(":42069")
	server.ListenAndServe()

}
