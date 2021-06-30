# Client_Test
Project, for load and view the products of a restaurant, this project works with go, dgraph and vue with html

# Dependencies 
All this project, it´s ejecuted on linux mint 2019 install Dgraph, API chi and VueJs 

For run dgraph database you need write on a terminal tthe follow instruction 
 $ docker run --rm -it -p "8080:8080" -p "9080:9080" -p "8000:8000" -v ~/dgraph:/dgraph "dgraph/standalone:v21.03.0"
 
you don0t forget init the module inside of the folder with instruction 
go mod init "namefolder" 
 
Open your navegator and access to localhost: 8000 
