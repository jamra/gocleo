gocleo
======

##A golang implementation of the Cleo search.

The Cleo search is explained here: [Linked in original article](http://engineering.linkedin.com/open-source/cleo-open-source-technology-behind-linkedins-typeahead-search)

The source for Jingwei Wu's version can be found here: [Jingwei's version](https://github.com/linkedin/cleo)

###Dependencies
[gorilla mux library](http://gorilla-web.appspot.com/pkg/mux)

###Instructions
This is a sample app:

      package main
   	import "cleo"
  
   	func main(){
   	  cleo.InitAndRun("w1_fixed.txt", "8080")
   	}

Run the program and navigate to localhost:8080/cleo/{query}

{query} is your search.  e.g.("tractor", "nightingale", "pizza")

###Your own corpus
You can have the search run off of your own corpus so long as each term is separated by a new line.  w1_fixed.txt is provided as an example.

###TODO
 - Add better configurability.  I want to make the scoring mechanism a passable function.