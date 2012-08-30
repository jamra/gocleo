gocleo
======

##A golang implementation of the Cleo search.

The Cleo search is explained here: [Linked in original article](http://engineering.linkedin.com/open-source/cleo-open-source-technology-behind-linkedins-typeahead-search)

The source for Jingwei Wu's version can be found here: [Jingwei's version](https://github.com/linkedin/cleo)

###Dependencies
[gorilla mux library](http://gorilla-web.appspot.com/pkg/mux)

###Instructions
  package main
  import "cleo"
  
  func main(){
   cleo.InitAndRun("w1_fixed.txt", "8080")
  }
Run the program and navigate to localhost:8080/cleo/{query}

{query} is your search.  e.g.("tractor", "nightingale", "pizza")

###TODO:  
 - Give a better explanation of the code.  
 - Split the web portion into a different file.  Perhaps "cleo_test.go".  