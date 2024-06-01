# Gutendex Library API
This was originally the first Assignmnet for [PROG2005 - Cloud Technologies](https://www.ntnu.edu/studies/courses/PROG2005#tab=omEmnet), at NTNU campus Gjøvik.
The original assignment code is hosted at on [Render](https://prog2005-cloud.onrender.com/) (it is using the free plan, so loading may take up to a minute).

The service communicates with several other APIs to retrieve information. The information depends on the endpoint provided in the URL by the user. The service queries the Gutenberg Library, REST Countries, and Language2Country APIs (self-hosted by teacher Chrisopher Frantz). There are three endpoints the service serves, which supply the user with different information. The endpoints are: `bookcount`, `readership`, and `status`. 

The service does not support any interactive front end, and Postman queries are required to serve endpoints. 

## bookcount
Bookhandler recieves an url in the form of `https://prog2005-cloud.onrender.com/librarystats/v1/bookhcount/?language={ISO_Code}`<br><br>
The `{ISO_Code}` of the parameter can be any 2-letter ISO code for a country. It extracts this information and queries the Gutenberg Library API to find the relevant information about this language.
<br><br>
If the url provided does not contain a parameter it throws and error asking for a language. The same message is provided if we are at the root `bookcount/` and if the parameters are empty `bookcount/?language=`.<br>
The handler function only handles GET requests and will provide a message to the client if the request is of another method. 
The handler will read from the Gutenberg library and provide the amount of books in a given language, how many unique authors have authored the books, and the fraction of books of the given language of the total amount of books.
 
## readership
Readershiphandler recieves a url in the form of `https://prog2005-cloud.onrender.com/librarystats/v1/readership/{ISO_code}/?{limit}`
<br><br>
The `{ISO_Code}` is a 2-letter ISO Code for a given country. The {limit} is the limit of the output specified by the client. <br>
If the `{ISO_Code}` is empty the server will throw an error and ask for a country code from the client. `{limit}` may be empty, and the server will supply the amount of countries it found for that language. 
<br><br>
The server queries the Gutenberg Library like the `bookcountHandler`and provides the same information of the number of books and authors for a given language. Additionally, it queries the Language2Country API to find the number of countries where the given `ISO_Code` language is an official language.<br>
It then queries the REST Countries API and returns the population of the countries it finds individually to the client, and assumes this is the potential readership of books in this language.

## status
The status endpoint resides at the url address: `https://prog2005-cloud.onrender.com/librarystats/v1/status/` <br><br>
Statushandler retrieves the status codes from all the APIs the service relies on, along with the version and the uptime of the service.
<br>
The errors in the endpoint are of log.Fatal so the program moves on. If there are errors on the GET requests we want to return the error message to the client and not return to the client because something failed. 
