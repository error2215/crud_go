# crud_go
Task description - create simple service for create/read/update/delete-ing items in Elasticsearch using REST API and create simple server-client service using gRPC and protobuf.
Every item has such fields:
  1) Id
  2) Date
  3) Data
  4) Description
  5) Author
  6) Title
    
crud_go/client/main.go - cliend side file, that makes request to server with id of item and receives item

crud/getData - protobuf and gRPC files

crud/serv - server side file, that has 3 fuctions:

    main() - main function, that receives requests from client
  
    GetDataById() - makes request to Elasticsearch client and returns item
  
    searchDataInES() - function is doing search in ES by id
  
templates/ - templates for CRUD service

/main.go - file with CRUD service, has 2 structs: 'Post' for create/update/delete functions and 'IndexPost' for viewing items.
