gocleo
======

##A golang implementation of the Cleo search.

The Cleo search is explained here: [Linked in original article](http://engineering.linkedin.com/open-source/cleo-open-source-technology-behind-linkedins-typeahead-search)

The source for Jingwei Wu's version can be found here: [Jingwei's version](https://github.com/linkedin/cleo)

Basically, this is a golang version of the original program.  The original program is written in Java.  I have included a corpus of words to search for.  I downloaded this corpus from http://www.wordfrequency.info/

###Algorithm overview
 - The algorithm starts out by searching for matches in the inverted index.  The inverted index contains a map of the word's prefix (up to 4 chars).  Each word prefix maps to an array of document ID, bloom filter tuples.  
 - The bloom filter of each candidate is compared against the query's bloom filter.  If it matches successfully, the candidate makes it to the next round.
 - The remaining words are scored by their [levenshtein distance](http://en.wikipedia.org/wiki/Levenshtein_distance) to the query, then normalized using the [Jaccard coefficient](http://en.wikipedia.org/wiki/Jaccard_index).
 - The final words are returned as JSON
 - You can also change how scoring works if you like. You just need to provide a function that conforms to
    func(s1, s2 string) (score float64)

###Instructions
This is a sample app:

    package main
   	import "github.com/jamra/gocleo"
  
   	func main(){
   	  cleo.InitAndRun("w1_fixed.txt", "8080", nil) //The last parameter is optional. Defaults to Levenshtein distance normalized by Jaccard coefficient
   	}

Run the program and navigate to localhost:8080/cleo/{query}

{query} is your search.  e.g.("tractor", "nightingale", "pizza")

###Your own corpus
You can have the search run off of your own corpus so long as each term is separated by a new line.  w1_fixed.txt is provided as an example.

###Setup
This should work with go get

    go get github.com/jamra/gocleo
###TODO
 - Give the user the ability to add and remove words from the index.
 - More robust Unit testing
